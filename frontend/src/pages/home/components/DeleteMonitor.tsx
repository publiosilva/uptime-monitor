import { useState } from 'react'
import { X, Trash2 } from 'lucide-react'
import type { Monitor } from '../../../services/monitor'
import { monitorService } from '../../../services/monitor'
import Button from '../../../components/Button'
import { ApiError } from '../../../lib/api'

export default function DeleteMonitor({
  onClose,
  monitor,
  onDeleted,
}: {
  onClose: () => void
  monitor: Monitor
  onDeleted?: () => void
}) {
  const [isDeleting, setIsDeleting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function onDelete() {
    try {
      setIsDeleting(true)
      setError(null)
      await monitorService.delete(monitor.id)
      onDeleted?.()
      onClose()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.')
    } finally {
      setIsDeleting(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative bg-[#0c1422] border border-[#1e3558] rounded-2xl w-full max-w-md shadow-2xl">
        <div className="flex items-center justify-between px-5 py-4 border-b border-[#1a2d4a]">
          <h2 className="text-sm font-semibold text-slate-100">Delete Monitor</h2>
          <button
            type="button"
            onClick={onClose}
            className="p-1 rounded text-[#3a5070] hover:text-slate-300 hover:bg-[#172240] transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
        <div className="p-5 flex flex-col gap-4">
          {error && (
            <p className="text-xs text-red-400 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2">
              {error}
            </p>
          )}
          <p className="text-sm text-slate-400">
            Delete <span className="text-slate-100 font-semibold">{monitor.name}</span>? This cannot be
            undone and all check history will be lost.
          </p>
          <p className="font-mono text-xs text-[#3a5070] bg-[#0d1829] px-3 py-2 rounded-lg border border-[#1e3558] truncate">
            {monitor.url}
          </p>
          <div className="flex gap-3">
            <Button variant="secondary" onClick={onClose} className="flex-1" type="button">
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={onDelete}
              disabled={isDeleting}
              className="flex-1"
              type="button"
            >
              <Trash2 className="w-3.5 h-3.5" />
              {isDeleting ? 'Deleting…' : 'Delete'}
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
