import type { InputHTMLAttributes } from 'react'

type InputProps = InputHTMLAttributes<HTMLInputElement> & {
  id: string
  label: string
  error?: string
}

export default function Input({ id, label, error, className, ...props }: InputProps) {
  const hasError = Boolean(error)

  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor={id} className="text-xs font-semibold text-slate-400 tracking-wide">
        {label}
      </label>
      <input
        id={id}
        aria-invalid={hasError}
        aria-describedby={hasError ? `${id}-error` : undefined}
        className={[
          'w-full px-3 py-2 text-sm bg-[#0d1829] border rounded-lg outline-none transition-colors',
          'text-slate-100 placeholder:text-[#3a5070]',
          hasError
            ? 'border-red-500/50'
            : 'border-[#1e3558] focus:border-cyan-500/50 focus:ring-1 focus:ring-cyan-500/15',
          className,
        ].filter(Boolean).join(' ')}
        {...props}
      />
      {error && (
        <p id={`${id}-error`} className="text-xs text-red-400">
          {error}
        </p>
      )}
    </div>
  )
}