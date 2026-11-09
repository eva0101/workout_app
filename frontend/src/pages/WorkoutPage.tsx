import {
  ArrowLeft,
  Check,
  CheckCircle2,
  ChevronRight,
  Circle,
  Dumbbell,
  LoaderCircle,
  Plus,
  Trophy,
  X,
} from 'lucide-react'
import { type FormEvent, useEffect, useMemo, useRef, useState } from 'react'
import {
  getExercises,
  type CatalogExercise,
  type Workout,
  type WorkoutExercise,
} from '../lib/api'

type WorkoutPageProps = {
  workout: Workout
  title: string
  isSavingSet: boolean
  isCompletingExercise: boolean
  isCompleting: boolean
  error: string
  onAddSet: (exerciseId: number, repsDone: number, weight: number | null) => Promise<boolean>
  onCompleteExercise: (exerciseId: number) => Promise<boolean>
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

function formatUnfinishedSummary(count: number) {
  if (count === 1) return '1 упражнение останется незавершённым в истории.'
  if (count >= 2 && count <= 4) return `${count} упражнения останутся незавершёнными в истории.`
  return `${count} упражнений останутся незавершёнными в истории.`
}

function imageSource(path: string | null | undefined): string | null {
  if (!path) return null
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  return /\.[a-z\d]+$/i.test(normalizedPath) ? normalizedPath : `${normalizedPath}.webp`
}

function getErrorMessage(error: unknown) {
  if (error instanceof TypeError) return 'Не удалось загрузить каталог упражнений. Проверьте подключение.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось загрузить данные упражнения.'
}

function ExerciseImage({
  exercise,
  catalogExercise,
  className,
}: {
  exercise: WorkoutExercise
  catalogExercise?: CatalogExercise
  className: string
}) {
  const [imageFailed, setImageFailed] = useState(false)
  const imagePath = imageSource(catalogExercise?.image_path)
  const exerciseName = exercise.title || exercise.name

  useEffect(() => {
    setImageFailed(false)
  }, [imagePath])

  return (
    <div className={className}>
      {imagePath && !imageFailed ? (
        <img
          src={imagePath}
          alt={exerciseName}
          onError={() => setImageFailed(true)}
        />
      ) : (
        <Dumbbell size={className.includes('large') ? 34 : 21} aria-hidden="true" />
      )}
    </div>
  )
}

export default function WorkoutPage({
  workout,
  title,
  isSavingSet,
  isCompletingExercise,
  isCompleting,
  error,
  onAddSet,
  onCompleteExercise,
  onComplete,
  onBack,
}: WorkoutPageProps) {
  const [selectedExerciseId, setSelectedExerciseId] = useState<number | null>(null)
  const [inputs, setInputs] = useState<Record<number, SetInput>>({})
  const [fatigueScore, setFatigueScore] = useState(5)
  const [isFinishConfirmationOpen, setIsFinishConfirmationOpen] = useState(false)
  const [isTechniqueExpanded, setIsTechniqueExpanded] = useState(false)
  const [catalogExercises, setCatalogExercises] = useState<CatalogExercise[]>([])
  const [isCatalogLoading, setIsCatalogLoading] = useState(true)
  const [catalogError, setCatalogError] = useState('')
  const actualSetsListRef = useRef<HTMLDivElement>(null)
  const previousActualSetsRef = useRef<{ exerciseId: number; count: number } | null>(null)

  useEffect(() => {
    let cancelled = false
    setIsCatalogLoading(true)
    setCatalogError('')
    getExercises()
      .then((items) => {
        if (!cancelled) setCatalogExercises(items)
      })
      .catch((loadError: unknown) => {
        if (!cancelled) setCatalogError(getErrorMessage(loadError))
      })
      .finally(() => {
        if (!cancelled) setIsCatalogLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  const completedCount = workout.exercises.filter((exercise) => Boolean(exercise.completed_at)).length
  const unfinishedCount = workout.exercises.length - completedCount
  const progress = workout.exercises.length === 0
    ? 0
    : Math.round((completedCount / workout.exercises.length) * 100)
  const selectedExercise = workout.exercises.find((exercise) => exercise.id === selectedExerciseId) ?? null
  const exerciseCatalog = useMemo(
    () => new Map(catalogExercises.map((exercise) => [exercise.id, exercise])),
    [catalogExercises],
  )

  useEffect(() => {
    if (!selectedExercise) {
      previousActualSetsRef.current = null
      return
    }

    const currentCount = selectedExercise.actual_sets.length
    const previous = previousActualSetsRef.current
    if (previous?.exerciseId === selectedExercise.id && currentCount > previous.count) {
      const list = actualSetsListRef.current
      list?.scrollTo({ top: list.scrollHeight, behavior: 'smooth' })
    }
    previousActualSetsRef.current = { exerciseId: selectedExercise.id, count: currentCount }
  }, [selectedExercise])

  async function submitSet(event: FormEvent<HTMLFormElement>, exerciseId: number) {
    event.preventDefault()
    const input = inputs[exerciseId] ?? { reps: '', weight: '' }
    const repsDone = Number(input.reps)
    const weight = input.weight.trim() === '' ? null : Number(input.weight)
    if (!Number.isInteger(repsDone) || repsDone < 0 ||
      (weight !== null && (!Number.isFinite(weight) || weight < 0))) {
      return
    }
    if (await onAddSet(exerciseId, repsDone, weight)) {
      setInputs((current) => ({ ...current, [exerciseId]: { reps: '', weight: '' } }))
    }
  }

  function requestFinishWorkout() {
    if (isCompleting || isSavingSet || isCompletingExercise) return
    if (unfinishedCount > 0) {
      setIsFinishConfirmationOpen(true)
      return
    }
    onComplete(fatigueScore)
  }

  function finishDespiteUncompletedExercises() {
    if (isCompleting) return
    setIsFinishConfirmationOpen(false)
    onComplete(fatigueScore)
  }

  function renderExerciseCard(exercise: WorkoutExercise, index: number) {
    const isCompleted = Boolean(exercise.completed_at)
    const catalogItem = exerciseCatalog.get(exercise.exercise_id)
    const planCount = exercise.planned_sets?.length ?? 0
    return (
      <li key={exercise.id}>
        <button
          className={`session-exercise-card${isCompleted ? ' session-exercise-card-completed' : ''}`}
          type="button"
          onClick={() => setSelectedExerciseId(exercise.id)}
          aria-label={`${exercise.title || exercise.name}, ${isCompleted ? 'выполнено' : 'не выполнено'}, ${planCount} плановых подходов`}
        >
          <ExerciseImage
            exercise={exercise}
            catalogExercise={catalogItem}
            className="session-exercise-image"
          />
          <span className="session-exercise-main">
            <span className="session-exercise-name">{exercise.title || exercise.name}</span>
            <span className="session-exercise-meta">
              {String(index + 1).padStart(2, '0')} <span>·</span> {planCount} {planCount === 1 ? 'подход' : planCount > 1 && planCount < 5 ? 'подхода' : 'подходов'}
              {exercise.actual_sets.length > 0 && <span>· {exercise.actual_sets.length} выполнено</span>}
            </span>
          </span>
          <span className={`session-exercise-status${isCompleted ? ' session-exercise-status-done' : ''}`}>
            {isCompleted ? <CheckCircle2 size={22} /> : <Circle size={21} />}
          </span>
          <ChevronRight className="session-exercise-chevron" size={18} />
        </button>
      </li>
    )
  }

  if (selectedExercise) {
    const catalogItem = exerciseCatalog.get(selectedExercise.exercise_id)
    const isCompleted = Boolean(selectedExercise.completed_at)
    const input = inputs[selectedExercise.id] ?? { reps: '', weight: '' }

    return (
      <section className="page-content workout-session-page session-exercise-detail-page">
        <button
          className="program-back-button"
          type="button"
          onClick={() => setSelectedExerciseId(null)}
        >
          <ArrowLeft size={16} /> Все упражнения
        </button>

        <header className="session-exercise-header">
          <div className="session-exercise-header-copy">
            <p className="eyebrow">{title || 'ТРЕНИРОВКА'} · УПРАЖНЕНИЕ</p>
            <h1>{selectedExercise.title || selectedExercise.name}</h1>
          </div>
          {isCompleted && (
            <span className="session-completed-badge"><CheckCircle2 size={15} /> Завершено</span>
          )}
        </header>

        {(error || catalogError) && (
          <div className="program-alert mt-4" role="alert">{error || catalogError}</div>
        )}

        <div className="session-exercise-layout">
          <div className="session-exercise-sets-column">
            <section className="session-sets-card session-actual-card">
              <div className="session-subsection-heading">
                <span className="session-subsection-icon session-subsection-icon-accent"><Check size={16} /></span>
                <div>
                  <p className="eyebrow">ВАШ РЕЗУЛЬТАТ</p>
                  <h2>Фактические подходы <span>{selectedExercise.actual_sets.length}</span></h2>
                </div>
              </div>
              <div className="session-actual-list-scroll" ref={actualSetsListRef}>
                {selectedExercise.actual_sets.length === 0 ? (
                  <p className="session-description-muted">Добавленные подходы появятся здесь.</p>
                ) : (
                  <ol className="session-set-list">
                    {selectedExercise.actual_sets.map((set) => (
                      <li className="session-set-row session-set-row-actual" key={set.id}>
                        <span className="session-set-index">{String(set.set_number).padStart(2, '0')}</span>
                        <span>{formatValue(set.reps_done)} <small>повт.</small></span>
                        <span>{formatValue(set.weight)} <small>кг</small></span>
                      </li>
                    ))}
                  </ol>
                )}
              </div>
              {!isCompleted && (
                <form className="session-set-form" onSubmit={(event) => void submitSet(event, selectedExercise.id)}>
                  <label>
                    <span>Повторения</span>
                    <input
                      type="number"
                      min="0"
                      step="1"
                      required
                      inputMode="numeric"
                      placeholder="Например, 10"
                      value={input.reps}
                      onChange={(event) => setInputs((current) => ({
                        ...current,
                        [selectedExercise.id]: { ...input, reps: event.target.value },
                      }))}
                    />
                  </label>
                  <label>
                    <span>Вес, кг</span>
                    <input
                      type="number"
                      min="0"
                      step="0.01"
                      inputMode="decimal"
                      placeholder="Например, 40"
                      value={input.weight}
                      onChange={(event) => setInputs((current) => ({
                        ...current,
                        [selectedExercise.id]: { ...input, weight: event.target.value },
                      }))}
                    />
                  </label>
                  <button
                    className="program-primary-button session-add-set-button"
                    type="submit"
                    disabled={isSavingSet || isCompletingExercise}
                  >
                    {isSavingSet ? <LoaderCircle className="spin" size={16} /> : <Plus size={16} />}
                    Завершить подход
                  </button>
                </form>
              )}
            </section>

            <section className="session-sets-card session-planned-card">
              <div className="session-subsection-heading">
                <span className="session-subsection-icon"><Trophy size={16} /></span>
                <div>
                  <p className="eyebrow">ИЗ ПЛАНА</p>
                  <h2>Плановые подходы <span>{selectedExercise.planned_sets.length}</span></h2>
                </div>
              </div>
              {selectedExercise.planned_sets.length === 0 ? (
                <p className="session-description-muted">Плановые подходы не заданы.</p>
              ) : (
                <ol className="session-set-list">
                  {selectedExercise.planned_sets.map((set) => (
                    <li className="session-set-row session-set-row-planned" key={set.id}>
                      <span className="session-set-index">{String(set.set_number).padStart(2, '0')}</span>
                      <span>{formatValue(set.repetitions)} <small>повт.</small></span>
                      <span>{formatValue(set.weight)} <small>кг</small></span>
                    </li>
                  ))}
                </ol>
              )}
            </section>

            <div className="session-detail-actions">
              {isCompleted ? (
                <button className="program-secondary-button" type="button" onClick={() => setSelectedExerciseId(null)}>
                  <CheckCircle2 size={16} /> К списку упражнений
                </button>
              ) : (
                <button
                  className="program-primary-button"
                  type="button"
                  disabled={isCompletingExercise || isSavingSet}
                  onClick={() => void onCompleteExercise(selectedExercise.id).then((success) => {
                    if (success) setSelectedExerciseId(null)
                  })}
                >
                  {isCompletingExercise
                    ? <LoaderCircle className="spin" size={16} />
                    : <CheckCircle2 size={16} />}
                  Завершить упражнение
                </button>
              )}
              <button className="program-secondary-button" type="button" onClick={() => setSelectedExerciseId(null)}>
                Вернуться к списку
              </button>
            </div>
          </div>

          <div className="session-exercise-technique-column">
            <ExerciseImage
              exercise={selectedExercise}
              catalogExercise={catalogItem}
              className="session-exercise-image session-detail-exercise-image"
            />
            <section className={`session-technique-card panel-card${isTechniqueExpanded ? ' session-technique-expanded' : ''}`}>
              <div className="session-technique-desktop-heading">
                <span className="session-subsection-heading">
                  <span className="session-subsection-icon"><Dumbbell size={16} /></span>
                  <span>
                    <span className="eyebrow">ТЕХНИКА</span>
                    <span className="session-technique-title">Техника выполнения</span>
                  </span>
                </span>
              </div>
              <button
                className="session-technique-toggle"
                type="button"
                aria-expanded={isTechniqueExpanded}
                onClick={() => setIsTechniqueExpanded((expanded) => !expanded)}
              >
                <span className="session-subsection-heading">
                  <span className="session-subsection-icon"><Dumbbell size={16} /></span>
                  <span>
                    <span className="eyebrow">ТЕХНИКА</span>
                    <span className="session-technique-title">Техника выполнения</span>
                  </span>
                </span>
                <ChevronRight className="session-technique-chevron" size={18} />
              </button>
              <div className={`session-technique-content${isTechniqueExpanded ? ' session-technique-content-expanded' : ''}`}>
                <div className="session-technique-content-inner">
                  {isCatalogLoading ? (
                    <p className="session-description-muted"><LoaderCircle className="spin inline-block" size={15} /> Загружаем данные упражнения...</p>
                  ) : catalogError ? (
                    <p className="session-description-muted">Описание и изображение пока недоступны.</p>
                  ) : (
                    <p className="session-technique-description">
                      {catalogItem?.description || 'Описание техники для этого упражнения не добавлено.'}
                    </p>
                  )}
                </div>
              </div>
            </section>
          </div>
        </div>
      </section>
    )
  }

  return (
    <section className="page-content workout-session-page">
      <button className="program-back-button" type="button" onClick={onBack}>
        <ArrowLeft size={16} /> К обзору
      </button>

      <header className="session-overview-header">
        <div>
          <p className="eyebrow">ВАША ТРЕНИРОВКА</p>
          <h1>{title || workout.training_day_name || `Тренировка #${workout.id}`}</h1>
          <div className="session-started-at">
            <span>Начата</span>
            <time dateTime={workout.started_at}>
              {new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', hour: '2-digit', minute: '2-digit' }).format(new Date(workout.started_at))}
            </time>
          </div>
        </div>
        <div className="session-progress-ring" aria-label={`Выполнено ${completedCount} из ${workout.exercises.length} упражнений`}>
          <span>{completedCount}<small>/{workout.exercises.length}</small></span>
          <span className="session-progress-caption">упражнений</span>
        </div>
      </header>

      {(error || catalogError) && (
        <div className="program-alert mt-4" role="alert">{error || catalogError}</div>
      )}
      {catalogError && (
        <button
          className="program-secondary-button mt-3"
          type="button"
          onClick={() => {
            setIsCatalogLoading(true)
            setCatalogError('')
            void getExercises()
              .then(setCatalogExercises)
              .catch((loadError: unknown) => setCatalogError(getErrorMessage(loadError)))
              .finally(() => setIsCatalogLoading(false))
          }}
        >
          Повторить загрузку изображений
        </button>
      )}

      <section className="session-progress-card panel-card" aria-label="Прогресс тренировки">
        <div className="session-progress-copy">
          <div>
            <p className="eyebrow">ПРОГРЕСС ТРЕНИРОВКИ</p>
            <h2>{completedCount} из {workout.exercises.length} упражнений</h2>
          </div>
          <span>{progress}%</span>
        </div>
        <div
          className="session-progress-track"
          role="progressbar"
          aria-valuenow={completedCount}
          aria-valuemin={0}
          aria-valuemax={Math.max(1, workout.exercises.length)}
          aria-label="Завершённые упражнения"
        >
          <span style={{ width: `${progress}%` }} />
        </div>
      </section>

      <section className="session-exercises-section" aria-labelledby="session-exercises-title">
        <div className="session-list-heading">
          <div>
            <p className="eyebrow">ПЛАН ТРЕНИРОВКИ</p>
            <h2 id="session-exercises-title">Упражнения <span>{workout.exercises.length}</span></h2>
          </div>
          <p>Выберите любое упражнение, чтобы начать</p>
        </div>

        {workout.exercises.length === 0 ? (
          <div className="session-empty panel-card">
            <span className="session-empty-icon"><Dumbbell size={23} /></span>
            <h3>В этой тренировке нет упражнений</h3>
            <p>Завершите тренировку или вернитесь к обзору.</p>
          </div>
        ) : (
          <ol className="session-exercise-list">
            {workout.exercises.map(renderExerciseCard)}
          </ol>
        )}
      </section>

      <section className="session-finish-card panel-card">
        <div className="session-finish-copy">
          <span className="session-finish-icon"><Trophy size={19} /></span>
          <div>
            <p className="eyebrow">ФИНИШ</p>
            <h2>Завершить тренировку</h2>
            <p>
              {unfinishedCount === 0
                ? 'Все упражнения отмечены как выполненные.'
                : `${unfinishedCount} ${unfinishedCount === 1 ? 'упражнение осталось' : 'упражнений осталось'} незавершёнными.`}
            </p>
          </div>
        </div>
        <label className="session-fatigue-field">
          <span>Самочувствие</span>
          <select
            value={fatigueScore}
            onChange={(event) => setFatigueScore(Number(event.target.value))}
            disabled={isCompleting}
          >
            {Array.from({ length: 11 }, (_, score) => (
              <option key={score} value={score}>{score} / 10</option>
            ))}
          </select>
        </label>
        <button
          className="program-primary-button session-finish-button"
          type="button"
          disabled={isCompleting || isSavingSet || isCompletingExercise}
          onClick={requestFinishWorkout}
        >
          {isCompleting ? <LoaderCircle className="spin" size={16} /> : <Check size={16} />}
          Завершить тренировку
        </button>
      </section>

      {isFinishConfirmationOpen && (
        <div
          className="program-modal-backdrop session-confirmation-backdrop"
          onClick={(event) => {
            if (event.target === event.currentTarget && !isCompleting) setIsFinishConfirmationOpen(false)
          }}
          role="presentation"
        >
          <section
            aria-labelledby="finish-confirmation-title"
            aria-modal="true"
            className="program-create-dialog session-confirmation-dialog"
            role="dialog"
          >
            <div className="session-confirmation-icon"><Trophy size={21} /></div>
            <button
              className="program-icon-action session-confirmation-close"
              type="button"
              aria-label="Закрыть диалог"
              disabled={isCompleting}
              onClick={() => setIsFinishConfirmationOpen(false)}
            >
              <X size={18} />
            </button>
            <p className="eyebrow">ЗАВЕРШЕНИЕ ТРЕНИРОВКИ</p>
            <h2 id="finish-confirmation-title">Остались упражнения</h2>
            <p className="session-confirmation-copy">
              {formatUnfinishedSummary(unfinishedCount)} Вы можете продолжить или завершить тренировку сейчас.
            </p>
            <div className="session-confirmation-actions">
              <button
                className="program-secondary-button"
                type="button"
                disabled={isCompleting}
                onClick={() => setIsFinishConfirmationOpen(false)}
              >
                Продолжить тренировку
              </button>
              <button
                className="program-primary-button"
                type="button"
                disabled={isCompleting}
                onClick={finishDespiteUncompletedExercises}
              >
                {isCompleting ? <LoaderCircle className="spin" size={15} /> : <Check size={15} />}
                Всё равно завершить
              </button>
            </div>
          </section>
        </div>
      )}
    </section>
  )
}
