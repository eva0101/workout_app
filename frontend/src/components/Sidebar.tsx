import {
  BarChart3,
  Dumbbell,
  LayoutDashboard,
  Target,
} from 'lucide-react'

const navItems = [
  { id: 'Обзор', label: 'Главная', icon: LayoutDashboard },
  { id: 'Программы', label: 'Программы', icon: Target },
  { id: 'Тренировки', label: 'История', icon: Dumbbell },
  { id: 'Статистика', label: 'Статистика', icon: BarChart3 },
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
          {navItems.map(({ id, label, icon: Icon }) => (
            <button
              className={`nav-link ${activeItem === id ? 'nav-link-active' : ''}`}
              key={id}
              onClick={() => onNavigate(id)}
              type="button"
            >
              <Icon size={18} strokeWidth={1.8} />
              <span>{label}</span>
            </button>
          ))}
        </nav>
      </aside>

      <nav className="mobile-nav lg:hidden" aria-label="Мобильная навигация">
        {navItems.map(({ id, label, icon: Icon }) => (
          <button
            aria-current={activeItem === id ? 'page' : undefined}
            className={`mobile-nav-link ${activeItem === id ? 'mobile-nav-active' : ''}`}
            key={id}
            onClick={() => onNavigate(id)}
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
