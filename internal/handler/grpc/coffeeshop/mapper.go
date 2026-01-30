package coffeeshop

import (
	entity "gopher-cafe/internal/entity/coffeeshop"

	pb "github.com/rexyajaib/gopher-cafe/pkg/gen/go/v1"
)

func toEntityDrink(d pb.DrinkType) entity.DrinkType {
	switch d {
	case pb.DrinkType_DRINK_TYPE_ESPRESSO:
		return entity.DrinkEspresso
	case pb.DrinkType_DRINK_TYPE_LATTE:
		return entity.DrinkLatte
	case pb.DrinkType_DRINK_TYPE_FRAPPE:
		return entity.DrinkFrappe
	case pb.DrinkType_DRINK_TYPE_MATCHA:
		return entity.DrinkMatcha
	default:
		return entity.DrinkUnspecified
	}
}

func toPbEquipment(e entity.EquipmentType) pb.EquipmentType {
	switch e {
	case entity.EquipGrinder:
		return pb.EquipmentType_EQUIPMENT_TYPE_GRINDER
	case entity.EquipEspressoMachine:
		return pb.EquipmentType_EQUIPMENT_TYPE_ESPRESSO_MACHINE
	case entity.EquipMilkSteamer:
		return pb.EquipmentType_EQUIPMENT_TYPE_MILK_STEAMER
	case entity.EquipBlender:
		return pb.EquipmentType_EQUIPMENT_TYPE_BLENDER
	case entity.EquipWhisk:
		return pb.EquipmentType_EQUIPMENT_TYPE_WHISK
	default:
		return pb.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED
	}
}
