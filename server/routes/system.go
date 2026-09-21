package routes

import (
	"net/http"
	"server/messages"

	"github.com/gofiber/fiber/v2"
)

var SysInfo messages.SystemInformation

func registerSystemRoutes(api fiber.Router) {
	api.Get("/system/sysinfo", GetSysInfo)
}

func GetSysInfo(c *fiber.Ctx) error {
	var response messages.SysInfoResponse
	response.Success = true
	response.SysInfo = SysInfo
	return c.Status(http.StatusOK).JSON(response)
}
