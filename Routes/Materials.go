package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func MaterialRoute(route fiber.Router) {
	route.Post("/new", Controllers.MaterialNew)
	route.Get("/get_all", Controllers.MaterialGetAll)
	route.Put("/set_status/:id/:new_status", Controllers.MaterialSetStatus)
	route.Put("/modify", Controllers.MaterialModify)
	route.Get("/get_all/populated", Controllers.MaterialGetAllPopulated)
}
