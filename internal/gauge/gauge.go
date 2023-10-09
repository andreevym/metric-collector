package gauge

import (
	"strconv"

	"github.com/andreevym/metric-collector/internal/repository"
)

func Store(s repository.Storage, metricName string, metricValue string) error {
	err := s.Create(metricName, metricValue)
	return err
}

func Validate(metricValue string) error {
	_, err := strconv.ParseFloat(metricValue, 64)
	return err
}
