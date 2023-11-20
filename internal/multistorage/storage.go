package multistorage

import (
	"errors"
	"log"
	"os"
	"reflect"
	"time"

	"fmt"

	"github.com/andreevym/metric-collector/internal/backup"
	"github.com/andreevym/metric-collector/internal/config/serverconfig"
	"github.com/andreevym/metric-collector/internal/counter"
	"github.com/andreevym/metric-collector/internal/gauge"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/repository"
	"go.uber.org/zap"
)

type Storage struct {
	counterStorage    repository.Storage
	gaugeStorage      repository.Storage
	counterBackupPath string
	gaugeBackupPath   string
	cfg               *serverconfig.ServerConfig
}

func NewStorage(counterStorage repository.Storage, gaugeStorage repository.Storage, cfg *serverconfig.ServerConfig) (*Storage, error) {
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

	s := &Storage{
		counterStorage: counterStorage,
		gaugeStorage:   gaugeStorage,
		cfg:            cfg,
	}

	if cfg != nil && cfg.FileStoragePath != "" {
		err := os.MkdirAll(cfg.FileStoragePath+"/", 0777)
		if err != nil {
			panic(err)
		}
		s.counterBackupPath = cfg.FileStoragePath + "/counter.backup"
		_, err = os.Create(s.counterBackupPath)
		if err != nil {
			log.Fatalf("create file: %v", err)
		}
		s.gaugeBackupPath = cfg.FileStoragePath + "/gauge.backup"
		_, err = os.Create(s.gaugeBackupPath)
		if err != nil {
			log.Fatalf("create file: %v", err)
		}
		if ok, _ := isDirectory(cfg.FileStoragePath); !ok {
			return nil, fmt.Errorf("storage path need to be directory %s", cfg.FileStoragePath)
		}
	}
	if cfg != nil && cfg.Restore {
		err := s.Restore()
		if err != nil {
			panic(err)
		}
	}
	return s, nil
}

func (s Storage) GaugeStorage() repository.Storage {
	return s.gaugeStorage
}

func (s Storage) CounterStorage() repository.Storage {
	return s.counterStorage
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func (s Storage) Restore() error {
	if s.counterBackupPath == "" {
		return errors.New("storage counterBackupPath path can't be empty")
	}
	if s.gaugeBackupPath == "" {
		return errors.New("storage gaugeBackupPath path can't be empty")
	}

	data, err := backup.Load(s.counterBackupPath)
	s.counterStorage.UpdateData(data)
	if err != nil {
		return err
	}
	data, err = backup.Load(s.gaugeBackupPath)
	s.gaugeStorage.UpdateData(data)
	if err != nil {
		return err
	}
	return nil
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
		if storage.cfg != nil && storage.counterBackupPath != "" && storage.cfg.StoreInterval > 0 {
			time.AfterFunc(storage.cfg.StoreInterval, func() {
				err = backup.Save(storage.counterBackupPath, storage.CounterStorage().Data())
				if err != nil {
					logger.Log.Error("problem to save backup ", zap.Error(err))
				}
			})
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
			time.AfterFunc(storage.cfg.StoreInterval, func() {
				err = backup.Save(storage.gaugeBackupPath, storage.GaugeStorage().Data())
				if err != nil {
					logger.Log.Error("problem to save backup ", zap.Error(err))
				}
			})
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
