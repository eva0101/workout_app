import { ArrowRight, Dumbbell, Eye, EyeOff, LockKeyhole, Moon, Sun, UserRound } from 'lucide-react'
import { type FormEvent, useMemo, useState } from 'react'
import { loginUser, registerUser, setStoredToken } from '../lib/api'

type AuthMode = 'login' | 'register'

type AuthPageProps = {
  theme: 'light' | 'dark'
  onToggleTheme: () => void
  onAuthSuccess: () => void
}

export default function AuthPage({ theme, onToggleTheme, onAuthSuccess }: AuthPageProps) {
  const [mode, setMode] = useState<AuthMode>('login')
  const [form, setForm] = useState({ login: '', password: '' })
  const [showPassword, setShowPassword] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [successMessage, setSuccessMessage] = useState('')

  const title = useMemo(
    () => (mode === 'login' ? 'Вход в аккаунт' : 'Создать аккаунт'),
    [mode],
  )

  const subtitle = useMemo(
    () =>
      mode === 'login'
        ? 'Продолжай свой путь к сильному телу.'
        : 'Создай профиль, чтобы начать тренировки.',
    [mode],
  )

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setSuccessMessage('')

    const login = form.login.trim()
    const password = form.password.trim()

    if (login.length < 3) {
      setError('Логин должен содержать не менее 3 символов.')
      return
    }

    if (password.length < 8) {
      setError('Пароль должен содержать не менее 8 символов.')
      return
    }

    try {
      setIsSubmitting(true)

      if (mode === 'login') {
        const response = await loginUser({ login, password })
        setStoredToken(response.token)
        onAuthSuccess()
        return
      }

      await registerUser({ login, password })
      setSuccessMessage('Регистрация успешна. Теперь войдите в аккаунт.')
      setForm({ login: '', password: '' })
      setMode('login')
    } catch (authError) {
      const message =
        authError instanceof Error && authError.message
          ? authError.message
          : 'Не удалось выполнить запрос. Проверьте соединение.'

      setError(message)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="auth-shell" data-theme={theme}>
      <header className="auth-topbar">
        <div className="auth-brand" aria-label="Workout App brand">
          <div className="brand-mark" aria-hidden="true">
            <Dumbbell size={19} strokeWidth={2.2} />
          </div>
          <div>
            <div className="auth-brand-name">Workout App</div>
            <div className="auth-brand-caps">training system</div>
          </div>
        </div>

        <button
          className="icon-button theme-toggle"
          type="button"
          onClick={onToggleTheme}
          aria-label={theme === 'light' ? 'Включить тёмную тему' : 'Включить светлую тему'}
          title={theme === 'light' ? 'Тёмная тема' : 'Светлая тема'}
        >
          {theme === 'light' ? <Moon size={18} strokeWidth={1.9} /> : <Sun size={18} strokeWidth={1.9} />}
        </button>
      </header>

      <section className="auth-panel" aria-live="polite">
        <div className="auth-panel-header">
          <span className="eyebrow auth-eyebrow">Workout App</span>
          <h1>{title}</h1>
          <p>{subtitle}</p>
        </div>

        <form className="auth-form" onSubmit={handleSubmit}>
          <label className="auth-field" htmlFor="login">
            <span>Логин</span>
            <div className="auth-input-wrap">
              <UserRound size={17} aria-hidden="true" />
              <input
                id="login"
                name="login"
                type="text"
                autoComplete="username"
                value={form.login}
                onChange={(event) => setForm((current) => ({ ...current, login: event.target.value }))}
                placeholder="Введите логин"
                maxLength={32}
              />
            </div>
          </label>

          <label className="auth-field" htmlFor="password">
            <span>Пароль</span>
            <div className="auth-input-wrap">
              <LockKeyhole size={17} aria-hidden="true" />
              <input
                id="password"
                name="password"
                type={showPassword ? 'text' : 'password'}
                autoComplete={mode === 'login' ? 'current-password' : 'new-password'}
                value={form.password}
                onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))}
                placeholder="Введите пароль"
                maxLength={64}
              />
              <button
                className="auth-visibility-toggle"
                type="button"
                onClick={() => setShowPassword((shown) => !shown)}
                aria-label={showPassword ? 'Скрыть пароль' : 'Показать пароль'}
              >
                {showPassword ? <EyeOff size={17} /> : <Eye size={17} />}
              </button>
            </div>
          </label>

          {error ? <div className="auth-message auth-message-error">{error}</div> : null}
          {successMessage ? <div className="auth-message auth-message-success">{successMessage}</div> : null}

          <button className="auth-submit" type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Проверка...' : mode === 'login' ? 'Войти' : 'Создать аккаунт'}
            <ArrowRight size={16} strokeWidth={2} />
          </button>
        </form>

        <div className="auth-switcher">
          <span>{mode === 'login' ? 'Нет аккаунта?' : 'Уже есть аккаунт?'}</span>
          <button
            type="button"
            onClick={() => {
              setMode((current) => (current === 'login' ? 'register' : 'login'))
              setError('')
              setSuccessMessage('')
            }}
          >
            {mode === 'login' ? 'Регистрация' : 'Войти'}
          </button>
        </div>
      </section>
    </div>
  )
}
