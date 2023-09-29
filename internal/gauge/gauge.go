package gauge

import (
	"strconv"
)

type Storage interface {
	Create(key string, val string) error
}

func StoreGauge(s Storage, metricName string, metricValue string) error {
	// validate metric value argument for 'gauge' metric type
	_, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return err
	}

	err = s.Create(metricName, metricValue)
	if err != nil {
		return err
	}

	return nil
}
