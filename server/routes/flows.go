package routes

import (
	"net/http"
	"server/devices"
	"server/flows"
	"server/messages"

	"github.com/gofiber/fiber/v2"
)

func registerFlowRoutes(api fiber.Router) {
	api.Get("/flows/white/allow", GetAllWhiteAllowFlows)
	api.Get("/flows/white/block", GetAllWhiteBlockFlows)
	api.Get("/flows/grey", GetAllGreyFlows)
	api.Get("/flows/black", GetAllBlackFlows)
	api.Get("/flows/ethertypes", GetAllEtherTypes)
}

func GetAllWhiteAllowFlows(c *fiber.Ctx) error {
	var response messages.RequestResponse
	response.Success = devices.DefaultAddressPresent
	return c.Status(http.StatusOK).JSON(response)
}

func GetAllWhiteBlockFlows(c *fiber.Ctx) error {
	var response messages.RequestResponse
	response.Success = devices.DefaultAddressPresent
	return c.Status(http.StatusOK).JSON(response)
}

func GetAllGreyFlows(c *fiber.Ctx) error {
	var response messages.EntryListResponse
	response.Success = devices.DefaultAddressPresent
	response.EntriesList = flows.GetAllEntriesAsList()
	return c.Status(http.StatusOK).JSON(response)
}

func GetAllBlackFlows(c *fiber.Ctx) error {
	var response messages.RequestResponse
	response.Success = devices.DefaultAddressPresent
	return c.Status(http.StatusOK).JSON(response)
}

func GetAllEtherTypes(c *fiber.Ctx) error {
	var response messages.EtherTypeListResponse
	response.Success = true
	response.EtherTypeList = flows.GetAllEtherTypesAsList()
	return c.Status(http.StatusOK).JSON(response)
}
