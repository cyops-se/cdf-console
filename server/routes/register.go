package routes

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(api fiber.Router) {
	registerDataRoutes(api)
	registerFlowRoutes(api)
	registerDeviceRoutes(api)
	registerSystemRoutes(api)
	registerHelpRoutes(api)
}
