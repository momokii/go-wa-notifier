package utils

import "github.com/gofiber/fiber/v2"

type MessageResponseSuccess struct {
	Error   bool   `json:"error" example:"false" default:"false"`
	Message string `json:"message"`
}

type MessageResponseError struct {
	Error   bool   `json:"error" example:"true" default:"true"`
	Message string `json:"message"`
}

type DataResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   bool        `json:"error" example:"false"`
}

func ResponseMessage(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(&MessageResponseSuccess{
		Error:   false,
		Message: message,
	})
}

func ResponseWitData(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(&DataResponse{
		Data:    data,
		Message: message,
		Error:   false,
	})
}

func ResponseError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(&MessageResponseError{
		Error:   true,
		Message: message,
	})
}
