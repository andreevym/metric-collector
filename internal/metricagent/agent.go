package metricagent

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type Agent struct {
	Address        string
	PollDuration   time.Duration
	ReportDuration time.Duration
	LiveTime       time.Duration
	SecretKey      string
	CryptoKey      string
	RateLimit      int
}

func NewAgent(
	secretKey string,
	cryptoKey string,
	address string,
	pollDuration time.Duration,
	reportDuration time.Duration,
	liveTime time.Duration,
	rateLimit int,
) *Agent {
	return &Agent{
		Address:        address,
		PollDuration:   pollDuration,
		ReportDuration: reportDuration,
		LiveTime:       liveTime,
		SecretKey:      secretKey,
		CryptoKey:      cryptoKey,
		RateLimit:      rateLimit,
	}
}

func (a Agent) Run() error {
	// время жизни клиента для сбора метрик
	ctx, cancelFunc := context.WithTimeout(context.Background(), a.LiveTime)
	defer cancelFunc()
	metricsCh, err := collectMetric(ctx, a.PollDuration, a.RateLimit)
	if err != nil {
		return fmt.Errorf("failed to collect metric: %w", err)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < a.RateLimit; i++ {
		wg.Add(1)
		go func() {
			// откладываем уменьшение счетчика в WaitGroup, когда завершится горутина
			defer wg.Done()
			sendMetric(ctx, metricsCh, a.SecretKey, a.CryptoKey, a.ReportDuration, a.Address)
		}()
	}
	wg.Wait()
	return nil
}
