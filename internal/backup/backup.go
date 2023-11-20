package backup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andreevym/metric-collector/internal/logger"
)

func Load(filename string) (map[string][]string, error) {
	_, err := os.Stat(filename)
	if err != nil {
		logger.Log.Error(err.Error())
		return map[string][]string{}, nil
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, fmt.Errorf("can't read backup file %s: %w", filename, err)
	}
	if len(bytes) == 0 {
		return map[string][]string{}, nil
	}

	m, err := marshal(bytes)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, fmt.Errorf("can't marshal data from backup file %s: %w", filename, err)
	}
	return m, nil
}

func Save(filename string, data map[string][]string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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
	Storage map[string][]string `json:"storage"`
}

func marshal(data []byte) (map[string][]string, error) {
	v := file{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return v.Storage, nil
}

func unmarshal(m map[string][]string) ([]byte, error) {
	return json.Marshal(file{m})
}
