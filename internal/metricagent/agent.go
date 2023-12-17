package metricagent

import (
	"context"
	"time"
)

type Agent struct {
	Address        string
	PollDuration   time.Duration
	ReportDuration time.Duration
	LiveTime       time.Duration
	SecretKey      string
}

func NewAgent(
	secretKey string,
	address string,
	pollDuration time.Duration,
	reportDuration time.Duration,
	liveTime time.Duration,
) *Agent {
	return &Agent{
		Address:        address,
		PollDuration:   pollDuration,
		ReportDuration: reportDuration,
		LiveTime:       liveTime,
		SecretKey:      secretKey,
	}
}

func (a Agent) Start() {
	ctx := context.Background()
	tickerPoll := time.NewTicker(a.PollDuration)
	go pollLastMemStatByTicker(tickerPoll)

	tickerReport := time.NewTicker(a.ReportDuration)
	go sendLastMemStats(ctx, a.SecretKey, tickerReport, a.Address)

	// время жизни клиента для сбора метрик
	time.Sleep(a.LiveTime)
}
