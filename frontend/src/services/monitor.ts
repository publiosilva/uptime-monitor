import { gql } from '@apollo/client'
import { apiRequest } from '../lib/api'

export const GET_MONITOR_STATS = gql`
    query GetMonitorStats($monitor_id: ID!) {
        getMonitorStats(monitor_id: $monitor_id) {
            id
            name
            url
            method
            timeout
            frequency
            isActive
            isUp
            createdAt
            stats24h {
                uptimePercentage
                averageLatencyMs
                history {
                    id
                    statusCode
                    latencyMs
                    isUp
                    createdAt
                }
            }
        }
    }
`

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

export type UpdateMonitorInput = {
    name: string
    url: string
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD'
    frequency: number
    timeout: number
    is_active: boolean
}

export type MonitorStats = {
    id: string
    name: string
    url: string
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD'
    timeout: number
    frequency: number
    isActive: boolean
    isUp: boolean
    createdAt: string
    stats24h: MonitorStats24h
}

export type MonitorStats24h = {
    uptimePercentage: number
    averageLatencyMs: number
    history: Heartbeat[]
}

export type Heartbeat = {
    id: string
    statusCode: number
    latencyMs: number
    isUp: boolean
    createdAt: string
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
    },

    async get(id: string): Promise<Monitor> {
        return await apiRequest<Monitor>(`/api/v1/monitors/${id}`, {
            method: 'GET',
        })
    },

    async update(id: string, data: UpdateMonitorInput): Promise<Monitor> {
        return await apiRequest<Monitor>(`/api/v1/monitors/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        })
    },
}
