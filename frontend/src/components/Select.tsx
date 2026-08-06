import type { SelectHTMLAttributes } from 'react'
import { ChevronDown } from 'lucide-react'

type SelectProps = SelectHTMLAttributes<HTMLSelectElement> & {
  id: string
  label: string
  error?: string
}

export default function Select({ id, label, error, className, children, ...props }: SelectProps) {
  const hasError = Boolean(error)

  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor={id} className="text-xs font-semibold text-slate-400 tracking-wide">
        {label}
      </label>
      <div className="relative">
        <select
          id={id}
          aria-invalid={hasError}
          aria-describedby={hasError ? `${id}-error` : undefined}
          className={[
            'w-full px-3 py-2 text-sm bg-[#0d1829] border rounded-lg outline-none transition-colors appearance-none cursor-pointer',
            'text-slate-100',
            hasError
              ? 'border-red-500/50'
              : 'border-[#1e3558] focus:border-cyan-500/50 focus:ring-1 focus:ring-cyan-500/15',
            className,
          ].filter(Boolean).join(' ')}
          {...props}
        >
          {children}
        </select>
        <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-[#3a5070] pointer-events-none" />
      </div>
      {error && (
        <p id={`${id}-error`} className="text-xs text-red-400">
          {error}
        </p>
      )}
    </div>
  )
}
