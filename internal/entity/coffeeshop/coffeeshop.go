package coffeeshop

import "time"

type EquipmentType int

const (
	EquipGrinder EquipmentType = iota
	EquipEspressoMachine
	EquipMilkSteamer
	EquipBlender
	EquipWhisk
)

type DrinkType int

const (
	DrinkUnspecified DrinkType = iota
	DrinkEspresso
	DrinkLatte
	DrinkFrappe
	DrinkMatcha
)

type RecipeStep struct {
	Equipment EquipmentType
	Duration  time.Duration
}

var Recipes = map[DrinkType][]RecipeStep{
	DrinkEspresso: {
		{EquipGrinder, 5 * time.Millisecond},
		{EquipEspressoMachine, 8 * time.Millisecond},
	},
	DrinkLatte: {
		{EquipGrinder, 5 * time.Millisecond},
		{EquipEspressoMachine, 8 * time.Millisecond},
		{EquipMilkSteamer, 15 * time.Millisecond},
	},
	DrinkFrappe: {
		{EquipGrinder, 5 * time.Millisecond},
		{EquipBlender, 12 * time.Millisecond},
	},
	DrinkMatcha: {
		{EquipGrinder, 5 * time.Millisecond},
		{EquipMilkSteamer, 15 * time.Millisecond},
		{EquipWhisk, 3 * time.Millisecond},
	},
}

type Order struct {
	ID    int64
	Drink DrinkType
}

type StepExecution struct {
	Equipment   EquipmentType
	StartTimeMs int64
	EndTimeMs   int64
}

type OrderResult struct {
	OrderID int64
	Steps   []StepExecution
}
