import type { InputHTMLAttributes } from 'react'

type ToggleProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> & {
  id: string
  label: string
  error?: string
}

export default function Toggle({ id, label, error, className, ...props }: ToggleProps) {
  const hasError = Boolean(error)

  return (
    <div className="flex flex-col gap-1.5">
      <div className="flex items-center justify-between gap-3">
        <span className="text-xs font-semibold text-slate-400 tracking-wide">
          {label}
        </span>
        <label
          htmlFor={id}
          className={[
            'relative inline-flex h-5 w-9 shrink-0 cursor-pointer items-center',
            className,
          ].filter(Boolean).join(' ')}
        >
          <input
            id={id}
            type="checkbox"
            role="switch"
            aria-invalid={hasError}
            aria-describedby={hasError ? `${id}-error` : undefined}
            className="peer sr-only"
            {...props}
          />
          <span
            className={[
              'absolute inset-0 rounded-full transition-colors',
              'bg-[var(--switch-background)] peer-checked:bg-cyan-500',
              'peer-focus-visible:ring-1 peer-focus-visible:ring-cyan-500/15',
              hasError ? 'ring-1 ring-red-500/50' : '',
            ].filter(Boolean).join(' ')}
          />
          <span className="pointer-events-none absolute left-0.5 top-0.5 h-4 w-4 rounded-full bg-white transition-transform peer-checked:translate-x-4" />
        </label>
      </div>
      {error && (
        <p id={`${id}-error`} className="text-xs text-red-400">
          {error}
        </p>
      )}
    </div>
  )
}
