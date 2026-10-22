import { ArrowUpRight, Clock3, Flame, Trophy } from 'lucide-react'

export default function WeeklyStats() {
  return (
    <section aria-label="Статистика за неделю" className="grid gap-4 sm:grid-cols-2">
      <article className="metric-card">
        <div className="flex items-start justify-between">
          <div className="metric-icon metric-icon-accent">
            <Flame size={18} strokeWidth={1.8} />
          </div>
          <span className="trend-pill">
            <ArrowUpRight size={12} />
            +1
          </span>
        </div>
        <div className="mt-5 flex items-end justify-between gap-3">
          <div>
            <p className="text-[11px] font-medium text-muted">Тренировок за неделю</p>
            <p className="metric-value mt-1 text-[26px] font-bold leading-none tracking-[-0.055em] text-ink">
              3<span className="ml-1 text-sm font-medium tracking-normal text-subtle">/ 4</span>
            </p>
          </div>
          <div className="flex h-[44px] items-end gap-1">
            {[17, 26, 13, 31, 19, 13, 13].map((height, index) => (
              <span
                className={`weekly-bar w-[5px] rounded-full ${index < 3 ? 'bg-accent' : ''}`}
                key={index}
                style={{ height }}
              />
            ))}
          </div>
        </div>
      </article>

      <article className="metric-card">
        <div className="flex items-start justify-between">
          <div className="metric-icon metric-icon-accent-soft">
            <Clock3 size={18} strokeWidth={1.8} />
          </div>
          <div className="metric-icon metric-icon-muted">
            <Trophy size={16} strokeWidth={1.8} />
          </div>
        </div>
        <div className="mt-5">
          <p className="text-[11px] font-medium text-muted">Общее время</p>
          <p className="metric-value mt-1 text-[26px] font-bold leading-none tracking-[-0.055em] text-ink">
            2 ч 46 мин
          </p>
          <p className="mt-2 text-[10px] text-muted">Суммарная длительность тренировок</p>
        </div>
      </article>
    </section>
  )
}
