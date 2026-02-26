/**
 * 应用路由：/（占卜主界面）、/login、/register。
 * 根路径需 token，否则重定向到 /login；已登录访问 /login、/register 重定向到 /。
 * 使用 PageTransition 做 200–300ms 淡入+短距滑动。
 */

import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { getToken } from './auth'
import { FortuneMain } from './App'
import Login from './Login'
import Register from './Register'

/** 无 token 时重定向到 /login */
function Protected({ children }: { children: React.ReactNode }) {
  if (!getToken()) return <Navigate to="/login" replace />
  return <>{children}</>
}

/** 已登录时重定向到 /，避免重复登录/注册 */
function GuestOnly({ children }: { children: React.ReactNode }) {
  if (getToken()) return <Navigate to="/" replace />
  return <>{children}</>
}

/** 页面切换时淡入 + 短距上滑，约 250ms */
function PageTransition({ children }: { children: React.ReactNode }) {
  const [active, setActive] = React.useState(false)
  React.useEffect(() => {
    const id = requestAnimationFrame(() => setActive(true))
    return () => cancelAnimationFrame(id)
  }, [])
  return (
    <div className={`page-transition ${active ? 'page-enter-active' : ''}`}>
      {children}
    </div>
  )
}

export function AppRoutes() {
  return (
    <Routes>
      <Route
        path="/"
        element={
          <Protected>
            <PageTransition>
              <FortuneMain />
            </PageTransition>
          </Protected>
        }
      />
      <Route
        path="/login"
        element={
          <GuestOnly>
            <PageTransition>
              <Login />
            </PageTransition>
          </GuestOnly>
        }
      />
      <Route
        path="/register"
        element={
          <GuestOnly>
            <PageTransition>
              <Register />
            </PageTransition>
          </GuestOnly>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
