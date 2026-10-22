import { useState } from 'react'
import DashboardHeader from './components/DashboardHeader'
import NextWorkoutCard from './components/NextWorkoutCard'
import RecentWorkouts from './components/RecentWorkouts'
import Sidebar from './components/Sidebar'
import WeeklyStats from './components/WeeklyStats'
import WorkoutDialog from './components/WorkoutDialog'
import { recentWorkouts } from './data'
import { getStoredToken, removeStoredToken } from './lib/api'
import AuthPage from './pages/AuthPage'

type Theme = 'light' | 'dark'

function getInitialTheme(): Theme {
  const savedTheme = localStorage.getItem('workout-app-theme')
  if (savedTheme === 'light' || savedTheme === 'dark') return savedTheme
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function getIsAuthenticated(): boolean {
  return Boolean(getStoredToken())
}

export default function App() {
  const [activeItem, setActiveItem] = useState('Обзор')
  const [dialogOpen, setDialogOpen] = useState(false)
  const [theme, setTheme] = useState<Theme>(getInitialTheme)
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(getIsAuthenticated)

  function toggleTheme() {
    const nextTheme = theme === 'light' ? 'dark' : 'light'
    localStorage.setItem('workout-app-theme', nextTheme)
    setTheme(nextTheme)
  }

  function handleLogout() {
    removeStoredToken()
    setIsAuthenticated(false)
  }

  if (!isAuthenticated) {
    return <AuthPage theme={theme} onToggleTheme={toggleTheme} onAuthSuccess={() => setIsAuthenticated(true)} />
  }

  return (
    <div className="app-shell" data-theme={theme}>
      <Sidebar activeItem={activeItem} onNavigate={setActiveItem} />
      <main className="main-content">
        <DashboardHeader
          theme={theme}
          onToggleTheme={toggleTheme}
          onLogout={handleLogout}
        />
        <div className="dashboard-content">
          <div className="section-heading">
            <div>
              <p className="eyebrow">ТВОЙ ПЛАН</p>
              <h2 className="mt-1 text-[17px] font-semibold tracking-[-0.035em] text-ink">
                Неделя в движении
              </h2>
            </div>
            <button className="week-select" type="button">
              Эта неделя <span aria-hidden="true">⌄</span>
            </button>
          </div>

          <div className="dashboard-grid">
            <div className="flex flex-col gap-4">
              <NextWorkoutCard onStart={() => setDialogOpen(true)} />
              <WeeklyStats />
            </div>
            <aside className="weekly-panel">
              <div className="flex items-start justify-between">
                <div>
                  <p className="text-[11px] font-medium text-muted">Твой ритм</p>
                  <h3 className="mt-1 text-sm font-semibold text-ink">План на неделю</h3>
                </div>
                <span className="streak-badge">75%</span>
              </div>
              <div className="mt-5">
                <div className="flex items-end justify-between">
                  <span className="text-[29px] font-semibold leading-none tracking-[-0.06em] text-ink">
                    3<span className="ml-1 text-sm font-medium tracking-normal text-subtle">/ 4</span>
                  </span>
                  <span className="pb-0.5 text-[10px] font-medium text-muted">цель выполнена</span>
                </div>
                <div className="progress-track mt-4 h-1.5 overflow-hidden rounded-full">
                  <div className="h-full w-3/4 rounded-full bg-accent" />
                </div>
              </div>
              <div className="mt-5 border-t border-line pt-4">
                <WorkoutProgress />
              </div>
              <div className="encouragement mt-4 rounded-xl px-3.5 py-3">
                <p className="text-xs font-semibold text-ink">Ещё одна — и цель твоя</p>
                <p className="mt-1 text-[11px] leading-4 text-muted">
                  Ты стабильно двигаешься вперёд. Так держать!
                </p>
              </div>
            </aside>
          </div>

          <div className="mt-6">
            <RecentWorkouts workouts={recentWorkouts} />
          </div>
          <p className="demo-note">Демонстрационные данные · история тренировок не сохраняется</p>
        </div>
      </main>
      <WorkoutDialog open={dialogOpen} onClose={() => setDialogOpen(false)} />
    </div>
  )
}

function WorkoutProgress() {
  const days = [
    { label: 'ПН', completed: true },
    { label: 'ВТ', completed: true },
    { label: 'СР', completed: false },
    { label: 'ЧТ', completed: true },
    { label: 'ПТ', completed: false },
    { label: 'СБ', completed: false },
    { label: 'ВС', completed: false },
  ]

  return (
    <div className="flex justify-between">
      {days.map(({ label, completed }) => (
        <div className="flex flex-col items-center gap-2" key={label}>
          <span className={`day-marker ${completed ? 'day-marker-done' : ''}`}>
            {completed ? '✓' : ''}
          </span>
          <span className="text-[9px] font-semibold text-subtle">{label}</span>
        </div>
      ))}
    </div>
  )
}
