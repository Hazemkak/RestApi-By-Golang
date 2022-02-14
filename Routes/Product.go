package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductRoute(route fiber.Router) {
	route.Post("/new", Controllers.ProductNew)
	route.Post("/get_all", Controllers.ProductGetAll)
	route.Post("/get_all_populated", Controllers.ProductGetAllPopulated)
	route.Put("/set_status/:id/:new_status", Controllers.ProductSetStatus)
	route.Post("/get_categories", Controllers.ProductGetDistinctCategories)
	route.Put("/modify/:id", Controllers.ProductModify)
	route.Post("/add_BOM/:id", Controllers.ProductAddBOMNew)
}
