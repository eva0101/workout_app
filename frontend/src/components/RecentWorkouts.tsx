import { ArrowRight, Check, Clock3, Dumbbell, MoreHorizontal } from 'lucide-react'
import type { Workout } from '../data'

type RecentWorkoutsProps = {
  workouts: Workout[]
}

const toneClasses = {
  sage: 'workout-icon-sage',
  sand: 'workout-icon-sand',
  blue: 'workout-icon-blue',
}

export default function RecentWorkouts({ workouts }: RecentWorkoutsProps) {
  return (
    <section className="panel-card">
      <div className="flex items-center justify-between px-5 pb-4 pt-5 sm:px-6">
        <div>
          <h2 className="text-[15px] font-semibold tracking-[-0.025em] text-ink">
            Последние тренировки
          </h2>
          <p className="mt-1 text-xs text-muted">Твоя работа складывается в результат</p>
        </div>
        <button className="text-button" type="button">
          Все тренировки <ArrowRight size={14} />
        </button>
      </div>

      {workouts.length > 0 ? (
        <div className="divide-y divide-line">
          {workouts.map((workout) => (
            <article
              className="workout-row"
              key={workout.id}
            >
              <div className={`workout-icon ${toneClasses[workout.tone]}`}>
                <Dumbbell size={17} strokeWidth={1.8} />
              </div>
              <div className="min-w-0 flex-1">
                <h3 className="truncate text-[13px] font-semibold text-ink">{workout.name}</h3>
                <p className="mt-1 text-[11px] text-muted">{workout.date}</p>
              </div>
              <div className="hidden items-center gap-1.5 text-[11px] text-muted sm:flex">
                <Dumbbell size={13} />
                {workout.exercises} упр.
              </div>
              <div className="flex items-center gap-1.5 text-[11px] text-muted">
                <Clock3 size={13} />
                {workout.duration}
              </div>
              <span className="completion-check" title="Тренировка завершена">
                <Check size={12} strokeWidth={2.5} />
              </span>
              <button className="row-more" type="button" aria-label={`Действия: ${workout.name}`}>
                <MoreHorizontal size={18} />
              </button>
            </article>
          ))}
        </div>
      ) : (
        <div className="empty-state">
          <div className="empty-icon">
            <Dumbbell size={21} strokeWidth={1.6} />
          </div>
          <h3 className="mt-3 text-sm font-semibold text-ink">Здесь появится твой прогресс</h3>
          <p className="mt-1 max-w-[250px] text-xs leading-5 text-muted">
            Заверши первую тренировку — и она сохранится в твоей истории.
          </p>
        </div>
      )}
    </section>
  )
}
