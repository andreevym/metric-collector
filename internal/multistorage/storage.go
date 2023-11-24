package multistorage

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/andreevym/metric-collector/internal/counter"
	"github.com/andreevym/metric-collector/internal/gauge"
	"github.com/andreevym/metric-collector/internal/storage"
)

type MetricManager struct {
	counterStorage storage.Storage
	gaugeStorage   storage.Storage
}

type Backup interface {
	Restore() error
	Backup() error
}

func NewMetricManager(
	counterStorage storage.Storage,
	gaugeStorage storage.Storage,
) (*MetricManager, error) {
	if counterStorage == nil ||
		(reflect.ValueOf(counterStorage).Kind() == reflect.Ptr && reflect.ValueOf(counterStorage).IsNil()) {
		return nil, errors.New("counter storage can't be nil")
	}

	if gaugeStorage == nil ||
		(reflect.ValueOf(gaugeStorage).Kind() == reflect.Ptr && reflect.ValueOf(gaugeStorage).IsNil()) {
		return nil, errors.New("gauge storage can't be nil")
	}

	return &MetricManager{
		counterStorage: counterStorage,
		gaugeStorage:   gaugeStorage,
	}, nil
}

const (
	MetricTypeGauge   string = "gauge"
	MetricTypeCounter string = "counter"
)

func (storage *MetricManager) SaveMetric(
	metricName string,
	metricType string,
	metricValue string,
) (string, error) {
	var newVal string
	switch metricType {
	case MetricTypeCounter:
		err := counter.Validate(metricValue)
		if err != nil {
			return "", err
		}
		newVal, err = counter.Store(storage.counterStorage, metricName, metricValue)
		if err != nil {
			return "", err
		}
		if b, ok := storage.counterStorage.(Backup); ok {
			err = b.Backup()
			if err != nil {
				return "", err
			}
		}
	case MetricTypeGauge:
		err := gauge.Validate(metricValue)
		if err != nil {
			return "", err
		}
		newVal, err = gauge.Store(storage.gaugeStorage, metricName, metricValue)
		if err != nil {
			return "", err
		}
		if b, ok := storage.gaugeStorage.(Backup); ok {
			err = b.Backup()
			if err != nil {
				return "", err
			}
		}
	default:
		return "", errors.New("metric type not found")
	}

	return newVal, nil
}

func (storage *MetricManager) GetMetric(metricType string, metricName string) (string, error) {
	switch metricType {
	case MetricTypeCounter:
		return counter.Get(storage.counterStorage, metricName)
	case MetricTypeGauge:
		return gauge.Get(storage.gaugeStorage, metricName)
	default:
		return "", fmt.Errorf("metric type '%s' not found", metricType)
	}
}

func (storage *MetricManager) SaveMetrics(metricType string, kvMap map[string]*storage.Metric) error {
	var err error
	switch metricType {
	case MetricTypeCounter:
		err = counter.StoreAll(storage.counterStorage, kvMap)
		if err != nil {
			return err
		}
		if b, ok := storage.counterStorage.(Backup); ok {
			err = b.Backup()
			if err != nil {
				return err
			}
		}
	case MetricTypeGauge:
		err = gauge.StoreAll(storage.gaugeStorage, kvMap)
		if err != nil {
			return err
		}
		if b, ok := storage.gaugeStorage.(Backup); ok {
			err = b.Backup()
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("metric type not found")
	}

	return nil
}
