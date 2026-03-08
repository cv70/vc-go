import { OpenAPI } from './client/core/OpenAPI'

const DEFAULT_BASE_URL = 'http://localhost:8888'
const TOKEN_KEY = 'vc_token'

export function configureApiClient(baseUrl?: string): void {
  OpenAPI.BASE = baseUrl ?? import.meta.env.VITE_API_BASE_URL ?? DEFAULT_BASE_URL

  OpenAPI.TOKEN = async () => {
    const token = window.localStorage.getItem(TOKEN_KEY)
    return token ?? ''
  }

  OpenAPI.HEADERS = async () => ({
    'X-Request-Id': crypto.randomUUID(),
  })
}

export function setAuthToken(token: string): void {
  window.localStorage.setItem(TOKEN_KEY, token)
}

export function clearAuthToken(): void {
  window.localStorage.removeItem(TOKEN_KEY)
}

export function handleApiErrorStatus(status: number): void {
  if (status === 401) {
    clearAuthToken()
  }
}
