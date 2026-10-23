import { ArrowLeft, Dumbbell, FolderPlus, LoaderCircle, Pencil, Plus, RefreshCw, Trash2, X } from 'lucide-react'
import { type FormEvent, useEffect, useState } from 'react'
import {
  createProgram,
  createTrainingDay,
  deleteProgram,
  deleteTrainingDay,
  getProgram,
  getPrograms,
  getTrainingDay,
  type Program,
  type ProgramDetails,
  type TrainingDayDetails,
  type TrainingDaySummary,
  renameProgram,
} from '../lib/api'

type ProgramsPageProps = {
  selectedProgramId: number | null
  selectedDayId: number | null
  onOpenProgram: (id: number) => void
  onOpenDay: (id: number) => void
  onBackToPrograms: () => void
  onBackToProgram: () => void
}

function getErrorMessage(error: unknown): string {
  const status = typeof error === 'object' && error !== null && 'status' in error
    ? (error as { status?: number }).status
    : undefined

  if (status === 401) return 'Сессия завершена. Войдите в аккаунт снова.'
  if (status === 404) return 'Программа или тренировочный день не найдены. Обновите страницу и попробуйте снова.'
  if (status === 409) return 'Не удалось выполнить действие из-за конфликта данных.'
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте, запущен ли backend.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось выполнить запрос. Попробуйте ещё раз.'
}

function formatDate(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Дата неизвестна'
  return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' }).format(date)
}

