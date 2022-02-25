package utils

import "github.com/gofiber/fiber/v2"

func ErrorMsg(c *fiber.Ctx, err string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": true,
		"msg":   err,
	})
}

func SuccessMsg(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   msg,
	})
}
