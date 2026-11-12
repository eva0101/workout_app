import {
  CalendarDays,
  Check,
  ChevronRight,
  Clock3,
  Dumbbell,
  LoaderCircle,
  Play,
  RotateCcw,
} from 'lucide-react'
import { useCallback, useEffect, useMemo, useState } from 'react'
import WeeklyStats from '../components/WeeklyStats'
import WorkoutCalendarDialog from '../components/WorkoutCalendarDialog'
import {
  getPrograms,
  getWorkoutHistory,
  getScheduleWithDayNames,
  getWorkoutStatisticsDashboard,
  type Program,
  type ScheduledWorkoutWithDayName,
  type Workout,
  type WorkoutHistoryItem,
  type WorkoutStatisticsDashboard,
} from '../lib/api'

type DashboardPageProps = {
  onWelcomeMessageChange: (message: DashboardWelcomeMessage) => void
  activeWorkout: Workout | null
  activeWorkoutError: string
  isActiveWorkoutLoading: boolean
  isStartingWorkout: boolean
  onStartScheduledWorkout: (scheduleId: number, title: string) => Promise<void>
  onResumeWorkout: () => void
  onOpenPrograms: () => void
}

export type DashboardWelcomeMessage = {
  title: string
  description: string
}

const WELCOME_MESSAGES = {
  welcome: {
    title: 'Добро пожаловать',
    description: 'Создайте свою первую программу тренировок и начните путь к новым результатам.',
  },
  rest: {
    title: 'Сегодня день отдыха',
    description: 'Отдыхайте и набирайтесь сил.',
  },
  workout: {
    title: 'Время двигаться вперёд',
    description: 'Ваш тренировочный план готов. Продолжайте работать над собой.',
  },
  completed: {
    title: 'Отличная работа!',
    description: 'Сегодняшняя тренировка завершена. Время восстановиться.',
  },
} satisfies Record<string, DashboardWelcomeMessage>

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

function getWeekDates(today: string) {
  const [year, month, day] = today.split('-').map(Number)
  const date = new Date(year, month - 1, day)
  const mondayOffset = (date.getDay() + 6) % 7
  const monday = addDays(today, -mondayOffset)
  return Array.from({ length: 7 }, (_, index) => addDays(monday, index))
}

function getCurrentWeekSummary(statistics: WorkoutStatisticsDashboard, weekDates: string[], today: string) {
  const currentWeekDates = new Set(weekDates.filter((date) => date <= today))
  const workoutsByDate = new Map(
    statistics.regularity.data.map(({ date, workouts_count }) => [date, workouts_count]),
  )
  const loadByDate = new Map(
    statistics.load_dynamics.data.map(({ date, completed_sets, duration_seconds }) => [
      date,
      { completedSets: completed_sets, durationSeconds: duration_seconds },
    ]),
  )

  return [...currentWeekDates].reduce(
    (summary, date) => {
      summary.workoutsCount += workoutsByDate.get(date) ?? 0
      const load = loadByDate.get(date)
      summary.completedSets += load?.completedSets ?? 0
      summary.totalDurationSeconds += load?.durationSeconds ?? 0
      return summary
    },
    { workoutsCount: 0, completedSets: 0, totalDurationSeconds: 0 },
  )
}

