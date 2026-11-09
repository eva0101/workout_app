export const AUTH_TOKEN_KEY = 'workout-app-token'

export type AuthRequest = {
  login: string
  password: string
}

export type AuthTokenResponse = {
  token: string
}

export type RegisterUserResponse = {
  id: string
  login: string
  created_at: string
}

export type Program = {
  id: number
  name: string
  started_at: string
}

export type ProgramDetails = {
  id: number
  name: string
  training_days: TrainingDaySummary[] | null
}

export type TrainingDaySummary = {
  id: number
  program_id: number
  day_number: number
  name: string
}

export type TrainingDayDetails = {
  id: number
  day_number: number
  name: string
  exercises: Array<{
    id: number
    exercise_id: number
    sets: number
    planned_sets: PlannedSet[]
  }>
}

export type PlannedSet = {
  set_number: number
  repetitions: number | null
  weight: number | null
}

export type MuscleGroup = {
  id: number
  name: string
  title: string
  image_path: string | null
}

export type CatalogExercise = {
  id: number
  name: string
  title: string
  muscle_group_id: number
  description: string | null
  image_path: string | null
}

export type CreatedTrainingDayExercise = {
  id: number
  exercise_id: number
  name_exercise: string
  sets: number
  planned_sets: PlannedSet[]
}

export type ScheduledWorkout = {
  id: number
  scheduled_date: string
  training_day_id: number
  day_number: number
  program_name: string
  created_at?: string
}

export type ScheduledWorkoutWithDayName = ScheduledWorkout & {
  training_day_name: string
}

export type Workout = {
  id: number
  training_day_name?: string | null
  status: string
  started_at: string
  completed_at: string | null
  exercises: WorkoutExercise[]
}

export type WorkoutExercise = {
  id: number
  exercise_id: number
  name: string
  title: string
  order: number
  completed_at: string | null
  planned_sets: WorkoutPlannedSet[]
  actual_sets: WorkoutActualSet[]
}

export type WorkoutPlannedSet = {
  id: number
  set_number: number
  repetitions: number | null
  weight: number | null
}

export type WorkoutActualSet = {
  id: number
  set_number: number
  reps_done: number | null
  weight: number | null
}

export type ExecutedWorkoutSet = {
  workout_exercise_id: number
  set_number: number
  reps_done: number | null
  weight: number | null
}

export type CompletedWorkoutExercise = {
  id: number
  completed_at: string
}

export type StartedWorkout = {
  id?: number
  ID?: number
  training_day_id?: number
  TrainingDayID?: number
  status?: string
  Status?: string
  begin_at?: string
  BeginAt?: string
}

export type ApiError = Error & {
  status?: number
}

export const UNAUTHORIZED_EVENT = 'workout-app:unauthorized'

function clearAuthToken() {
  localStorage.removeItem(AUTH_TOKEN_KEY)
}

export function getStoredToken() {
  return localStorage.getItem(AUTH_TOKEN_KEY)
}

export function setStoredToken(token: string) {
  localStorage.setItem(AUTH_TOKEN_KEY, token)
}

export function removeStoredToken() {
  clearAuthToken()
}

async function parsePayload(response: Response) {
  const text = await response.text()

  if (!text) {
    return null
  }

  try {
    return JSON.parse(text)
  } catch {
    return text
  }
}

export async function apiRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  const token = getStoredToken()
  const headers = new Headers(init.headers ?? {})

  if (!(init.body instanceof FormData) && !headers.has('Content-Type') && init.body !== undefined) {
    headers.set('Content-Type', 'application/json')
  }

  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(`/api${path}`, {
    ...init,
    headers,
  })

  const payload = await parsePayload(response)

  if (response.status === 401) {
    clearAuthToken()
    if (token) {
      window.dispatchEvent(new Event(UNAUTHORIZED_EVENT))
    }
  }

  if (!response.ok) {
    const message =
      typeof payload === 'string'
        ? payload
        : payload && typeof payload === 'object' && 'message' in payload
          ? String((payload as { message?: string }).message)
          : payload && typeof payload === 'object' && 'error' in payload
            ? String((payload as { error?: string }).error)
            : 'Ошибка запроса'

    const error = new Error(message || 'Ошибка запроса') as ApiError
    error.status = response.status
    throw error
  }

  return payload as T
}

