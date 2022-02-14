package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func ProdStagesRoute(route fiber.Router) {
	route.Post("/new", Controllers.ProdStagesNew)
	route.Get("/get_all", Controllers.ProdStagesGetAll)
	route.Put("/modify", Controllers.ProdStagesModify)
	route.Put("/set_status/:id/:new_status", Controllers.ProdStagesSetStatus)

}
