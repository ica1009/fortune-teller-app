#!/usr/bin/env bash
set -e
ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

echo "Starting backend on :8081 ..."
(cd backend && go run ./cmd/server) &
BACKEND_PID=$!

echo "Starting frontend on :5174 ..."
(cd frontend && npm run dev) &
FRONTEND_PID=$!

echo "Backend PID: $BACKEND_PID  Frontend PID: $FRONTEND_PID"
echo "API: http://localhost:8081  App: http://localhost:5174"
echo "Press Ctrl+C to stop both."
trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT TERM
wait
