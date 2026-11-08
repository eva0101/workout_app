import { CalendarDays, ChevronLeft, ChevronRight, Dumbbell, LoaderCircle, Play, X } from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'
import type { ScheduledWorkoutWithDayName } from '../lib/api'

type WorkoutCalendarDialogProps = {
  open: boolean
  workouts: ScheduledWorkoutWithDayName[]
  loadError: string
  isLoading: boolean
  isStarting: boolean
  hasActiveWorkout: boolean
  initialDate: string
  onClose: () => void
  onRetry: () => void
  onOpenPrograms: () => void
  onStartWorkout: (scheduleId: number, title: string) => Promise<void>
}

const weekdays = ['ПН', 'ВТ', 'СР', 'ЧТ', 'ПТ', 'СБ', 'ВС']
const monthFormatter = new Intl.DateTimeFormat('ru-RU', { month: 'long', year: 'numeric' })
const dateFormatter = new Intl.DateTimeFormat('ru-RU', {
  day: 'numeric',
  month: 'long',
  year: 'numeric',
})

function toDateKey(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function parseDateKey(value: string): Date {
  const [year, month, day] = value.split('-').map(Number)
  return new Date(year, month - 1, day)
}

function formatDateKey(value: string): string {
  return dateFormatter.format(parseDateKey(value))
}

function getErrorMessage(error: unknown): string {
  if (typeof error === 'object' && error !== null && 'status' in error) {
    const status = (error as { status?: number }).status
    if (status === 404) return 'Запись расписания уже недоступна. Обновите календарь.'
    if (status === 409) return 'У вас уже есть активная тренировка. Продолжите её вместо запуска новой.'
  }
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте подключение.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось начать тренировку. Попробуйте ещё раз.'
}

export default function WorkoutCalendarDialog({
  open,
  workouts,
  loadError,
  isLoading,
  isStarting,
  hasActiveWorkout,
  initialDate,
  onClose,
  onRetry,
  onOpenPrograms,
  onStartWorkout,
}: WorkoutCalendarDialogProps) {
  const [visibleMonth, setVisibleMonth] = useState(() => {
    const initial = parseDateKey(initialDate)
    return new Date(initial.getFullYear(), initial.getMonth(), 1)
  })
  const [selectedDate, setSelectedDate] = useState(initialDate)
  const [startError, setStartError] = useState('')

  useEffect(() => {
    if (!open) return
    const initial = parseDateKey(initialDate)
    setVisibleMonth(new Date(initial.getFullYear(), initial.getMonth(), 1))
    setSelectedDate(initialDate)
    setStartError('')
  }, [open, initialDate])

  const days = useMemo(() => {
    const year = visibleMonth.getFullYear()
    const month = visibleMonth.getMonth()
    const offset = (new Date(year, month, 1).getDay() + 6) % 7
    const total = new Date(year, month + 1, 0).getDate()
    return [
      ...Array.from({ length: offset }, () => null),
      ...Array.from({ length: total }, (_, index) => index + 1),
    ]
  }, [visibleMonth])

  if (!open) return null

  const workoutsByDate = new Map<string, ScheduledWorkoutWithDayName[]>()
  for (const workout of workouts) {
    const items = workoutsByDate.get(workout.scheduled_date) ?? []
    items.push(workout)
    workoutsByDate.set(workout.scheduled_date, items)
  }
  const selectedWorkouts = workoutsByDate.get(selectedDate) ?? []
  const hasAnyWorkouts = workouts.length > 0

  async function startWorkout(workout: ScheduledWorkoutWithDayName) {
    if (hasActiveWorkout || isStarting) return
    setStartError('')
    const title = `${workout.program_name} — ${workout.training_day_name}`
    try {
      await onStartWorkout(workout.id, title)
      onClose()
    } catch (error) {
      setStartError(getErrorMessage(error))
    }
  }

  return (
    <div
      className="schedule-modal-backdrop"
      onClick={(event) => {
        if (event.target === event.currentTarget && !isStarting) onClose()
      }}
      role="presentation"
    >
      <section
        aria-labelledby="workout-calendar-title"
        aria-modal="true"
        className="schedule-dialog workout-calendar-dialog"
        role="dialog"
      >
        <header className="schedule-dialog-header">
          <div>
            <p className="eyebrow">ПЛАН ТРЕНИРОВОК</p>
            <h2 id="workout-calendar-title">Календарь</h2>
          </div>
          <button
            className="program-icon-action"
            type="button"
            aria-label="Закрыть календарь"
            disabled={isStarting}
            onClick={onClose}
          >
            <X size={19} />
          </button>
        </header>

        <div className="schedule-calendar-month">
          <button
            className="schedule-month-button"
            type="button"
            aria-label="Предыдущий месяц"
            onClick={() => setVisibleMonth((month) => new Date(month.getFullYear(), month.getMonth() - 1, 1))}
          >
            <ChevronLeft size={18} />
          </button>
          <h3>{monthFormatter.format(visibleMonth)}</h3>
          <button
            className="schedule-month-button"
            type="button"
            aria-label="Следующий месяц"
            onClick={() => setVisibleMonth((month) => new Date(month.getFullYear(), month.getMonth() + 1, 1))}
          >
            <ChevronRight size={18} />
          </button>
        </div>

        <div className="schedule-calendar-grid" role="grid" aria-label="Запланированные тренировки">
          {weekdays.map((weekday) => (
            <span className="schedule-weekday" key={weekday} role="columnheader">{weekday}</span>
          ))}
          {days.map((day, index) => {
            if (day === null) {
              return <span className="schedule-calendar-empty" key={`empty-${index}`} aria-hidden="true" />
            }
            const date = new Date(visibleMonth.getFullYear(), visibleMonth.getMonth(), day)
            const key = toDateKey(date)
            const isBooked = (workoutsByDate.get(key)?.length ?? 0) > 0
            return (
              <button
                aria-label={`${formatDateKey(key)}${isBooked ? ', есть тренировка' : ''}`}
                aria-pressed={selectedDate === key}
                className={[
                  'schedule-date-button',
                  selectedDate === key ? 'schedule-date-selected' : '',
                  isBooked ? 'schedule-date-booked' : '',
                ].filter(Boolean).join(' ')}
                key={key}
                onClick={() => setSelectedDate(key)}
                role="gridcell"
                type="button"
              >
                {day}
                {isBooked && <span className="schedule-booked-dot" aria-hidden="true" />}
              </button>
            )
          })}
        </div>

        <section className="workout-calendar-results" aria-live="polite">
          <h3>{formatDateKey(selectedDate)}</h3>
          {isLoading ? (
            <div className="program-state"><LoaderCircle className="spin" size={19} /> Загружаем календарь...</div>
          ) : loadError ? (
            <div className="program-alert" role="alert">
              <span>{loadError}</span>
              <button className="program-secondary-button" type="button" onClick={onRetry}>Повторить</button>
            </div>
          ) : selectedWorkouts.length > 0 ? (
            <ul className="workout-calendar-list">
              {selectedWorkouts.map((workout) => {
                const title = `${workout.program_name} — ${workout.training_day_name}`
                return (
                  <li className="workout-calendar-item" key={workout.id}>
                    <span className="workout-icon workout-icon-sage"><Dumbbell size={17} /></span>
                    <span className="min-w-0 flex-1">
                      <span className="workout-calendar-title">{title}</span>
                      <span className="workout-calendar-subtitle">Запланированная тренировка</span>
                    </span>
                    <button
                      className="program-primary-button"
                      type="button"
                      disabled={hasActiveWorkout || isStarting}
                      onClick={() => void startWorkout(workout)}
                    >
                      {isStarting
                        ? <LoaderCircle className="spin" size={15} />
                        : <Play size={15} fill="currentColor" />}
                      Начать
                    </button>
                  </li>
                )
              })}
            </ul>
          ) : !hasAnyWorkouts && !loadError ? (
            <div className="workout-calendar-empty">
              <CalendarDays size={20} />
              <p>В календаре пока нет запланированных тренировок.</p>
              <button className="program-secondary-button" type="button" onClick={onOpenPrograms}>
                Перейти к программам
              </button>
            </div>
          ) : (
            <p className="workout-calendar-empty-copy">На эту дату тренировок нет. Выберите день с отметкой в календаре.</p>
          )}
          {hasActiveWorkout && (
            <p className="workout-calendar-note">Сначала продолжите и завершите активную тренировку.</p>
          )}
          {startError && <p className="program-alert mt-3" role="alert">{startError}</p>}
        </section>
      </section>
    </div>
  )
}
