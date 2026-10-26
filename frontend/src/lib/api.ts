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
  }>
}

export type ScheduledWorkout = {
  id: number
  scheduled_date: string
  training_day_id: number
  day_number: number
  program_name: string
  created_at?: string
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

export async function deleteScheduledWorkout(id: number): Promise<void> {
  await apiRequest<void>(`/workouts/schedule/${id}`, { method: 'DELETE' })
}
