import { Dumbbell, LoaderCircle, Plus, Trash2, X } from 'lucide-react'
import { useState } from 'react'
import {
  deleteTrainingDayExercise,
  type CatalogExercise,
  type TrainingDayDetails,
  type PlannedSet,
  updateTrainingDayExerciseSets,
} from '../lib/api'

type TrainingExerciseDialogProps = {
  exercise: TrainingDayDetails['exercises'][number]
  catalogExercise: CatalogExercise | undefined
  onClose: () => void
  onDeleted: (exerciseId: number) => void
  onSaved: (exerciseId: number, plannedSets: PlannedSet[]) => void
}

type SetDraft = {
  id: number
  repetitions: string
  weight: string
}

function imageSource(path: string | null): string | null {
  if (!path) return null
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  return /\.[a-z\d]+$/i.test(normalizedPath) ? normalizedPath : `${normalizedPath}.webp`
}

function getErrorMessage(error: unknown): string {
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте подключение.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось выполнить запрос. Попробуйте ещё раз.'
}

export default function TrainingExerciseDialog({
  exercise,
  catalogExercise,
  onClose,
  onDeleted,
  onSaved,
}: TrainingExerciseDialogProps) {
  const [sets, setSets] = useState<SetDraft[]>(() => exercise.planned_sets.map((set) => ({
    id: set.set_number,
    repetitions: set.repetitions === null ? '' : String(set.repetitions),
    weight: set.weight === null ? '' : String(set.weight),
  })))
  const [error, setError] = useState('')
  const [isDeleting, setIsDeleting] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const imagePath = imageSource(catalogExercise?.image_path ?? null)
  const exerciseName = catalogExercise?.title ?? `Упражнение #${exercise.exercise_id}`

  const updatedSets: PlannedSet[] = sets.map((set, index) => ({
    set_number: index + 1,
    repetitions: set.repetitions === '' ? null : Number(set.repetitions),
    weight: set.weight === '' ? null : Number(set.weight),
  }))
  const hasChanges = updatedSets.length !== exercise.planned_sets.length ||
    updatedSets.some((set, index) => (
      set.repetitions !== exercise.planned_sets[index]?.repetitions ||
      set.weight !== exercise.planned_sets[index]?.weight
    ))
  const isBusy = isDeleting || isSaving

  function addSet() {
    setSets((current) => [
      ...current,
      {
        id: Math.max(0, ...current.map(({ id }) => id)) + 1,
        repetitions: '',
        weight: '',
      },
    ])
  }

  function updateSet(id: number, field: 'repetitions' | 'weight', value: string) {
    setSets((current) => current.map((set) => set.id === id ? { ...set, [field]: value } : set))
  }

  async function handleSave() {
    if (!hasChanges || isBusy) return

    setIsSaving(true)
    setError('')
    try {
      await updateTrainingDayExerciseSets(exercise.id, updatedSets)
      onSaved(exercise.id, updatedSets)
      onClose()
    } catch (saveError) {
      setError(getErrorMessage(saveError))
    } finally {
      setIsSaving(false)
    }
  }

  async function handleDelete() {
    if (!window.confirm(
      `Удалить «${exerciseName}» из плана тренировочного дня? История выполненных тренировок не удаляется. Если упражнение уже использовалось в тренировке, сервер может запретить удаление.`,
    )) return

    setIsDeleting(true)
    setError('')
    try {
      await deleteTrainingDayExercise(exercise.id)
      onDeleted(exercise.id)
    } catch (deleteError) {
      setError(getErrorMessage(deleteError))
      setIsDeleting(false)
    }
  }

  return (
    <div
      className="program-modal-backdrop exercise-catalog-backdrop"
      onClick={(event) => {
        if (event.target === event.currentTarget && !isBusy) onClose()
      }}
      role="presentation"
    >
      <section
        aria-labelledby="training-exercise-title"
        aria-modal="true"
        className="program-create-dialog exercise-catalog-dialog exercise-detail-dialog"
        role="dialog"
      >
        <div className="program-create-heading exercise-catalog-heading">
          <div className="exercise-catalog-title-wrap">
            <div>
              <p className="eyebrow">УПРАЖНЕНИЕ В ПЛАНЕ</p>
              <h2 id="training-exercise-title">{exerciseName}</h2>
            </div>
          </div>
          <button
            className="program-icon-action"
            type="button"
            aria-label="Закрыть"
            onClick={onClose}
            disabled={isBusy}
          >
            <X size={18} />
          </button>
        </div>

        <div className="exercise-detail-layout">
          <div className="exercise-detail-art exercise-catalog-image">
            {imagePath ? (
              <img src={imagePath} alt={exerciseName} onError={(event) => event.currentTarget.remove()} />
            ) : (
              <Dumbbell size={25} aria-hidden="true" />
            )}
          </div>
          <div className="exercise-detail-copy">
            <p className="eyebrow">ТЕХНИКА ВЫПОЛНЕНИЯ</p>
            <p className="exercise-detail-description">
              {catalogExercise
                ? catalogExercise.description || 'Описание техники пока не добавлено.'
                : 'Не удалось загрузить данные упражнения из каталога.'}
            </p>
          </div>
        </div>

        <div className="planned-sets-heading">
          <div>
            <p className="eyebrow">ПЛАН ТРЕНИРОВКИ</p>
            <h3>Подходы <span>{sets.length}</span></h3>
          </div>
          <button
            className="program-secondary-button"
            type="button"
            onClick={addSet}
            disabled={isBusy || sets.length >= 100}
          >
            <Plus size={16} /> Добавить подход
          </button>
        </div>

        {sets.length > 0 ? (
          <ol className="training-exercise-detail-sets">
            {sets.map((set, index) => (
              <li key={set.id}>
                <span className="training-exercise-set-number">
                  <span className="training-exercise-set-prefix">Подход </span>
                  {index + 1}
                </span>
                <label>
                  <span>Повторения</span>
                  <input
                    type="number"
                    min="0"
                    step="1"
                    inputMode="numeric"
                    placeholder="—"
                    value={set.repetitions}
                    onChange={(event) => updateSet(set.id, 'repetitions', event.target.value)}
                    disabled={isBusy}
                  />
                </label>
                <label>
                  <span>Вес, кг</span>
                  <input
                    type="number"
                    min="0"
                    max="999.99"
                    step="0.1"
                    inputMode="decimal"
                    placeholder="—"
                    value={set.weight}
                    onChange={(event) => updateSet(set.id, 'weight', event.target.value)}
                    disabled={isBusy}
                  />
                </label>
                <button
                  className="program-icon-action program-delete-action planned-set-remove"
                  type="button"
                  aria-label={`Удалить подход ${index + 1}`}
                  onClick={() => setSets((current) => current.filter(({ id }) => id !== set.id))}
                  disabled={isBusy}
                >
                  <Trash2 size={16} />
                </button>
              </li>
            ))}
          </ol>
        ) : (
          <p className="planned-sets-hint">Подходов нет. Упражнение останется в тренировочном дне.</p>
        )}

        {error && <div className="program-alert" role="alert">{error}</div>}
        <div className="training-exercise-edit-actions">
          <button
            className="program-secondary-button program-delete-action"
            type="button"
            onClick={() => void handleDelete()}
            disabled={isBusy}
          >
            {isDeleting ? <LoaderCircle className="spin" size={16} /> : <Trash2 size={16} />}
            Удалить упражнение
          </button>
          <div>
            <button
              className="program-secondary-button"
              type="button"
              onClick={onClose}
              disabled={isBusy}
            >
              Закрыть
            </button>
            <button
              className="program-primary-button"
              type="button"
              onClick={() => void handleSave()}
              disabled={!hasChanges || isBusy}
            >
              {isSaving ? <LoaderCircle className="spin" size={16} /> : null}
              {isSaving ? 'Сохраняем...' : 'Сохранить изменения'}
            </button>
          </div>
        </div>
      </section>
    </div>
  )
}
