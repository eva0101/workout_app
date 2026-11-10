import { ArrowRight, Dumbbell, LoaderCircle, RotateCcw } from 'lucide-react'
import { useEffect, useState } from 'react'
import {
  getWorkoutHistory,
  type Workout,
  type WorkoutHistoryItem,
} from '../lib/api'

const HISTORY_PAGE_SIZE = 20

type WorkoutHistoryPageProps = {
  activeWorkout: Workout | null
  activeWorkoutError: string
  isActiveWorkoutLoading: boolean
  onRefreshActiveWorkout: () => void
  onResumeWorkout: () => void
  onOpenWorkout: (workout: WorkoutHistoryItem) => void
  onGoToOverview: () => void
}

function getErrorMessage(error: unknown) {
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте подключение.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось загрузить историю тренировок.'
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(value))
}

function formatDay(value: string) {
  return new Intl.DateTimeFormat('ru-RU', { day: '2-digit' }).format(new Date(value))
}

function formatMonth(value: string) {
  return new Intl.DateTimeFormat('ru-RU', { month: 'short' })
    .format(new Date(value))
    .replace('.', '')
}

function formatTime(value: string) {
  return new Intl.DateTimeFormat('ru-RU', { hour: '2-digit', minute: '2-digit' })
    .format(new Date(value))
}

function formatDuration(seconds: number) {
  const totalMinutes = Math.floor(seconds / 60)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  if (hours > 0) return `${hours} ч ${minutes} мин`
  if (minutes > 0) return `${minutes} мин`
  return `${seconds} сек`
}

function formatWeight(weight: number) {
  return `${new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 2 }).format(weight)} кг`
}

function getWellbeing(score: number | null | undefined) {
  if (score === null || score === undefined) return null
  if (score >= 1 && score <= 3) return { label: 'Плохое', tone: 'poor' }
  if (score <= 6 && score >= 4) return { label: 'Среднее', tone: 'average' }
  if (score >= 7 && score <= 8) return { label: 'Хорошее', tone: 'good' }
  if (score >= 9 && score <= 10) return { label: 'Отличное', tone: 'excellent' }
  return null
}

