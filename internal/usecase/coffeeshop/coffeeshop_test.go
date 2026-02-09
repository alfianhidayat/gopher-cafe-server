package coffeeshop

import (
	entity "gopher-cafe/internal/entity/coffeeshop"
	"gopher-cafe/internal/worker"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tests := []struct {
		name      string
		baristas  int
		orders    []entity.Order
		want      []entity.OrderResult
		wantStats struct {
			totalRequests int64
			totalOrders   int64
			p90RequestsMs int64
		}
	}{
		{
			name:     "success",
			baristas: 1,
			orders: []entity.Order{
				{
					ID:    1,
					Drink: entity.DrinkEspresso,
				},
			},
			want: []entity.OrderResult{
				{
					OrderID: 1,
					Steps: []entity.StepExecution{
						{
							Equipment:   entity.EquipGrinder,
							StartTimeMs: 0,
							EndTimeMs:   0,
						},
						{
							Equipment:   entity.EquipEspressoMachine,
							StartTimeMs: 0,
							EndTimeMs:   0,
						},
					},
				},
			},
			wantStats: struct {
				totalRequests int64
				totalOrders   int64
				p90RequestsMs int64
			}{totalRequests: int64(1), totalOrders: int64(1), p90RequestsMs: int64(1)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ew := worker.EquipmentWorkers

			manager := worker.NewEquipPoolManager(uint8(len(ew)))
			for k, v := range ew {
				manager.Register(k, v)
			}
			manager.StartAll()

			metrics := entity.NewOrderMetrics(10000)

			usecase := NewCoffeeshopUsecase(manager, metrics)
			results := usecase.ExecuteBrew(t.Context(), test.orders, test.baristas)

			assert.Len(t, results, len(test.want))

			for i, r := range results {
				want := test.want[i]
				assert.Equal(t, want.OrderID, r.OrderID)
				for i, s := range r.Steps {
					assert.Equal(t, want.Steps[i].Equipment, s.Equipment)
				}
			}

			//wantStats := test.wantStats
			//assert.Equal(t, wantStats.totalRequests, usecase.totalRequests)
			//assert.Equal(t, wantStats.totalOrders, usecase.totalOrders)
		})
	}
}
