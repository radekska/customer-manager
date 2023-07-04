package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
)

func jsonContentValidator() fiber.Handler {
	methodsToValidate := []string{fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch}
	return func(ctx *fiber.Ctx) error {
		if !slices.Contains(methodsToValidate, ctx.Method()) {
			return ctx.Next()
		}
		if contentType := ctx.Get("Content-Type"); contentType != fiber.MIMEApplicationJSON {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf(
					"invalid content-type header specified: '%s', allowed: '%s'",
					contentType,
					fiber.MIMEApplicationJSON,
				),
			})
		}
		return ctx.Next()
	}
}
