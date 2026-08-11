import type { Heartbeat, Monitor, MonitorStats } from '../../../services/monitor'

export type Status = 'up' | 'down' | 'paused'

export function cn(...classes: (string | false | null | undefined)[]) {
  return classes.filter(Boolean).join(' ')
}

export function getStatus(monitor: Pick<MonitorStats, 'isActive' | 'isUp'>): Status {
  if (!monitor.isActive) return 'paused'
  return monitor.isUp ? 'up' : 'down'
}

export function formatFrequency(seconds: number): string {
  if (seconds < 60) return `${seconds}s`
  if (seconds < 3600) {
    const minutes = seconds / 60
    return Number.isInteger(minutes) ? `${minutes}m` : `${minutes.toFixed(1)}m`
  }
  const hours = seconds / 3600
  return Number.isInteger(hours) ? `${hours}h` : `${hours.toFixed(1)}h`
}

export function formatRelativeTime(iso: string | undefined): string {
  if (!iso) return '—'
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return '—'

  const seconds = Math.max(0, Math.floor((Date.now() - date.getTime()) / 1000))
  if (seconds < 60) return `${seconds}s ago`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
  return `${Math.floor(seconds / 86400)}d ago`
}

export function formatCheckTime(iso: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return '—'
  return date.toLocaleTimeString('en-US', { hour12: false })
}

export function toDeleteMonitor(monitor: MonitorStats): Monitor {
  return {
    id: monitor.id,
    name: monitor.name,
    url: monitor.url,
    method: monitor.method,
    timeout: monitor.timeout,
    frequency: monitor.frequency,
    is_active: monitor.isActive,
    is_up: monitor.isUp,
    created_at: monitor.createdAt,
  }
}

export function chronologicalHistory(history: Heartbeat[]): Heartbeat[] {
  return [...history].reverse()
}
