package counter

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/andreevym/metric-collector/internal/storage"
)

func Store(s storage.Storage, metricName string, metricValue string) (string, error) {
	metricValues, err := s.Read(metricName)
	if err != nil && !errors.Is(err, storage.ErrValueNotFound) {
		return "", err
	}
	if len(metricValues) == 0 {
		err = s.Create(metricName, metricValue)
		if err != nil {
			return "", err
		}
		return metricValue, nil
	}

	existsMetricVal, err := strconv.ParseInt(metricValues, 10, 64)
	if err != nil {
		return "", err
	}
	v, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return "", err
	}
	newVal := strconv.FormatInt(existsMetricVal+v, 10)
	err = s.Update(metricName, newVal)
	if err != nil {
		return "", err
	}
	return newVal, err
}

func StoreAll(s storage.Storage, kvMap map[string]string) error {
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

	if len(v) == 0 {
		return "", fmt.Errorf("can't find metric by name %s", metricName)
	}

	return v, nil
}
