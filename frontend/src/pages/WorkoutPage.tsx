import { ArrowLeft, Check, Dumbbell, LoaderCircle, Plus } from 'lucide-react'
import { type FormEvent, useState } from 'react'
import type { Workout } from '../lib/api'

type WorkoutPageProps = {
  workout: Workout
  title: string
  isSavingSet: boolean
  isCompleting: boolean
  error: string
  onAddSet: (exerciseId: number, repsDone: number, weight: number | null) => Promise<boolean>
  onComplete: (fatigueScore: number) => void
  onBack: () => void
}

type SetInput = {
  reps: string
  weight: string
}

function formatValue(value: number | null) {
  return value === null ? '—' : String(value)
}

export default function WorkoutPage({
  workout,
  title,
  isSavingSet,
  isCompleting,
  error,
  onAddSet,
  onComplete,
  onBack,
}: WorkoutPageProps) {
  const [inputs, setInputs] = useState<Record<number, SetInput>>({})
  const [fatigueScore, setFatigueScore] = useState(5)

  async function submitSet(event: FormEvent<HTMLFormElement>, exerciseId: number) {
    event.preventDefault()
    const input = inputs[exerciseId] ?? { reps: '', weight: '' }
    const repsDone = Number(input.reps)
    const weight = input.weight.trim() === '' ? null : Number(input.weight)
    if (!Number.isInteger(repsDone) || repsDone < 0 || (weight !== null && (!Number.isFinite(weight) || weight < 0))) {
      return
    }
    if (await onAddSet(exerciseId, repsDone, weight)) {
      setInputs((current) => ({ ...current, [exerciseId]: { reps: '', weight: '' } }))
    }
  }

  return (
    <section className="page-content programs-page">
      <button className="program-back-button" type="button" onClick={onBack}>
        <ArrowLeft size={16} /> К обзору
      </button>

      <div className="program-detail-heading">
        <div className="program-detail-icon"><Dumbbell size={23} /></div>
        <div>
          <p className="eyebrow">ВЫПОЛНЕНИЕ ТРЕНИРОВКИ · #{workout.id}</p>
          <h2>{title || `Тренировка #${workout.id}`}</h2>
        </div>
      </div>
      {error && <div className="program-alert mt-4" role="alert">{error}</div>}

      {workout.exercises.length === 0 ? (
        <div className="panel-card mt-5 p-5 text-sm text-muted">
          В снимке этой тренировки нет упражнений.
        </div>
      ) : (
        <div className="mt-5 flex flex-col gap-4">
          {workout.exercises.map((exercise, index) => {
            const input = inputs[exercise.id] ?? { reps: '', weight: '' }
            return (
              <article className="panel-card p-5 sm:p-6" key={exercise.id}>
                <div className="flex items-start gap-3">
                  <span className="exercise-number">{String(index + 1).padStart(2, '0')}</span>
                  <div className="min-w-0">
                    <h3 className="text-base font-semibold text-ink">{exercise.title || exercise.name}</h3>
                    <p className="mt-1 text-xs text-muted">Упражнение из сохранённого плана</p>
                  </div>
                </div>

                <div className="mt-5">
                  <h4 className="text-xs font-semibold uppercase tracking-wide text-muted">Плановые подходы</h4>
                  {exercise.planned_sets.length === 0 ? (
                    <p className="mt-2 text-sm text-muted">Плановые подходы не заданы.</p>
                  ) : (
                    <ul className="mt-2 flex flex-wrap gap-2">
                      {exercise.planned_sets.map((set) => (
                        <li className="rounded-lg bg-[var(--surface-soft)] px-3 py-2 text-xs text-ink" key={set.id}>
                          {set.set_number}. {formatValue(set.repetitions)} повт. · {formatValue(set.weight)} кг
                        </li>
                      ))}
                    </ul>
                  )}
                </div>

                <div className="mt-5 border-t border-line pt-4">
                  <h4 className="text-xs font-semibold uppercase tracking-wide text-muted">Выполненные подходы</h4>
                  {exercise.actual_sets.length === 0 ? (
                    <p className="mt-2 text-sm text-muted">Пока нет записанных подходов.</p>
                  ) : (
                    <ul className="mt-2 flex flex-wrap gap-2">
                      {exercise.actual_sets.map((set) => (
                        <li className="inline-flex items-center gap-1.5 rounded-lg bg-[var(--surface-soft)] px-3 py-2 text-xs text-ink" key={set.id}>
                          <Check size={13} className="text-accent" />
                          {set.set_number}. {formatValue(set.reps_done)} повт. · {formatValue(set.weight)} кг
                        </li>
                      ))}
                    </ul>
                  )}
                </div>

                <form className="mt-5 flex flex-wrap items-end gap-3" onSubmit={(event) => submitSet(event, exercise.id)}>
                  <label className="flex flex-col gap-1 text-xs font-medium text-muted">
                    Повторения
                    <input
                      className="w-28 rounded-lg border border-line bg-[var(--surface)] px-3 py-2 text-sm text-ink"
                      type="number"
                      min="0"
                      step="1"
                      required
                      value={input.reps}
                      onChange={(event) => setInputs((current) => ({
                        ...current,
                        [exercise.id]: { ...input, reps: event.target.value },
                      }))}
                    />
                  </label>
                  <label className="flex flex-col gap-1 text-xs font-medium text-muted">
                    Вес, кг
                    <input
                      className="w-28 rounded-lg border border-line bg-[var(--surface)] px-3 py-2 text-sm text-ink"
                      type="number"
                      min="0"
                      step="0.01"
                      value={input.weight}
                      onChange={(event) => setInputs((current) => ({
                        ...current,
                        [exercise.id]: { ...input, weight: event.target.value },
                      }))}
                    />
                  </label>
                  <button
                    className="program-secondary-button"
                    type="submit"
                    disabled={isSavingSet || isCompleting}
                  >
                    {isSavingSet ? <LoaderCircle className="spin" size={15} /> : <Plus size={15} />}
                    Добавить подход
                  </button>
                </form>
              </article>
            )
          })}
        </div>
      )}

      <section className="panel-card mt-5 flex flex-wrap items-end justify-between gap-4 p-5 sm:p-6">
        <label className="flex flex-col gap-1 text-xs font-medium text-muted">
          Самочувствие после тренировки
          <select
            className="min-h-10 rounded-lg border border-line bg-[var(--surface)] px-3 text-sm text-ink"
            value={fatigueScore}
            onChange={(event) => setFatigueScore(Number(event.target.value))}
          >
            {Array.from({ length: 11 }, (_, score) => (
              <option key={score} value={score}>{score} / 10</option>
            ))}
          </select>
        </label>
        <button
          className="program-primary-button"
          type="button"
          disabled={isCompleting || isSavingSet}
          onClick={() => onComplete(fatigueScore)}
        >
          {isCompleting ? <LoaderCircle className="spin" size={16} /> : <Check size={16} />}
          Завершить тренировку
        </button>
      </section>
    </section>
  )
}
