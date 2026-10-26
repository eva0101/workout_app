import { useEffect, useState } from 'react'
import DashboardHeader from './components/DashboardHeader'
import NextWorkoutCard from './components/NextWorkoutCard'
import RecentWorkouts from './components/RecentWorkouts'
import Sidebar from './components/Sidebar'
import WeeklyStats from './components/WeeklyStats'
import WorkoutDialog from './components/WorkoutDialog'
import { recentWorkouts } from './data'
import { getStoredToken, removeStoredToken, UNAUTHORIZED_EVENT } from './lib/api'
import AuthPage from './pages/AuthPage'
import ProgramsPage from './pages/ProgramsPage'

type Theme = 'light' | 'dark'

function getInitialTheme(): Theme {
  const savedTheme = localStorage.getItem('workout-app-theme')
  if (savedTheme === 'light' || savedTheme === 'dark') return savedTheme
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function getIsAuthenticated(): boolean {
  return Boolean(getStoredToken())
}

type AppLocation = {
  activeItem: string
  programId: number | null
  dayId: number | null
}

function readAppLocation(): AppLocation {
  const params = new URLSearchParams(window.location.search)
  const programId = Number(params.get('program'))
  const dayId = Number(params.get('day'))
  const isProgramRoute = params.get('section') === 'programs' && Number.isInteger(programId) && programId > 0

  return {
    activeItem: params.get('section') === 'programs' ? 'Программы' : 'Обзор',
    programId: isProgramRoute ? programId : null,
    dayId: isProgramRoute && Number.isInteger(dayId) && dayId > 0 ? dayId : null,
  }
}

function writeAppLocation(programId: number | null, dayId: number | null, keepProgramsSection = false) {
  const url = new URL(window.location.href)
  url.search = ''

  if (programId !== null) {
    url.searchParams.set('section', 'programs')
    url.searchParams.set('program', String(programId))
    if (dayId !== null) url.searchParams.set('day', String(dayId))
  } else if (keepProgramsSection) {
    url.searchParams.set('section', 'programs')
  }

  window.history.pushState(null, '', url)
}

export default function App() {
  const [initialLocation] = useState(readAppLocation)
  const [activeItem, setActiveItem] = useState(initialLocation.activeItem)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [theme, setTheme] = useState<Theme>(getInitialTheme)
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(getIsAuthenticated)
  const [selectedProgramId, setSelectedProgramId] = useState<number | null>(initialLocation.programId)
  const [selectedDayId, setSelectedDayId] = useState<number | null>(initialLocation.dayId)

  useEffect(() => {
    function handleUnauthorized() {
      setIsAuthenticated(false)
      setActiveItem('Обзор')
      setSelectedProgramId(null)
      setSelectedDayId(null)
    }

    window.addEventListener(UNAUTHORIZED_EVENT, handleUnauthorized)
    return () => window.removeEventListener(UNAUTHORIZED_EVENT, handleUnauthorized)
  }, [])

  useEffect(() => {
    function handlePopState() {
      const location = readAppLocation()
      setActiveItem(location.activeItem)
      setSelectedProgramId(location.programId)
      setSelectedDayId(location.dayId)
    }

    window.addEventListener('popstate', handlePopState)
    return () => window.removeEventListener('popstate', handlePopState)
  }, [])

  function toggleTheme() {
    const nextTheme = theme === 'light' ? 'dark' : 'light'
    localStorage.setItem('workout-app-theme', nextTheme)
    setTheme(nextTheme)
  }

  function handleLogout() {
    removeStoredToken()
    setIsAuthenticated(false)
    setActiveItem('Обзор')
    setSelectedProgramId(null)
    setSelectedDayId(null)
    writeAppLocation(null, null)
  }

  if (!isAuthenticated) {
    return <AuthPage theme={theme} onToggleTheme={toggleTheme} onAuthSuccess={() => setIsAuthenticated(true)} />
  }

  return (
    <div className="app-shell" data-theme={theme}>
      <Sidebar
        activeItem={activeItem}
        onNavigate={(item) => {
          setActiveItem(item)
          setSelectedProgramId(null)
          setSelectedDayId(null)
          writeAppLocation(null, null, item === 'Программы')
        }}
      />
      <main className="main-content">
        <DashboardHeader
          theme={theme}
          onToggleTheme={toggleTheme}
          onLogout={handleLogout}
        />
        {activeItem === 'Программы' ? (
          <ProgramsPage
            selectedProgramId={selectedProgramId}
            selectedDayId={selectedDayId}
            onOpenProgram={(programId) => {
              setSelectedProgramId(programId)
              setSelectedDayId(null)
              writeAppLocation(programId, null)
            }}
            onOpenDay={(dayId) => {
              if (selectedProgramId === null) return
              setSelectedDayId(dayId)
              writeAppLocation(selectedProgramId, dayId)
            }}
            onBackToPrograms={() => {
              setSelectedProgramId(null)
              setSelectedDayId(null)
              writeAppLocation(null, null, true)
            }}
            onBackToProgram={() => {
              setSelectedDayId(null)
              if (selectedProgramId !== null) writeAppLocation(selectedProgramId, null)
            }}
          />
        ) : <div className="dashboard-content">
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
        </div>}
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
