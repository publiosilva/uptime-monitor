import { apiRequest } from '../lib/api'
import { clearToken, getToken, setToken } from '../lib/token'

type LoginResponse = {
    token: string
}

type LoginCredentials = {
    email: string
    password: string
}

export const authService = {
    getToken,
    setToken,
    clearToken,

    isAuthenticated(): boolean {
        return getToken() !== null
    },

    async login({ email, password }: LoginCredentials): Promise<string> {
        const { token } = await apiRequest<LoginResponse>('/api/v1/auth/login', {
            method: 'POST',
            body: JSON.stringify({ email, password }),
        })

        setToken(token)
        return token
    },

    logout(): void {
        clearToken()
    },
}
