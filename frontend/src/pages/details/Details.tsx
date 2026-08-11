import { useState } from 'react'
import { useQuery } from '@apollo/client/react'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, Pause, Play, Trash2 } from 'lucide-react'
import type { MonitorStats } from '../../services/monitor'
import { GET_MONITOR_STATS, monitorService } from '../../services/monitor'
import Button from '../../components/Button'
import { ApiError } from '../../lib/api'
import DeleteMonitor from '../home/components/DeleteMonitor'
import UrlCard from './components/UrlCard'
import UptimeTimeline from './components/UptimeTimeline'
import RecentChecks from './components/RecentChecks'
import ResponseLatencyChart from './components/ResponseLatencyChart'
import { StatusBadge, StatusDot, StatCard } from './components/Status'
import {
  formatFrequency,
  formatRelativeTime,
  getStatus,
  toDeleteMonitor,
} from './components/helpers'

type GetMonitorStatsData = {
  getMonitorStats: MonitorStats
}

function DetailsSkeleton() {
  return (
    <div className="min-h-screen bg-[#070c18]">
      <header className="sticky top-0 z-40 bg-[#070c18]/90 backdrop-blur border-b border-[#131f38]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 h-14 flex items-center gap-3 animate-pulse">
          <div className="h-4 w-16 bg-[#1e3050] rounded-md" />
          <div className="h-4 w-px bg-[#1a2d4a]" />
          <div className="h-4 w-40 bg-[#1e3050] rounded-md" />
        </div>
      </header>
      <main className="max-w-6xl mx-auto px-4 sm:px-6 py-6 flex flex-col gap-5 animate-pulse">
        <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl h-24" />
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl h-24" />
          ))}
        </div>
        <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl h-48" />
        <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl h-32" />
        <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl h-64" />
      </main>
    </div>
  )
}

