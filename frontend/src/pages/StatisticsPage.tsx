import { Activity, BarChart3, Clock3, Dumbbell, LoaderCircle, RotateCcw, Weight } from 'lucide-react'
import { useEffect, useMemo, useRef, useState } from 'react'
import {
  getWorkoutStatisticsDashboard,
  type WorkoutStatisticsDashboard,
  type WorkoutStatisticsPeriod,
} from '../lib/api'

const periods: Array<{ value: WorkoutStatisticsPeriod; label: string }> = [
  { value: '7d', label: '7 дней' },
  { value: '30d', label: '30 дней' },
  { value: '3m', label: '3 месяца' },
  { value: '1y', label: 'Год' },
  { value: 'all', label: 'Всё время' },
]

const muscleGroups = [
  { name: 'chest', title: 'Грудь' },
  { name: 'back', title: 'Спина' },
  { name: 'legs', title: 'Ноги' },
  { name: 'biceps', title: 'Бицепс' },
  { name: 'triceps', title: 'Трицепс' },
  { name: 'shoulders', title: 'Плечи' },
  { name: 'abs', title: 'Пресс' },
]

type LoadMetric = 'weight' | 'sets' | 'reps' | 'time'

const loadMetrics: Array<{ value: LoadMetric; label: string }> = [
  { value: 'weight', label: 'Вес' },
  { value: 'sets', label: 'Подходы' },
  { value: 'reps', label: 'Повторения' },
  { value: 'time', label: 'Время' },
]

type ChartPoint = {
  date: string
  value: number
  tooltip: string
}

function safeNumber(value: unknown) {
  return typeof value === 'number' && Number.isFinite(value) && value > 0 ? value : 0
}

function formatNumber(value: number, maximumFractionDigits = 0) {
  return new Intl.NumberFormat('ru-RU', { maximumFractionDigits }).format(safeNumber(value))
}

function formatDuration(seconds: number) {
  const totalMinutes = Math.floor(safeNumber(seconds) / 60)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  if (hours > 0) return `${formatNumber(hours)} ч ${formatNumber(minutes)} мин`
  return `${formatNumber(minutes)} мин`
}

function formatDateLabel(value: string, interval: 'day' | 'week' | 'month') {
  const isMonthBucket = interval === 'month' && /^\d{4}-\d{2}$/.test(value)
  const dateValue = isMonthBucket ? `${value}-01` : value
  const date = new Date(`${dateValue}T00:00:00Z`)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('ru-RU', isMonthBucket
    ? { month: 'short', year: '2-digit', timeZone: 'UTC' }
    : interval === 'week' || interval === 'month'
      ? { day: 'numeric', month: 'short', timeZone: 'UTC' }
      : { day: 'numeric', timeZone: 'UTC' },
  ).format(date)
}

function getErrorMessage(error: unknown) {
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте подключение.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось загрузить статистику. Попробуйте ещё раз.'
}

function getLoadMetricValue(item: WorkoutStatisticsDashboard['load_dynamics']['data'][number], metric: LoadMetric) {
  switch (metric) {
    case 'weight':
      return safeNumber(item.total_weight)
    case 'sets':
      return safeNumber(item.completed_sets)
    case 'reps':
      return safeNumber(item.actual_reps)
    case 'time':
      return safeNumber(item.duration_seconds)
  }
}

function addInitialZeroPoint(
  points: ChartPoint[],
  interval: 'day' | 'week' | 'month',
  metric: LoadMetric,
): ChartPoint[] {
  const firstActualIndex = points.findIndex((point) => point.value > 0)
  if (firstActualIndex === -1) return points

  const firstActualPoint = points[firstActualIndex]
  const firstDate = interval === 'month' && /^\d{4}-\d{2}$/.test(firstActualPoint.date)
    ? `${firstActualPoint.date}-01`
    : firstActualPoint.date
  const previousDate = new Date(`${firstDate}T00:00:00.000Z`)
  if (Number.isNaN(previousDate.getTime())) return points
  previousDate.setUTCDate(previousDate.getUTCDate() - 1)

  const date = previousDate.toISOString().slice(0, 10)
  const zeroPoint: ChartPoint = {
    date,
    value: 0,
    tooltip: `${formatDateLabel(date, 'week')}: ${metric === 'time' ? formatDuration(0) : `0 ${getMetricUnit(metric, 0)}`}`,
  }
  return [
    ...points.slice(0, firstActualIndex),
    zeroPoint,
    ...points.slice(firstActualIndex),
  ]
}

