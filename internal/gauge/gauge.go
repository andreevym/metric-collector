package gauge

import (
	"strconv"

	"github.com/andreevym/metric-collector/internal/storage"
)

func Store(s storage.Storage, metricName string, metricValue string) (string, error) {
	err := s.Create(metricName, metricValue)
	return metricValue, err
}

func StoreAll(s storage.Storage, kvMap map[string]*storage.Metric) error {
	if len(kvMap) != 0 {
		return s.CreateAll(kvMap)
	}
	return nil
}

func Validate(metricValue string) error {
	_, err := strconv.ParseFloat(metricValue, 64)
	return err
}

func Get(s storage.Storage, metricName string) (string, error) {
	v, err := s.Read(metricName)
	if err != nil {
		return "", err
	}

	return v, nil
}
