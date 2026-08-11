import { useCallback, useState } from 'react'
import { Check, Copy } from 'lucide-react'
import type { MonitorStats } from '../../../services/monitor'
import { MethodBadge } from './Status'
import { cn, formatFrequency } from './helpers'

export default function UrlCard({ monitor }: { monitor: MonitorStats }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = useCallback(() => {
    navigator.clipboard.writeText(monitor.url).catch(() => {})
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }, [monitor.url])

  const meta: [string, string, boolean][] = [
    ['Frequency', formatFrequency(monitor.frequency), true],
    ['Timeout', `${monitor.timeout}s`, true],
    ['Created', monitor.createdAt, false],
    ['Status', monitor.isActive ? 'Active' : 'Paused', false],
  ]

  return (
    <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl p-4">
      <div className="flex items-center gap-3 mb-3">
        <MethodBadge method={monitor.method} />
        <span className="font-mono text-sm text-slate-300 truncate flex-1">{monitor.url}</span>
        <button
          type="button"
          onClick={handleCopy}
          className="flex items-center gap-1.5 text-xs text-[#3a5070] hover:text-slate-300 transition-colors flex-shrink-0 px-2 py-1 rounded hover:bg-[#172240]"
        >
          {copied ? (
            <>
              <Check className="w-3 h-3 text-emerald-400" />
              <span className="text-emerald-400">Copied</span>
            </>
          ) : (
            <>
              <Copy className="w-3 h-3" />
              Copy
            </>
          )}
        </button>
      </div>
      <div className="flex flex-wrap gap-x-5 gap-y-1.5 text-xs">
        {meta.map(([label, value, mono]) => (
          <span key={label} className="text-[#3a5070]">
            {label}:{' '}
            <span className={cn(mono && 'font-mono', 'text-slate-400')}>{value}</span>
          </span>
        ))}
      </div>
    </div>
  )
}
