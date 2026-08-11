import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Eye, MoreVertical, Pause, Play, Trash2 } from 'lucide-react'
import type { Monitor } from '../../../services/monitor'
import { monitorService } from '../../../services/monitor'
import DeleteMonitor from './DeleteMonitor'

type Status = 'up' | 'down' | 'paused'

function cn(...classes: (string | false | null | undefined)[]) {
  return classes.filter(Boolean).join(' ')
}

function getStatus(monitor: Monitor): Status {
  if (!monitor.is_active) return 'paused'
  return monitor.is_up ? 'up' : 'down'
}

function formatFrequency(seconds: number): string {
  if (seconds < 60) return `${seconds}s`
  if (seconds < 3600) {
    const minutes = seconds / 60
    return Number.isInteger(minutes) ? `${minutes}m` : `${minutes.toFixed(1)}m`
  }
  const hours = seconds / 3600
  return Number.isInteger(hours) ? `${hours}h` : `${hours.toFixed(1)}h`
}

function StatusDot({ status }: { status: Status }) {
  const colors = { up: 'bg-emerald-400', down: 'bg-red-400', paused: 'bg-slate-500' }
  return (
    <span className="relative flex h-2 w-2 flex-shrink-0">
      {status === 'down' && (
        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-70" />
      )}
      <span className={cn('relative inline-flex rounded-full h-2 w-2', colors[status])} />
    </span>
  )
}

function StatusBadge({ status }: { status: Status }) {
  const cfg = {
    up: { label: 'Operational', cls: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/25' },
    down: { label: 'Down', cls: 'bg-red-500/10 text-red-400 border-red-500/25' },
    paused: { label: 'Paused', cls: 'bg-slate-500/10 text-slate-500 border-slate-500/25' },
  }
  const { label, cls } = cfg[status]
  return (
    <span className={cn('inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-semibold border', cls)}>
      <StatusDot status={status} />
      {label}
    </span>
  )
}

function MethodBadge({ method }: { method: Monitor['method'] }) {
  const colors: Record<Monitor['method'], string> = {
    GET: 'text-emerald-400 bg-emerald-400/10',
    POST: 'text-sky-400 bg-sky-400/10',
    PUT: 'text-amber-400 bg-amber-400/10',
    PATCH: 'text-orange-400 bg-orange-400/10',
    DELETE: 'text-red-400 bg-red-400/10',
    HEAD: 'text-violet-400 bg-violet-400/10',
  }
  return (
    <span className={cn('font-mono text-xs px-1.5 py-0.5 rounded font-semibold tracking-wide', colors[method])}>
      {method}
    </span>
  )
}

export default function MonitorCard({
  monitor,
  onDeleted,
  onUpdated,
}: {
  monitor: Monitor
  onDeleted?: () => void
  onUpdated?: () => void
}) {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [isDeleteMonitorOpen, setIsDeleteMonitorOpen] = useState(false)
  const [isTogglingActive, setIsTogglingActive] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const navigate = useNavigate()
  const status = getStatus(monitor)

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setIsMenuOpen(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  async function toggleActive(e: React.MouseEvent) {
    e.stopPropagation()
    setIsMenuOpen(false)
    try {
      setIsTogglingActive(true)
      await monitorService.update(monitor.id, {
        name: monitor.name,
        url: monitor.url,
        method: monitor.method,
        frequency: monitor.frequency,
        timeout: monitor.timeout,
        is_active: !monitor.is_active,
      })
      onUpdated?.()
    } finally {
      setIsTogglingActive(false)
    }
  }

  return (
    <>
      <div
        onClick={() => navigate(`/details/${monitor.id}`)}
        className={cn(
          'bg-[#0c1422] border rounded-xl p-4 flex flex-col gap-3 cursor-pointer',
          'transition-all duration-150 hover:border-[#2a4470] hover:bg-[#0d1628]',
          status === 'down' ? 'border-red-500/30 hover:border-red-500/50' : 'border-[#1a2d4a]',
        )}
      >
        <div className="flex items-start justify-between gap-2">
          <div className="flex items-center gap-2 min-w-0">
            <StatusDot status={status} />
            <span className="font-semibold text-sm text-slate-100 truncate leading-tight">
              {monitor.name}
            </span>
          </div>
          <div className="flex items-center gap-2 flex-shrink-0">
            <span
              className={cn(
                'text-xs px-2 py-0.5 rounded-full font-semibold',
                monitor.is_active
                  ? 'bg-cyan-500/10 text-cyan-400'
                  : 'bg-slate-600/20 text-slate-500',
              )}
            >
              {monitor.is_active ? 'Active' : 'Paused'}
            </span>
            <div className="relative" ref={menuRef}>
              <button
                type="button"
                className="p-1 rounded text-slate-600 hover:text-slate-300 hover:bg-[#172240] transition-colors"
                onClick={(e) => {
                  e.stopPropagation()
                  setIsMenuOpen(!isMenuOpen)
                }}
              >
                <MoreVertical className="w-3.5 h-3.5" />
              </button>
              {isMenuOpen && (
                <div className="absolute right-0 top-7 z-50 bg-[#0d1829] border border-[#1e3558] rounded-xl shadow-2xl min-w-[140px] py-1">
                  <button
                    type="button"
                    className="w-full px-3 py-2 text-left text-xs text-slate-300 hover:bg-[#172240] flex items-center gap-2"
                    onClick={(e) => {
                      e.stopPropagation()
                      setIsMenuOpen(false)
                      navigate(`/details/${monitor.id}`)
                    }}
                  >
                    <Eye className="w-3 h-3" /> View
                  </button>
                  <button
                    type="button"
                    disabled={isTogglingActive}
                    className="w-full px-3 py-2 text-left text-xs text-slate-300 hover:bg-[#172240] flex items-center gap-2 disabled:opacity-40"
                    onClick={toggleActive}
                  >
                    {monitor.is_active ? (
                      <>
                        <Pause className="w-3 h-3" /> Pause
                      </>
                    ) : (
                      <>
                        <Play className="w-3 h-3" /> Resume
                      </>
                    )}
                  </button>
                  <button
                    type="button"
                    className="w-full px-3 py-2 text-left text-xs text-red-400 hover:bg-red-500/10 flex items-center gap-2"
                    onClick={(e) => {
                      e.stopPropagation()
                      setIsMenuOpen(false)
                      setIsDeleteMonitorOpen(true)
                    }}
                  >
                    <Trash2 className="w-3 h-3" /> Delete
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <MethodBadge method={monitor.method} />
          <span className="font-mono text-xs text-[#3a5878] truncate">{monitor.url}</span>
        </div>

        <div className="flex items-center justify-end">
          <StatusBadge status={status} />
        </div>

        <div className="flex justify-between text-xs text-[#2e4560]">
          <span>Every {formatFrequency(monitor.frequency)}</span>
        </div>
      </div>

      {isDeleteMonitorOpen && (
        <DeleteMonitor
          monitor={monitor}
          onClose={() => setIsDeleteMonitorOpen(false)}
          onDeleted={onDeleted}
        />
      )}
    </>
  )
}
