package util

import (
	"time"
)

func Schedule(seconds time.Duration) {
	for {
		go RefreshPolicy()
		time.Sleep(seconds * time.Second)
	}
}
