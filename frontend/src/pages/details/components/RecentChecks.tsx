import { CheckCircle, XCircle } from 'lucide-react'
import type { Heartbeat } from '../../../services/monitor'
import { formatCheckTime } from './helpers'

export default function RecentChecks({ history }: { history: Heartbeat[] }) {
  return (
    <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl overflow-hidden">
      <div className="px-4 py-3 border-b border-[#1a2d4a] flex items-center justify-between">
        <h3 className="text-sm font-semibold text-slate-200">Recent Checks</h3>
        <span className="text-xs text-[#3a5070] font-mono">{history.length} results</span>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-xs">
          <thead>
            <tr className="border-b border-[#131f38]">
              {['Time', 'Status', 'Code', 'Latency', 'Error'].map((header) => (
                <th key={header} className="px-4 py-2.5 text-left text-[#3a5070] font-semibold">
                  {header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {history.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-4 py-6 text-center text-[#3a5070]">
                  No checks yet
                </td>
              </tr>
            ) : (
              history.map((check) => (
                <tr
                  key={check.id}
                  className="border-b border-[#0f1c2e] hover:bg-[#0d1829] transition-colors"
                >
                  <td className="px-4 py-2.5 font-mono text-[#2e4560]">
                    {formatCheckTime(check.createdAt)}
                  </td>
                  <td className="px-4 py-2.5">
                    {check.isUp ? (
                      <span className="flex items-center gap-1.5 text-emerald-400">
                        <CheckCircle className="w-3 h-3" />
                        Up
                      </span>
                    ) : (
                      <span className="flex items-center gap-1.5 text-red-400">
                        <XCircle className="w-3 h-3" />
                        Down
                      </span>
                    )}
                  </td>
                  <td className="px-4 py-2.5 font-mono">
                    <span className={check.statusCode >= 200 && check.statusCode < 400 ? 'text-emerald-400' : 'text-red-400'}>
                      {check.statusCode || '—'}
                    </span>
                  </td>
                  <td className="px-4 py-2.5 font-mono text-slate-400">{check.latencyMs}ms</td>
                  <td className="px-4 py-2.5 text-[#3a5070]">
                    {check.isUp ? '—' : 'Service Unavailable'}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
