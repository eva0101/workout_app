package workout_transport_http

import (
	"encoding/json"
	"net/http"
	core_domain "workout_app/internal/core/domain"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetStatisticsDashboard godoc
// @Summary     Получение статистики тренировок для вкладки «Статистика»
// @Description Возвращает общую сводку, регулярность, динамику нагрузки и группы мышц по завершённым тренировкам пользователя
// @Tags        workout
// @Produce     json
// @Security    BearerAuth
// @Param       period query string false "Период статистики" Enums(7d, 30d, 3m, 1y, all) default(30d)
// @Success     200 {object} WorkoutStatisticsDashboardDTOResponse "Статистика тренировок"
// @Failure     400 {string} string "Неизвестный период"
// @Failure     401 {string} string "Unauthorized"
// @Failure     500 {string} string "Internal server error"
// @Router      /workouts/statistics/dashboard [get]
func (h *WorkoutHTTPHandler) GetStatisticsDashboard(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}
	statistics, err := h.workoutService.GetStatisticsDashboard(r.Context(), userID, period)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	response := toWorkoutStatisticsDashboardResponse(statistics)
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.log.Error("failed to encode workout statistics dashboard", zap.Error(err))
	}
}

func toWorkoutStatisticsDashboardResponse(
	statistics core_domain.WorkoutStatisticsDashboard,
) WorkoutStatisticsDashboardDTOResponse {
	response := WorkoutStatisticsDashboardDTOResponse{
		Summary: WorkoutStatisticsSummaryDTOResponse{
			WorkoutsCount:        statistics.Summary.WorkoutsCount,
			TotalDurationSeconds: statistics.Summary.TotalDurationSeconds,
			CompletedSets:        statistics.Summary.CompletedSets,
			TotalWeight:          statistics.Summary.TotalWeight,
		},
		Regularity: WorkoutStatisticsRegularityDTOResponse{
			Interval: statistics.Interval,
			Data:     make([]WorkoutStatisticsRegularityItemDTOResponse, 0, len(statistics.Regularity)),
		},
		LoadDynamics: WorkoutStatisticsLoadDynamicsDTOResponse{
			Interval: statistics.Interval,
			Data:     make([]WorkoutStatisticsLoadItemDTOResponse, 0, len(statistics.LoadDynamics)),
		},
		MuscleGroups: make([]WorkoutStatisticsMuscleGroupDTOResponse, 0, len(statistics.MuscleGroups)),
	}
	for _, entry := range statistics.Regularity {
		response.Regularity.Data = append(response.Regularity.Data, WorkoutStatisticsRegularityItemDTOResponse{
			Date:          entry.Date,
			WorkoutsCount: entry.WorkoutsCount,
		})
	}
	for _, entry := range statistics.LoadDynamics {
		response.LoadDynamics.Data = append(response.LoadDynamics.Data, WorkoutStatisticsLoadItemDTOResponse{
			Date:            entry.Date,
			TotalWeight:     entry.TotalWeight,
			CompletedSets:   entry.CompletedSets,
			ActualReps:      entry.ActualReps,
			DurationSeconds: entry.DurationSeconds,
		})
	}
	for _, group := range statistics.MuscleGroups {
		response.MuscleGroups = append(response.MuscleGroups, WorkoutStatisticsMuscleGroupDTOResponse{
			Name:          group.Name,
			Title:         group.Title,
			CompletedSets: group.CompletedSets,
		})
	}
	return response
}