function getMetricUnit(metric: LoadMetric, maximumValue: number) {
  if (metric === 'weight') return 'кг'
  if (metric === 'sets') return 'подходов'
  if (metric === 'reps') return 'повторений'
  return maximumValue >= 3600 ? 'ч' : 'мин'
}

function getDisplayValue(value: number, metric: LoadMetric | 'regularity', unit: string) {
  const displayValue = metric === 'time' || metric === 'regularity'
    ? value / (unit === 'ч' ? 3600 : 60)
    : value
  return formatNumber(displayValue, metric === 'weight' ? 1 : 0)
}

function StatisticsBarChart({
  points,
  interval,
  unit,
  metric,
  chartType = 'bar',
}: {
  points: ChartPoint[]
  interval: 'day' | 'week' | 'month'
  unit: string
  metric: LoadMetric | 'regularity'
  chartType?: 'bar' | 'line'
}) {
  const chartHeight = 220
  const chartContainerRef = useRef<HTMLDivElement>(null)
  const [chartWidth, setChartWidth] = useState(320)
  useEffect(() => {
    const container = chartContainerRef.current
    if (!container) return

    const updateWidth = () => setChartWidth(Math.max(1, container.clientWidth))
    updateWidth()
    const observer = new ResizeObserver(updateWidth)
    observer.observe(container)
    return () => observer.disconnect()
  }, [])

  const plotLeft = 38
  const plotRight = chartWidth - 12
  const plotTop = 15
  const plotBottom = 174
  const plotHeight = plotBottom - plotTop
  const plotWidth = plotRight - plotLeft
  const maximumValue = Math.max(0, ...points.map(({ value }) => safeNumber(value)))
  const scaleMaximum = maximumValue > 0 ? maximumValue * 1.1 : 1
  const barWidth = Math.min(24, (plotWidth / Math.max(points.length, 1)) * 0.58)
  const labelCapacity = Math.max(1, Math.floor(plotWidth / 34))
  const labelStep = Math.max(1, Math.ceil(points.length / labelCapacity))
  const axisValues = [0, scaleMaximum / 2, scaleMaximum]
  const pointCoordinates = points.map((point, index) => ({
    x: plotLeft + (index + 0.5) * (plotWidth / points.length),
    y: plotBottom - (safeNumber(point.value) / scaleMaximum) * plotHeight,
  }))
  const linePath = pointCoordinates
    .map(({ x, y }, index) => `${index === 0 ? 'M' : 'L'} ${x} ${y}`)
    .join(' ')

  if (points.length === 0) {
    return <div className="statistics-chart-empty">Нет данных за выбранный период</div>
  }

  return (
    <div className="statistics-chart-scroll" ref={chartContainerRef}>
      <svg
        aria-label={`${chartType === 'line' ? 'Линейный' : 'Столбчатый'} график, единицы измерения: ${unit}`}
        className="statistics-chart"
        height={chartHeight}
        role="img"
        style={{ width: '100%', height: chartHeight }}
        viewBox={`0 0 ${chartWidth} ${chartHeight}`}
        width="100%"
      >
        {axisValues.map((value, index) => {
          const y = plotBottom - (value / scaleMaximum) * plotHeight
          const axisLabel = metric === 'time' || metric === 'regularity'
            ? getDisplayValue(value, metric, unit)
            : formatNumber(value, metric === 'weight' ? 1 : 0)
          return (
            <g key={`${index}-${axisLabel}`}>
              <line className="statistics-chart-gridline" x1={plotLeft} x2={plotRight} y1={y} y2={y} />
              <text className="statistics-chart-axis-label" textAnchor="end" x={plotLeft - 7} y={y + 3}>
                {axisLabel}
              </text>
            </g>
          )
        })}
        <text className="statistics-chart-unit" textAnchor="end" x={plotRight} y={chartHeight - 4}>
          {unit}
        </text>
        {chartType === 'line' && (
          <path className="statistics-chart-line" d={linePath} />
        )}
        {points.map((point, index) => {
          const value = safeNumber(point.value)
          const { x, y } = pointCoordinates[index]
          const height = (value / scaleMaximum) * plotHeight
          const label = formatDateLabel(point.date, interval)
          return (
            <g key={`${point.date}-${index}`}>
              <title>{point.tooltip}</title>
              {chartType === 'line' ? (
                <circle className="statistics-chart-point" cx={x} cy={y} r={4} />
              ) : (
                <rect
                  className="statistics-chart-bar"
                  height={height}
                  rx={4}
                  width={barWidth}
                  x={x - barWidth / 2}
                  y={plotBottom - height}
                />
              )}
              {(index % labelStep === 0 || index === points.length - 1) && (
                <text className="statistics-chart-date" textAnchor="middle" x={x} y={plotBottom + 20}>
                  {label}
                </text>
              )}
            </g>
          )
        })}
      </svg>
    </div>
  )
}

