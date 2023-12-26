package metricagent

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"

	"github.com/andreevym/metric-collector/internal/storage"
)

func mapMemStatToMetrics(stats *runtime.MemStats) ([]*storage.Metric, error) {
	metrics := make([]*storage.Metric, 0)
	metrics = mustAppendGaugeMetricFloat64(metrics, "RandomValue", rand.Float64())
	metrics = mustAppendGaugeMetricUint64(metrics, "Alloc", stats.Alloc)
	metrics = mustAppendGaugeMetricUint64(metrics, "BuckHashSys", stats.BuckHashSys)
	metrics = mustAppendGaugeMetricUint64(metrics, "Frees", stats.Frees)
	metrics = mustAppendGaugeMetricFloat64(metrics, "GCCPUFraction", stats.GCCPUFraction)
	metrics = mustAppendGaugeMetricUint64(metrics, "GCSys", stats.GCSys)
	metrics = mustAppendGaugeMetricUint64(metrics, "HeapAlloc", stats.HeapAlloc)
	metrics = mustAppendGaugeMetricUint64(metrics, "HeapIdle", stats.HeapIdle)
	metrics = mustAppendGaugeMetricUint64(metrics, "HeapInuse", stats.HeapInuse)
	metrics = mustAppendGaugeMetricUint64(metrics, "HeapObjects", stats.HeapObjects)
	metrics = mustAppendGaugeMetricUint64(metrics, "HeapReleased", stats.HeapReleased)
	metrics = mustAppendGaugeMetricUint64(metrics, "HeapSys", stats.HeapSys)
	metrics = mustAppendGaugeMetricUint64(metrics, "LastGC", stats.LastGC)
	metrics = mustAppendGaugeMetricUint64(metrics, "Lookups", stats.Lookups)
	metrics = mustAppendGaugeMetricUint64(metrics, "MCacheInuse", stats.MCacheInuse)
	metrics = mustAppendGaugeMetricUint64(metrics, "MCacheSys", stats.MCacheSys)
	metrics = mustAppendGaugeMetricUint64(metrics, "MSpanInuse", stats.MSpanInuse)
	metrics = mustAppendGaugeMetricUint64(metrics, "MSpanSys", stats.MSpanSys)
	metrics = mustAppendGaugeMetricUint64(metrics, "Mallocs", stats.Mallocs)
	metrics = mustAppendGaugeMetricUint64(metrics, "NextGC", stats.NextGC)
	metrics = mustAppendGaugeMetricUint32(metrics, "NumForcedGC", stats.NumForcedGC)
	metrics = mustAppendGaugeMetricUint32(metrics, "NumGC", stats.NumGC)
	metrics = mustAppendGaugeMetricUint64(metrics, "OtherSys", stats.OtherSys)
	metrics = mustAppendGaugeMetricUint64(metrics, "PauseTotalNs", stats.PauseTotalNs)
	metrics = mustAppendGaugeMetricUint64(metrics, "StackInuse", stats.StackInuse)
	metrics = mustAppendGaugeMetricUint64(metrics, "StackSys", stats.StackSys)
	metrics = mustAppendGaugeMetricUint64(metrics, "Sys", stats.Sys)
	metrics = mustAppendGaugeMetricUint64(metrics, "TotalAlloc", stats.TotalAlloc)
	return metrics, nil
}

func mustAppendGaugeMetricUint64(metrics []*storage.Metric, id string, v uint64) []*storage.Metric {
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	if err != nil {
		panic(err)
	}
	metrics = append(metrics, &storage.Metric{
		ID:    id,
		MType: storage.MTypeGauge,
		Value: &f,
	})
	return metrics
}

func mustAppendGaugeMetricUint32(metrics []*storage.Metric, id string, v uint32) []*storage.Metric {
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 32)
	if err != nil {
		panic(err)
	}
	metrics = append(metrics, &storage.Metric{
		ID:    id,
		MType: storage.MTypeGauge,
		Value: &f,
	})
	return metrics
}

func mustAppendGaugeMetricFloat64(metrics []*storage.Metric, id string, f float64) []*storage.Metric {
	metrics = append(metrics, &storage.Metric{
		ID:    id,
		MType: storage.MTypeGauge,
		Value: &f,
	})
	return metrics
}
