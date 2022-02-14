package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func PriceListRoute(route fiber.Router) {
	route.Put("/set_status/:id/:new_status", Controllers.PriceListSetStatus)
	route.Put("/modify/:id", Controllers.PriceListModify)
	route.Post("/get_all", Controllers.PriceListGetAll)
	route.Post("/get_all_populated", Controllers.PriceListGetAllPopulated)
	route.Post("/new", Controllers.PriceListNew)
}
