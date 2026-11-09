package workout_transport_http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"
	core_middleware "workout_app/internal/core/middleware"
	core_http_server "workout_app/internal/core/transport/http/server"

	"github.com/google/uuid"
)

type workoutReadServiceStub struct {
	WorkoutService
	workout core_domain.WorkoutDetails
	err     error
}

func (s workoutReadServiceStub) GetActiveWorkout(
	context.Context,
	uuid.UUID,
) (core_domain.WorkoutDetails, error) {
	return s.workout, s.err
}

func (s workoutReadServiceStub) GetWorkout(
	context.Context,
	uuid.UUID,
	int,
) (core_domain.WorkoutDetails, error) {
	return s.workout, s.err
}

func TestGetWorkoutRoutesAreAuthenticated(t *testing.T) {
	handler := NewWorkoutHTTPHandler(nil, nil)
	routes := handler.Routes()
	router := core_http_server.NewRouter()
	router.RegisterRoutes(func(next http.Handler) http.Handler { return next }, routes...)

	want := map[string]bool{
		http.MethodGet + " /workouts/active": false,
		http.MethodGet + " /workouts/{id}":   false,
	}
	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, ok := want[key]; ok {
			want[key] = route.Auth
		}
	}
	for route, authenticated := range want {
		if !authenticated {
			t.Errorf("route %s must require authentication", route)
		}
	}
}

func TestWorkoutResponseIncludesNestedSetsAndEmptyArrays(t *testing.T) {
	repetitions, planWeight, repsDone, actualWeight := 10, 25.0, 9, 22.5
	workout := core_domain.WorkoutDetails{
		ID: 7, Status: "completed",
		Exercises: []core_domain.WorkoutExerciseSnapshot{
			{
				ID: 12, ExerciseID: 3, Name: "press", Title: "Press", Order: 1,
				PlannedSets: []core_domain.WorkoutPlannedSet{
					{ID: 20, SetNumber: 1, Repetitions: &repetitions, Weight: &planWeight},
				},
				ActualSets: []core_domain.WorkoutActualSet{
					{ID: 30, SetNumber: 1, RepsDone: &repsDone, Weight: &actualWeight},
				},
			},
			{
				ID: 13, ExerciseID: 4, Name: "row", Title: "Row", Order: 2,
				PlannedSets: []core_domain.WorkoutPlannedSet{},
				ActualSets:  []core_domain.WorkoutActualSet{},
			},
		},
	}
	handler := NewWorkoutHTTPHandler(workoutReadServiceStub{workout: workout}, nil)
	request := httptest.NewRequest(http.MethodGet, "/workouts/7", nil)
	request.SetPathValue("id", "7")
	request = request.WithContext(context.WithValue(request.Context(), core_middleware.UserIDKey, uuid.New()))
	response := httptest.NewRecorder()

	handler.GetWorkout(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}

	var body WorkoutDTOResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode workout response: %v", err)
	}
	if body.ID != 7 || body.Status != "completed" || len(body.Exercises) != 2 {
		t.Fatalf("unexpected workout response: %+v", body)
	}
	if body.Exercises[0].ID != 12 || len(body.Exercises[0].PlannedSets) != 1 || len(body.Exercises[0].ActualSets) != 1 {
		t.Fatalf("first exercise or its sets missing: %+v", body.Exercises[0])
	}
	if body.Exercises[1].PlannedSets == nil || body.Exercises[1].ActualSets == nil {
		t.Fatalf("empty set arrays serialized as null: %+v", body.Exercises[1])
	}
}

func TestGetActiveWorkoutReturns404WhenMissing(t *testing.T) {
	handler := NewWorkoutHTTPHandler(
		workoutReadServiceStub{err: core_errors.ErrActiveWorkoutNotFound},
		nil,
	)
	request := httptest.NewRequest(http.MethodGet, "/workouts/active", nil)
	request = request.WithContext(context.WithValue(request.Context(), core_middleware.UserIDKey, uuid.New()))
	response := httptest.NewRecorder()

	handler.GetActiveWorkout(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}
}