function getLocalDateKey(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

async function getCurrentWeekHistory(weekStart: string) {
  const history: WorkoutHistoryItem[] = []
  let offset = 0

  while (true) {
    const page = await getWorkoutHistory(100, offset)
    history.push(...page)
    if (page.length < 100 || getLocalDateKey(page[page.length - 1].completed_at) < weekStart) {
      return history
    }
    offset += page.length
  }
}

function formatWorkoutDate(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Дата недоступна'
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

function formatWorkoutDuration(seconds: number) {
  const totalSeconds = Math.max(0, Math.floor(seconds))
  const totalMinutes = Math.floor(totalSeconds / 60)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  if (hours > 0) return `${hours} ч ${minutes} мин`
  if (minutes > 0) return `${minutes} мин`
  return `${totalSeconds} сек`
}

function formatScheduleDate(value: string) {
  const [year, month, day] = value.split('-').map(Number)
  return new Intl.DateTimeFormat('ru-RU', { weekday: 'short', day: 'numeric', month: 'short' })
    .format(new Date(year, month - 1, day))
}

function getDashboardDataError(error: unknown, fallback: string) {
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте подключение.'
  if (error instanceof Error && error.message) return error.message
  return fallback
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
  onWelcomeMessageChange,
  activeWorkout,
  activeWorkoutError,
  isActiveWorkoutLoading,
  isStartingWorkout,
  onStartScheduledWorkout,
  onResumeWorkout,
  onOpenPrograms,
}: DashboardPageProps) {
  const today = getLocalDate()
  const weekDates = useMemo(() => getWeekDates(today), [today])
  const [scheduledWorkouts, setScheduledWorkouts] = useState<ScheduledWorkoutWithDayName[]>([])
  const [isScheduleLoading, setIsScheduleLoading] = useState(true)
  const [scheduleError, setScheduleError] = useState('')
  const [programCount, setProgramCount] = useState<number | null>(null)
  const [isProgramsLoading, setIsProgramsLoading] = useState(true)
  const [programsError, setProgramsError] = useState('')
  const [startError, setStartError] = useState('')
  const [isCalendarOpen, setIsCalendarOpen] = useState(false)
  const [calendarDate, setCalendarDate] = useState(today)
  const [statistics, setStatistics] = useState<WorkoutStatisticsDashboard | null>(null)
  const [isStatisticsLoading, setIsStatisticsLoading] = useState(true)
  const [statisticsError, setStatisticsError] = useState('')
  const [statisticsRetry, setStatisticsRetry] = useState(0)
  const [latestWorkout, setLatestWorkout] = useState<WorkoutHistoryItem | null>(null)
  const [recentHistory, setRecentHistory] = useState<WorkoutHistoryItem[]>([])
  const [isLatestWorkoutLoading, setIsLatestWorkoutLoading] = useState(true)
  const [latestWorkoutError, setLatestWorkoutError] = useState('')
  const [historyRetry, setHistoryRetry] = useState(0)

  const loadSchedule = useCallback(async (programs?: Program[]) => {
    setIsScheduleLoading(true)
    setScheduleError('')
    try {
      setScheduledWorkouts(await getScheduleWithDayNames('0001-01-01', '9999-12-31', programs))
    } catch (error) {
      setScheduleError(getErrorMessage(error))
    } finally {
      setIsScheduleLoading(false)
    }
  }, [])

  useEffect(() => {
    let cancelled = false
    getPrograms()
      .then((programs) => {
        if (cancelled) return
        setProgramCount(programs.length)
        void loadSchedule(programs)
      })
      .catch((error: unknown) => {
        if (!cancelled) {
          setProgramsError(getDashboardDataError(error, 'Не удалось определить список программ.'))
          setScheduleError(getDashboardDataError(error, 'Не удалось загрузить расписание.'))
          setIsScheduleLoading(false)
        }
      })
      .finally(() => {
        if (!cancelled) setIsProgramsLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [loadSchedule])

  useEffect(() => {
    let cancelled = false
    setIsStatisticsLoading(true)
    setStatisticsError('')
    getWorkoutStatisticsDashboard('7d')
      .then((result) => {
        if (!cancelled) setStatistics(result)
      })
      .catch((error: unknown) => {
        if (!cancelled) {
          setStatisticsError(getDashboardDataError(error, 'Не удалось загрузить статистику.'))
        }
      })
      .finally(() => {
        if (!cancelled) setIsStatisticsLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [statisticsRetry])

  useEffect(() => {
    let cancelled = false
    setIsLatestWorkoutLoading(true)
    setLatestWorkoutError('')
    getCurrentWeekHistory(weekDates[0])
      .then((items) => {
        if (!cancelled) {
          setRecentHistory(items)
          setLatestWorkout(items[0] ?? null)
        }
      })
      .catch((error: unknown) => {
        if (!cancelled) {
          setLatestWorkoutError(getDashboardDataError(error, 'Не удалось загрузить историю тренировок.'))
        }
      })
      .finally(() => {
        if (!cancelled) setIsLatestWorkoutLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [historyRetry, weekDates])

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
  const weeklySummary = useMemo(
    () => statistics ? getCurrentWeekSummary(statistics, weekDates, today) : null,
    [statistics, weekDates, today],
  )
  const completedWorkoutDates = useMemo(
    () => new Set(recentHistory.map(({ completed_at }) => getLocalDateKey(completed_at))),
    [recentHistory],
  )
  const hasCompletedWorkoutToday = recentHistory.some(
    ({ completed_at }) => getLocalDateKey(completed_at) === today,
  )
  const canStart = !activeWorkout && !isActiveWorkoutLoading && !activeWorkoutError && !isStartingWorkout

  useEffect(() => {
    if (
      isProgramsLoading ||
      programCount === null ||
      isScheduleLoading ||
      scheduleError ||
      isLatestWorkoutLoading ||
      latestWorkoutError ||
      isActiveWorkoutLoading ||
      activeWorkoutError
    ) return

    const welcomeMessage = hasCompletedWorkoutToday
      ? WELCOME_MESSAGES.completed
      : activeWorkout || todayWorkouts.length > 0
        ? WELCOME_MESSAGES.workout
        : programCount === 0
          ? WELCOME_MESSAGES.welcome
          : WELCOME_MESSAGES.rest

    onWelcomeMessageChange(welcomeMessage)
  }, [
    activeWorkout,
    activeWorkoutError,
    hasCompletedWorkoutToday,
    isActiveWorkoutLoading,
    isLatestWorkoutLoading,
    isProgramsLoading,
    isScheduleLoading,
    latestWorkoutError,
    onWelcomeMessageChange,
    programCount,
    scheduleError,
    todayWorkouts,
  ])

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
    <section className="dashboard-content dashboard-home-content">
      <div className="section-heading dashboard-home-heading">
        <div>
          <p className="eyebrow">ТВОИ ТРЕНИРОВКИ</p>
          <h1>Главная</h1>
          <p>Твой план, активность и результаты тренировок в одном месте.</p>
        </div>
      </div>
      {programsError && <p className="program-alert mb-0" role="alert">{programsError}</p>}
      {activeWorkoutError && <p className="program-alert mb-4" role="alert">{activeWorkoutError}</p>}

      <WeeklyStats
        summary={weeklySummary}
        isLoading={isStatisticsLoading}
        error={statisticsError}
        onRetry={() => setStatisticsRetry((current) => current + 1)}
      />

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
                    ? 'Время тренировки'
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
                    {upcomingWorkouts.length > 0
                      ? 'На сегодня тренировка не запланирована. Выберите будущую тренировку, чтобы начать её сегодня.'
                      : 'В календаре нет предстоящих тренировок. Сначала запланируйте занятие, чтобы оно появилось здесь.'}
                  </p>
                  {upcomingWorkouts.length > 0 ? (
                    <button
                      className="program-secondary-button mt-5"
                      type="button"
                      onClick={() => openCalendar(upcomingWorkouts[0].scheduled_date)}
                    >
                      <CalendarDays size={15} /> Выбрать тренировку
                    </button>
                  ) : (
                    <button className="program-primary-button mt-5" type="button" onClick={onOpenPrograms}>
                      Начать планирование <ChevronRight size={15} />
                    </button>
                  )}
                </>
              )}
              {startError && <p className="program-alert mt-4" role="alert">{startError}</p>}
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

      <div className="dashboard-bottom-grid">
        <section className="dashboard-week-activity panel-card" aria-labelledby="dashboard-week-activity-title">
            <div className="dashboard-widget-heading">
              <div>
                <p className="eyebrow">ВАША АКТИВНОСТЬ</p>
                <h3 id="dashboard-week-activity-title">Активность за неделю</h3>
              </div>
              {isStatisticsLoading && (
                <LoaderCircle aria-label="Загружаем активность" className="spin text-muted" size={18} />
              )}
            </div>
            {statisticsError ? (
              <p className="dashboard-data-loading" role="status">Активность временно недоступна.</p>
            ) : isStatisticsLoading ? (
              <p className="dashboard-data-loading">Загружаем активность...</p>
            ) : (
              <div aria-label="Дни текущей недели" className="dashboard-week-days" role="list">
                {weekDates.map((date, index) => {
                  const isToday = date === today
                  const isFuture = date > today
                  const isComplete = completedWorkoutDates.has(date)
                  const stateLabel = isComplete
                    ? 'тренировка завершена'
                    : isToday
                      ? 'сегодня, тренировка не завершена'
                      : isFuture
                        ? 'будущий день'
                        : 'нет завершённой тренировки'
                  return (
                    <div
                      aria-label={`${['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'][index]}: ${stateLabel}`}
                      className="dashboard-week-day"
                      key={date}
                      role="listitem"
                    >
                      <span
                        className={[
                          'dashboard-week-day-marker',
                          isComplete && 'is-complete',
                          isToday && !isComplete && 'is-today',
                          isFuture && 'is-future',
                        ].filter(Boolean).join(' ')}
                      >
                        {isComplete && <Check aria-hidden="true" size={15} strokeWidth={2.5} />}
                      </span>
                      <span className={isToday ? 'is-today' : ''}>
                        {['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'][index]}
                      </span>
                    </div>
                  )
                })}
              </div>
            )}
            {statisticsError && (
              <button
                className="text-button dashboard-week-retry"
                type="button"
                onClick={() => setStatisticsRetry((current) => current + 1)}
              >
                <RotateCcw size={13} /> Повторить загрузку
              </button>
            )}
        </section>

        <section className="dashboard-latest-workout panel-card" aria-labelledby="dashboard-latest-title">
            <div className="dashboard-widget-heading">
              <div>
                <p className="eyebrow">ВАША ИСТОРИЯ</p>
                <h3 id="dashboard-latest-title">Последняя тренировка</h3>
              </div>
              {isLatestWorkoutLoading && (
                <LoaderCircle aria-label="Загружаем историю" className="spin text-muted" size={18} />
              )}
            </div>
            {latestWorkoutError ? (
              <div className="dashboard-latest-state" role="alert">
                <p>{latestWorkoutError}</p>
                <button
                  className="text-button"
                  type="button"
                  onClick={() => setHistoryRetry((current) => current + 1)}
                >
                  <RotateCcw size={13} /> Повторить
                </button>
              </div>
            ) : isLatestWorkoutLoading ? (
              <p className="dashboard-latest-state">Загружаем последнюю тренировку...</p>
            ) : latestWorkout ? (
              <article className="dashboard-latest-content">
                <div className="dashboard-latest-icon"><Dumbbell size={19} /></div>
                <div className="dashboard-latest-details">
                  <h4>{latestWorkout.training_day_name || 'Без названия'}</h4>
                  <p>{formatWorkoutDate(latestWorkout.completed_at)}</p>
                  <div className="dashboard-latest-metrics">
                    <span><Clock3 size={14} /> {formatWorkoutDuration(latestWorkout.duration_seconds)}</span>
                    <span><Dumbbell size={14} /> {latestWorkout.completed_sets} подходов</span>
                    <span className="dashboard-latest-completion">
                      {latestWorkout.completion_percent === null
                        ? 'План: —'
                        : `План: ${latestWorkout.completion_percent}%`}
                    </span>
                  </div>
                </div>
            </article>
          ) : (
            <div className="dashboard-latest-empty">
              <span className="dashboard-latest-icon"><Dumbbell size={19} /></span>
              <p>Завершённых тренировок пока нет.</p>
            </div>
          )}
        </section>
      </div>

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
