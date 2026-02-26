import { useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { getToken, clearToken, getApiBase, getUsernameFromToken } from './auth'
import { AppRoutes } from './AppRoutes'

interface FortuneItem {
  category: string
  title: string
  content: string
  hint?: string
}

interface CategoryItem {
  id: string
  label: string
}

const CATEGORY_LABELS: Record<string, string> = {
  love: '姻缘',
  career: '事业',
  health: '健康',
  wealth: '财运',
  general: '综合',
}

export function FortuneMain() {
  const navigate = useNavigate()
  const [fortune, setFortune] = useState<FortuneItem | null>(null)
  const [categories, setCategories] = useState<CategoryItem[]>([])
  const [selectedCategory, setSelectedCategory] = useState<string>('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [drawn, setDrawn] = useState(false)
  const apiBase = getApiBase()
  const token = getToken()

  const handleLogout = () => {
    clearToken()
    navigate('/login')
  }

  const fetchCategories = useCallback(() => {
    fetch(`${apiBase}/api/categories`)
      .then((r) => (r.ok ? r.json() : Promise.reject(new Error(r.statusText))))
      .then((data: { categories?: CategoryItem[] }) => {
        setCategories(data?.categories ?? [])
      })
      .catch(() => setCategories([]))
  }, [apiBase])

  if (categories.length === 0 && !loading) fetchCategories()

  const readErrorFromResponse = (r: Response): Promise<never> =>
    r.json().then((d: { error?: string }) => Promise.reject(new Error(d?.error || r.statusText))).catch(() => Promise.reject(new Error(r.statusText)))

  const draw = () => {
    setError(null)
    setLoading(true)
    setDrawn(false)
    const url = selectedCategory
      ? `${apiBase}/api/fortune?category=${encodeURIComponent(selectedCategory)}`
      : `${apiBase}/api/fortune`
    fetch(url)
      .then((r) => (r.ok ? r.json() : readErrorFromResponse(r)))
      .then((data: FortuneItem) => {
        setFortune(data)
        setDrawn(true)
      })
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }

  const drawAI = () => {
    setError(null)
    setLoading(true)
    setDrawn(false)
    const category = selectedCategory || 'general'
    fetch(`${apiBase}/api/fortune/ai`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ category }),
    })
      .then((r) => {
        if (r.ok) return r.json()
        return readErrorFromResponse(r)
      })
      .then((data: FortuneItem) => {
        setFortune(data)
        setDrawn(true)
      })
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }

  const label = fortune ? CATEGORY_LABELS[fortune.category] ?? fortune.category : ''
  const username = token ? getUsernameFromToken(token) : null

  return (
    <div className="app">
      <header className="header">
        <h1 className="title">🔮 高级算命</h1>
        <p className="subtitle">心诚则灵 · 抽签占卜</p>
        {username && (
          <p className="user-bar">
            <span>{username}</span>
            <button type="button" onClick={handleLogout}>退出</button>
          </p>
        )}
      </header>

      <section className="controls">
        {categories.length > 0 && (
          <div className="category-wrap">
            <label htmlFor="cat">运势类别</label>
            <select
              id="cat"
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="category-select"
            >
              <option value="">随机</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.label}
                </option>
              ))}
            </select>
          </div>
        )}
        <button
          type="button"
          className="draw-btn"
          onClick={draw}
          disabled={loading}
        >
          {loading ? '占卜中…' : '抽签占卜'}
        </button>
        <button
          type="button"
          className="draw-btn draw-btn-ai"
          onClick={drawAI}
          disabled={loading}
        >
          AI 占卜
        </button>
      </section>

      {error && (
        <div className="error" role="alert">
          {error}
        </div>
      )}

      {fortune && drawn && !loading && (
        <article className={`card ${drawn ? 'card-visible' : ''}`}>
          <span className="card-category">{label}</span>
          <h2 className="card-title">{fortune.title}</h2>
          <p className="card-content">{fortune.content}</p>
          {fortune.hint && (
            <p className="card-hint">※ {fortune.hint}</p>
          )}
        </article>
      )}

      <footer className="footer">
        <p>仅供娱乐 · 理性看待</p>
      </footer>
    </div>
  )
}

export default function App() {
  return <AppRoutes />
}
