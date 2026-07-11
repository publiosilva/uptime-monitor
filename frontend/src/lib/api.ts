import { redirectToSignIn } from './navigation'
import { clearToken, getToken } from './token'

export class ApiError extends Error {
    status: number

    constructor(status: number, message: string) {
        super(message)
        this.name = 'ApiError'
        this.status = status
    }
}

type ErrorResponse = {
    error?: string
}

function isAuthRequest(path: string): boolean {
    return path.includes('/auth/login') || path.includes('/auth/register')
}

function handleUnauthorized(path: string): void {
    if (isAuthRequest(path)) return

    clearToken()
    redirectToSignIn()
}

export async function apiRequest<T>(path: string, options?: RequestInit): Promise<T> {
    const API_BASE = import.meta.env.VITE_API_URL ?? ''
    const token = getToken()
    const headers = new Headers(options?.headers)

    if (!headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json')
    }

    if (token) {
        headers.set('Authorization', `Bearer ${token}`)
    }

    const response = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers,
    })

    const data: T & ErrorResponse = await response.json().catch(() => ({} as T & ErrorResponse))

    if (response.status === 401) {
        handleUnauthorized(path)
        throw new ApiError(response.status, data.error ?? 'Unauthorized')
    }

    if (!response.ok) {
        throw new ApiError(response.status, data.error ?? 'Request failed')
    }

    return data
}
