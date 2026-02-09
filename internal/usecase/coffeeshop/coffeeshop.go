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
	u.metrics.RecordTotalRequests(len(orders))

	orderInputChan := make(chan OrderInput)
	defer close(orderInputChan)

	for i := range baristas {
		go func() {
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
					logger.Debugf("Baristas: %d executing order: %d", i, input.data.ID)
					res, err := u.processOrder(input.data)
					if err != nil {
						input.resultChannel <- entity.OrderResult{OrderID: input.data.ID}
						logger.Errorf("Baristas: %d processing order %d failed: %s", i, input.data.ID, err)
						return
					}
					input.resultChannel <- res
				}
			}
		}()
	}

	wg := sync.WaitGroup{}

	logger.Debugf("Start submit orders")

	results := make([]entity.OrderResult, 0, len(orders))

	for _, order := range orders {
		orderResultChan := make(chan entity.OrderResult)

		logger.Debugf("Order %d submitted", order.ID)

		select {
		case orderInputChan <- OrderInput{
			data:          order,
			resultChannel: orderResultChan,
		}:
		case <-ctx.Done():
			return results
		}

		wg.Add(1)
		go func() {
			defer close(orderResultChan)
			defer wg.Done()
			select {
			case <-ctx.Done():
				logger.Debugf("OrderResult %d got context done", order.ID)
				return
			case orderResult, ok := <-orderResultChan:
				if !ok {
					logger.Debugf("OrderResult %d got channel closed", order.ID)
					return
				}
				results = append(results, orderResult)
				u.recordOrderStats(orderResult)
			}
		}()
	}

	wg.Wait()

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
	return u.metrics.GetTotalRequests(),
		u.metrics.GetTotalOrders(),
		u.metrics.GetP90Duration()
}
