package controller

import (
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/lib/schedule_filter"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

// GetScheduleOptions handles the request to retrieve schedule options based on the provided filters.
//
// @Summary Get schedule options
// @Description Retrieve schedule options based on filters like branch, employee, service, and time.
// @Tags Schedule
// @Accept json
// @Produce json
// @Param branch_id query string false "Filter by branch ID"
// @Param employee_id query string false "Filter by employee ID"
// @Param service_id query string false "Filter by service ID"
// @Param get query string false "Specify what to retrieve: 'services', 'branches', 'employees', or 'time_slots'"
// @Success 200 {object} schedule_filter.ScheduleOptions "Successful operation"
// @Failure 400 {object} lib.ErrorResponse "Bad Request"
// @Failure 500 {object} lib.ErrorResponse "Internal Server Error"
// @Router /schedule/options [get]
func GetScheduleOptions(c *fiber.Ctx) error {
	var err error
	tx, def, err := database.ContextTransaction(c)
	if err != nil {
		return err
	}
	defer def(err)

	sf, err := ScheduleFilter.NewFromContext(tx, c)
	if err != nil {
		return err
	}

	options, err := sf.GetScheduleOptions()
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).Send(200, options)
}

func Schedule(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		GetScheduleOptions,
	})
}
