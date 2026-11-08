import { useCallback, useEffect, useState } from 'react'
import { Dumbbell } from 'lucide-react'
import DashboardHeader from './components/DashboardHeader'
import Sidebar from './components/Sidebar'
import AuthPage from './pages/AuthPage'
import DashboardPage from './pages/DashboardPage'
import ProgramsPage from './pages/ProgramsPage'
import WorkoutPage from './pages/WorkoutPage'
import {
  completeWorkout,
  executeWorkoutSet,
  getActiveWorkout,
  getStoredToken,
  getWorkout,
  removeStoredToken,
  startScheduledWorkout,
  UNAUTHORIZED_EVENT,
  type ApiError,
  type Workout,
} from './lib/api'

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
    activeItem: params.get('section') === 'programs'
      ? 'Программы'
      : params.get('section') === 'workouts'
        ? 'Тренировки'
        : 'Обзор',
    programId: isProgramRoute ? programId : null,
    dayId: isProgramRoute && Number.isInteger(dayId) && dayId > 0 ? dayId : null,
  }
}

function writeAppLocation(
  programId: number | null,
  dayId: number | null,
  keepProgramsSection = false,
  keepWorkoutsSection = false,
) {
  const url = new URL(window.location.href)
  url.search = ''

  if (programId !== null) {
    url.searchParams.set('section', 'programs')
    url.searchParams.set('program', String(programId))
    if (dayId !== null) url.searchParams.set('day', String(dayId))
  } else if (keepProgramsSection) {
    url.searchParams.set('section', 'programs')
  } else if (keepWorkoutsSection) {
    url.searchParams.set('section', 'workouts')
  }

  window.history.pushState(null, '', url)
}

function getErrorMessage(error: unknown): string {
  const status = typeof error === 'object' && error !== null && 'status' in error
    ? (error as ApiError).status
    : undefined
  if (status === 401) return 'Сессия завершена. Войдите в аккаунт снова.'
  if (status === 404) return 'Тренировка не найдена или больше не активна.'
  if (status === 409) return 'У вас уже есть активная тренировка. Продолжите её вместо запуска новой.'
  if (error instanceof TypeError) return 'Не удалось связаться с сервером. Проверьте подключение.'
  if (error instanceof Error && error.message) return error.message
  return 'Не удалось выполнить запрос. Попробуйте ещё раз.'
}

function getWorkoutID(started: { id?: number; ID?: number }): number | null {
  const id = started.id ?? started.ID
  return Number.isInteger(id) && Number(id) > 0 ? Number(id) : null
}

