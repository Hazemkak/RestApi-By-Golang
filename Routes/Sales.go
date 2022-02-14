package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func SalesRoute(route fiber.Router) {
	route.Post("/new", Controllers.SalesNew)
	route.Post("/get_all", Controllers.SalesGetAll)
	route.Post("/get_all_populated", Controllers.SalesGetAllPopulated)
	route.Put("/modify/:id", Controllers.SalesModify)
	route.Put("/set_status/:type/:id/:new_status", Controllers.SetStatus)
	route.Post("new_sales_delivery", Controllers.SalesDeliveryNew)
	route.Post("modify_sales_delivery/:id", Controllers.SalesDeliveryModify)
	route.Put("/converted_to_sales/:id", Controllers.SetConvertedToSales)
}
