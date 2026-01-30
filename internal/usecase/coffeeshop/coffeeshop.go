package coffeeshop

import (
	"sync/atomic"
	"time"

	entity "gopher-cafe/internal/entity/coffeeshop"
)

type CoffeeshopUsecase struct {
	totalRequests int64
	totalOrders   int64
	p90RequestsMs int64
}

func NewCoffeeshopUsecase() *CoffeeshopUsecase {
	return &CoffeeshopUsecase{}
}

func (u *CoffeeshopUsecase) ExecuteBrew(orders []entity.Order, baristas int) []entity.OrderResult {
	atomic.AddInt64(&u.totalRequests, 1)

	results := make([]entity.OrderResult, 0, len(orders))

	for _, order := range orders {
		recipe := entity.Recipes[order.Drink]
		res := entity.OrderResult{OrderID: order.ID}

		for _, step := range recipe {
			startStep := time.Now().UnixMilli()
			time.Sleep(step.Duration)
			endStep := time.Now().UnixMilli()

			res.Steps = append(res.Steps, entity.StepExecution{
				Equipment:   step.Equipment,
				StartTimeMs: startStep,
				EndTimeMs:   endStep,
			})
		}
		results = append(results, res)
		u.recordOrderStats(res)
	}
	return results
}

func (u *CoffeeshopUsecase) recordOrderStats(res entity.OrderResult) {
	atomic.AddInt64(&u.totalOrders, 1)
	if len(res.Steps) > 0 {
		duration := res.Steps[len(res.Steps)-1].EndTimeMs - res.Steps[0].StartTimeMs
		atomic.AddInt64(&u.p90RequestsMs, duration)
	}
}

func (u *CoffeeshopUsecase) GetStats() (int64, int64, int64) {
	return atomic.LoadInt64(&u.totalRequests),
		atomic.LoadInt64(&u.totalOrders),
		atomic.LoadInt64(&u.p90RequestsMs)
}
