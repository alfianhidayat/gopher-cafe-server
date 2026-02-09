package coffeeshop

import entity "gopher-cafe/internal/entity/coffeeshop"

type OrderInput struct {
	data          entity.Order
	resultChannel chan entity.OrderResult
}
