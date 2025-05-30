package mJSON

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type WorkSchedule struct {
	Monday    []WorkRange `json:"monday"`
	Tuesday   []WorkRange `json:"tuesday"`
	Wednesday []WorkRange `json:"wednesday"`
	Thursday  []WorkRange `json:"thursday"`
	Friday    []WorkRange `json:"friday"`
	Saturday  []WorkRange `json:"saturday"`
	Sunday    []WorkRange `json:"sunday"`
}

type WorkRange struct {
	Start    string      `json:"start"`
	End      string      `json:"end"`
	BranchID uuid.UUID   `json:"branch_id"`
	Services []uuid.UUID `json:"services"` // List of service IDs
}

// Implement driver.Valuer
func (ws WorkSchedule) Value() (driver.Value, error) {
	return json.Marshal(ws)
}

// Implement sql.Scanner
func (ws *WorkSchedule) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan WorkSchedule: expected []byte")
	}

	return json.Unmarshal(bytes, ws)
}

func (ws *WorkSchedule) IsEmpty() bool {
	if ws == nil {
		return true
	}
	return len(ws.Monday) == 0 && len(ws.Tuesday) == 0 && len(ws.Wednesday) == 0 &&
		len(ws.Thursday) == 0 && len(ws.Friday) == 0 && len(ws.Saturday) == 0 &&
		len(ws.Sunday) == 0
}

func (ws *WorkSchedule) GetRangesForDay(day time.Weekday) []WorkRange {
	if ws == nil {
		return nil
	}
	switch day {
	case time.Monday:
		return ws.Monday
	case time.Tuesday:
		return ws.Tuesday
	case time.Wednesday:
		return ws.Wednesday
	case time.Thursday:
		return ws.Thursday
	case time.Friday:
		return ws.Friday
	case time.Saturday:
		return ws.Saturday
	case time.Sunday:
		return ws.Sunday
	default:
		return nil
	}
}

func (ws *WorkSchedule) GetAllRanges() []WorkRange {
	if ws == nil {
		return nil
	}
	ranges := []WorkRange{}
	ranges = append(ranges, ws.Monday...)
	ranges = append(ranges, ws.Tuesday...)
	ranges = append(ranges, ws.Wednesday...)
	ranges = append(ranges, ws.Thursday...)
	ranges = append(ranges, ws.Friday...)
	ranges = append(ranges, ws.Saturday...)
	ranges = append(ranges, ws.Sunday...)
	return ranges
}