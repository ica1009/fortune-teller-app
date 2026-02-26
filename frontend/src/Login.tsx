import { useState } from 'react'
import { getApiBase, setToken } from './auth'
import './AuthForm.css'

interface LoginProps {
  onSuccess: () => void
  onSwitchRegister: () => void
}

export default function Login({ onSuccess, onSwitchRegister }: LoginProps) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const res = await fetch(`${getApiBase()}/api/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      })
      const data = await res.json().catch(() => ({}))
      if (!res.ok) {
        setError(data?.error ?? res.statusText ?? '登录失败')
        return
      }
      if (data.token) {
        setToken(data.token)
        onSuccess()
      } else {
        setError('未返回 token')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '请求失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-form">
      <h1 className="auth-title">🔮 高级算命</h1>
      <p className="auth-subtitle">登录</p>
      <form onSubmit={handleSubmit}>
        <label>
          用户名
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            autoComplete="username"
          />
        </label>
        <label>
          密码
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            autoComplete="current-password"
          />
        </label>
        {error && <p className="auth-error">{error}</p>}
        <button type="submit" disabled={loading}>
          {loading ? '登录中…' : '登录'}
        </button>
      </form>
      <p className="auth-switch">
        还没有账号？ <button type="button" onClick={onSwitchRegister}>注册</button>
      </p>
    </div>
  )
}
