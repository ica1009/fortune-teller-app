/** 鉴权：token 存 localStorage，API 基地址从环境读取。 */

const TOKEN_KEY = 'fortune_teller_token'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

export function getApiBase(): string {
  if (import.meta.env.DEV) return ''
  return (import.meta.env.VITE_API_BASE ?? '').replace(/\/$/, '')
}

export function getUsernameFromToken(token: string): string | null {
  try {
    const payload = token.split('.')[1]
    if (!payload) return null
    const decoded = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
    return decoded.username ?? null
  } catch {
    return null
  }
}
