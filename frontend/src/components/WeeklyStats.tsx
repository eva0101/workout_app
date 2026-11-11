import { Activity, Clock3, Dumbbell, LoaderCircle, RotateCcw } from 'lucide-react'

type WeeklyStatsProps = {
  summary: {
    workoutsCount: number
    totalDurationSeconds: number
    completedSets: number
  } | null
  isLoading: boolean
  error: string
  onRetry: () => void
}

function formatNumber(value: number) {
  return new Intl.NumberFormat('ru-RU').format(value)
}

function formatDuration(seconds: number) {
  const totalMinutes = Math.floor(seconds / 60)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  if (hours > 0) return `${formatNumber(hours)} ч ${formatNumber(minutes)} мин`
  return `${formatNumber(minutes)} мин`
}

export default function WeeklyStats({ summary, isLoading, error, onRetry }: WeeklyStatsProps) {
  return (
    <section aria-label="Краткая сводка за неделю" className="dashboard-week-summary">
      <div className="dashboard-week-summary-heading">
        <div className="dashboard-week-summary-title">
          <p className="eyebrow">ЭТА НЕДЕЛЯ</p>
          <h3>Краткая сводка</h3>
        </div>
        {isLoading && <LoaderCircle aria-label="Загружаем статистику" className="spin text-muted" size={18} />}
      </div>

      {error ? (
        <div className="dashboard-data-error" role="alert">
          <span>{error}</span>
          <button className="text-button" type="button" onClick={onRetry}>
            <RotateCcw size={13} /> Повторить
          </button>
        </div>
      ) : summary ? (
        <div className="dashboard-week-metrics">
          <article className="dashboard-week-metric">
            <div className="dashboard-week-metric-value">
              <span className="dashboard-week-metric-icon" aria-hidden="true"><Activity size={17} /></span>
              <strong>{formatNumber(summary.workoutsCount)}</strong>
            </div>
            <p>Завершённые тренировки</p>
          </article>
          <article className="dashboard-week-metric">
            <div className="dashboard-week-metric-value">
              <span className="dashboard-week-metric-icon" aria-hidden="true"><Clock3 size={17} /></span>
              <strong>{formatDuration(summary.totalDurationSeconds)}</strong>
            </div>
            <p>Время тренировок</p>
          </article>
          <article className="dashboard-week-metric">
            <div className="dashboard-week-metric-value">
              <span className="dashboard-week-metric-icon" aria-hidden="true"><Dumbbell size={17} /></span>
              <strong>{formatNumber(summary.completedSets)}</strong>
            </div>
            <p>Выполненные подходы</p>
          </article>
        </div>
      ) : (
        <p className="dashboard-data-loading">Загружаем сводку за неделю...</p>
      )}
    </section>
  )
}
