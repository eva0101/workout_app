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

func TestCreateExerciseResponseContainsOptionalPlannedSetFields(t *testing.T) {
	repetitions := 10
	weight := 40.5

	response := CreateExerciseDTOResponse{
		ID:           4,
		ExerciseID:   12,
		NameExercise: "barbell bench press",
		Sets:         1,
		PlannedSets: ToPlannedSetResponse([]core_domain.PlannedSet{{
			SetNumber:   1,
			Repetitions: &repetitions,
			Weight:      &weight,
		}}),
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal created exercise: %v", err)
	}

	jsonValue := string(data)
	for _, field := range []string{
		`"exercise_id":12`,
		`"planned_sets":[{"set_number":1,"repetitions":10,"weight":40.5}]`,
	} {
		if !strings.Contains(jsonValue, field) {
			t.Errorf("expected JSON to contain %s, got %s", field, jsonValue)
		}
	}
}

func TestPatchExerciseSetsRequestUsesRepsAndOptionalWeight(t *testing.T) {
	var request PatchExerciseSetsDTORequest
	if err := json.Unmarshal([]byte(`{"sets":[{"reps":10},{"weight":40.5}]}`), &request); err != nil {
		t.Fatalf("unmarshal patch exercise sets request: %v", err)
	}
	if len(request.Sets) != 2 {
		t.Fatalf("expected 2 sets, got %d", len(request.Sets))
	}
	if request.Sets[0].Reps == nil || *request.Sets[0].Reps != 10 || request.Sets[0].Weight != nil {
		t.Fatalf("unexpected first set: %+v", request.Sets[0])
	}
	if request.Sets[1].Reps != nil || request.Sets[1].Weight == nil || *request.Sets[1].Weight != 40.5 {
		t.Fatalf("unexpected second set: %+v", request.Sets[1])
	}

	if err := json.Unmarshal([]byte(`{"sets":[]}`), &request); err != nil {
		t.Fatalf("unmarshal empty sets request: %v", err)
	}
	if request.Sets == nil || len(request.Sets) != 0 {
		t.Fatalf("expected a non-nil empty set list, got %#v", request.Sets)
	}
}

func TestCatalogResponsesUseRequestedJSONFields(t *testing.T) {
	imagePath := "/exercises/groups/chest.webp"
	description := "Technique"

	groupJSON, err := json.Marshal(ToMuscleGroupResponse([]core_domain.MuscleGroup{{
		ID:        1,
		Name:      "chest",
		Title:     "Грудь",
		ImagePath: &imagePath,
	}}))
	if err != nil {
		t.Fatalf("marshal muscle groups: %v", err)
	}
	for _, field := range []string{`"id":1`, `"name":"chest"`, `"title":"Грудь"`, `"image_path":"/exercises/groups/chest.webp"`} {
		if !strings.Contains(string(groupJSON), field) {
			t.Errorf("expected muscle group JSON to contain %s, got %s", field, groupJSON)
		}
	}

	exerciseJSON, err := json.Marshal(ToCatalogExerciseResponse([]core_domain.CatalogExercise{{
		ID:            12,
		Name:          "barbell bench press",
		Title:         "Жим штанги лёжа",
		MuscleGroupID: 1,
		Description:   &description,
		ImagePath:     &imagePath,
	}}))
	if err != nil {
		t.Fatalf("marshal exercises: %v", err)
	}
	for _, field := range []string{`"id":12`, `"name":"barbell bench press"`, `"title":"Жим штанги лёжа"`, `"muscle_group_id":1`, `"description":"Technique"`, `"image_path":"/exercises/groups/chest.webp"`} {
		if !strings.Contains(string(exerciseJSON), field) {
			t.Errorf("expected exercise JSON to contain %s, got %s", field, exerciseJSON)
		}
	}
}
