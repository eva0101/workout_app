import { CalendarDays, ChevronRight, Dumbbell, LoaderCircle, Play, RotateCcw } from 'lucide-react'
import { useCallback, useEffect, useMemo, useState } from 'react'
import WorkoutCalendarDialog from '../components/WorkoutCalendarDialog'
import {
  getScheduleWithDayNames,
  type ScheduledWorkoutWithDayName,
  type Workout,
} from '../lib/api'

type DashboardPageProps = {
  activeWorkout: Workout | null
  activeWorkoutError: string
  isActiveWorkoutLoading: boolean
  isStartingWorkout: boolean
  onStartScheduledWorkout: (scheduleId: number, title: string) => Promise<void>
  onResumeWorkout: () => void
  onOpenPrograms: () => void
}

function getLocalDate() {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function addDays(dateKey: string, days: number) {
  const [year, month, day] = dateKey.split('-').map(Number)
  const date = new Date(year, month - 1, day + days)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function formatScheduleDate(value: string) {
  const [year, month, day] = value.split('-').map(Number)
  return new Intl.DateTimeFormat('ru-RU', { weekday: 'short', day: 'numeric', month: 'short' })
    .format(new Date(year, month - 1, day))
}

function getErrorMessage(error: unknown) {
  if (typeof error === 'object' && error !== null && 'status' in error) {
    const status = (error as { status?: number }).status
    if (status === 404) return 'Расписание не найдено. Обновите страницу и попробуйте снова.'
  }
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте подключение.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось загрузить расписание.'
}

export default function DashboardPage({
  activeWorkout,
  activeWorkoutError,
  isActiveWorkoutLoading,
  isStartingWorkout,
  onStartScheduledWorkout,
  onResumeWorkout,
  onOpenPrograms,
}: DashboardPageProps) {
  const today = getLocalDate()
  const [scheduledWorkouts, setScheduledWorkouts] = useState<ScheduledWorkoutWithDayName[]>([])
  const [isScheduleLoading, setIsScheduleLoading] = useState(true)
  const [scheduleError, setScheduleError] = useState('')
  const [startError, setStartError] = useState('')
  const [isCalendarOpen, setIsCalendarOpen] = useState(false)
  const [calendarDate, setCalendarDate] = useState(today)

  const loadSchedule = useCallback(async () => {
    setIsScheduleLoading(true)
    setScheduleError('')
    try {
      setScheduledWorkouts(await getScheduleWithDayNames('0001-01-01', '9999-12-31'))
    } catch (error) {
      setScheduleError(getErrorMessage(error))
    } finally {
      setIsScheduleLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadSchedule()
  }, [loadSchedule])

  const todayWorkouts = useMemo(
    () => scheduledWorkouts.filter((workout) => workout.scheduled_date === today),
    [scheduledWorkouts, today],
  )
  const upcomingWorkouts = useMemo(
    () => scheduledWorkouts
      .filter((workout) => workout.scheduled_date > today)
      .sort((left, right) => left.scheduled_date.localeCompare(right.scheduled_date))
      .slice(0, 4),
    [scheduledWorkouts, today],
  )
  const nextSevenDaysCount = scheduledWorkouts.filter(
    (workout) => workout.scheduled_date >= today && workout.scheduled_date <= addDays(today, 6),
  ).length
  const canStart = !activeWorkout && !isActiveWorkoutLoading && !activeWorkoutError && !isStartingWorkout

  async function startWorkout(workout: ScheduledWorkoutWithDayName) {
    if (!canStart) return
    setStartError('')
    try {
      await onStartScheduledWorkout(
        workout.id,
        `${workout.program_name} — ${workout.training_day_name}`,
      )
    } catch (error) {
      setStartError(getErrorMessage(error))
    }
  }

  function openCalendar(date = today) {
    setCalendarDate(date)
    setIsCalendarOpen(true)
  }

  return (
    <section className="dashboard-content">
      <div className="section-heading dashboard-section-heading">
        <div>
          <p className="eyebrow">ВАШ ПЛАН</p>
          <h2 className="mt-1 text-[20px] font-semibold tracking-[-0.035em] text-ink">Обзор</h2>
        </div>
        <button className="text-button" type="button" onClick={() => openCalendar()}>
          Календарь <CalendarDays size={15} />
        </button>
      </div>

      {activeWorkoutError && <p className="program-alert mb-4" role="alert">{activeWorkoutError}</p>}

      <div className="dashboard-overview-grid">
        <section className="dashboard-hero-card" aria-labelledby="dashboard-today-title">
          <div className="dashboard-hero-copy">
            <p className="eyebrow">
              {activeWorkout ? 'ТРЕНИРОВКА В ПРОЦЕССЕ' : 'СЕГОДНЯ'}
            </p>
            <h2 id="dashboard-today-title">
              {activeWorkout
                ? activeWorkout.training_day_name || 'Продолжите тренировку'
                : todayWorkouts.length > 0
                  ? 'Пора двигаться'
                  : 'План на сегодня'}
            </h2>

            {isActiveWorkoutLoading && !activeWorkout ? (
              <p className="dashboard-hero-description">
                <LoaderCircle className="spin inline-block" size={15} /> Проверяем активную тренировку...
              </p>
            ) : activeWorkout ? (
              <>
                <p className="dashboard-hero-description">
                  Тренировка сохранена. Продолжите с того места, где остановились.
                </p>
                <button className="program-primary-button mt-5" type="button" onClick={onResumeWorkout}>
                  <RotateCcw size={15} /> Продолжить тренировку
                </button>
              </>
            ) : isScheduleLoading ? (
              <p className="dashboard-hero-description">
                <LoaderCircle className="spin inline-block" size={15} /> Загружаем расписание...
              </p>
            ) : scheduleError ? (
              <>
                <p className="dashboard-hero-description">Не удалось проверить план тренировок.</p>
                <button className="program-secondary-button mt-4" type="button" onClick={() => void loadSchedule()}>
                  Повторить загрузку
                </button>
              </>
            ) : todayWorkouts.length > 0 ? (
              <ul className="dashboard-today-list">
                {todayWorkouts.map((workout) => (
                  <li className="dashboard-today-item" key={workout.id}>
                    <span className="min-w-0">
                      <span className="dashboard-today-title">
                        {workout.program_name} — {workout.training_day_name}
                      </span>
                      <span className="dashboard-hero-description">Тренировка запланирована на сегодня</span>
                    </span>
                    <button
                      className="program-primary-button"
                      type="button"
                      disabled={!canStart}
                      onClick={() => void startWorkout(workout)}
                    >
                      {isStartingWorkout
                        ? <LoaderCircle className="spin" size={15} />
                        : <Play size={15} fill="currentColor" />}
                      Начать
                    </button>
                  </li>
                ))}
              </ul>
            ) : (
              <>
                <p className="dashboard-hero-description">
                  {scheduledWorkouts.length > 0
                    ? 'На сегодня ничего не запланировано. Вы можете посмотреть другие даты.'
                    : 'В календаре пока нет тренировок. Сначала запланируйте занятие, чтобы оно появилось здесь.'}
                </p>
                <div className="mt-5 flex flex-wrap gap-2">
                  <button className="program-secondary-button" type="button" onClick={() => openCalendar()}>
                    <CalendarDays size={15} />
                    {scheduledWorkouts.length > 0 ? 'Выбрать другую тренировку' : 'Открыть календарь'}
                  </button>
                  {scheduledWorkouts.length === 0 && (
                    <button className="text-button" type="button" onClick={onOpenPrograms}>
                      Перейти к программам <ChevronRight size={15} />
                    </button>
                  )}
                </div>
              </>
            )}
            {startError && <p className="program-alert mt-4" role="alert">{startError}</p>}
          </div>
          <div className="dashboard-hero-art" aria-hidden="true">
            <span className="dashboard-art-ring dashboard-art-ring-outer" />
            <span className="dashboard-art-ring dashboard-art-ring-inner" />
            <span className="dashboard-art-icon"><Dumbbell size={39} strokeWidth={1.4} /></span>
          </div>
        </section>

        <section className="dashboard-side-card panel-card" aria-labelledby="dashboard-upcoming-title">
          <div className="dashboard-widget-heading">
            <div>
              <p className="eyebrow">ВАШ КАЛЕНДАРЬ</p>
              <h3 id="dashboard-upcoming-title">Ближайшие тренировки</h3>
            </div>
            <button
              className="program-icon-action"
              type="button"
              aria-label="Открыть календарь"
              onClick={() => openCalendar()}
            >
              <CalendarDays size={17} />
            </button>
          </div>
          {isScheduleLoading ? (
            <div className="program-state"><LoaderCircle className="spin" size={19} /> Загружаем...</div>
          ) : scheduleError ? (
            <p className="dashboard-widget-empty">Расписание пока недоступно.</p>
          ) : upcomingWorkouts.length > 0 ? (
            <ul className="dashboard-upcoming-list">
              {upcomingWorkouts.map((workout) => (
                <li key={workout.id}>
                  <button
                    className="dashboard-upcoming-button"
                    type="button"
                    onClick={() => openCalendar(workout.scheduled_date)}
                  >
                    <span className="dashboard-upcoming-date">{formatScheduleDate(workout.scheduled_date)}</span>
                    <span className="dashboard-upcoming-name">
                      {workout.program_name} — {workout.training_day_name}
                    </span>
                    <ChevronRight size={15} />
                  </button>
                </li>
              ))}
            </ul>
          ) : (
            <p className="dashboard-widget-empty">Следующие занятия появятся здесь после планирования в календаре.</p>
          )}
          <div className="dashboard-widget-footer">
            <span>{nextSevenDaysCount} тренировок на ближайшие 7 дней</span>
            <button className="text-button" type="button" onClick={() => openCalendar()}>
              Все даты <ChevronRight size={14} />
            </button>
          </div>
        </section>
      </div>

      {activeWorkout && todayWorkouts.length > 0 && (
        <p className="dashboard-active-note">
          Сегодня запланировано занятий: {todayWorkouts.length}. Завершите текущую тренировку, чтобы начать следующую.
        </p>
      )}

      <WorkoutCalendarDialog
        open={isCalendarOpen}
        workouts={scheduledWorkouts}
        loadError={scheduleError}
        isLoading={isScheduleLoading}
        isStarting={isStartingWorkout}
        hasActiveWorkout={Boolean(activeWorkout)}
        initialDate={calendarDate}
        onClose={() => setIsCalendarOpen(false)}
        onRetry={() => void loadSchedule()}
        onOpenPrograms={() => {
          setIsCalendarOpen(false)
          onOpenPrograms()
        }}
        onStartWorkout={onStartScheduledWorkout}
      />
    </section>
  )
}
