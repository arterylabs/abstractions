package function_utils

import (
	"time"
)

func WrapInRetries(attempts int, sleep time.Duration, fn func() error) {
	for i := 0; ; i++ {
		err := fn()
		if err == nil {
			break
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleep)
	}
}
