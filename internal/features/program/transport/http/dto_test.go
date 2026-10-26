package program_transport_http

import (
	"encoding/json"
	"strings"
	"testing"

	core_domain "workout_app/internal/core/domain"
)

func TestTrainingDaySummaryResponseUsesJSONNames(t *testing.T) {
	response := ToTrainingDaySummaryResponse(core_domain.TrainingDays{
		ID:        4,
		ProgramID: 2,
		DayNumber: 1,
		Name:      "День спины",
	})

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal training day summary: %v", err)
	}

	jsonValue := string(data)
	for _, field := range []string{`"id":4`, `"program_id":2`, `"day_number":1`, `"name":"День спины"`} {
		if !strings.Contains(jsonValue, field) {
			t.Errorf("expected JSON to contain %s, got %s", field, jsonValue)
		}
	}
}

func TestTrainingDayResponseIncludesName(t *testing.T) {
	response := ToTrainingDayResponse(core_domain.TrainingDay{
		ID:        4,
		DayNumber: 1,
		Name:      "День спины",
	})

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal training day: %v", err)
	}

	jsonValue := string(data)
	for _, field := range []string{`"id":4`, `"day_number":1`, `"name":"День спины"`, `"exercises":[]`} {
		if !strings.Contains(jsonValue, field) {
			t.Errorf("expected JSON to contain %s, got %s", field, jsonValue)
		}
	}
}
