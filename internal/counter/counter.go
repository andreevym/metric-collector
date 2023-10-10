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
		err = s.Create(metricName, metricValue)
		if err != nil {
			return err
		}
	} else {
		err = s.Update(metricName, []string{metricValue})
		if err != nil {
			return err
		}
	}

	return nil
}

func Validate(metricValue string) error {
	_, err := strconv.ParseInt(metricValue, 10, 64)
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