export default function WorkoutHistoryPage({
  activeWorkout,
  activeWorkoutError,
  isActiveWorkoutLoading,
  onRefreshActiveWorkout,
  onResumeWorkout,
  onOpenWorkout,
  onGoToOverview,
}: WorkoutHistoryPageProps) {
  const [history, setHistory] = useState<WorkoutHistoryItem[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isLoadingMore, setIsLoadingMore] = useState(false)
  const [hasMore, setHasMore] = useState(false)
  const [historyError, setHistoryError] = useState('')

  useEffect(() => {
    let cancelled = false
    setIsLoading(true)
    setHistoryError('')
    getWorkoutHistory(HISTORY_PAGE_SIZE)
      .then((items) => {
        if (cancelled) return
        setHistory(items)
        setHasMore(items.length === HISTORY_PAGE_SIZE)
      })
      .catch((error: unknown) => {
        if (!cancelled) setHistoryError(getErrorMessage(error))
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  async function loadMore() {
    if (isLoadingMore || !hasMore) return
    setIsLoadingMore(true)
    setHistoryError('')
    try {
      const nextPage = await getWorkoutHistory(HISTORY_PAGE_SIZE, history.length)
      setHistory((current) => [...current, ...nextPage])
      setHasMore(nextPage.length === HISTORY_PAGE_SIZE)
    } catch (error) {
      setHistoryError(getErrorMessage(error))
    } finally {
      setIsLoadingMore(false)
    }
  }

  function retryHistory() {
    if (history.length > 0) {
      void loadMore()
      return
    }

    setIsLoading(true)
    setHistoryError('')
    void getWorkoutHistory(HISTORY_PAGE_SIZE)
      .then((items) => {
        setHistory(items)
        setHasMore(items.length === HISTORY_PAGE_SIZE)
      })
      .catch((error: unknown) => setHistoryError(getErrorMessage(error)))
      .finally(() => setIsLoading(false))
  }

  const isCompletelyEmpty = !isLoading &&
    !isActiveWorkoutLoading &&
    !activeWorkout &&
    history.length === 0 &&
    !historyError

  return (
    <section className="dashboard-content workout-history-page">
      <div className="section-heading workout-history-heading">
        <div>
          <p className="eyebrow">ВАША АКТИВНОСТЬ</p>
          <h1>Тренировки</h1>
          <p className="workout-history-intro">Завершённые занятия и ваш прогресс</p>
        </div>
      </div>

      {activeWorkout && (
        <section className="workout-active-card panel-card" aria-label="Активная тренировка">
          <div className="workout-history-card-icon"><Dumbbell size={19} /></div>
          <div className="workout-active-copy">
            <p className="eyebrow">АКТИВНАЯ ТРЕНИРОВКА</p>
            <h2>{activeWorkout.training_day_name || `Тренировка #${activeWorkout.id}`}</h2>
            <p>
              Начата {formatDate(activeWorkout.started_at)}
              <span> · </span>
              {activeWorkout.exercises.length} упражнений
            </p>
          </div>
          <button className="program-primary-button workout-resume-button" type="button" onClick={onResumeWorkout}>
            Продолжить <ArrowRight size={16} />
          </button>
        </section>
      )}

      {activeWorkoutError && (
        <div className="program-alert mt-4" role="alert">
          {activeWorkoutError}
          <button className="text-button workout-inline-retry" type="button" onClick={onRefreshActiveWorkout}>
            <RotateCcw size={14} /> Повторить
          </button>
        </div>
      )}

      <div className="workout-history-list-heading">
        <div>
          <p className="eyebrow">АРХИВ</p>
          <h2>История завершённых тренировок</h2>
        </div>
      </div>

      {historyError && (
        <div className="program-alert mt-4" role="alert">
          {historyError}
          {!isLoadingMore && (
            <button className="text-button workout-inline-retry" type="button" onClick={retryHistory}>
              <RotateCcw size={14} /> Повторить
            </button>
          )}
        </div>
      )}

      {isLoading ? (
        <div className="workout-history-loading panel-card" role="status">
          <LoaderCircle className="spin" size={20} /> Загружаем историю тренировок...
        </div>
      ) : history.length > 0 ? (
        <>
          <ol className="workout-history-list">
            {history.map((item) => {
              const name = item.training_day_name || `Тренировка #${item.id}`
              const percent = item.completion_percent
              const wellbeing = getWellbeing(item.fatigue_score)
              const progressTone = percent === null
                ? 'neutral'
                : percent < 50
                  ? 'low'
                  : percent < 80
                    ? 'medium'
                    : 'high'
              return (
                <li key={item.id}>
                  <button
                    className={`workout-history-card panel-card workout-history-card-${progressTone}`}
                    type="button"
                    onClick={() => onOpenWorkout(item)}
                    aria-label={`Открыть тренировку ${name}, завершена ${formatDate(item.completed_at)}`}
                  >
                    <span className="workout-history-card-top">
                      <span className="workout-history-date-block" aria-label={formatDate(item.completed_at)}>
                        <span>{formatDay(item.completed_at)}</span>
                        <small>{formatMonth(item.completed_at)}</small>
                      </span>
                      <span className="workout-history-card-main">
                        <span className="workout-history-card-topline">
                          <span className="workout-history-card-title">{name}</span>
                          <span className="workout-history-completion" aria-label={percent === null ? 'Нет данных' : `Выполнение плана ${percent}%`}>
                            {percent === null ? '—' : `${percent}%`}
                          </span>
                        </span>
                        <span className="workout-history-card-subtitle">
                          Тренировка <span>·</span> завершена в{' '}
                          <time dateTime={item.completed_at}>{formatTime(item.completed_at)}</time>
                        </span>
                        {wellbeing && (
                          <span className="workout-history-wellbeing">
                            <span className="workout-history-wellbeing-label">Самочувствие</span>
                            <span className={`workout-history-wellbeing-badge workout-history-wellbeing-${wellbeing.tone}`}>
                              {wellbeing.label}
                            </span>
                          </span>
                        )}
                      </span>
                    </span>
                    <span
                      className="workout-history-progress"
                      role="progressbar"
                      aria-label="Выполнение плана"
                      aria-valuenow={percent ?? undefined}
                      aria-valuetext={percent === null ? 'Нет данных' : `${percent}%`}
                      aria-valuemin={0}
                      aria-valuemax={100}
                    >
                      <span style={{ width: `${percent ?? 0}%` }} />
                    </span>
                    <span className="workout-history-metrics">
                      <span><small>Длительность</small><strong>{formatDuration(item.duration_seconds)}</strong></span>
                      <span><small>Подходы</small><strong>{item.completed_sets}</strong></span>
                      <span><small>Повторения</small><strong>{item.actual_reps}</strong></span>
                      <span><small>Общий вес</small><strong>{formatWeight(item.total_weight)}</strong></span>
                    </span>
                  </button>
                </li>
              )
            })}
          </ol>
          {hasMore && (
            <div className="workout-history-more">
              <button
                className="program-secondary-button"
                type="button"
                disabled={isLoadingMore}
                onClick={() => void loadMore()}
              >
                {isLoadingMore && <LoaderCircle className="spin" size={16} />}
                {isLoadingMore ? 'Загружаем...' : 'Загрузить ещё'}
              </button>
            </div>
          )}
        </>
      ) : isCompletelyEmpty ? (
        <div className="workout-history-empty panel-card">
          <span className="session-empty-icon"><Dumbbell size={23} /></span>
          <h3>История пока пуста</h3>
          <p>Здесь появятся завершённые тренировки. Перейдите на обзор, чтобы запланировать занятие.</p>
          <button className="program-secondary-button mt-4" type="button" onClick={onGoToOverview}>
            Перейти на вкладку «Обзор»
          </button>
        </div>
      ) : !historyError ? (
        <div className="workout-history-empty panel-card">
          <p>Завершённых тренировок пока нет.</p>
        </div>
      ) : null}
    </section>
  )
}
