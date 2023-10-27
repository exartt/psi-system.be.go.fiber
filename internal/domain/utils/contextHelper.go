package utils

import "github.com/gofiber/fiber/v2"

func GetPsychologistIDFromContext(c *fiber.Ctx) (uint, error) {
	rawPsychologistID := c.Locals("psychologistID")
	if rawPsychologistID == nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Psychologist ID not found in context")
	}

	floatPsychologistID, ok := rawPsychologistID.(float64)
	if !ok {
		return 0, fiber.NewError(fiber.StatusInternalServerError, "Could not cast psychologistID to float64")
	}

	psychologistID := uint(floatPsychologistID)

	return psychologistID, nil
}
