import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Activity, LogOut, Plus } from 'lucide-react'
import Logo from '../../components/Logo'
import Button from '../../components/Button'
import AddMonitor from './components/AddMonitor'
import MonitorCard from './components/MonitorCard'
import type { Monitor } from '../../services/monitor'
import { monitorService } from '../../services/monitor'
import { authService } from '../../services/auth'

function cn(...classes: (string | false | null | undefined)[]) {
  return classes.filter(Boolean).join(' ')
}

function StatCard({
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

function SkeletonCard() {
  return (
    <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl p-4 flex flex-col gap-3 animate-pulse">
      <div className="flex items-center gap-2">
        <div className="w-2 h-2 rounded-full bg-[#1e3050]" />
        <div className="h-3.5 w-28 bg-[#1e3050] rounded-md" />
        <div className="ml-auto h-5 w-14 bg-[#1e3050] rounded-full" />
      </div>
      <div className="h-3 w-full bg-[#1a2840] rounded-md" />
      <div className="flex justify-end">
        <div className="h-5 w-24 bg-[#1a2840] rounded-full" />
      </div>
      <div className="h-3 w-16 bg-[#172035] rounded-md" />
    </div>
  )
}

export default function Home() {
  const navigate = useNavigate()
  const [isAddMonitorOpen, setIsAddMonitorOpen] = useState(false)
  const [monitors, setMonitors] = useState<Monitor[]>([])
  const [isLoading, setIsLoading] = useState(true)

  const qtdMonitors = monitors.length
  const qtdMonitorsUp = monitors.filter((monitor) => monitor.is_active && monitor.is_up).length
  const qtdMonitorsDown = monitors.filter((monitor) => monitor.is_active && !monitor.is_up).length

  function loadMonitors() {
    return monitorService.list().then(setMonitors)
  }

  useEffect(() => {
    let cancelled = false
    monitorService
      .list()
      .then((data) => {
        if (!cancelled) setMonitors(data)
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  function handleLogout() {
    authService.logout()
    navigate('/signin')
  }

  return (
    <div className="min-h-screen bg-[#070c18]">
      <header className="sticky top-0 z-40 bg-[#070c18]/90 backdrop-blur border-b border-[#131f38]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between gap-4">
          <div className="flex items-center gap-2.5">
            <div className="scale-[0.78] origin-left">
              <Logo />
            </div>
            <span className="text-base font-semibold text-slate-100 tracking-tight">uptime</span>
            <span className="hidden sm:block text-xs text-[#2e4560] font-mono ml-1">v1.0.0</span>
          </div>
          <div className="flex items-center gap-2">
            <Button size="sm" onClick={() => setIsAddMonitorOpen(true)}>
              <Plus className="w-3.5 h-3.5" />
              <span className="hidden sm:block">Add Monitor</span>
            </Button>
            <button
              type="button"
              onClick={handleLogout}
              className="flex items-center gap-2 px-2.5 py-1.5 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-[#111d33] transition-colors"
            >
              <LogOut className="w-3.5 h-3.5" />
              <span className="hidden sm:block text-xs">Sign out</span>
            </button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 py-6 flex flex-col gap-6">
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
          <StatCard label="Total" value={qtdMonitors} sub="monitors configured" accent="slate" />
          <StatCard label="Operational" value={qtdMonitorsUp} sub="online now" accent="emerald" />
          <StatCard
            label="Down"
            value={qtdMonitorsDown}
            sub={qtdMonitorsDown > 0 ? 'requires attention' : 'all systems go'}
            accent={qtdMonitorsDown > 0 ? 'red' : 'slate'}
          />
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <h2 className="text-sm font-semibold text-slate-300">Monitors</h2>
            <span className="text-xs text-[#2e4560] font-mono">({qtdMonitors})</span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-60" />
              <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-400" />
            </span>
            <span className="text-xs text-[#2e4560]">Live checks running</span>
          </div>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
            {Array.from({ length: 4 }).map((_, i) => (
              <SkeletonCard key={i} />
            ))}
          </div>
        ) : monitors.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 gap-4">
            <div className="w-16 h-16 rounded-2xl bg-[#0c1422] border border-[#1a2d4a] flex items-center justify-center">
              <Activity className="w-7 h-7 text-[#2e4560]" />
            </div>
            <div className="text-center">
              <p className="text-sm font-semibold text-slate-400">No monitors configured</p>
              <p className="text-xs text-[#2e4560] mt-1">
                Add your first URL to start tracking uptime
              </p>
            </div>
            <Button onClick={() => setIsAddMonitorOpen(true)}>
              <Plus className="w-3.5 h-3.5" /> Add your first monitor
            </Button>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
            {monitors.map((monitor) => (
              <MonitorCard
                key={monitor.id}
                monitor={monitor}
                onDeleted={loadMonitors}
                onUpdated={loadMonitors}
              />
            ))}
          </div>
        )}
      </main>

      {isAddMonitorOpen && (
        <AddMonitor
          onClose={() => setIsAddMonitorOpen(false)}
          onCreated={loadMonitors}
        />
      )}
    </div>
  )
}
