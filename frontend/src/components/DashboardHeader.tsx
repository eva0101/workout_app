import { ChevronDown, Moon, Sun } from 'lucide-react'

type DashboardHeaderProps = {
  theme: 'light' | 'dark'
  onToggleTheme: () => void
  onLogout: () => void
}

export default function DashboardHeader({ theme, onToggleTheme, onLogout }: DashboardHeaderProps) {
  const dateLabel = new Intl.DateTimeFormat('ru-RU', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
  }).format(new Date())

  return (
    <header className="dashboard-header">
      <div>
        <p className="eyebrow">{dateLabel}</p>
        <h1 className="mt-1 text-[22px] font-semibold tracking-[-0.04em] text-ink sm:text-[25px]">
          Добро пожаловать <span aria-hidden="true">✦</span>
        </h1>
        <p className="mt-1 text-sm text-muted">
          Планируйте занятия и двигайтесь в своём ритме.
        </p>
      </div>

      <div className="flex items-center gap-3">
        <button
          className="icon-button theme-toggle"
          type="button"
          onClick={onToggleTheme}
          aria-label={theme === 'light' ? 'Включить тёмную тему' : 'Включить светлую тему'}
          title={theme === 'light' ? 'Тёмная тема' : 'Светлая тема'}
        >
          {theme === 'light' ? <Moon size={18} strokeWidth={1.9} /> : <Sun size={18} strokeWidth={1.9} />}
        </button>
        <button className="profile-button" type="button" aria-label="Аккаунт">
          <span className="avatar">W</span>
          <span className="hidden text-left sm:block">
            <span className="block text-xs font-semibold text-ink">Аккаунт</span>
            <span className="mt-0.5 block text-[11px] text-muted">Workoutapp</span>
          </span>
          <ChevronDown size={15} className="hidden text-muted sm:block" />
        </button>
        <button className="text-button logout-button" type="button" onClick={onLogout}>
          Выйти
        </button>
      </div>
    </header>
  )
}
