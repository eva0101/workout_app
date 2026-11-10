package workout_transport_http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_middleware "workout_app/internal/core/middleware"

	"github.com/google/uuid"
)

type workoutHistoryServiceStub struct {
	WorkoutService
	history         []core_domain.WorkoutHistoryItem
	requestedUser   uuid.UUID
	requestedLimit  int
	requestedOffset int
}

func (s *workoutHistoryServiceStub) GetWorkoutHistory(
	_ context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]core_domain.WorkoutHistoryItem, error) {
	s.requestedUser = userID
	s.requestedLimit = limit
	s.requestedOffset = offset
	return s.history, nil
}

func TestGetWorkoutHistoryUsesDefaultPaginationAndReturnsItems(t *testing.T) {
	name := "Snapshot day"
	percent, fatigueScore := 80, 7
	startedAt := time.Date(2026, 10, 15, 14, 0, 0, 0, time.UTC)
	completedAt := startedAt.Add(100 * time.Second)
	userID := uuid.New()
	service := &workoutHistoryServiceStub{
		history: []core_domain.WorkoutHistoryItem{{
			ID: 9, TrainingDayName: &name, StartedAt: startedAt, CompletedAt: completedAt,
			DurationSeconds: 100, CompletedSets: 2, TotalWeight: 160,
			ActualReps: 12, PlannedReps: 15, CompletionPercent: &percent, FatigueScore: &fatigueScore,
		}},
	}
	handler := NewWorkoutHTTPHandler(service, nil)
	request := httptest.NewRequest(http.MethodGet, "/workouts/history", nil)
	request = request.WithContext(context.WithValue(request.Context(), core_middleware.UserIDKey, userID))
	response := httptest.NewRecorder()

	handler.GetWorkoutHistory(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.requestedUser != userID ||
		service.requestedLimit != core_domain.WorkoutHistoryDefaultLimit ||
		service.requestedOffset != 0 {
		t.Fatalf("service received user=%s limit=%d offset=%d",
			service.requestedUser, service.requestedLimit, service.requestedOffset)
	}

	var body []WorkoutHistoryDTOResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode workout history response: %v", err)
	}
	if len(body) != 1 || body[0].ID != 9 || body[0].TrainingDayName == nil ||
		*body[0].TrainingDayName != name || body[0].CompletionPercent == nil ||
		*body[0].CompletionPercent != 80 || body[0].DurationSeconds != 100 ||
		body[0].FatigueScore == nil || *body[0].FatigueScore != 7 {
		t.Fatalf("unexpected workout history response: %+v", body)
	}
}

func TestWorkoutHistoryResponseSerializesMissingFatigueScoreAsNull(t *testing.T) {
	response, err := json.Marshal(ToWorkoutHistoryResponse(core_domain.WorkoutHistoryItem{ID: 9}))
	if err != nil {
		t.Fatalf("marshal workout history response: %v", err)
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(response, &body); err != nil {
		t.Fatalf("decode workout history response: %v", err)
	}
	if value, ok := body["fatigue_score"]; !ok || string(value) != "null" {
		t.Fatalf("fatigue_score = %s, want null", value)
	}
}

func TestGetWorkoutHistoryParsesPaginationAndRejectsInvalidQuery(t *testing.T) {
	service := &workoutHistoryServiceStub{}
	handler := NewWorkoutHTTPHandler(service, nil)
	userID := uuid.New()
	request := httptest.NewRequest(http.MethodGet, "/workouts/history?limit=12&offset=24", nil)
	request = request.WithContext(context.WithValue(request.Context(), core_middleware.UserIDKey, userID))
	response := httptest.NewRecorder()

	handler.GetWorkoutHistory(response, request)
	if response.Code != http.StatusOK || service.requestedLimit != 12 || service.requestedOffset != 24 {
		t.Fatalf("pagination not passed through: status=%d limit=%d offset=%d",
			response.Code, service.requestedLimit, service.requestedOffset)
	}

	for _, query := range []string{"?limit=abc", "?offset=-1"} {
		request := httptest.NewRequest(http.MethodGet, "/workouts/history"+query, nil)
		request = request.WithContext(context.WithValue(request.Context(), core_middleware.UserIDKey, userID))
		response := httptest.NewRecorder()
		handler.GetWorkoutHistory(response, request)
		if response.Code != http.StatusBadRequest {
			t.Errorf("query %q returned status %d, want 400", query, response.Code)
		}
	}
}

func TestGetWorkoutHistoryRequiresAuthentication(t *testing.T) {
	handler := NewWorkoutHTTPHandler(&workoutHistoryServiceStub{}, nil)
	request := httptest.NewRequest(http.MethodGet, "/workouts/history", nil)
	response := httptest.NewRecorder()

	handler.GetWorkoutHistory(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}
