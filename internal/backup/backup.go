package backup

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(filename string) (map[string][]string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't read backup file %s: %w", filename, err)
	}

	m, err := marshal(bytes)
	if err != nil {
		return nil, fmt.Errorf("can't marshal data from backup file %s: %w", filename, err)
	}
	return m, err
}

func Save(filename string, data map[string][]string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	bytes, err := unmarshal(data)
	if err != nil {
		return fmt.Errorf("can't unmarshal data for backup %w", err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("can't write file backup %w", err)
	}
	return err
}

type file struct {
	Storage map[string][]string `json:"storage"`
}

func marshal(data []byte) (map[string][]string, error) {
	v := file{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v.Storage, nil
}

func unmarshal(m map[string][]string) ([]byte, error) {
	return json.Marshal(file{m})
}