export async function loginUser(data: AuthRequest): Promise<AuthTokenResponse> {
  return apiRequest<AuthTokenResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function registerUser(data: AuthRequest): Promise<RegisterUserResponse> {
  return apiRequest<RegisterUserResponse>('/auth/register', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function getPrograms(): Promise<Program[]> {
  return apiRequest<Program[]>('/programs')
}

export async function getProgram(id: number): Promise<ProgramDetails> {
  return apiRequest<ProgramDetails>(`/programs/${id}`)
}

export async function createProgram(name: string): Promise<Program> {
  return apiRequest<Program>('/programs', {
    method: 'POST',
    body: JSON.stringify({ name }),
  })
}

export async function renameProgram(id: number, name: string): Promise<{ name: string }> {
  return apiRequest<{ name: string }>(`/programs/${id}`, {
    method: 'PATCH',
    body: JSON.stringify({ name }),
  })
}

export async function deleteProgram(id: number): Promise<void> {
  await apiRequest<void>(`/programs/${id}`, { method: 'DELETE' })
}

export async function createTrainingDay(programId: number, name: string): Promise<TrainingDaySummary> {
  return apiRequest<TrainingDaySummary>(`/programs/${programId}/days`, {
    method: 'POST',
    body: JSON.stringify({ name }),
  })
}

export async function getTrainingDay(id: number): Promise<TrainingDayDetails> {
  return apiRequest<TrainingDayDetails>(`/days/${id}`)
}

export async function getMuscleGroups(): Promise<MuscleGroup[]> {
  return apiRequest<MuscleGroup[]>('/muscle-groups')
}

export async function getExercises(muscleGroupId?: number): Promise<CatalogExercise[]> {
  const params = new URLSearchParams()
  if (muscleGroupId !== undefined) params.set('muscle_group_id', String(muscleGroupId))
  const query = params.size > 0 ? `?${params.toString()}` : ''
  return apiRequest<CatalogExercise[]>(`/exercises${query}`)
}

export async function createTrainingDayExercise(
  trainingDayId: number,
  exerciseId: number,
  sets: Array<Pick<PlannedSet, 'repetitions' | 'weight'>>,
): Promise<CreatedTrainingDayExercise> {
  return apiRequest<CreatedTrainingDayExercise>(`/days/${trainingDayId}/exercises`, {
    method: 'POST',
    body: JSON.stringify({
      exercise_id: exerciseId,
      sets,
    }),
  })
}

export async function deleteTrainingDayExercise(id: number): Promise<void> {
  await apiRequest<void>(`/days/exercises/${id}`, { method: 'DELETE' })
}

export async function updateTrainingDayExerciseSets(
  id: number,
  sets: Array<Pick<PlannedSet, 'repetitions' | 'weight'>>,
): Promise<void> {
  await apiRequest<void>(`/days/exercises/${id}`, {
    method: 'PATCH',
    body: JSON.stringify({
      sets: sets.map(({ repetitions, weight }) => ({
        reps: repetitions,
        weight,
      })),
    }),
  })
}

export async function deleteTrainingDay(id: number): Promise<void> {
  await apiRequest<void>(`/days/${id}`, { method: 'DELETE' })
}

export async function createSchedule(
  trainingDayId: number,
  dates: string[],
): Promise<ScheduledWorkout[]> {
  return apiRequest<ScheduledWorkout[]>(`/days/${trainingDayId}/schedule`, {
    method: 'POST',
    body: JSON.stringify({ dates }),
  })
}

export async function getSchedule(from: string, to: string): Promise<ScheduledWorkout[]> {
  const params = new URLSearchParams({ from, to })
  return apiRequest<ScheduledWorkout[]>(`/workouts/schedule?${params.toString()}`)
}

export async function getScheduleWithDayNames(
  from: string,
  to: string,
): Promise<ScheduledWorkoutWithDayName[]> {
  const scheduledWorkouts = await getSchedule(from, to)
  if (scheduledWorkouts.length === 0) return []

  const programs = await getPrograms()
  const programDetails = await Promise.all(programs.map(({ id }) => getProgram(id)))
  const trainingDayNames = new Map(
    programDetails.flatMap((program) =>
      (program.training_days ?? []).map((day) => [day.id, day.name] as const),
    ),
  )

  return scheduledWorkouts.map((workout) => {
    const trainingDayName = trainingDayNames.get(workout.training_day_id)
    if (!trainingDayName) {
      throw new Error(`Не удалось определить название тренировочного дня для записи №${workout.id}.`)
    }
    return { ...workout, training_day_name: trainingDayName }
  })
}

export async function startScheduledWorkout(scheduleId: number): Promise<StartedWorkout> {
  return apiRequest<StartedWorkout>(`/workouts/schedule/${scheduleId}/start`, {
    method: 'POST',
  })
}

export async function getActiveWorkout(): Promise<Workout> {
  return apiRequest<Workout>('/workouts/active')
}

export async function getWorkout(id: number): Promise<Workout> {
  return apiRequest<Workout>(`/workouts/${id}`)
}

export async function executeWorkoutSet(
  workoutExerciseId: number,
  repsDone: number,
  weight: number | null,
): Promise<ExecutedWorkoutSet> {
  return apiRequest<ExecutedWorkoutSet>(`/workouts/exercises/${workoutExerciseId}/sets`, {
    method: 'POST',
    body: JSON.stringify({ reps_done: repsDone, weight }),
  })
}

export async function completeWorkoutExercise(workoutExerciseId: number): Promise<CompletedWorkoutExercise> {
  return apiRequest<CompletedWorkoutExercise>(`/workouts/exercises/${workoutExerciseId}/complete`, {
    method: 'POST',
  })
}

export async function completeWorkout(id: number, fatigueScore: number): Promise<void> {
  await apiRequest<void>(`/workouts/${id}/complete`, {
    method: 'POST',
    body: JSON.stringify({ fatigue_score: fatigueScore }),
  })
}

export async function deleteScheduledWorkout(id: number): Promise<void> {
  await apiRequest<void>(`/workouts/schedule/${id}`, { method: 'DELETE' })
}
