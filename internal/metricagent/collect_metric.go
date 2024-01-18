package metricagent

import (
	"context"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
)

func collectMetric(
	ctx context.Context,
	pollDuration time.Duration,
	rateLimit int,
) (chan []*storage.Metric, error) {
	// создаем буферизованный канал для отправки результатов
	metricsCh := make(chan []*storage.Metric, rateLimit)

	go func() {
		defer close(metricsCh)

		// pollCountAtomic (тип counter) — счётчик, увеличивающийся на 1
		// при каждом обновлении метрики из пакета runtime (на каждый pollInterval — см. ниже).
		var pollCountAtomic atomic.Int64
		ticker := time.NewTicker(pollDuration)
		for t := range ticker.C {
			metrics := make([]*storage.Metric, 0)

			logger.Logger().Debug("pollLastMemStatByTicker", zap.String("ticker", t.String()))
			memStats := runtime.MemStats{}
			runtime.ReadMemStats(&memStats)

			statMetrics, err := mapMemStatToMetrics(&memStats)
			if err != nil {
				return
			}
			metrics = append(metrics, statMetrics...)

			pollCountAtomic.Add(1)
			pollCount := pollCountAtomic.Load()
			metrics = append(metrics, &storage.Metric{
				ID:    "PollCount",
				MType: storage.MTypeCounter,
				Delta: &pollCount,
			})

			total, free, err := Memory()
			if err != nil {
				logger.Logger().Error("failed to get memory statistics",
					zap.Error(err),
				)
				return
			}
			metrics = append(metrics, &storage.Metric{
				ID:    "TotalMemory",
				MType: storage.MTypeGauge,
				Value: total,
			})

			metrics = append(metrics, &storage.Metric{
				ID:    "FreeMemory",
				MType: storage.MTypeGauge,
				Value: free,
			})
			cpuUtilization, err := CPUUtilization()
			if err != nil {
				logger.Logger().Error("failed to get cpu utilization",
					zap.Error(err),
				)
				return
			}
			metrics = append(metrics, &storage.Metric{
				ID:    "PollCount",
				MType: storage.MTypeGauge,
				Value: cpuUtilization,
			})

			select {
			case <-ctx.Done():
				return
			case metricsCh <- metrics:
			}
		}
	}()

	return metricsCh, nil
}

// Memory received virtual memory
// Total amount of RAM on this system
// Free is the kernel's notion of free memory; RAM chips whose bits nobody
// cares about the value of right now. For a human consumable number,
// Available is what you really want.
func Memory() (*float64, *float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, nil, err
	}
	totalFloat := float64(v.Total)
	freeFloat := float64(v.Free)
	return &totalFloat, &freeFloat, nil
}

// CPUUtilization точное количество — по числу CPU, определяемому во время исполнения
// CPU utilization indicates the amount of load handled by individual processor
// cores to run various programs on a computer. Definition.
func CPUUtilization() (*float64, error) {
	percent, err := cpu.Percent(time.Millisecond, false)
	if err != nil {
		return nil, err
	}
	return &percent[0], nil
}