export default function ProgramsPage({
  selectedProgramId,
  selectedDayId,
  onOpenProgram,
  onOpenDay,
  onBackToPrograms,
  onBackToProgram,
}: ProgramsPageProps) {
  const [programs, setPrograms] = useState<Program[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [pageError, setPageError] = useState('')
  const [formError, setFormError] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [name, setName] = useState('')
  const [isSaving, setIsSaving] = useState(false)
  const [details, setDetails] = useState<ProgramDetails | null>(null)
  const [isDetailsLoading, setIsDetailsLoading] = useState(false)
  const [detailsError, setDetailsError] = useState('')
  const [dayDetails, setDayDetails] = useState<TrainingDayDetails | null>(null)
  const [isDayLoading, setIsDayLoading] = useState(false)
  const [dayError, setDayError] = useState('')
  const [dayActionError, setDayActionError] = useState('')
  const [isDaySaving, setIsDaySaving] = useState(false)
  const [isDayDeleting, setIsDayDeleting] = useState(false)

  async function loadPrograms() {
    setIsLoading(true)
    setPageError('')
    try {
      setPrograms(await getPrograms())
    } catch (error) {
      setPageError(getErrorMessage(error))
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    void loadPrograms()
  }, [])

  useEffect(() => {
    if (selectedProgramId === null) {
      setDetails(null)
      setDetailsError('')
      return
    }

    let cancelled = false
    setIsDetailsLoading(true)
    setDetailsError('')
    setDetails(null)
    getProgram(selectedProgramId)
      .then((programDetails) => {
        if (!cancelled) setDetails(programDetails)
      })
      .catch((error: unknown) => {
        if (!cancelled) setDetailsError(getErrorMessage(error))
      })
      .finally(() => {
        if (!cancelled) setIsDetailsLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [selectedProgramId])

  useEffect(() => {
    if (selectedDayId === null) {
      setDayDetails(null)
      setDayError('')
      return
    }

    let cancelled = false
    setIsDayLoading(true)
    setDayError('')
    setDayDetails(null)
    getTrainingDay(selectedDayId)
      .then((trainingDay) => {
        if (!cancelled) setDayDetails(trainingDay)
      })
      .catch((error: unknown) => {
        if (!cancelled) setDayError(getErrorMessage(error))
      })
      .finally(() => {
        if (!cancelled) setIsDayLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [selectedDayId])

  async function handleCreateTrainingDay() {
    if (!details) return
    setIsDaySaving(true)
    setDayActionError('')
    try {
      const createdDay: TrainingDaySummary = await createTrainingDay(details.id)
      setDetails((current) => current
        ? { ...current, training_days: [...(current.training_days ?? []), createdDay] }
        : current,
      )
    } catch (error) {
      setDayActionError(getErrorMessage(error))
    } finally {
      setIsDaySaving(false)
    }
  }

  async function handleDeleteTrainingDay(day: TrainingDaySummary) {
    if (!details) return
    if (!window.confirm(`Удалить День ${day.DayNumber}? Это действие нельзя отменить.`)) return

    setIsDayDeleting(true)
    setDayActionError('')
    try {
      await deleteTrainingDay(day.ID)
      setDetails((current) => current
        ? {
            ...current,
            training_days: (current.training_days ?? []).filter(({ ID }) => ID !== day.ID),
          }
        : current,
      )
    } catch (error) {
      setDayActionError(getErrorMessage(error))
    } finally {
      setIsDayDeleting(false)
    }
  }

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const trimmedName = name.trim()
    if (!trimmedName) {
      setFormError('Введите название программы.')
      return
    }

    setIsSaving(true)
    setFormError('')
    try {
      await createProgram(trimmedName)
      setName('')
      setIsCreateOpen(false)
      await loadPrograms()
    } catch (error) {
      setFormError(getErrorMessage(error))
    } finally {
      setIsSaving(false)
    }
  }

  async function handleRename(event: FormEvent<HTMLFormElement>, id: number) {
    event.preventDefault()
    const trimmedName = name.trim()
    if (!trimmedName) {
      setFormError('Введите название программы.')
      return
    }

    setIsSaving(true)
    setFormError('')
    try {
      await renameProgram(id, trimmedName)
      setPrograms((current) => current.map((program) =>
        program.id === id ? { ...program, name: trimmedName } : program,
      ))
      if (details?.id === id) setDetails({ ...details, name: trimmedName })
      setEditingId(null)
      setName('')
    } catch (error) {
      setFormError(getErrorMessage(error))
    } finally {
      setIsSaving(false)
    }
  }

  async function handleDelete(program: Program) {
    if (!window.confirm(`Удалить программу «${program.name}»? Это действие нельзя отменить.`)) return

    setPageError('')
    try {
      await deleteProgram(program.id)
      setPrograms((current) => current.filter(({ id }) => id !== program.id))
      if (selectedProgramId === program.id) onBackToPrograms()
    } catch (error) {
      setPageError(getErrorMessage(error))
    }
  }

  if (selectedProgramId !== null) {
    const trainingDays = details?.training_days ?? []
    const lastDayNumber = trainingDays.reduce((last, day) => Math.max(last, day.DayNumber), 0)

    if (selectedDayId !== null) {
      return (
        <section className="page-content programs-page">
          <button className="program-back-button" type="button" onClick={onBackToProgram}>
            <ArrowLeft size={16} /> К программе
          </button>
          {isDayLoading || isDetailsLoading || (!details && !detailsError) ? (
            <div className="program-state"><LoaderCircle className="spin" size={22} /> Загружаем тренировочный день...</div>
          ) : dayError ? (
            <div className="program-state program-error" role="alert">
              <p>{dayError}</p>
              <button className="program-secondary-button" type="button" onClick={onBackToProgram}>Вернуться к программе</button>
            </div>
          ) : detailsError ? (
            <div className="program-state program-error" role="alert">
              <p>{detailsError}</p>
              <button className="program-secondary-button" type="button" onClick={onBackToProgram}>Вернуться к программе</button>
            </div>
          ) : dayDetails ? (
            <>
              <div className="program-detail-heading">
                <div className="program-detail-icon"><Dumbbell size={23} /></div>
                <div>
                  <p className="eyebrow">{details?.name}</p>
                  <h2>День {dayDetails.DayNumber}</h2>
                </div>
              </div>
              <div className="training-day-content">
                <div className="training-day-section-heading">
                  <div>
                    <p className="eyebrow">ПЛАН ДНЯ</p>
                    <h3>Упражнения</h3>
                  </div>
                  <button className="program-primary-button program-unavailable-button" type="button" disabled title="Добавление упражнений появится позже">
                    <Plus size={16} /> Добавить упражнение · недоступно
                  </button>
                </div>
                {dayDetails.Exercises.length === 0 ? (
                  <div className="program-empty-state training-day-empty">
                    <div className="program-empty-icon"><Dumbbell size={22} /></div>
                    <h3>В этом тренировочном дне пока нет упражнений</h3>
                    <p>Добавление упражнений будет доступно на следующем этапе.</p>
                  </div>
                ) : (
                  <ol className="training-exercise-list">
                    {dayDetails.Exercises.map((exercise, index) => (
                      <li className="training-exercise-card" key={exercise.ID}>
                        <span className="exercise-number">{String(index + 1).padStart(2, '0')}</span>
                        <span className="training-exercise-name">Упражнение #{exercise.ExerciseID}</span>
                        <span className="training-exercise-sets">{exercise.Sets} подходов</span>
                      </li>
                    ))}
                  </ol>
                )}
              </div>
            </>
          ) : (
            <div className="program-state">Не удалось загрузить тренировочный день.</div>
          )}
        </section>
      )
    }

    return (
      <section className="page-content programs-page">
        <button className="program-back-button" type="button" onClick={onBackToPrograms}>
          <ArrowLeft size={16} /> Все программы
        </button>
        {isDetailsLoading ? (
          <div className="program-state"><LoaderCircle className="spin" size={22} /> Загружаем программу...</div>
        ) : detailsError ? (
          <div className="program-state program-error" role="alert">
            <p>{detailsError}</p>
            <button className="program-secondary-button" type="button" onClick={onBackToPrograms}>К списку программ</button>
          </div>
        ) : details ? (
          <>
            <div className="program-detail-heading">
              <div className="program-detail-icon"><Dumbbell size={23} /></div>
              <div>
                <p className="eyebrow">ПРОГРАММА ТРЕНИРОВОК</p>
                <h2>{details.name}</h2>
              </div>
            </div>
            <div className="training-day-section-heading">
              <div>
                <p className="eyebrow">СТРУКТУРА ПРОГРАММЫ</p>
                <h3>Тренировочные дни</h3>
              </div>
              <button
                className="program-primary-button"
                type="button"
                onClick={() => void handleCreateTrainingDay()}
                disabled={isDaySaving || isDayDeleting}
              >
                <Plus size={17} /> {isDaySaving ? 'Добавление...' : 'Добавить тренировочный день'}
              </button>
            </div>
            {dayActionError && <div className="program-alert" role="alert">{dayActionError}</div>}
            {trainingDays.length === 0 ? (
              <div className="program-empty-state">
                <div className="program-empty-icon"><FolderPlus size={23} /></div>
                <h3>В программе пока нет тренировочных дней</h3>
                <p>Добавь первый день, чтобы начать наполнять программу.</p>
              </div>
            ) : (
              <div className="training-day-grid">
                {trainingDays.map((day) => (
                  <article className="training-day-card" key={day.ID}>
                    <button
                      className="training-day-card-open"
                      type="button"
                      onClick={() => onOpenDay(day.ID)}
                    >
                      <span className="program-card-icon"><Dumbbell size={20} /></span>
                      <span className="training-day-card-label">День {day.DayNumber}</span>
                      <span className="training-day-card-hint">Открыть упражнения</span>
                    </button>
                    {day.DayNumber === lastDayNumber && (
                      <button
                        className="program-icon-action program-delete-action training-day-delete"
                        type="button"
                        aria-label={`Удалить День ${day.DayNumber}`}
                        title={isDayDeleting ? 'Удаление...' : 'Удалить последний день'}
                        onClick={() => void handleDeleteTrainingDay(day)}
                        disabled={isDayDeleting || isDaySaving}
                      >
                        <Trash2 size={16} />
                      </button>
                    )}
                  </article>
                ))}
              </div>
            )}
          </>
        ) : null}
      </section>
    )
  }

  return (
    <section className="page-content programs-page">
      <div className="programs-heading">
        <div>
          <p className="eyebrow">ТВОЙ ПЛАН</p>
          <h2>Программы тренировок</h2>
          <p className="programs-subtitle">Создавай и настраивай свои планы силовых тренировок.</p>
        </div>
        <button
          className="program-primary-button"
          type="button"
          onClick={() => {
            setIsCreateOpen((open) => !open)
            setEditingId(null)
            setName('')
            setFormError('')
          }}
        >
          <Plus size={17} /> Создать программу
        </button>
      </div>

      {isCreateOpen && (
        <div
          className="program-modal-backdrop"
          onClick={(event) => {
            if (event.target === event.currentTarget && !isSaving) setIsCreateOpen(false)
          }}
          role="presentation"
        >
          <section
            aria-labelledby="create-program-title"
            aria-modal="true"
            className="program-create-dialog"
            role="dialog"
          >
            <div className="program-create-heading">
              <div>
                <p className="eyebrow">НОВАЯ ПРОГРАММА</p>
                <h2 id="create-program-title">Создать программу</h2>
              </div>
              <button
                className="program-icon-action"
                type="button"
                aria-label="Закрыть"
                disabled={isSaving}
                onClick={() => setIsCreateOpen(false)}
              >
                <X size={18} />
              </button>
            </div>
            <form className="program-create-form" onSubmit={handleCreate}>
              <label htmlFor="new-program-name">Название программы</label>
              <input
                id="new-program-name"
                autoFocus
                maxLength={100}
                value={name}
                onChange={(event) => setName(event.target.value)}
                placeholder="Например, Силовая база"
              />
              <p className="program-form-hint">Название должно содержать от 3 до 100 символов.</p>
              {formError && <p className="program-inline-error" role="alert">{formError}</p>}
              <div className="program-create-actions">
                <button className="program-secondary-button" type="button" disabled={isSaving} onClick={() => setIsCreateOpen(false)}>
                  Отмена
                </button>
                <button className="program-primary-button" type="submit" disabled={isSaving}>
                  {isSaving ? 'Сохранение...' : 'Сохранить'}
                </button>
              </div>
            </form>
          </section>
        </div>
      )}

      {pageError && (
        <div className="program-alert" role="alert">
          <span>{pageError}</span>
          <button className="program-secondary-button" type="button" onClick={() => void loadPrograms()}>
            <RefreshCw size={14} /> Повторить
          </button>
        </div>
      )}

      {isLoading ? (
        <div className="program-state"><LoaderCircle className="spin" size={22} /> Загружаем программы...</div>
      ) : pageError ? null : programs.length === 0 ? (
        <div className="program-empty-state">
          <div className="program-empty-icon"><FolderPlus size={23} /></div>
          <h3>Пока нет программ</h3>
          <p>Создай первую программу, чтобы подготовить план тренировок.</p>
          <button className="program-primary-button" type="button" onClick={() => setIsCreateOpen(true)}>
            <Plus size={17} /> Создать программу
          </button>
        </div>
      ) : (
        <div className="program-card-grid">
          {programs.map((program) => (
            <article className="program-card" key={program.id}>
              {editingId === program.id ? (
                <form className="program-edit-form" onSubmit={(event) => void handleRename(event, program.id)}>
                  <label htmlFor={`program-name-${program.id}`}>Новое название</label>
                  <input
                    id={`program-name-${program.id}`}
                    autoFocus
                    maxLength={100}
                    value={name}
                    onChange={(event) => setName(event.target.value)}
                  />
                  {formError && <p className="program-inline-error" role="alert">{formError}</p>}
                  <div className="program-card-actions">
                    <button className="program-primary-button" type="submit" disabled={isSaving}>
                      {isSaving ? 'Сохранение...' : 'Сохранить'}
                    </button>
                    <button className="program-secondary-button" type="button" onClick={() => { setEditingId(null); setFormError('') }}>
                      Отмена
                    </button>
                  </div>
                </form>
              ) : (
                <>
                  <button className="program-card-open" type="button" onClick={() => onOpenProgram(program.id)}>
                    <span className="program-card-icon"><Dumbbell size={20} /></span>
                    <span className="program-card-title">{program.name}</span>
                    <span className="program-card-date">Создана {formatDate(program.started_at)}</span>
                  </button>
                  <div className="program-card-actions">
                    <button
                      className="program-icon-action"
                      type="button"
                      aria-label={`Переименовать ${program.name}`}
                      title="Переименовать"
                      onClick={() => {
                        setEditingId(program.id)
                        setName(program.name)
                        setFormError('')
                      }}
                    >
                      <Pencil size={16} />
                    </button>
                    <button
                      className="program-icon-action program-delete-action"
                      type="button"
                      aria-label={`Удалить ${program.name}`}
                      title="Удалить"
                      onClick={() => void handleDelete(program)}
                    >
                      <Trash2 size={16} />
                    </button>
                  </div>
                </>
              )}
            </article>
          ))}
        </div>
      )}
    </section>
  )
}
