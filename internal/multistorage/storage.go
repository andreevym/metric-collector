package multistorage

import (
	"errors"
	"reflect"

	"github.com/andreevym/metric-collector/internal/counter"
	"github.com/andreevym/metric-collector/internal/gauge"
	"github.com/andreevym/metric-collector/internal/repository"

	"fmt"
)

type Storage struct {
	counterStorage repository.Storage
	gaugeStorage   repository.Storage
}

func NewStorage(counterStorage repository.Storage,
	gaugeStorage repository.Storage) (*Storage, error) {
	if counterStorage == nil ||
		(reflect.ValueOf(counterStorage).Kind() == reflect.Ptr && reflect.ValueOf(counterStorage).IsNil()) {
		return nil, errors.New("counter storage can't be nil")
	}

	if gaugeStorage == nil ||
		(reflect.ValueOf(gaugeStorage).Kind() == reflect.Ptr && reflect.ValueOf(gaugeStorage).IsNil()) {
		return nil, errors.New("gauge storage can't be nil")
	}

	return &Storage{
		counterStorage: counterStorage,
		gaugeStorage:   gaugeStorage,
	}, nil
}

func (s Storage) GaugeStorage() repository.Storage {
	return s.gaugeStorage
}

func (s Storage) CounterStorage() repository.Storage {
	return s.counterStorage
}

const (
	MetricTypeGauge   string = "gauge"
	MetricTypeCounter string = "counter"
)

func SaveMetric(storage *Storage, metricName string, metricType string, metricValue string) (string, error) {
	var newVal string
	switch metricType {
	case MetricTypeCounter:
		err := counter.Validate(metricValue)
		if err != nil {
			return "", err
		}
		newVal, err = counter.Store(storage.CounterStorage(), metricName, metricValue)
		if err != nil {
			return "", err
		}
	case MetricTypeGauge:
		err := gauge.Validate(metricValue)
		if err != nil {
			return "", err
		}
		newVal, err = gauge.Store(storage.GaugeStorage(), metricName, metricValue)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("metric type not found")
	}

	return newVal, nil
}

func GetMetric(storage *Storage, metricType string, metricName string) (string, error) {
	switch metricType {
	case MetricTypeCounter:
		return counter.Get(storage.CounterStorage(), metricName)
	case MetricTypeGauge:
		v, err := gauge.Get(storage.GaugeStorage(), metricName)
		if err != nil {
			return "", err
		}
		if len(v) != 0 {
			return v[len(v)-1], nil
		} else {
			return "", nil
		}
	default:
		return "", fmt.Errorf("metric type '%s' not found", metricType)
	}
}
