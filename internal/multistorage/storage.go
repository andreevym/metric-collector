package multistorage

import (
	"errors"
	"os"
	"reflect"
	"time"

	"fmt"

	"github.com/andreevym/metric-collector/internal/backup"
	"github.com/andreevym/metric-collector/internal/config/serverconfig"
	"github.com/andreevym/metric-collector/internal/counter"
	"github.com/andreevym/metric-collector/internal/gauge"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"go.uber.org/zap"
)

type MetricStorage struct {
	counterStorage    storage.Storage
	gaugeStorage      storage.Storage
	counterBackupPath string
	gaugeBackupPath   string
	cfg               *serverconfig.ServerConfig
}

type Restorable interface {
	UpdateData(data map[string][]string)
	Data() map[string][]string
}

func NewMetricStorage(counterStorage storage.Storage, gaugeStorage storage.Storage, cfg *serverconfig.ServerConfig) (*MetricStorage, error) {
	if counterStorage == nil ||
		(reflect.ValueOf(counterStorage).Kind() == reflect.Ptr && reflect.ValueOf(counterStorage).IsNil()) {
		return nil, errors.New("counter storage can't be nil")
	}

	if gaugeStorage == nil ||
		(reflect.ValueOf(gaugeStorage).Kind() == reflect.Ptr && reflect.ValueOf(gaugeStorage).IsNil()) {
		return nil, errors.New("gauge storage can't be nil")
	}

	if cfg == nil {
		return nil, errors.New("server config can't be nil")
	}

	s := &MetricStorage{
		counterStorage: counterStorage,
		gaugeStorage:   gaugeStorage,
		cfg:            cfg,
	}

	if cfg.FileStoragePath != "" {
		err := os.MkdirAll(cfg.FileStoragePath+"/", 0777)
		if err != nil {
			panic(err)
		}
		s.counterBackupPath = cfg.FileStoragePath + "/counter.backup"
		s.gaugeBackupPath = cfg.FileStoragePath + "/gauge.backup"
		if ok, _ := isDirectory(cfg.FileStoragePath); !ok {
			return nil, fmt.Errorf("storage path need to be directory %s", cfg.FileStoragePath)
		}
	}
	if cfg.Restore {
		err := s.Restore()
		if err != nil {
			panic(err)
		}
	}
	return s, nil
}

func (storage *MetricStorage) GaugeStorage() storage.Storage {
	return storage.gaugeStorage
}

func (storage *MetricStorage) CounterStorage() storage.Storage {
	return storage.counterStorage
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func (storage *MetricStorage) Restore() error {
	if r, ok := storage.counterStorage.(Restorable); ok {
		if storage.counterBackupPath == "" {
			return errors.New("storage counterBackupPath path can't be empty")
		}
		data, err := backup.Load(storage.counterBackupPath)
		if err != nil {
			return err
		}
		r.UpdateData(data)
	}
	if r, ok := storage.gaugeStorage.(Restorable); ok {
		if storage.gaugeBackupPath == "" {
			return errors.New("storage gaugeBackupPath path can't be empty")
		}
		data, err := backup.Load(storage.gaugeBackupPath)
		if err != nil {
			return err
		}
		r.UpdateData(data)
	}
	return nil
}

const (
	MetricTypeGauge   string = "gauge"
	MetricTypeCounter string = "counter"
)

func (storage *MetricStorage) SaveMetric(
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
		newVal, err = counter.Store(storage.CounterStorage(), metricName, metricValue)
		if err != nil {
			return "", err
		}
		if storage.cfg != nil && storage.counterBackupPath != "" && storage.cfg.StoreInterval > 0 {
			if r, ok := storage.gaugeStorage.(Restorable); ok {
				time.AfterFunc(storage.cfg.StoreInterval, func() {
					err = backup.Save(storage.counterBackupPath, r.Data())
					if err != nil {
						logger.Log.Error("problem to save backup ", zap.Error(err))
					}
				})
			}
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
		if storage.cfg != nil && storage.gaugeBackupPath != "" && storage.cfg.StoreInterval > 0 {
			if r, ok := storage.gaugeStorage.(Restorable); ok {
				time.AfterFunc(storage.cfg.StoreInterval, func() {
					err = backup.Save(storage.gaugeBackupPath, r.Data())
					if err != nil {
						logger.Log.Error("problem to save backup ", zap.Error(err))
					}
				})
			}
		}
	default:
		return "", errors.New("metric type not found")
	}

	return newVal, nil
}

func (storage *MetricStorage) GetMetric(metricType string, metricName string) (string, error) {
	switch metricType {
	case MetricTypeCounter:
		return counter.Get(storage.CounterStorage(), metricName)
	case MetricTypeGauge:
		return gauge.Get(storage.GaugeStorage(), metricName)
	default:
		return "", fmt.Errorf("metric type '%s' not found", metricType)
	}
}
