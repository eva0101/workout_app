import { ChevronLeft, ChevronRight, X } from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'

type ScheduleDialogProps = {
  open: boolean
  bookedDates: string[]
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
  bookedDates,
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

  useEffect(() => {
    if (open) setSelectedDates([])
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
  const bookedDateSet = new Set(bookedDates)
  const selectedDateSet = new Set(selectedDates)

  function toggleDate(date: Date) {
    const key = toDateKey(date)
    if (key < todayKey || bookedDateSet.has(key)) return
    setSelectedDates((current) =>
      current.includes(key)
        ? current.filter((value) => value !== key)
        : [...current, key].sort(),
    )
  }

  function changeMonth(delta: number) {
    setVisibleMonth((current) => new Date(current.getFullYear(), current.getMonth() + delta, 1))
  }

  function handleClose() {
    if (isSaving) return
    setSelectedDates([])
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
            const isBooked = bookedDateSet.has(key)
            const isSelected = selectedDateSet.has(key)
            return (
              <button
                aria-label={`${formatDateKey(key)}${isBooked ? ', уже запланировано' : ''}`}
                aria-pressed={isSelected}
                className={[
                  'schedule-date-button',
                  isSelected ? 'schedule-date-selected' : '',
                  isPast || isBooked ? 'schedule-date-disabled' : '',
                ].filter(Boolean).join(' ')}
                disabled={isPast || isBooked || isSaving}
                key={key}
                onClick={() => toggleDate(date)}
                role="gridcell"
                title={isBooked ? 'Тренировка на этот день уже запланирована' : undefined}
                type="button"
              >
                {day}
                {isBooked && <span className="schedule-booked-dot" aria-hidden="true" />}
              </button>
            )
          })}
        </div>

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
