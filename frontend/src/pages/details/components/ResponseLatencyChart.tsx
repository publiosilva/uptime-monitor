import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
} from 'recharts'
import type { TooltipContentProps } from 'recharts'
import { Zap } from 'lucide-react'
import type { Heartbeat } from '../../../services/monitor'
import { chronologicalHistory } from './helpers'

type ChartPoint = {
  time: string
  ms: number
}

function formatTickTime(iso: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
}

function toChartData(history: Heartbeat[]): ChartPoint[] {
  return chronologicalHistory(history).map((heartbeat) => ({
    time: formatTickTime(heartbeat.createdAt),
    ms: heartbeat.latencyMs,
  }))
}

function ChartTooltip({ active, payload, label }: TooltipContentProps) {
  if (!active || !payload?.length) return null

  return (
    <div className="bg-[#0d1829] border border-[#1e3558] rounded-lg px-3 py-2 shadow-2xl">
      <p className="text-xs text-[#3a5070] font-mono mb-0.5">{String(label ?? '')}</p>
      <p className="text-sm font-mono font-semibold text-cyan-400">{payload[0].value}ms</p>
    </div>
  )
}

export default function ResponseLatencyChart({ history }: { history: Heartbeat[] }) {
  const data = toChartData(history)

  return (
    <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl p-4 sm:p-5">
      <div className="flex items-center justify-between mb-5">
        <h3 className="text-sm font-semibold text-slate-200 flex items-center gap-2">
          <Zap className="w-3.5 h-3.5 text-cyan-500" />
          Response Latency
        </h3>
        <span className="text-xs text-[#3a5070] font-mono">Last 24 hours</span>
      </div>

      {data.length === 0 ? (
        <p className="text-xs text-[#3a5070] py-10 text-center">No latency data yet</p>
      ) : (
        <ResponsiveContainer width="100%" height={150}>
          <LineChart data={data} margin={{ top: 4, right: 4, left: -24, bottom: 0 }}>
            <CartesianGrid strokeDasharray="2 4" stroke="#131f38" vertical={false} />
            <XAxis
              dataKey="time"
              tick={{ fill: '#2e4560', fontSize: 9, fontFamily: "'JetBrains Mono', monospace" }}
              tickLine={false}
              axisLine={false}
              minTickGap={40}
            />
            <YAxis
              tick={{ fill: '#2e4560', fontSize: 9, fontFamily: "'JetBrains Mono', monospace" }}
              tickLine={false}
              axisLine={false}
              unit="ms"
            />
            <Tooltip content={ChartTooltip} />
            <Line
              type="monotone"
              dataKey="ms"
              stroke="#06b6d4"
              strokeWidth={1.5}
              dot={false}
              activeDot={{ r: 3, fill: '#06b6d4', strokeWidth: 0 }}
            />
          </LineChart>
        </ResponsiveContainer>
      )}
    </div>
  )
}
