package serverconfig

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Address         string
	LogLevel        string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	DatabaseDsn     string
}

type ServerEnvConfig struct {
	Address  string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
	// StoreInterval переменная окружения STORE_INTERVAL — интервал времени в секундах,
	// по истечении которого текущие показания сервера сохраняются на диск (по умолчанию 300 секунд,
	// значение 0 делает запись синхронной).
	StoreInterval string `env:"STORE_INTERVAL"`
	// FileStoragePath переменная окружения FILE_STORAGE_PATH — полное имя файла,
	// куда сохраняются текущие значения (по умолчанию /tmp/metrics-db.json,
	// пустое значение отключает функцию записи на диск).
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	// Restore переменная окружения RESTORE — булево значение (true/false),
	// определяющее, загружать или нет ранее сохранённые значения из указанного
	// файла при старте сервера (по умолчанию true).
	Restore string `env:"RESTORE"`
	// DatabaseDsn Строка с адресом подключения к БД должна получаться из переменной окружения DATABASE_DSN
	// или флага командной строки -d.
	DatabaseDsn string `env:"DATABASE_DSN"`
}

var (
	// flagRunAddr адрес и порт для запуска сервера, аргумент -a со значением :8080 по умолчанию
	flagRunAddr string
	// flagLogLevel уровень логирования агента
	flagLogLevel string
	// flagStoreInterval флаг -i — интервал времени в секундах,
	// по истечении которого текущие показания сервера сохраняются на диск (по умолчанию 300 секунд,
	// значение 0 делает запись синхронной).
	flagStoreInterval int
	// flagFileStoragePath флаг -f — полное имя файла,
	// куда сохраняются текущие значения (по умолчанию /tmp/metrics-db.json,
	// пустое значение отключает функцию записи на диск).
	flagFileStoragePath string
	// flagRestore флаг -r — булево значение (true/false),
	// определяющее, загружать или нет ранее сохранённые значения из указанного
	// файла при старте сервера (по умолчанию true).
	flagRestore bool
	// flagDatabaseDsn флаг командной строки -d,
	// Строка с адресом подключения к БД
	flagDatabaseDsn string
)

func Flags() (*ServerConfig, error) {
	flag.StringVar(&flagRunAddr, "a", ":8081", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.IntVar(&flagStoreInterval, "i", 300, "STORE INTERVAL")
	flag.StringVar(&flagFileStoragePath, "f", "/tmp/metricserver", "file storage path")
	flag.BoolVar(&flagRestore, "r", true, "restore")
	flag.StringVar(&flagDatabaseDsn, "d", "", "postgres connection DATABASE_DSN")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	var config ServerEnvConfig

	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	resultConfig := ServerConfig{}

	if config.Address == "" {
		resultConfig.Address = flagRunAddr
	} else {
		resultConfig.Address = config.Address
	}

	// Логирование, по умолчанию info
	if config.LogLevel == "" {
		resultConfig.LogLevel = flagLogLevel
	} else {
		resultConfig.LogLevel = config.LogLevel
	}

	// Логирование, по умолчанию info
	if config.DatabaseDsn == "" {
		resultConfig.DatabaseDsn = flagDatabaseDsn
	} else {
		resultConfig.DatabaseDsn = config.DatabaseDsn
	}

	if config.StoreInterval == "" {
		resultConfig.StoreInterval = time.Second * time.Duration(flagStoreInterval)
	} else {
		v, err := strconv.ParseInt(config.StoreInterval, 10, 32)
		if err != nil {
			panic(fmt.Errorf("problem setup config.StoreInterval %w", err))
		}
		resultConfig.StoreInterval = time.Second * time.Duration(v)
	}

	// Логирование, по умолчанию info
	if config.FileStoragePath == "" {
		resultConfig.FileStoragePath = flagFileStoragePath
	} else {
		resultConfig.FileStoragePath = config.FileStoragePath
	}

	// Логирование, по умолчанию info
	if config.Restore == "" {
		resultConfig.Restore = flagRestore
	} else {
		restore, err := strconv.ParseBool(config.Restore)
		if err != nil {
			logger.Log.Fatal(
				"can't parse env RESTORE",
				zap.String("RESTORE", config.Restore),
				zap.Error(err),
			)
		}
		resultConfig.Restore = restore
	}

	return &resultConfig, nil
}
