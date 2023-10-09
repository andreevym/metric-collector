package counter

import (
	"strconv"

	"github.com/andreevym/metric-collector/internal/repository"
)

func Store(s repository.Storage, metricName string, metricValue string) error {
	metricValues, err := s.Read(metricName)
	if err != nil {
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
