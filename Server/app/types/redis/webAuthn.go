package redis

import "time"

type Challenge struct {
	Challenge string        `redis:"challenge"`
	OutTime   time.Duration `redis:"out_time"`
}
