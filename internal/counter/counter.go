package counter

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/andreevym/metric-collector/internal/repository"
	"github.com/andreevym/metric-collector/internal/storage/mem"
)

func Store(s repository.Storage, metricName string, metricValue string) error {
	metricValues, err := s.Read(metricName)
	if err != nil && !errors.Is(err, mem.ErrValueNotFound) {
		return err
	}
	if len(metricValues) == 0 {
		return s.Create(metricName, metricValue)
	}

	existsMetricVal, err := strconv.ParseFloat(metricValues[0], 64)
	if err != nil {
		return err
	}
	v, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return err
	}
	newVal := fmt.Sprintf("%v", existsMetricVal+v)
	return s.Update(metricName, []string{newVal})
}

func Validate(metricValue string) error {
	_, err := strconv.ParseFloat(metricValue, 64)
	return err
}

func Get(s repository.Storage, metricName string) (string, error) {
	v, err := s.Read(metricName)
	if err != nil {
		return "", err
	}

	if len(v) == 0 {
		return "", fmt.Errorf("can't find metric by name %s", metricName)
	}

	return v[0], nil
}
