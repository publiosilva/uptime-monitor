import { apiRequest } from '../lib/api'

export type Monitor = {
    id: string
    name: string
    url: string
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD'
    timeout: number
    frequency: number
    is_active: boolean
    is_up: boolean
    created_at: string
}

export type CreateMonitorInput = {
    name: string
    url: string
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD'
    frequency: number
    is_active: boolean
}

export const monitorService = {
    async create(data: CreateMonitorInput): Promise<Monitor> {
        return await apiRequest<Monitor>('/api/v1/monitors', {
            method: 'POST',
            body: JSON.stringify(data),
        })
    },

    async list(): Promise<Monitor[]> {
        return await apiRequest<Monitor[]>('/api/v1/monitors', {
            method: 'GET',
        })
    },

    async delete(id: string): Promise<void> {
        return await apiRequest<void>(`/api/v1/monitors/${id}`, {
            method: 'DELETE',
        })
    }
}