function StatisticsLoading() {
  return (
    <div className="statistics-loading panel-card" role="status">
      <LoaderCircle className="spin" size={20} />
      Загружаем статистику...
    </div>
  )
}

export default function StatisticsPage() {
  const [period, setPeriod] = useState<WorkoutStatisticsPeriod>('7d')
  const [statistics, setStatistics] = useState<WorkoutStatisticsDashboard | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [retryNumber, setRetryNumber] = useState(0)
  const [loadMetric, setLoadMetric] = useState<LoadMetric>('weight')

  useEffect(() => {
    let cancelled = false
    setIsLoading(true)
    setError('')
    getWorkoutStatisticsDashboard(period)
      .then((response) => {
        if (!cancelled) setStatistics(response)
      })
      .catch((requestError: unknown) => {
        if (!cancelled) setError(getErrorMessage(requestError))
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [period, retryNumber])

  const muscleGroupValues = useMemo(() => {
    const groups = new Map(
      (Array.isArray(statistics?.muscle_groups) ? statistics.muscle_groups : [])
        .map((group) => [group.name, group]),
    )
    return muscleGroups.map((group) => {
      const apiGroup = groups.get(group.name)
      return {
        name: group.name,
        title: apiGroup?.title || group.title,
        completedSets: safeNumber(apiGroup?.completed_sets),
      }
    })
  }, [statistics])
  const maxMuscleSets = Math.max(0, ...muscleGroupValues.map(({ completedSets }) => completedSets))

  const regularityPoints = useMemo<ChartPoint[]>(() => {
    const regularity = statistics?.regularity
    const loadData = statistics?.load_dynamics?.data
    if (!regularity || !Array.isArray(regularity.data)) return []
    const durationByDate = new Map(
      (Array.isArray(loadData) ? loadData : [])
        .map((item) => [item.date, safeNumber(item.duration_seconds)]),
    )
    return regularity.data.map((item) => {
      const seconds = durationByDate.get(item.date) ?? 0
      return {
        date: item.date,
        value: seconds,
        tooltip: `${formatDateLabel(item.date, regularity.interval)}: ${formatDuration(seconds)}`,
      }
    })
  }, [statistics])
  const loadPoints = useMemo<ChartPoint[]>(() => {
    const dynamics = statistics?.load_dynamics
    if (!dynamics || !Array.isArray(dynamics.data)) return []
    const points = dynamics.data.map((item) => {
      const value = getLoadMetricValue(item, loadMetric)
      return {
        date: item.date,
        value,
        tooltip: `${formatDateLabel(item.date, dynamics.interval)}: ${loadMetric === 'time' ? formatDuration(value) : `${formatNumber(value, loadMetric === 'weight' ? 1 : 0)} ${getMetricUnit(loadMetric, value)}`}`,
      }
    })
    return period === 'all'
      ? addInitialZeroPoint(points, dynamics.interval, loadMetric)
      : points
  }, [loadMetric, period, statistics])

  const regularityMaxSeconds = Math.max(0, ...regularityPoints.map(({ value }) => value))
  const regularityUnit = regularityMaxSeconds >= 3600 ? 'ч' : 'мин'
  const maxLoadValue = Math.max(0, ...loadPoints.map(({ value }) => value))
  const loadUnit = getMetricUnit(loadMetric, maxLoadValue)
  const summary = statistics?.summary

  return (
    <section className="dashboard-content statistics-page">
      <div className="section-heading statistics-heading">
        <div>
          <p className="eyebrow">ВАШ ПРОГРЕСС</p>
          <h1>Статистика</h1>
          <p className="statistics-intro">Результаты завершённых тренировок</p>
        </div>
      </div>

      <div className="statistics-periods" aria-label="Период статистики">
        {periods.map((item) => (
          <button
            aria-pressed={period === item.value}
            className={`statistics-period-button ${period === item.value ? 'statistics-period-active' : ''}`}
            key={item.value}
            onClick={() => setPeriod(item.value)}
            type="button"
          >
            {item.label}
          </button>
        ))}
      </div>

      {error && (
        <div className="program-alert statistics-error" role="alert">
          <span>{error}</span>
          <button
            className="text-button statistics-retry"
            onClick={() => setRetryNumber((count) => count + 1)}
            type="button"
          >
            <RotateCcw size={14} /> Повторить
          </button>
        </div>
      )}

      {isLoading ? (
        <StatisticsLoading />
      ) : error ? null : (
        <>
          <section aria-label="Общая сводка" className="statistics-summary-grid">
            <article className="statistics-summary-card panel-card">
              <span className="statistics-summary-icon"><Dumbbell size={18} /></span>
              <p>Тренировки</p>
              <strong>{formatNumber(safeNumber(summary?.workouts_count))}</strong>
            </article>
            <article className="statistics-summary-card panel-card">
              <span className="statistics-summary-icon"><Clock3 size={18} /></span>
              <p>Общее время</p>
              <strong>{formatDuration(safeNumber(summary?.total_duration_seconds))}</strong>
            </article>
            <article className="statistics-summary-card panel-card">
              <span className="statistics-summary-icon"><Activity size={18} /></span>
              <p>Выполненные подходы</p>
              <strong>{formatNumber(safeNumber(summary?.completed_sets))}</strong>
            </article>
            <article className="statistics-summary-card panel-card">
              <span className="statistics-summary-icon"><Weight size={18} /></span>
              <p>Общий поднятый вес</p>
              <strong>{formatNumber(safeNumber(summary?.total_weight), 2)} кг</strong>
            </article>
          </section>

          <div className="statistics-charts-grid">
            <section className="statistics-panel panel-card" aria-labelledby="statistics-regularity-title">
              <div className="statistics-panel-heading">
                <div>
                  <p className="eyebrow">АКТИВНОСТЬ</p>
                  <h2 id="statistics-regularity-title">Регулярность тренировок</h2>
                </div>
                <span className="statistics-chart-heading-unit">
                  {regularityUnit === 'ч' ? 'Время, ч' : 'Время, мин'}
                </span>
              </div>
              <div aria-hidden="true" className="statistics-chart-controls statistics-chart-controls-placeholder" />
              <StatisticsBarChart
                interval={statistics?.regularity.interval ?? 'day'}
                metric="regularity"
                points={regularityPoints}
                unit={regularityUnit}
              />
            </section>

            <section className="statistics-panel panel-card" aria-labelledby="statistics-load-title">
              <div className="statistics-panel-heading statistics-load-heading">
                <div>
                  <p className="eyebrow">СИЛА И ОБЪЁМ</p>
                  <h2 id="statistics-load-title">Динамика нагрузки</h2>
                </div>
                <BarChart3 className="statistics-load-icon" size={18} />
              </div>
              <div className="statistics-chart-controls">
                <div className="statistics-metric-switch" aria-label="Метрика графика нагрузки">
                  {loadMetrics.map((metric) => (
                    <button
                      aria-pressed={loadMetric === metric.value}
                      className={`statistics-metric-button ${loadMetric === metric.value ? 'statistics-metric-active' : ''}`}
                      key={metric.value}
                      onClick={() => setLoadMetric(metric.value)}
                      type="button"
                    >
                      {metric.label}
                    </button>
                  ))}
                </div>
              </div>
              <StatisticsBarChart
                chartType="line"
                interval={statistics?.load_dynamics.interval ?? 'day'}
                metric={loadMetric}
                points={loadPoints}
                unit={loadUnit}
              />
            </section>
          </div>

          <section className="statistics-panel statistics-muscle-panel panel-card" aria-labelledby="statistics-muscle-title">
            <div className="statistics-panel-heading">
              <div>
                <p className="eyebrow">РАСПРЕДЕЛЕНИЕ НАГРУЗКИ</p>
                <h2 id="statistics-muscle-title">Мышечные группы</h2>
              </div>
              <span className="statistics-chart-heading-unit">Выполненные подходы</span>
            </div>
            <div className="statistics-muscle-list">
              {muscleGroupValues.map((group) => {
                const percentage = maxMuscleSets > 0
                  ? Math.min(100, (group.completedSets / maxMuscleSets) * 100)
                  : 0
                return (
                  <div className="statistics-muscle-row" key={group.name}>
                    <span className="statistics-muscle-title">{group.title}</span>
                    <div
                      aria-label={`${group.completedSets} подходов`}
                      className="statistics-muscle-track"
                      role="img"
                    >
                      <span className="statistics-muscle-fill" style={{ width: `${percentage}%` }} />
                    </div>
                    <span className="statistics-muscle-count">{formatNumber(group.completedSets)}</span>
                  </div>
                )
              })}
            </div>
          </section>
        </>
      )}
    </section>
  )
}
