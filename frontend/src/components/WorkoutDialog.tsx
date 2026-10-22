import { Check, Clock3, Dumbbell, X } from 'lucide-react'
import { nextWorkout } from '../data'

type WorkoutDialogProps = {
  open: boolean
  onClose: () => void
}

export default function WorkoutDialog({ open, onClose }: WorkoutDialogProps) {
  if (!open) return null

  return (
    <div
      className="dialog-backdrop"
      onClick={(event) => {
        if (event.target === event.currentTarget) onClose()
      }}
      role="presentation"
    >
      <section
        aria-labelledby="workout-dialog-title"
        aria-modal="true"
        className="workout-dialog"
        role="dialog"
      >
        <div className="flex items-start justify-between">
          <div>
            <p className="eyebrow">ПЛАН НА СЕГОДНЯ</p>
            <h2 className="mt-2 text-xl font-semibold tracking-[-0.04em] text-ink" id="workout-dialog-title">
              {nextWorkout.name}
            </h2>
            <div className="mt-2 flex gap-4 text-xs text-muted">
              <span className="inline-flex items-center gap-1.5">
                <Clock3 size={14} /> {nextWorkout.duration}
              </span>
              <span className="inline-flex items-center gap-1.5">
                <Dumbbell size={14} /> {nextWorkout.exercises.length} упражнения
              </span>
            </div>
          </div>
          <button className="icon-button" type="button" onClick={onClose} aria-label="Закрыть">
            <X size={17} />
          </button>
        </div>

        <ol className="mt-6 space-y-2">
          {nextWorkout.exercises.map((exercise, index) => (
            <li className="exercise-row" key={exercise}>
              <span className="exercise-number">{String(index + 1).padStart(2, '0')}</span>
              <span className="flex-1 text-sm font-medium text-ink">{exercise}</span>
              <Check size={15} className="text-subtle" />
            </li>
          ))}
        </ol>

        <div className="dialog-note mt-5 rounded-xl px-4 py-3 text-xs leading-5 text-muted">
          Это демонстрационный план. Тренировка не сохраняется — подключение к серверу пока не настроено.
        </div>
        <button className="primary-button mt-5 w-full justify-center" type="button" onClick={onClose}>
          Понятно
        </button>
      </section>
    </div>
  )
}
