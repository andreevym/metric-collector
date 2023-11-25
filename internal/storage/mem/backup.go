package mem

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
)

func Load(filename string) (map[string]*storage.Metric, error) {
	_, err := os.Stat(filename)
	if err != nil {
		logger.Log.Error(err.Error())
		return map[string]*storage.Metric{}, nil
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, fmt.Errorf("can't read backup file %s: %w", filename, err)
	}
	if len(bytes) == 0 {
		return map[string]*storage.Metric{}, nil
	}

	m, err := marshal(bytes)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, fmt.Errorf("can't marshal data from backup file %s: %w", filename, err)
	}
	return m, nil
}

func Save(filename string, data map[string]*storage.Metric) error {
	file, err := os.Create(filename)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	bytes, err := unmarshal(data)
	if err != nil {
		logger.Log.Error(err.Error())
		return fmt.Errorf("can't unmarshal data for backup %w", err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		logger.Log.Error(err.Error())
		return fmt.Errorf("can't write file backup %w", err)
	}
	return nil
}

type file struct {
	Storage map[string]*storage.Metric `json:"storage"`
}

func marshal(data []byte) (map[string]*storage.Metric, error) {
	v := file{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return v.Storage, nil
}

func unmarshal(m map[string]*storage.Metric) ([]byte, error) {
	return json.Marshal(file{m})
}
