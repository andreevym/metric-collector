package counter

import (
	"strconv"
)

type Storage interface {
	Create(key string, val string) error
	Read(key string) ([]string, error)
	Update(key string, val []string) error
}

func StoreCounter(s Storage, metricName string, metricValue string) error {
	// validate metric value argument for 'counter' metric type
	_, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return err
	}
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
