import { ChevronLeft, ChevronRight, X } from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'
import type { ScheduledWorkout } from '../lib/api'

type ScheduleDialogProps = {
  open: boolean
  scheduledWorkouts: ScheduledWorkout[]
  trainingDayId: number
  isSaving: boolean
  error: string
  onClose: () => void
  onSave: (dates: string[]) => void
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

function startOfToday(): Date {
  const now = new Date()
  return new Date(now.getFullYear(), now.getMonth(), now.getDate())
}

function formatDateKey(value: string): string {
  const [year, month, day] = value.split('-').map(Number)
  return dateFormatter.format(new Date(year, month - 1, day))
}

export default function ScheduleDialog({
  open,
  scheduledWorkouts,
  trainingDayId,
  isSaving,
  error,
  onClose,
  onSave,
}: ScheduleDialogProps) {
  const [visibleMonth, setVisibleMonth] = useState(() => {
    const today = startOfToday()
    return new Date(today.getFullYear(), today.getMonth(), 1)
  })
  const [selectedDates, setSelectedDates] = useState<string[]>([])
  const [pendingConfirmationDate, setPendingConfirmationDate] = useState<string | null>(null)
  const [selectionWarning, setSelectionWarning] = useState('')

  useEffect(() => {
    if (open) {
      setSelectedDates([])
      setPendingConfirmationDate(null)
      setSelectionWarning('')
    }
  }, [open])

  const days = useMemo(() => {
    const year = visibleMonth.getFullYear()
    const month = visibleMonth.getMonth()
    const firstDayOffset = (new Date(year, month, 1).getDay() + 6) % 7
    const numberOfDays = new Date(year, month + 1, 0).getDate()
    return [
      ...Array.from({ length: firstDayOffset }, () => null),
      ...Array.from({ length: numberOfDays }, (_, index) => index + 1),
    ]
  }, [visibleMonth])

  if (!open) return null

  const todayKey = toDateKey(startOfToday())
  const workoutsByDate = new Map<string, ScheduledWorkout[]>()
  for (const workout of scheduledWorkouts) {
    const items = workoutsByDate.get(workout.scheduled_date) ?? []
    items.push(workout)
    workoutsByDate.set(workout.scheduled_date, items)
  }
  const selectedDateSet = new Set(selectedDates)

  function toggleDate(date: Date) {
    const key = toDateKey(date)
    if (key < todayKey) return
    if (selectedDateSet.has(key)) {
      setSelectedDates((current) => current.filter((value) => value !== key))
      setSelectionWarning('')
      setPendingConfirmationDate(null)
      return
    }

    const existingWorkouts = workoutsByDate.get(key) ?? []
    if (existingWorkouts.some((workout) => workout.training_day_id === trainingDayId)) {
      setSelectionWarning(`На ${formatDateKey(key)} уже запланирован этот тренировочный день. Повторно добавить его на ту же дату нельзя.`)
      setPendingConfirmationDate(null)
      return
    }
    if (existingWorkouts.length > 0) {
      setPendingConfirmationDate(key)
      setSelectionWarning('')
      return
    }

    setSelectedDates((current) =>
      [...current, key].sort(),
    )
    setPendingConfirmationDate(null)
    setSelectionWarning('')
  }

  function confirmSecondWorkout() {
    if (!pendingConfirmationDate) return
    setSelectedDates((current) => [...current, pendingConfirmationDate].sort())
    setPendingConfirmationDate(null)
    setSelectionWarning('')
  }

  function changeMonth(delta: number) {
    setVisibleMonth((current) => new Date(current.getFullYear(), current.getMonth() + delta, 1))
  }

  function handleClose() {
    if (isSaving) return
    setSelectedDates([])
    setPendingConfirmationDate(null)
    setSelectionWarning('')
    onClose()
  }

  return (
    <div
      className="schedule-modal-backdrop"
      onClick={(event) => {
        if (event.target === event.currentTarget) handleClose()
      }}
      role="presentation"
    >
      <section
        aria-labelledby="schedule-dialog-title"
        aria-modal="true"
        className="schedule-dialog"
        role="dialog"
      >
        <header className="schedule-dialog-header">
          <div>
            <p className="eyebrow">КАЛЕНДАРЬ ТРЕНИРОВОК</p>
            <h2 id="schedule-dialog-title">Запланировать</h2>
          </div>
          <button
            className="program-icon-action"
            type="button"
            aria-label="Закрыть календарь"
            disabled={isSaving}
            onClick={handleClose}
          >
            <X size={19} />
          </button>
        </header>

        <div className="schedule-calendar-month">
          <button
            className="schedule-month-button"
            type="button"
            aria-label="Предыдущий месяц"
            onClick={() => changeMonth(-1)}
          >
            <ChevronLeft size={18} />
          </button>
          <h3>{monthFormatter.format(visibleMonth)}</h3>
          <button
            className="schedule-month-button"
            type="button"
            aria-label="Следующий месяц"
            onClick={() => changeMonth(1)}
          >
            <ChevronRight size={18} />
          </button>
        </div>

        <div className="schedule-calendar-grid" role="grid" aria-label="Выбор дат">
          {weekdays.map((weekday) => (
            <span className="schedule-weekday" key={weekday} role="columnheader">{weekday}</span>
          ))}
          {days.map((day, index) => {
            if (day === null) {
              return <span className="schedule-calendar-empty" key={`empty-${index}`} aria-hidden="true" />
            }

            const date = new Date(visibleMonth.getFullYear(), visibleMonth.getMonth(), day)
            const key = toDateKey(date)
            const isPast = key < todayKey
            const existingWorkouts = workoutsByDate.get(key) ?? []
            const isBooked = existingWorkouts.length > 0
            const isSelected = selectedDateSet.has(key)
            return (
              <button
                aria-label={`${formatDateKey(key)}${isBooked ? ', уже запланировано' : ''}`}
                aria-pressed={isSelected}
                className={[
                  'schedule-date-button',
                  isSelected ? 'schedule-date-selected' : '',
                  isBooked ? 'schedule-date-booked' : '',
                  isPast ? 'schedule-date-disabled' : '',
                ].filter(Boolean).join(' ')}
                disabled={isPast || isSaving}
                key={key}
                onClick={() => toggleDate(date)}
                role="gridcell"
                title={isBooked ? 'На эту дату уже запланирована тренировка' : undefined}
                type="button"
              >
                {day}
                {isBooked && <span className="schedule-booked-dot" aria-hidden="true" />}
              </button>
            )
          })}
        </div>

        {pendingConfirmationDate && (
          <div className="schedule-overlap-confirmation" role="alertdialog" aria-live="polite">
            <p>
              На {formatDateKey(pendingConfirmationDate)} уже запланирована тренировка.
              Вы пытаетесь добавить вторую тренировку на этот день.
            </p>
            <div>
              <button
                className="program-secondary-button"
                type="button"
                disabled={isSaving}
                onClick={() => setPendingConfirmationDate(null)}
              >
                Отменить
              </button>
              <button
                className="program-primary-button"
                type="button"
                disabled={isSaving}
                onClick={confirmSecondWorkout}
              >
                Всё равно добавить
              </button>
            </div>
          </div>
        )}

        {selectionWarning && <p className="schedule-error" role="alert">{selectionWarning}</p>}
        <div className="schedule-selection-summary" aria-live="polite">
          <span>
            {selectedDates.length === 0
              ? 'Выберите одну или несколько дат'
              : `Выбрано дат: ${selectedDates.length}`}
          </span>
          {selectedDates.length > 0 && (
            <span className="schedule-selection-preview">
              {selectedDates.slice(0, 2).map(formatDateKey).join(', ')}
              {selectedDates.length > 2 ? ` и ещё ${selectedDates.length - 2}` : ''}
            </span>
          )}
        </div>

        {error && <p className="schedule-error" role="alert">{error}</p>}

        <footer className="schedule-dialog-actions">
          <button className="program-secondary-button" type="button" disabled={isSaving} onClick={handleClose}>
            Отмена
          </button>
          <button
            className="program-primary-button"
            type="button"
            disabled={selectedDates.length === 0 || isSaving}
            onClick={() => onSave(selectedDates)}
          >
            {isSaving ? 'Сохраняем...' : `Сохранить${selectedDates.length ? ` · ${selectedDates.length}` : ''}`}
          </button>
        </footer>
      </section>
    </div>
  )
}
