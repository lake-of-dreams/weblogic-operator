package retry

import (
	"k8s.io/apimachinery/pkg/util/wait"
)

// Retry executes the provided function repeatedly, retrying until the function
// returns done = true, errors, or exceeds the given timeout.
func Retry(backoff wait.Backoff, fn wait.ConditionFunc) error {
	var lastErr error
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		done, err := fn()
		if err != nil {
			lastErr = err
		}
		return done, err
	})
	if err == wait.ErrWaitTimeout {
		err = lastErr
	}
	return err
}
