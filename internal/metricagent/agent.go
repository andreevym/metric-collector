package metricagent

import (
	"time"
)

type Agent struct {
	Address        string
	PollDuration   time.Duration
	ReportDuration time.Duration
	LiveTime       time.Duration
}

func NewAgent(
	address string,
	pollDuration time.Duration,
	reportDuration time.Duration,
	liveTime time.Duration) *Agent {
	return &Agent{
		Address:        address,
		PollDuration:   pollDuration,
		ReportDuration: reportDuration,
		LiveTime:       liveTime,
	}
}

func (a Agent) Start() {
	tickerPoll := time.NewTicker(a.PollDuration)
	go pollLastMemStatByTicker(tickerPoll)

	tickerReport := time.NewTicker(a.ReportDuration)
	go sendByTickerAndAddress(tickerReport, a.Address)

	// время жизни клиента для сбора метрик
	time.Sleep(a.LiveTime)
}
