package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func UnitsOfMeasurementRoute(route fiber.Router) {
	route.Post("/new", Controllers.UnitsOfMeasurementNew)
	route.Post("/add_relations/:to_id", Controllers.UnitsOfMeasurementAddRelations)
	route.Get("/get_all", Controllers.UnitsOfMeasurementGetAll)
	route.Get("/get_all/populated", Controllers.UnitsOfMeasurementGetAllPopulated)
	route.Get("/get_categories", Controllers.UnitsOfMeasurementGetDistinctCategories)
	route.Put("/set_status/:id/:new_status", Controllers.UnitsOfMeasurementSetStatus)
	route.Put("/set_status/:id/:embed_id/:new_status", Controllers.UnitsOfMeasurementSetRelationStatus)
	route.Post("/convert", Controllers.UnitsOfMeasurementConvertEP)
	route.Post("/get_names", Controllers.UnitsOfMeasurementGetDistinctUOM)
}
