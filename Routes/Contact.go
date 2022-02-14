package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func ContactRoute(route fiber.Router) {
	route.Put("/set_status/:id/:new_status", Controllers.ContactSetStatus)
	route.Put("/modify/:id", Controllers.ContactModify)
	route.Post("/get_all", Controllers.ContactGetAll)
	route.Post("/get_all_populated", Controllers.ContactGetAllPopulated)
	route.Post("/new", Controllers.ContactNew)
}
