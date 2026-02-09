package coffeeshop

import (
	entity "gopher-cafe/internal/entity/coffeeshop"
	"gopher-cafe/internal/worker"
	"testing"
)

func BenchmarkExecuteBrew(b *testing.B) {
	ew := worker.EquipmentWorkers

	manager := worker.NewEquipPoolManager(uint8(len(ew)))
	for k, v := range ew {
		manager.Register(k, v)
	}

	metrics := entity.NewOrderMetrics(10000)

	usecase := NewCoffeeshopUsecase(manager, metrics)

	orders := []entity.Order{
		{
			ID:    1,
			Drink: entity.DrinkLatte,
		},
		{
			ID:    2,
			Drink: entity.DrinkEspresso,
		},
		{
			ID:    3,
			Drink: entity.DrinkMatcha,
		},
		{
			ID:    4,
			Drink: entity.DrinkFrappe,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		usecase.ExecuteBrew(b.Context(), orders, 3)
	}
}
