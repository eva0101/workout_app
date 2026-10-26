import { ArrowLeft, CalendarDays, Dumbbell, FolderPlus, LoaderCircle, Pencil, Plus, RefreshCw, Trash2, X } from 'lucide-react'
import { type FormEvent, useEffect, useState } from 'react'
import ScheduleDialog from '../components/ScheduleDialog'
import {
  createProgram,
  createSchedule,
  createTrainingDay,
  deleteProgram,
  deleteScheduledWorkout,
  deleteTrainingDay,
  getProgram,
  getPrograms,
  getSchedule,
  getTrainingDay,
  type Program,
  type ProgramDetails,
  type ScheduledWorkout,
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
  if (status === 400) return 'Запрос содержит некорректные данные. Проверьте выбранные даты.'
  if (status === 404) return 'Запись или тренировочный день не найдены. Обновите страницу и попробуйте снова.'
  if (status === 409) return 'Одна или несколько дат уже запланированы для этого тренировочного дня.'
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте, запущен ли backend.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось выполнить запрос. Попробуйте ещё раз.'
}

function formatDate(value: string): string {
  const match = /^(\d{4})-(\d{2})-(\d{2})(?:$|T)/.exec(value)
  if (!match || (value.includes('T') && Number.isNaN(Date.parse(value)))) {
    return 'Дата неизвестна'
  }

  const [, yearText, monthText, dayText] = match
  const year = Number(yearText)
  const month = Number(monthText)
  const day = Number(dayText)
  const date = new Date(year, month - 1, day)
  if (
    Number.isNaN(date.getTime()) ||
    date.getFullYear() !== year ||
    date.getMonth() !== month - 1 ||
    date.getDate() !== day
  ) {
    return 'Дата неизвестна'
  }
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
  const [dayFormError, setDayFormError] = useState('')
  const [dayName, setDayName] = useState('')
  const [isDaySaving, setIsDaySaving] = useState(false)
  const [isDayDeleting, setIsDayDeleting] = useState(false)
  const [isCreateDayOpen, setIsCreateDayOpen] = useState(false)
  const [scheduledWorkouts, setScheduledWorkouts] = useState<ScheduledWorkout[]>([])
  const [isScheduleLoading, setIsScheduleLoading] = useState(false)
  const [scheduleError, setScheduleError] = useState('')
  const [isScheduleDialogOpen, setIsScheduleDialogOpen] = useState(false)
  const [isScheduleSaving, setIsScheduleSaving] = useState(false)
  const [scheduleSaveError, setScheduleSaveError] = useState('')
  const [scheduleDeleteError, setScheduleDeleteError] = useState('')
  const [deletingScheduleId, setDeletingScheduleId] = useState<number | null>(null)

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

  useEffect(() => {
    if (selectedDayId === null) {
      setScheduledWorkouts([])
      setScheduleError('')
      return
    }

    let cancelled = false
    setIsScheduleLoading(true)
    setScheduleError('')
    getSchedule('0001-01-01', '9999-12-31')
      .then((items) => {
        if (!cancelled) {
          setScheduledWorkouts(items.filter((item) => item.training_day_id === selectedDayId))
        }
      })
      .catch((error: unknown) => {
        if (!cancelled) setScheduleError(getErrorMessage(error))
      })
      .finally(() => {
        if (!cancelled) setIsScheduleLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [selectedDayId])

  async function refreshSchedule() {
    if (selectedDayId === null) return
    setIsScheduleLoading(true)
    setScheduleError('')
    try {
      const items = await getSchedule('0001-01-01', '9999-12-31')
      setScheduledWorkouts(items.filter((item) => item.training_day_id === selectedDayId))
    } catch (error) {
      setScheduleError(getErrorMessage(error))
    } finally {
      setIsScheduleLoading(false)
    }
  }

  async function handleSaveSchedule(dates: string[]) {
    if (selectedDayId === null) return
    setIsScheduleSaving(true)
    setScheduleSaveError('')
    try {
      await createSchedule(selectedDayId, dates)
      setIsScheduleDialogOpen(false)
      await refreshSchedule()
    } catch (error) {
      setScheduleSaveError(getErrorMessage(error))
    } finally {
      setIsScheduleSaving(false)
    }
  }

  async function handleDeleteScheduledWorkout(item: ScheduledWorkout) {
    if (!window.confirm(`Удалить тренировку из расписания на ${formatDate(item.scheduled_date)}?`)) return
    setDeletingScheduleId(item.id)
    setScheduleDeleteError('')
    try {
      await deleteScheduledWorkout(item.id)
      setScheduledWorkouts((current) => current.filter(({ id }) => id !== item.id))
    } catch (error) {
      setScheduleDeleteError(getErrorMessage(error))
    } finally {
      setDeletingScheduleId(null)
    }
  }

  async function handleCreateTrainingDay(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!details) return
    const trimmedName = dayName.trim()
    const nameLength = Array.from(trimmedName).length
    if (nameLength < 3 || nameLength > 100) {
      setDayFormError('Название должно содержать от 3 до 100 символов.')
      return
    }

    setIsDaySaving(true)
    setDayActionError('')
    setDayFormError('')
    try {
      const createdDay = await createTrainingDay(details.id, trimmedName)
      setDetails((current) => current
        ? { ...current, training_days: [...(current.training_days ?? []), createdDay] }
        : current,
      )
      setDayName('')
      setIsCreateDayOpen(false)
    } catch (error) {
      const message = getErrorMessage(error)
      setDayActionError(message)
      setDayFormError(message)
    } finally {
      setIsDaySaving(false)
    }
  }

  async function handleDeleteTrainingDay(day: TrainingDaySummary) {
    if (!details) return
    if (!window.confirm(`Удалить «${day.name}»? Это действие нельзя отменить.`)) return

    setIsDayDeleting(true)
    setDayActionError('')
    try {
      await deleteTrainingDay(day.id)
      setDetails((current) => current
        ? {
            ...current,
            training_days: (current.training_days ?? []).filter(({ id }) => id !== day.id),
          }
        : current,
      )
    } catch (error) {
      const status = typeof error === 'object' && error !== null && 'status' in error
        ? (error as { status?: number }).status
        : undefined
      setDayActionError(
        status === 409
          ? 'Нельзя удалить этот день: с ним связана история тренировок.'
          : getErrorMessage(error),
      )
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
                  <p className="eyebrow">{details?.name} · ДЕНЬ {dayDetails.day_number}</p>
                  <h2>{dayDetails.name}</h2>
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
                {dayDetails.exercises.length === 0 ? (
                  <div className="program-empty-state training-day-empty">
                    <div className="program-empty-icon"><Dumbbell size={22} /></div>
                    <h3>В этом тренировочном дне пока нет упражнений</h3>
                    <p>Добавление упражнений будет доступно на следующем этапе.</p>
                  </div>
                ) : (
                  <ol className="training-exercise-list">
                    {dayDetails.exercises.map((exercise, index) => (
                      <li className="training-exercise-card" key={exercise.id}>
                        <span className="exercise-number">{String(index + 1).padStart(2, '0')}</span>
                        <span className="training-exercise-name">Упражнение #{exercise.exercise_id}</span>
                        <span className="training-exercise-sets">{exercise.sets} подходов</span>
                      </li>
                    ))}
                  </ol>
                )}
              </div>
              <section className="scheduled-workouts-section" aria-labelledby="scheduled-workouts-title">
                <div className="scheduled-workouts-heading">
                  <div>
                    <p className="eyebrow">КАЛЕНДАРЬ</p>
                    <h3 id="scheduled-workouts-title">Запланированные даты</h3>
                  </div>
                  <button
                    className="program-primary-button"
                    type="button"
                    disabled={isScheduleLoading || isScheduleSaving}
                    onClick={() => {
                      setScheduleSaveError('')
                      setIsScheduleDialogOpen(true)
                    }}
                  >
                    <CalendarDays size={16} /> Запланировать
                  </button>
                </div>

                {scheduleSaveError && <div className="program-alert" role="alert">{scheduleSaveError}</div>}
                {scheduleDeleteError && <div className="program-alert" role="alert">{scheduleDeleteError}</div>}
                {scheduleError ? (
                  <div className="program-alert" role="alert">
                    <span>{scheduleError}</span>
                    <button className="program-secondary-button" type="button" onClick={() => void refreshSchedule()}>
                      <RefreshCw size={14} /> Повторить
                    </button>
                  </div>
                ) : isScheduleLoading ? (
                  <div className="program-state schedule-state"><LoaderCircle className="spin" size={19} /> Загружаем расписание...</div>
                ) : scheduledWorkouts.length === 0 ? (
                  <div className="scheduled-workouts-empty">
                    <CalendarDays size={19} />
                    <span>Для этого дня пока нет запланированных дат.</span>
                  </div>
                ) : (
                  <ul className="scheduled-workouts-list">
                    {[...scheduledWorkouts]
                      .sort((left, right) => left.scheduled_date.localeCompare(right.scheduled_date))
                      .map((item) => (
                        <li className="scheduled-workout-item" key={item.id}>
                          <span className="scheduled-workout-date">
                            <CalendarDays size={17} />
                            {formatDate(item.scheduled_date)}
                          </span>
                          <button
                            className="program-icon-action program-delete-action"
                            type="button"
                            aria-label={`Удалить дату ${formatDate(item.scheduled_date)}`}
                            title="Удалить из расписания"
                            disabled={deletingScheduleId === item.id}
                            onClick={() => void handleDeleteScheduledWorkout(item)}
                          >
                            {deletingScheduleId === item.id
                              ? <LoaderCircle className="spin" size={16} />
                              : <Trash2 size={16} />}
                          </button>
                        </li>
                      ))}
                  </ul>
                )}
              </section>
            </>
          ) : (
            <div className="program-state">Не удалось загрузить тренировочный день.</div>
          )}
          <ScheduleDialog
            open={isScheduleDialogOpen}
            bookedDates={scheduledWorkouts.map((item) => item.scheduled_date)}
            isSaving={isScheduleSaving}
            error={scheduleSaveError}
            onClose={() => {
              setIsScheduleDialogOpen(false)
              setScheduleSaveError('')
            }}
            onSave={(dates) => void handleSaveSchedule(dates)}
          />
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
                onClick={() => {
                  setDayName('')
                  setDayFormError('')
                  setDayActionError('')
                  setIsCreateDayOpen(true)
                }}
                disabled={isDaySaving || isDayDeleting}
              >
                <Plus size={17} /> Добавить тренировочный день
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
                  <article className="training-day-card" key={day.id}>
                    <button
                      className="training-day-card-open"
                      type="button"
                      onClick={() => onOpenDay(day.id)}
                    >
                      <span className="program-card-icon"><Dumbbell size={20} /></span>
                      <span className="training-day-card-label">{day.name}</span>
                      <span className="training-day-card-hint">Открыть упражнения</span>
                    </button>
                    <button
                      className="program-icon-action program-delete-action training-day-delete"
                      type="button"
                      aria-label={`Удалить ${day.name}`}
                      title={isDayDeleting ? 'Удаление...' : 'Удалить тренировочный день'}
                      onClick={() => void handleDeleteTrainingDay(day)}
                      disabled={isDayDeleting || isDaySaving}
                    >
                      <Trash2 size={16} />
                    </button>
                  </article>
                ))}
              </div>
            )}
          </>
        ) : null}
        {isCreateDayOpen && (
          <div
            className="program-modal-backdrop"
            onClick={(event) => {
              if (event.target === event.currentTarget && !isDaySaving) setIsCreateDayOpen(false)
            }}
            role="presentation"
          >
            <section
              aria-labelledby="create-training-day-title"
              aria-modal="true"
              className="program-create-dialog training-day-create-dialog"
              role="dialog"
            >
              <div className="program-create-heading">
                <div>
                  <p className="eyebrow">НОВЫЙ ДЕНЬ ПРОГРАММЫ</p>
                  <h2 id="create-training-day-title">Новый тренировочный день</h2>
                </div>
                <button
                  className="program-icon-action"
                  type="button"
                  aria-label="Закрыть"
                  disabled={isDaySaving}
                  onClick={() => setIsCreateDayOpen(false)}
                >
                  <X size={18} />
                </button>
              </div>
              <form className="program-create-form" onSubmit={(event) => void handleCreateTrainingDay(event)}>
                <label htmlFor="new-training-day-name">Название тренировочного дня</label>
                <input
                  id="new-training-day-name"
                  autoFocus
                  value={dayName}
                  onChange={(event) => setDayName(event.target.value)}
                  placeholder="Например, День спины"
                  disabled={isDaySaving}
                />
                <p className="program-form-hint">Название должно содержать от 3 до 100 символов.</p>
                {dayFormError && <p className="program-inline-error" role="alert">{dayFormError}</p>}
                <div className="program-create-actions">
                  <button
                    className="program-secondary-button"
                    type="button"
                    disabled={isDaySaving}
                    onClick={() => setIsCreateDayOpen(false)}
                  >
                    Отмена
                  </button>
                  <button className="program-primary-button" type="submit" disabled={isDaySaving}>
                    {isDaySaving ? 'Создание...' : 'Создать'}
                  </button>
                </div>
              </form>
            </section>
          </div>
        )}
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
