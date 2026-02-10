package coffeeshop

import (
	"context"
	"fmt"
	"gopher-cafe/internal/worker"
	"sync"
	"time"

	entity "gopher-cafe/internal/entity/coffeeshop"

	"github.com/ajaibid/coin-common-golang/logger"
)

type CoffeeshopUsecase struct {
	equipPoolManager *worker.EquipPoolManager
	metrics          *entity.OrderMetrics
}

func NewCoffeeshopUsecase(manager *worker.EquipPoolManager, metrics *entity.OrderMetrics) *CoffeeshopUsecase {
	return &CoffeeshopUsecase{
		equipPoolManager: manager,
		metrics:          metrics,
	}
}

func (u *CoffeeshopUsecase) ExecuteBrew(ctx context.Context, orders []entity.Order, baristas int) []entity.OrderResult {
	orderInputChan := make(chan entity.Order, len(orders))

	for _, order := range orders {
		orderInputChan <- order
		logger.Debugf("Order %d submitted", order.ID)
	}
	close(orderInputChan)

	orderResultChan := make(chan entity.OrderResult, len(orders))

	wg := sync.WaitGroup{}

	wg.Add(baristas)
	for i := range baristas {
		go func() {
			defer wg.Done()
			logger.Debugf("Baristas: %d start working", i)
			for {
				select {
				case <-ctx.Done():
					logger.Debugf("Baristas: %d got context done, %v", i, ctx.Err())
					return
				case input, ok := <-orderInputChan:
					if !ok {
						logger.Debugf("Baristas: %d got channel closed", i)
						return
					}
					logger.Debugf("Baristas: %d executing order: %d", i, input.ID)
					res, err := u.processOrder(input)
					if err != nil {
						logger.Errorf("Baristas: %d processing order %d failed: %s", i, input.ID, err)
						continue
					}
					orderResultChan <- res
					u.recordOrderStats(res)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(orderResultChan)
	}()

	results := make([]entity.OrderResult, 0, len(orders))

	for result := range orderResultChan {
		results = append(results, result)
	}

	if len(results) == len(orders) {
		u.metrics.RecordTotalRequests(1)
	}

	logger.Debugf("Finish execute brew : %d, %d", len(orders), baristas)
	return results
}

func (u *CoffeeshopUsecase) processOrder(order entity.Order) (entity.OrderResult, error) {
	var emptyResult entity.OrderResult

	recipe := entity.Recipes[order.Drink]
	res := entity.OrderResult{OrderID: order.ID}

	for _, step := range recipe {
		startStep := time.Now().UnixMilli()

		pool, err := u.equipPoolManager.GetWorkerPool(step.Equipment)
		if err != nil {
			return emptyResult, fmt.Errorf("get worker pool failed: %v", err)
		}

		_, err = pool.Submit(worker.Job{
			OrderID: order.ID,
			Timer:   step.Duration,
		})

		if err != nil {
			return emptyResult, fmt.Errorf("submit failed: %v", err)
		}

		endStep := time.Now().UnixMilli()

		res.Steps = append(res.Steps, entity.StepExecution{
			Equipment:   step.Equipment,
			StartTimeMs: startStep,
			EndTimeMs:   endStep,
		})
	}

	return res, nil
}

func (u *CoffeeshopUsecase) recordOrderStats(res entity.OrderResult) {
	u.metrics.RecordOrder(res)
}

func (u *CoffeeshopUsecase) GetStats() (int64, int64, int64) {
	return u.metrics.GetStats()
}
