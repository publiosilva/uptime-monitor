import type { NavigateFunction } from 'react-router-dom'

let navigateFn: NavigateFunction | null = null

export function setNavigate(fn: NavigateFunction): void {
    navigateFn = fn
}

export function redirectToSignIn(): void {
    navigateFn?.('/signin', { replace: true })
}
