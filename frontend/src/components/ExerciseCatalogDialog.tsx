import { ArrowLeft, ChevronRight, Dumbbell, LoaderCircle, Plus, Trash2, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import {
  createTrainingDayExercise,
  getExercises,
  getMuscleGroups,
  type CatalogExercise,
  type CreatedTrainingDayExercise,
  type MuscleGroup,
} from '../lib/api'

type ExerciseCatalogDialogProps = {
  trainingDayId: number
  onClose: () => void
  onAdded: (exercise: CreatedTrainingDayExercise, catalogExercise: CatalogExercise) => void
}

type SetDraft = {
  id: number
  repetitions: string
  weight: string
}

type View = 'groups' | 'exercises' | 'detail'

function getErrorMessage(error: unknown): string {
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте подключение.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось выполнить запрос. Попробуйте ещё раз.'
}

function imageSource(path: string | null): string | null {
  if (!path) return null
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  return /\.[a-z\d]+$/i.test(normalizedPath) ? normalizedPath : `${normalizedPath}.webp`
}

export default function ExerciseCatalogDialog({
  trainingDayId,
  onClose,
  onAdded,
}: ExerciseCatalogDialogProps) {
  const [view, setView] = useState<View>('groups')
  const [groups, setGroups] = useState<MuscleGroup[]>([])
  const [selectedGroup, setSelectedGroup] = useState<MuscleGroup | null>(null)
  const [exercises, setExercises] = useState<CatalogExercise[]>([])
  const [selectedExercise, setSelectedExercise] = useState<CatalogExercise | null>(null)
  const [sets, setSets] = useState<SetDraft[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState('')
  const [canRetryLoad, setCanRetryLoad] = useState(false)
  const [nextSetId, setNextSetId] = useState(1)

  useEffect(() => {
    let cancelled = false
    setIsLoading(true)
    setError('')
    getMuscleGroups()
      .then((items) => {
        if (!cancelled) setGroups(items)
      })
      .catch((requestError: unknown) => {
        if (!cancelled) {
          setError(getErrorMessage(requestError))
          setCanRetryLoad(true)
        }
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [])

  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape' && !isSaving) onClose()
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [isSaving, onClose])

  async function selectGroup(group: MuscleGroup) {
    setSelectedGroup(group)
    setSelectedExercise(null)
    setExercises([])
    setView('exercises')
    setIsLoading(true)
    setError('')
    setCanRetryLoad(false)
    try {
      setExercises(await getExercises(group.id))
    } catch (requestError) {
      setError(getErrorMessage(requestError))
      setCanRetryLoad(true)
    } finally {
      setIsLoading(false)
    }
  }

  async function retryLoading() {
    setIsLoading(true)
    setError('')
    setCanRetryLoad(false)
    try {
      if (view === 'groups') {
        setGroups(await getMuscleGroups())
      } else if (selectedGroup) {
        setExercises(await getExercises(selectedGroup.id))
      }
    } catch (requestError) {
      setError(getErrorMessage(requestError))
      setCanRetryLoad(true)
    } finally {
      setIsLoading(false)
    }
  }

  function addSet() {
    setSets((current) => [...current, { id: nextSetId, repetitions: '', weight: '' }])
    setNextSetId((current) => current + 1)
  }

  function updateSet(id: number, field: 'repetitions' | 'weight', value: string) {
    setSets((current) => current.map((set) => set.id === id ? { ...set, [field]: value } : set))
  }

  async function handleAddExercise() {
    if (!selectedExercise) return

    const plannedSets: Array<{ repetitions: number | null; weight: number | null }> = []
    for (const [index, set] of sets.entries()) {
      const repetitionsText = set.repetitions.trim()
      const weightText = set.weight.trim()
      const repetitions = repetitionsText === '' ? null : Number(repetitionsText)
      const weight = weightText === '' ? null : Number(weightText)

      if (repetitions !== null && (!Number.isInteger(repetitions) || repetitions < 0)) {
        setError(`В подходе ${index + 1} укажите целое неотрицательное количество повторений.`)
        setCanRetryLoad(false)
        return
      }
      if (weight !== null && (!Number.isFinite(weight) || weight < 0 || weight > 999.99)) {
        setError(`В подходе ${index + 1} укажите вес от 0 до 999,99 кг.`)
        setCanRetryLoad(false)
        return
      }

      plannedSets.push({ repetitions, weight })
    }

    setIsSaving(true)
    setError('')
    setCanRetryLoad(false)
    try {
      const created = await createTrainingDayExercise(trainingDayId, selectedExercise.id, plannedSets)
      onAdded(created, selectedExercise)
    } catch (requestError) {
      setError(getErrorMessage(requestError))
      setCanRetryLoad(false)
    } finally {
      setIsSaving(false)
    }
  }

  function renderImage(path: string | null, alt: string, className: string) {
    const src = imageSource(path)
    return (
      <div className={`${className} exercise-catalog-image`}>
        {src ? (
          <img src={src} alt={alt} onError={(event) => event.currentTarget.remove()} />
        ) : (
          <Dumbbell size={25} aria-hidden="true" />
        )}
      </div>
    )
  }

  return (
    <div
      className="program-modal-backdrop exercise-catalog-backdrop"
      onClick={(event) => {
        if (event.target === event.currentTarget && !isSaving) onClose()
      }}
      role="presentation"
    >
      <section
        aria-labelledby="exercise-catalog-title"
        aria-modal="true"
        className={`program-create-dialog exercise-catalog-dialog${view === 'detail' ? ' exercise-detail-dialog' : ''}`}
        role="dialog"
      >
        <div className="program-create-heading exercise-catalog-heading">
          <div className="exercise-catalog-title-wrap">
            {view !== 'groups' && (
              <button
                className="program-icon-action exercise-catalog-back"
                type="button"
                aria-label={view === 'detail' ? 'К списку упражнений' : 'К выбору групп'}
                disabled={isSaving}
                onClick={() => {
                  setError('')
                  setCanRetryLoad(false)
                  if (view === 'detail') {
                    setView('exercises')
                    setSets([])
                  } else {
                    setView('groups')
                    setSelectedGroup(null)
                    setExercises([])
                  }
                }}
              >
                <ArrowLeft size={17} />
              </button>
            )}
            <div>
              <p className="eyebrow">
                {view === 'groups' ? 'КАТАЛОГ УПРАЖНЕНИЙ' : selectedGroup?.title}
              </p>
              <h2 id="exercise-catalog-title">
                {view === 'groups'
                  ? 'Выберите группу мышц'
                  : view === 'exercises'
                    ? 'Выберите упражнение'
                    : selectedExercise?.name}
              </h2>
            </div>
          </div>
          <button
            className="program-icon-action"
            type="button"
            aria-label="Закрыть каталог"
            disabled={isSaving}
            onClick={onClose}
          >
            <X size={18} />
          </button>
        </div>

        {error && (
          <div className="program-alert exercise-catalog-error" role="alert">
            <span>{error}</span>
            {!isSaving && canRetryLoad && (
              <button className="program-secondary-button" type="button" onClick={() => void retryLoading()}>
                Повторить
              </button>
            )}
          </div>
        )}

        {isLoading ? (
          <div className="program-state exercise-catalog-state">
            <LoaderCircle className="spin" size={21} /> Загружаем каталог...
          </div>
        ) : error ? null : view === 'groups' ? (
          groups.length === 0 ? (
            <div className="program-empty-state exercise-catalog-empty">
              <Dumbbell size={24} />
              <h3>Группы мышц пока не добавлены</h3>
            </div>
          ) : (
            <div className="exercise-group-grid">
              {groups.map((group) => (
                <button
                  className="exercise-group-card"
                  key={group.id}
                  type="button"
                  onClick={() => void selectGroup(group)}
                >
                  {renderImage(group.image_path, group.title, 'exercise-group-art')}
                  <span>{group.title}</span>
                  <ChevronRight className="exercise-card-chevron" size={17} />
                </button>
              ))}
            </div>
          )
        ) : view === 'exercises' ? (
          exercises.length === 0 ? (
            <div className="program-empty-state exercise-catalog-empty">
              <Dumbbell size={24} />
              <h3>В этой группе пока нет упражнений</h3>
              <p>Выберите другую группу мышц или вернитесь позже.</p>
            </div>
          ) : (
            <div className="exercise-catalog-grid">
              {exercises.map((exercise) => (
                <button
                  className="exercise-catalog-card"
                  key={exercise.id}
                  type="button"
                  onClick={() => {
                    setSelectedExercise(exercise)
                    setSets([])
                    setError('')
                    setCanRetryLoad(false)
                    setView('detail')
                  }}
                >
                  {renderImage(exercise.image_path, exercise.title, 'exercise-card-art')}
                  <span>{exercise.title}</span>
                  <ChevronRight className="exercise-card-chevron" size={17} />
                </button>
              ))}
            </div>
          )
        ) : selectedExercise ? (
          <>
            <div className="exercise-detail-layout">
              {renderImage(selectedExercise.image_path, selectedExercise.title, 'exercise-detail-art')}
              <div className="exercise-detail-copy">
                <p className="eyebrow">ТЕХНИКА ВЫПОЛНЕНИЯ</p>
                <p className="exercise-detail-description">
                  {selectedExercise.description || 'Описание техники пока не добавлено.'}
                </p>
              </div>
            </div>

            <div className="planned-sets-heading">
              <div>
                <p className="eyebrow">ПЛАН ТРЕНИРОВКИ</p>
                <h3>Подходы <span>{sets.length}</span></h3>
              </div>
              <button className="program-secondary-button" type="button" onClick={addSet} disabled={sets.length >= 100}>
                <Plus size={16} /> Добавить подход
              </button>
            </div>

            {sets.length > 0 ? (
              <div className="planned-set-list">
                {sets.map((set, index) => (
                  <div className="planned-set-row" key={set.id}>
                    <span className="planned-set-number">{String(index + 1).padStart(2, '0')}</span>
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
                        disabled={isSaving}
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
                        disabled={isSaving}
                      />
                    </label>
                    <button
                      className="program-icon-action program-delete-action planned-set-remove"
                      type="button"
                      aria-label={`Удалить подход ${index + 1}`}
                      onClick={() => setSets((current) => current.filter(({ id }) => id !== set.id))}
                      disabled={isSaving}
                    >
                      <Trash2 size={16} />
                    </button>
                  </div>
                ))}
              </div>
            ) : (
              <p className="planned-sets-hint">Можно добавить упражнение без подходов, а заполнить план позже.</p>
            )}

            <div className="exercise-detail-actions">
              <button className="program-secondary-button" type="button" onClick={() => setView('exercises')} disabled={isSaving}>
                Назад к упражнениям
              </button>
              <button className="program-primary-button" type="button" onClick={() => void handleAddExercise()} disabled={isSaving}>
                {isSaving ? <LoaderCircle className="spin" size={16} /> : <Plus size={16} />}
                Добавить в тренировочный день
              </button>
            </div>
          </>
        ) : null}
      </section>
    </div>
  )
}
