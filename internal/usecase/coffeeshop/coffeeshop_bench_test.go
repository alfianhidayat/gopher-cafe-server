package coffeeshop

import (
	entity "gopher-cafe/internal/entity/coffeeshop"
	"gopher-cafe/internal/worker"
	"testing"
)

//BenchmarkExecuteBrew-8   	     100	 105996649 ns/op	    7335 B/op	     149 allocs/op
//BenchmarkExecuteBrew-8   	     100	 105234255 ns/op	    7897 B/op	     152 allocs/op
//BenchmarkExecuteBrew-8   	     100	 105521863 ns/op	    7318 B/op	     149 allocs/op

//BenchmarkExecuteBrew-8   	     100	 106223778 ns/op	    7499 B/op	     151 allocs/op
//BenchmarkExecuteBrew-8   	     100	 104450702 ns/op	    7646 B/op	      94 allocs/op
//BenchmarkExecuteBrew-8   	     100	 107174945 ns/op	    7355 B/op	     150 allocs/op

func BenchmarkExecuteBrew(b *testing.B) {
	ew := worker.EquipmentWorkers

	manager := worker.NewEquipPoolManager(uint8(len(ew)))
	for k, v := range ew {
		manager.Register(k, v)
	}
	manager.StartAll()

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
		{
			ID:    5,
			Drink: entity.DrinkMatcha,
		},
		{
			ID:    6,
			Drink: entity.DrinkEspresso,
		},
		{
			ID:    7,
			Drink: entity.DrinkMatcha,
		},
		{
			ID:    8,
			Drink: entity.DrinkLatte,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		usecase.ExecuteBrew(b.Context(), orders, 2)
	}
}
