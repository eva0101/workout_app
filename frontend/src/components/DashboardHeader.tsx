import { Moon, Sun } from 'lucide-react'
import { useEffect, useState } from 'react'

type DashboardHeaderProps = {
  theme: 'light' | 'dark'
  login: string | null
  isHomePage: boolean
  welcomeTitle: string
  welcomeDescription: string
  onToggleTheme: () => void
  onLogout: () => void
}

export default function DashboardHeader({
  theme,
  login,
  isHomePage,
  welcomeTitle,
  welcomeDescription,
  onToggleTheme,
  onLogout,
}: DashboardHeaderProps) {
  const [hasShownWelcome, setHasShownWelcome] = useState(isHomePage)

  useEffect(() => {
    if (isHomePage) setHasShownWelcome(true)
  }, [isHomePage])

  const dateLabel = new Intl.DateTimeFormat('ru-RU', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
  }).format(new Date())
  const displayedLogin = login?.trim() || 'Пользователь'
  const avatarLetter = Array.from(login?.trim() ?? '')[0]?.toLocaleUpperCase('ru-RU') || 'П'

  return (
    <header className="dashboard-header">
      <div className={hasShownWelcome ? 'welcome-content-enter' : undefined}>
        <p className="eyebrow">{dateLabel}</p>
        <h1 className="mt-1 text-[22px] font-semibold tracking-[-0.04em] text-ink sm:text-[25px]">
          {welcomeTitle} <span aria-hidden="true">✦</span>
        </h1>
        <p className="mt-1 text-sm text-muted">
          {welcomeDescription}
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
        <div className="profile-button">
          <span className="avatar">{avatarLetter}</span>
          <span className="hidden text-left sm:block">
            <span className="block text-xs font-semibold text-ink">Аккаунт</span>
            <span className="mt-0.5 block text-[11px] text-muted">{displayedLogin}</span>
          </span>
        </div>
        <button className="text-button logout-button" type="button" onClick={onLogout}>
          Выйти
        </button>
      </div>
    </header>
  )
}
