package parallel

import (
	"context"
	"sync"
	"time"
)

type TaskFunc[YieldData any, YieldResult any] func(itemCtx context.Context, data YieldData) (yieldResult YieldResult, yieldReturn bool, yieldErr error)

func Invoke[YieldData any, YieldResult any](ctx context.Context, datas []YieldData, task TaskFunc[YieldData, YieldResult], taskTimeout time.Duration) ([]YieldResult, []error) {
	var wg sync.WaitGroup
	var results []YieldResult
	var errors []error
	var resultLock sync.Mutex

	for _, data := range datas {
		wg.Add(1)

		itemCtx, cancel := context.WithTimeout(ctx, taskTimeout)
		defer cancel()

		go func(ctx context.Context, data YieldData) {
			defer wg.Done()

			result, yieldReturn, yieldErr := task(ctx, data)

			if yieldReturn {
				resultLock.Lock()
				defer resultLock.Unlock()
				results = append(results, result)
			}

			if yieldErr != nil {
				resultLock.Lock()
				defer resultLock.Unlock()
				errors = append(errors, yieldErr)
			}
		}(itemCtx, data)
	}

	wg.Wait()
	return results, errors
}
