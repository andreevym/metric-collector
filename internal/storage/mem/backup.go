package mem

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/avast/retry-go"
)

func Load(filename string) (map[string]*storage.Metric, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return map[string]*storage.Metric{}, nil
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, fmt.Errorf("can't read backup file %s: %w", filename, err)
	}
	if len(bytes) == 0 {
		return map[string]*storage.Metric{}, nil
	}

	m, err := marshal(bytes)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, fmt.Errorf("can't marshal data from backup file %s: %w", filename, err)
	}
	return m, nil
}

func Save(filename string, data map[string]*storage.Metric) error {
	var file *os.File
	var err error
	_ = retry.Do(
		func() error {
			file, err = os.Create(filename)
			// обработка ошибки доступа к файлу, который был заблокирован другим процессом.
			if err != nil && os.IsPermission(err) {
				return err
			}
			return nil
		},
	)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	bytes, err := unmarshal(data)
	if err != nil {
		logger.Logger().Error(err.Error())
		return fmt.Errorf("can't unmarshal data for backup %w", err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		logger.Logger().Error(err.Error())
		return fmt.Errorf("can't write file backup %w", err)
	}

	// Close the file when done
	return file.Close()
}

type metricStore struct {
	Metrics map[string]*storage.Metric `json:"storage"`
}

func marshal(data []byte) (map[string]*storage.Metric, error) {
	v := metricStore{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}
	return v.Metrics, nil
}

func unmarshal(m map[string]*storage.Metric) ([]byte, error) {
	return json.Marshal(metricStore{m})
}
