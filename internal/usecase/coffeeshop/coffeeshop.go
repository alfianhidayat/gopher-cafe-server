package coffeeshop

import (
	"context"
	"fmt"
	"gopher-cafe/internal/worker"
	"sync"
	"sync/atomic"
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
	orderInputChan := make(chan entity.OrderStep, len(orders))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, order := range orders {
		recipe := entity.Recipes[order.Drink]

		orderInputChan <- entity.OrderStep{
			Order:    order,
			Recipes:  recipe,
			Steps:    make([]entity.StepExecution, 0, len(recipe)),
			NextStep: uint8(0),
		}

		logger.Debugf("Order %d submitted", order.ID)
	}

	orderResultChan := make(chan entity.OrderResult, len(orders))

	wg := sync.WaitGroup{}
	completed := atomic.Int64{}

	wg.Add(baristas)
	for i := 0; i < baristas; i++ {
		go func() {
			defer wg.Done()
			logger.Debugf("Baristas: %d start working", i)
			for {
				select {
				case <-ctx.Done():
					logger.Debugf("Baristas: %d got context done", i)
					return
				case input, ok := <-orderInputChan:
					if !ok {
						logger.Debugf("Baristas: %d got channel closed", i)
						return
					}

					if len(input.Recipes) != len(input.Steps) {
						logger.Debugf("Baristas: %d executing order: %d, step: %d", i, input.Order.ID, input.NextStep)
						err := u.processStep(&input)
						if err != nil {
							logger.Errorf("Baristas: %d processing order %d failed: %s", i, input.Order.ID, err)
							continue
						}
						// requeue safely
						select {
						case orderInputChan <- input:
						case <-ctx.Done():
							return
						}
						continue
					}

					res := entity.OrderResult{
						OrderID: input.Order.ID,
						Steps:   input.Steps,
					}

					select {
					case orderResultChan <- res:
						u.recordOrderStats(res)
					case <-ctx.Done():
						return
					}

					// count completion
					if completed.Add(1) == int64(len(orders)) {
						// 🔴 shutdown trigger point
						cancel()              // stop workers
						close(orderInputChan) // unblock receivers
					}
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

func (u *CoffeeshopUsecase) processStep(input *entity.OrderStep) error {
	startStep := time.Now().UnixMilli()

	step := input.Recipes[input.NextStep]

	pool, err := u.equipPoolManager.GetWorkerPool(step.Equipment)
	if err != nil {
		return fmt.Errorf("get worker pool failed: %v", err)
	}

	err = pool.Submit(worker.Job{
		OrderID: input.Order.ID,
		Timer:   step.Duration,
	})

	if err != nil {
		return fmt.Errorf("submit failed: %v", err)
	}

	endStep := time.Now().UnixMilli()

	input.Steps = append(input.Steps, entity.StepExecution{
		Equipment:   step.Equipment,
		StartTimeMs: startStep,
		EndTimeMs:   endStep,
	})

	input.NextStep++

	return nil
}

func (u *CoffeeshopUsecase) recordOrderStats(res entity.OrderResult) {
	u.metrics.RecordOrder(res)
}

func (u *CoffeeshopUsecase) GetStats() (int64, int64, int64) {
	return u.metrics.GetStats()
}
