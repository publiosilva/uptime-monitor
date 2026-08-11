import { Clock } from 'lucide-react'
import type { Heartbeat } from '../../../services/monitor'
import { chronologicalHistory, cn } from './helpers'

export default function UptimeTimeline({ history }: { history: Heartbeat[] }) {
  const timeline = chronologicalHistory(history)

  return (
    <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl p-4 sm:p-5">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-semibold text-slate-200 flex items-center gap-2">
          <Clock className="w-3.5 h-3.5 text-cyan-500" />
          Uptime Timeline
        </h3>
        <div className="flex items-center gap-4 text-xs text-[#3a5070]">
          <span className="flex items-center gap-1.5">
            <span className="w-2 h-2 rounded-sm bg-emerald-500/65" />
            Up
          </span>
          <span className="flex items-center gap-1.5">
            <span className="w-2 h-2 rounded-sm bg-red-500/70" />
            Down
          </span>
        </div>
      </div>

      {timeline.length === 0 ? (
        <p className="text-xs text-[#3a5070] py-4 text-center">No checks in the last 24 hours</p>
      ) : (
        <>
          <div className="flex gap-px h-8">
            {timeline.map((heartbeat) => (
              <div
                key={heartbeat.id}
                title={heartbeat.isUp ? 'Operational' : 'Incident'}
                className={cn(
                  'flex-1 rounded-[1px] cursor-default hover:opacity-75 transition-opacity',
                  heartbeat.isUp ? 'bg-emerald-500/65' : 'bg-red-500/70',
                )}
              />
            ))}
          </div>
          <div className="flex justify-between mt-2 text-xs text-[#2a3e58] font-mono">
            <span>24h ago</span>
            <span>now</span>
          </div>
        </>
      )}
    </div>
  )
}
