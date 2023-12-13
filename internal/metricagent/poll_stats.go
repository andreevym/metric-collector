package metricagent

import (
	"runtime"
	"time"

	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
)

func pollLastMemStatByTicker(ticker *time.Ticker) {
	for t := range ticker.C {
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		lastMemStats = &memStats
		logger.Logger().Info("metric", zap.String("ticker", t.String()))
	}
}
