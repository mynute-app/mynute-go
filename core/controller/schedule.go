package controller

import (
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/lib/schedule_filter"

	"github.com/gofiber/fiber/v2"
)

func GetScheduleOptions(c *fiber.Ctx) error {
	branch_id := c.Query("branch_id")
	employee_id := c.Query("employee_id")
	service_id := c.Query("service_id")
	weekday := c.Query("weekday")
	get := c.Query("get")
	time := c.Query("time")

	tx, def, err := database.ContextTransaction(c)
	if err != nil {
		return err
	}
	defer def()

	ScheduleFilter := schedule_filter.NewScheduleFilter(tx)

	options, err := ScheduleFilter.GetScheduleOptions(get, branch_id, employee_id, service_id, weekday, time)
	if err != nil {
		return err
	}

	return c.JSON(options)
}
