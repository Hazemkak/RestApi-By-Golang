package Routes

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Controllers"

	"github.com/gofiber/fiber/v2"
)

func SettingRoute(route fiber.Router) {
	route.Post("/get_all", Controllers.SettingGetAll)
	route.Post("/new", Controllers.SettingNew)
}
