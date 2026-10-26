import { Bell, ChevronDown, Moon, Sun } from 'lucide-react'

type DashboardHeaderProps = {
  theme: 'light' | 'dark'
  onToggleTheme: () => void
  onLogout: () => void
}

export default function DashboardHeader({ theme, onToggleTheme, onLogout }: DashboardHeaderProps) {
  return (
    <header className="dashboard-header">
      <div>
        <p className="eyebrow">ВТОРНИК, 20 ОКТЯБРЯ</p>
        <h1 className="mt-1 text-[22px] font-semibold tracking-[-0.04em] text-ink sm:text-[25px]">
          Доброе утро, Алексей <span aria-hidden="true">✦</span>
        </h1>
        <p className="mt-1 text-sm text-muted">
          Отличный день, чтобы стать немного сильнее.
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
        <button className="icon-button relative hidden sm:grid" type="button" aria-label="Уведомления">
          <Bell size={18} strokeWidth={1.8} />
          <span className="notification-dot" />
        </button>
        <button className="profile-button" type="button" aria-label="Профиль Алексея">
          <span className="avatar">А</span>
          <span className="hidden text-left sm:block">
            <span className="block text-xs font-semibold text-ink">Алексей</span>
            <span className="mt-0.5 block text-[11px] text-muted">Мой профиль</span>
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
