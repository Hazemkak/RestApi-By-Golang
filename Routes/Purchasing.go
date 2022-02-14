package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func PurchasingRoute(route fiber.Router) {
	route.Post("/new", Controllers.PurchasingNew)
	route.Post("/get_all", Controllers.PurchasingGetAll)
	route.Post("/get_all_populated", Controllers.PurchasingGetAllPopulated)
	route.Put("/modify/:id", Controllers.PurchasingModify)
	route.Put("/set_status/:type/:id/:new_status", Controllers.PurchasingSetStatus)
	route.Put("/converted_to_delivery/:id", Controllers.SetConvertedToDelivery)
}
