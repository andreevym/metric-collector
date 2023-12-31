package utils

import (
	"time"

	"github.com/avast/retry-go"
)

// RetryDelayType кастомный тип ожидания, т.к интервалы между повторами должны увеличиваться: 1s, 3s, 5s.
func RetryDelayType(n uint, _ error, config *retry.Config) time.Duration {
	switch n {
	case 0:
		return 1 * time.Second
	case 1:
		return 3 * time.Second
	default:
		return 5 * time.Second
	}
}
