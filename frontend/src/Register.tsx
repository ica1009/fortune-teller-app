import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getApiBase } from './auth'
import './AuthForm.css'

export default function Register() {
  const navigate = useNavigate()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const res = await fetch(`${getApiBase()}/api/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      })
      const data = await res.json().catch(() => ({}))
      if (!res.ok) {
        setError(data?.error ?? res.statusText ?? '注册失败')
        return
      }
      navigate('/login')
    } catch (err) {
      setError(err instanceof Error ? err.message : '请求失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-form">
      <h1 className="auth-title">🔮 高级算命</h1>
      <p className="auth-subtitle">注册</p>
      <form onSubmit={handleSubmit}>
        <label>
          用户名（3–32 位字母、数字或下划线）
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            autoComplete="username"
          />
        </label>
        <label>
          密码（至少 6 位）
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            autoComplete="new-password"
          />
        </label>
        {error && <p className="auth-error">{error}</p>}
        <button type="submit" disabled={loading}>
          {loading ? '注册中…' : '注册'}
        </button>
      </form>
      <p className="auth-switch">
        已有账号？ <Link to="/login">去登录</Link>
      </p>
    </div>
  )
}
