import { useEffect } from 'react'
import { AlertCircle, CheckCircle, X } from 'lucide-react'

function cn(...classes: (string | false | null | undefined)[]) {
  return classes.filter(Boolean).join(' ')
}

export type ToastType = 'error' | 'success' | 'info'

export default function Toast({
  message,
  type,
  onDismiss,
}: {
  message: string
  type: ToastType
  onDismiss: () => void
}) {
  useEffect(() => {
    const t = setTimeout(onDismiss, 6000)
    return () => clearTimeout(t)
  }, [onDismiss])

  const cfg = {
    error: {
      icon: <AlertCircle className="w-4 h-4 text-red-400 flex-shrink-0" />,
      border: 'border-red-500/30',
    },
    success: {
      icon: <CheckCircle className="w-4 h-4 text-emerald-400 flex-shrink-0" />,
      border: 'border-emerald-500/30',
    },
    info: {
      icon: <AlertCircle className="w-4 h-4 text-cyan-400 flex-shrink-0" />,
      border: 'border-cyan-500/30',
    },
  }

  return (
    <div
      className={cn(
        'fixed bottom-5 right-5 z-50 flex items-center gap-3 bg-[#0d1829] border rounded-xl px-4 py-3 shadow-2xl max-w-xs sm:max-w-sm',
        cfg[type].border,
      )}
    >
      {cfg[type].icon}
      <span className="text-sm text-slate-200 flex-1 leading-snug">{message}</span>
      <button
        type="button"
        onClick={onDismiss}
        className="text-[#3a5070] hover:text-slate-300 ml-1 flex-shrink-0"
      >
        <X className="w-3.5 h-3.5" />
      </button>
    </div>
  )
}