export default function App() {
  const [initialLocation] = useState(readAppLocation)
  const [activeItem, setActiveItem] = useState(initialLocation.activeItem)
  const [theme, setTheme] = useState<Theme>(getInitialTheme)
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(getIsAuthenticated)
  const [selectedProgramId, setSelectedProgramId] = useState<number | null>(initialLocation.programId)
  const [selectedDayId, setSelectedDayId] = useState<number | null>(initialLocation.dayId)
  const [activeWorkout, setActiveWorkout] = useState<Workout | null>(null)
  const [currentWorkout, setCurrentWorkout] = useState<Workout | null>(null)
  const [currentWorkoutTitle, setCurrentWorkoutTitle] = useState('')
  const [isActiveWorkoutLoading, setIsActiveWorkoutLoading] = useState(false)
  const [isStartingWorkout, setIsStartingWorkout] = useState(false)
  const [isSavingSet, setIsSavingSet] = useState(false)
  const [isCompletingWorkout, setIsCompletingWorkout] = useState(false)
  const [workoutError, setWorkoutError] = useState('')
  const [activeWorkoutError, setActiveWorkoutError] = useState('')

  const refreshActiveWorkout = useCallback(async () => {
    setIsActiveWorkoutLoading(true)
    setActiveWorkoutError('')
    try {
      setActiveWorkout(await getActiveWorkout())
    } catch (error) {
      if (typeof error === 'object' && error !== null && 'status' in error && (error as ApiError).status === 404) {
        setActiveWorkout(null)
      } else {
        setActiveWorkoutError(getErrorMessage(error))
      }
    } finally {
      setIsActiveWorkoutLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!isAuthenticated) return

    let cancelled = false
    setIsActiveWorkoutLoading(true)
    setActiveWorkoutError('')
    getActiveWorkout()
      .then(async (workout) => {
        if (cancelled) return
        setActiveWorkout(workout)
        if (initialLocation.activeItem === 'Тренировки') {
          try {
            const details = await getWorkout(workout.id)
            if (!cancelled) {
              setCurrentWorkout(details)
              setCurrentWorkoutTitle(details.training_day_name || `Тренировка #${details.id}`)
            }
          } catch (error) {
            if (!cancelled) setActiveWorkoutError(getErrorMessage(error))
          }
        }
      })
      .catch((error: unknown) => {
        if (cancelled) return
        if (typeof error === 'object' && error !== null && 'status' in error && (error as ApiError).status === 404) {
          setActiveWorkout(null)
          if (initialLocation.activeItem === 'Тренировки') setActiveItem('Обзор')
        } else {
          setActiveWorkoutError(getErrorMessage(error))
        }
      })
      .finally(() => {
        if (!cancelled) setIsActiveWorkoutLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [initialLocation.activeItem, isAuthenticated])

  useEffect(() => {
    function handleUnauthorized() {
      setIsAuthenticated(false)
      setActiveItem('Обзор')
      setSelectedProgramId(null)
      setSelectedDayId(null)
      setActiveWorkout(null)
      setCurrentWorkout(null)
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
      if (location.activeItem === 'Тренировки' && activeWorkout) {
        setCurrentWorkout(activeWorkout)
        setCurrentWorkoutTitle(activeWorkout.training_day_name || `Тренировка #${activeWorkout.id}`)
      }
    }

    window.addEventListener('popstate', handlePopState)
    return () => window.removeEventListener('popstate', handlePopState)
  }, [activeWorkout])

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
    setActiveWorkout(null)
    setCurrentWorkout(null)
    writeAppLocation(null, null)
  }

  async function openWorkout(workoutId: number, title?: string) {
    setWorkoutError('')
    try {
      const workout = await getWorkout(workoutId)
      setCurrentWorkout(workout)
      setActiveWorkout(workout.status === 'in_progress' ? workout : activeWorkout)
      setCurrentWorkoutTitle(title || workout.training_day_name || `Тренировка #${workout.id}`)
      setActiveItem('Тренировки')
      setSelectedProgramId(null)
      setSelectedDayId(null)
      writeAppLocation(null, null, false, true)
    } catch (error) {
      setWorkoutError(getErrorMessage(error))
    }
  }

  async function handleStartScheduledWorkout(scheduleId: number, title: string) {
    if (activeWorkout) {
      throw new Error('У вас уже есть активная тренировка. Продолжите её вместо запуска новой.')
    }
    if (isActiveWorkoutLoading) {
      throw new Error('Проверяем состояние активной тренировки. Попробуйте ещё раз через мгновение.')
    }
    if (isStartingWorkout) {
      throw new Error('Запуск тренировки уже выполняется.')
    }
    setWorkoutError('')
    setIsStartingWorkout(true)
    try {
      const started = await startScheduledWorkout(scheduleId)
      const workoutId = getWorkoutID(started)
      if (workoutId === null) {
        throw new Error('Backend вернул некорректный ID созданной тренировки.')
      }
      await openWorkout(workoutId, title)
      await refreshActiveWorkout()
    } catch (error) {
      if (typeof error === 'object' && error !== null && 'status' in error && (error as ApiError).status === 409) {
        await refreshActiveWorkout()
      }
      throw error
    } finally {
      setIsStartingWorkout(false)
    }
  }

  async function handleAddSet(workoutExerciseId: number, repsDone: number, weight: number | null): Promise<boolean> {
    if (!currentWorkout || isSavingSet) return false
    setWorkoutError('')
    setIsSavingSet(true)
    try {
      await executeWorkoutSet(workoutExerciseId, repsDone, weight)
      try {
        const refreshed = await getWorkout(currentWorkout.id)
        setCurrentWorkout(refreshed)
        setActiveWorkout(refreshed)
      } catch (error) {
        setWorkoutError(`Подход сохранён, но не удалось обновить тренировку. ${getErrorMessage(error)}`)
      }
      return true
    } catch (error) {
      setWorkoutError(getErrorMessage(error))
      return false
    } finally {
      setIsSavingSet(false)
    }
  }

  async function handleCompleteWorkout(fatigueScore: number) {
    if (!currentWorkout || isCompletingWorkout) return
    setWorkoutError('')
    setIsCompletingWorkout(true)
    try {
      await completeWorkout(currentWorkout.id, fatigueScore)
      setCurrentWorkout(null)
      setActiveWorkout(null)
      setCurrentWorkoutTitle('')
      setActiveItem('Обзор')
      setSelectedProgramId(null)
      setSelectedDayId(null)
      writeAppLocation(null, null)
    } catch (error) {
      setWorkoutError(getErrorMessage(error))
    } finally {
      setIsCompletingWorkout(false)
    }
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
          setWorkoutError('')
          if (item === 'Тренировки' && activeWorkout) {
            void openWorkout(activeWorkout.id)
          } else {
            writeAppLocation(null, null, item === 'Программы', item === 'Тренировки')
          }
        }}
      />
      <main className="main-content">
        <DashboardHeader
          theme={theme}
          onToggleTheme={toggleTheme}
          onLogout={handleLogout}
        />
        {workoutError && !(activeItem === 'Тренировки' && currentWorkout) && (
          <div className="program-alert mx-5 mt-4 sm:mx-8 lg:mx-12" role="alert">
            {workoutError}
          </div>
        )}
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
        ) : activeItem === 'Тренировки' ? (
          currentWorkout ? (
            <WorkoutPage
              workout={currentWorkout}
              title={currentWorkoutTitle}
              isSavingSet={isSavingSet}
              isCompleting={isCompletingWorkout}
              error={workoutError}
              onAddSet={handleAddSet}
              onComplete={(fatigueScore) => void handleCompleteWorkout(fatigueScore)}
              onBack={() => {
                setActiveItem('Обзор')
                setWorkoutError('')
                writeAppLocation(null, null)
              }}
            />
          ) : (
            <section className="dashboard-content">
              <div className="panel-card empty-state">
                <div className="empty-icon"><Dumbbell size={21} /></div>
                <h2 className="mt-3 text-sm font-semibold text-ink">Нет активной тренировки</h2>
                <p className="mt-1 text-xs text-muted">Выберите запланированную тренировку в обзоре или календаре.</p>
                {activeWorkoutError && <p className="program-alert mt-4" role="alert">{activeWorkoutError}</p>}
                <button
                  className="program-secondary-button mt-4"
                  type="button"
                  onClick={() => {
                    setActiveItem('Обзор')
                    writeAppLocation(null, null)
                  }}
                >
                  К обзору
                </button>
              </div>
            </section>
          )
        ) : (
          <DashboardPage
            activeWorkout={activeWorkout}
            activeWorkoutError={activeWorkoutError}
            isActiveWorkoutLoading={isActiveWorkoutLoading}
            isStartingWorkout={isStartingWorkout}
            onStartScheduledWorkout={handleStartScheduledWorkout}
            onResumeWorkout={() => activeWorkout && void openWorkout(activeWorkout.id)}
            onOpenPrograms={() => {
              setActiveItem('Программы')
              setSelectedProgramId(null)
              setSelectedDayId(null)
              writeAppLocation(null, null, true)
            }}
          />
        )}
      </main>
    </div>
  )
}
