package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func InventoryRequestsRoute(route fiber.Router) {
	route.Post("/new", Controllers.InventoryDemandNew)
	route.Post("/get_all", Controllers.InventoryDemandGetAll)
	route.Post("/get_all_populated", Controllers.InventoryDemandGetAllPopulated)
	route.Put("/reject/:id", Controllers.InventoryDemandReject)
	route.Put("/order/:id", Controllers.InventoryDemandOrderPhase)
	route.Put("/delivery/:id", Controllers.InventoryDemandDeliveryPhase)
	route.Delete("/delete/:id", Controllers.InventoryDemandDelete)
} //2500 & 350
