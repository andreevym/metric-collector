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
	SecretKey       string
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
	DatabaseDsn string `env:"DATABASE_DSN"`
	SecretKey   string `env:"KEY"`
}

func Flags() (*ServerConfig, error) {
	cfg := ServerConfig{}
	flag.StringVar(&cfg.Address, "a", ":8080", "адрес и порт для запуска сервера")
	flag.StringVar(&cfg.LogLevel, "l", "info", "уровень логирования агента")
	flag.DurationVar(&cfg.StoreInterval, "i", 300*time.Second, "интервал времени в секундах "+
		"по истечении которого текущие показания сервера сохраняются на диск "+
		"(значение 0 делает запись синхронной).")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "полное имя файла "+
		"куда сохраняются текущие значения, пустое значение отключает функцию записи на диск.")
	flag.BoolVar(&cfg.Restore, "r", true, "определяющее, загружать или нет ранее сохранённые значения"+
		" из указанного файла при старте сервера")
	flag.StringVar(&cfg.DatabaseDsn, "d", "", "строка с адресом подключения к БД")
	flag.StringVar(&cfg.SecretKey, "k", "", "secret key, if variable is not empty will "+
		"make hash from request body and add header HashSHA256 for each http request")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	var envConfig ServerEnvConfig
	err := env.Parse(&envConfig)
	if err != nil {
		return nil, err
	}

	if envConfig.SecretKey != "" {
		cfg.SecretKey = envConfig.SecretKey
	}

	if envConfig.Address != "" {
		cfg.Address = envConfig.Address
	}

	if envConfig.LogLevel != "" {
		cfg.LogLevel = envConfig.LogLevel
	}

	if envConfig.DatabaseDsn != "" {
		cfg.DatabaseDsn = envConfig.DatabaseDsn
	}

	if envConfig.StoreInterval != "" {
		v, err := strconv.ParseInt(envConfig.StoreInterval, 10, 32)
		if err != nil {
			panic(fmt.Errorf("problem setup envConfig.StoreInterval %w", err))
		}
		cfg.StoreInterval = time.Second * time.Duration(v)
	}

	if envConfig.FileStoragePath != "" {
		cfg.FileStoragePath = envConfig.FileStoragePath
	}

	if envConfig.Restore != "" {
		restore, err := strconv.ParseBool(envConfig.Restore)
		if err != nil {
			logger.Logger().Fatal(
				"can't parse env RESTORE",
				zap.String("RESTORE", envConfig.Restore),
				zap.Error(err),
			)
		}
		cfg.Restore = restore
	}

	return &cfg, nil
}
