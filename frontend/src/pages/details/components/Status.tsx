import { cn, type Status } from './helpers'
import type { MonitorStats } from '../../../services/monitor'

export function StatusDot({ status }: { status: Status }) {
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

export function StatusBadge({ status }: { status: Status }) {
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

export function MethodBadge({ method }: { method: MonitorStats['method'] }) {
  const colors: Record<MonitorStats['method'], string> = {
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

export function StatCard({
  label,
  value,
  sub,
  accent = 'slate',
}: {
  label: string
  value: string | number
  sub?: string
  accent?: 'cyan' | 'emerald' | 'red' | 'slate'
}) {
  const accentCls = {
    cyan: 'text-cyan-400',
    emerald: 'text-emerald-400',
    red: 'text-red-400',
    slate: 'text-slate-300',
  }[accent]

  return (
    <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl p-4 flex flex-col gap-1">
      <span className="text-xs text-[#4a6080] font-semibold uppercase tracking-widest">{label}</span>
      <span className={cn('text-2xl font-mono font-semibold leading-tight', accentCls)}>{value}</span>
      {sub && <span className="text-xs text-[#3a5070]">{sub}</span>}
    </div>
  )
}
