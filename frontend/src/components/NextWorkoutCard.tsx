import { ArrowUpRight, Clock3, Dumbbell, Play } from 'lucide-react'

type NextWorkoutCardProps = {
  onStart: () => void
}

export default function NextWorkoutCard({ onStart }: NextWorkoutCardProps) {
  return (
    <section className="next-workout-card">
      <div className="relative z-10 max-w-[390px]">
        <div className="flex items-center gap-2">
          <span className="status-dot" />
          <span className="workout-eyebrow text-[10px] font-bold uppercase tracking-[0.16em]">
            СЛЕДУЮЩАЯ ТРЕНИРОВКА
          </span>
        </div>
        <h2 className="mt-5 text-[28px] font-semibold tracking-[-0.055em] text-ink sm:text-[32px]">
          Ноги и плечи
        </h2>
        <p className="workout-description mt-2 max-w-[320px] text-sm leading-6 text-muted">
          Сильная база для уверенного прогресса. Сегодня — твой день.
        </p>

        <div className="workout-meta mt-5 flex flex-wrap gap-x-5 gap-y-2 text-xs font-medium">
          <span className="inline-flex items-center gap-2">
            <Clock3 size={15} strokeWidth={1.9} />
            около 55 мин
          </span>
          <span className="inline-flex items-center gap-2">
            <Dumbbell size={15} strokeWidth={1.9} />
            6 упражнений
          </span>
        </div>

        <button className="primary-button mt-6" type="button" onClick={onStart}>
          <Play size={15} fill="currentColor" />
          Начать тренировку
          <ArrowUpRight size={16} />
        </button>
      </div>
      <div className="workout-art" aria-hidden="true">
        <div className="art-ring art-ring-outer" />
        <div className="art-ring art-ring-inner" />
        <div className="art-center">
          <Dumbbell size={47} strokeWidth={1.25} />
        </div>
        <div className="art-caption">СИЛА · ТРЕНИРОВКА B</div>
      </div>
    </section>
  )
}