export default function Details() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [showDelete, setShowDelete] = useState(false)
  const [isTogglingActive, setIsTogglingActive] = useState(false)
  const [actionError, setActionError] = useState<string | null>(null)

  const { data, loading, error, refetch } = useQuery<GetMonitorStatsData>(GET_MONITOR_STATS, {
    variables: { monitor_id: id },
    skip: !id,
  })

  async function toggleActive(monitor: MonitorStats) {
    try {
      setIsTogglingActive(true)
      setActionError(null)
      await monitorService.update(monitor.id, {
        name: monitor.name,
        url: monitor.url,
        method: monitor.method,
        frequency: monitor.frequency,
        timeout: monitor.timeout,
        is_active: !monitor.isActive,
      })
      await refetch()
    } catch (err) {
      setActionError(err instanceof ApiError ? err.message : 'Failed to update monitor')
    } finally {
      setIsTogglingActive(false)
    }
  }

  if (loading) {
    return <DetailsSkeleton />
  }

  if (error || !data?.getMonitorStats) {
    return (
      <div className="min-h-screen bg-[#070c18] flex flex-col">
        <header className="sticky top-0 z-40 bg-[#070c18]/90 backdrop-blur border-b border-[#131f38]">
          <div className="max-w-6xl mx-auto px-4 sm:px-6 h-14 flex items-center gap-3">
            <button
              type="button"
              onClick={() => navigate('/')}
              className="flex items-center gap-1.5 text-[#3a5070] hover:text-slate-200 transition-colors text-sm"
            >
              <ArrowLeft className="w-4 h-4" />
              <span className="hidden sm:block">Back</span>
            </button>
          </div>
        </header>
        <main className="max-w-6xl mx-auto px-4 sm:px-6 py-6">
          <div className="bg-red-500/10 border border-red-500/25 rounded-xl px-4 py-3 text-sm text-red-400">
            {error?.message ?? 'Monitor not found'}
          </div>
        </main>
      </div>
    )
  }

  const monitor = data.getMonitorStats
  const history = monitor.stats24h?.history ?? []
  const lastHeartbeat = history[0]
  const status = getStatus(monitor)
  const uptime = monitor.stats24h?.uptimePercentage ?? 0
  const avgLatency = Math.round(monitor.stats24h?.averageLatencyMs ?? 0)

  return (
    <div className="min-h-screen bg-[#070c18]">
      <header className="sticky top-0 z-40 bg-[#070c18]/90 backdrop-blur border-b border-[#131f38]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 h-14 flex items-center gap-3">
          <button
            type="button"
            onClick={() => navigate('/')}
            className="flex items-center gap-1.5 text-[#3a5070] hover:text-slate-200 transition-colors text-sm flex-shrink-0"
          >
            <ArrowLeft className="w-4 h-4" />
            <span className="hidden sm:block">Back</span>
          </button>
          <div className="h-4 w-px bg-[#1a2d4a]" />
          <div className="flex items-center gap-2.5 min-w-0 flex-1">
            <StatusDot status={status} />
            <span className="text-sm font-semibold text-slate-100 truncate">{monitor.name}</span>
            <StatusBadge status={status} />
          </div>
          <div className="flex items-center gap-2 flex-shrink-0">
            {monitor.isActive && (
              <div className="hidden sm:flex items-center gap-1.5 text-xs text-[#3a5070] px-2.5 py-1 rounded-full bg-[#0c1422] border border-[#1a2d4a]">
                <span className="w-1.5 h-1.5 rounded-full bg-cyan-400 animate-pulse" />
                Live
              </div>
            )}
            <Button
              variant="secondary"
              size="sm"
              type="button"
              disabled={isTogglingActive}
              onClick={() => toggleActive(monitor)}
            >
              {monitor.isActive ? (
                <>
                  <Pause className="w-3.5 h-3.5" />
                  <span className="hidden sm:block">{isTogglingActive ? 'Pausing…' : 'Pause'}</span>
                </>
              ) : (
                <>
                  <Play className="w-3.5 h-3.5" />
                  <span className="hidden sm:block">{isTogglingActive ? 'Resuming…' : 'Resume'}</span>
                </>
              )}
            </Button>
            <Button variant="destructive" size="sm" type="button" onClick={() => setShowDelete(true)}>
              <Trash2 className="w-3.5 h-3.5" />
              <span className="hidden sm:block">Delete</span>
            </Button>
          </div>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-4 sm:px-6 py-6 flex flex-col gap-5">
        {actionError && (
          <div className="bg-red-500/10 border border-red-500/25 rounded-xl px-4 py-3 text-sm text-red-400">
            {actionError}
          </div>
        )}

        <UrlCard monitor={monitor} />

        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
          <StatCard
            label="24h Uptime"
            value={`${uptime.toFixed(uptime % 1 === 0 ? 0 : 2)}%`}
            sub="last 24 hours"
            accent={uptime >= 99 ? 'emerald' : uptime >= 95 ? 'cyan' : 'red'}
          />
          <StatCard
            label="Avg Latency"
            value={`${avgLatency}ms`}
            sub="response time"
            accent="cyan"
          />
          <StatCard
            label="Last Check"
            value={formatRelativeTime(lastHeartbeat?.createdAt)}
            sub={`every ${formatFrequency(monitor.frequency)}`}
            accent="slate"
          />
          <StatCard
            label="Status Code"
            value={lastHeartbeat?.statusCode ?? '—'}
            sub={
              lastHeartbeat == null
                ? '—'
                : lastHeartbeat.isUp
                  ? 'OK'
                  : 'Service Unavailable'
            }
            accent={
              lastHeartbeat == null
                ? 'slate'
                : lastHeartbeat.isUp
                  ? 'emerald'
                  : 'red'
            }
          />
        </div>

        <ResponseLatencyChart history={history} />
        <UptimeTimeline history={history} />
        <RecentChecks history={history} />
      </main>

      {showDelete && (
        <DeleteMonitor
          monitor={toDeleteMonitor(monitor)}
          onClose={() => setShowDelete(false)}
          onDeleted={() => navigate('/')}
        />
      )}
    </div>
  )
}
