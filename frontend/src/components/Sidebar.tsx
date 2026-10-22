import {
  BarChart3,
  Dumbbell,
  LayoutDashboard,
  Settings2,
  Target,
} from 'lucide-react'

const navItems = [
  { label: 'Обзор', icon: LayoutDashboard },
  { label: 'Программы', icon: Target },
  { label: 'Тренировки', icon: Dumbbell },
  { label: 'Статистика', icon: BarChart3 },
]

type SidebarProps = {
  activeItem: string
  onNavigate: (label: string) => void
}

function Brand() {
  return (
    <div className="flex items-center gap-3 px-2">
      <div className="brand-mark" aria-hidden="true">
        <Dumbbell size={19} strokeWidth={2.4} />
      </div>
      <div className="leading-tight">
        <span className="block text-[15px] font-bold tracking-[-0.04em] text-ink">
          workout
        </span>
        <span className="block text-[10px] font-semibold uppercase tracking-[0.16em] text-muted">
          app
        </span>
      </div>
    </div>
  )
}

export default function Sidebar({ activeItem, onNavigate }: SidebarProps) {
  return (
    <>
      <aside className="sidebar hidden lg:flex">
        <Brand />

        <div className="mt-12 px-3 text-[10px] font-bold uppercase tracking-[0.16em] text-subtle">
          Меню
        </div>
        <nav className="mt-3 flex flex-col gap-1" aria-label="Основная навигация">
          {navItems.map(({ label, icon: Icon }) => (
            <button
              className={`nav-link ${activeItem === label ? 'nav-link-active' : ''}`}
              key={label}
              onClick={() => onNavigate(label)}
              type="button"
            >
              <Icon size={18} strokeWidth={1.8} />
              <span>{label}</span>
              {label === 'Обзор' && (
                <span className="nav-indicator" aria-hidden="true" />
              )}
            </button>
          ))}
        </nav>

        <div className="sidebar-bottom">
          <div className="plan-card">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-ink">Твоя цель</span>
              <Target size={16} className="text-accent" />
            </div>
            <p className="mt-2 text-xs leading-5 text-muted">
              4 тренировки в неделю — ты уже на верном пути.
            </p>
            <div className="sidebar-progress-track mt-3 h-1.5 overflow-hidden rounded-full">
              <div className="h-full w-3/4 rounded-full bg-accent" />
            </div>
            <p className="mt-2 text-[11px] font-medium text-muted">3 из 4 на этой неделе</p>
          </div>
          <button className="nav-link mt-3" type="button" onClick={() => onNavigate('Настройки')}>
            <Settings2 size={18} strokeWidth={1.8} />
            <span>Настройки</span>
          </button>
        </div>
      </aside>

      <nav className="mobile-nav lg:hidden" aria-label="Мобильная навигация">
        {navItems.map(({ label, icon: Icon }) => (
          <button
            aria-current={activeItem === label ? 'page' : undefined}
            className={`mobile-nav-link ${activeItem === label ? 'mobile-nav-active' : ''}`}
            key={label}
            onClick={() => onNavigate(label)}
            type="button"
          >
            <Icon size={19} strokeWidth={1.9} />
            <span>{label}</span>
          </button>
        ))}
      </nav>
    </>
  )
}
