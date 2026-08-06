import { useState } from 'react'
import { X } from 'lucide-react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import Button from '../../../components/Button'
import Input from '../../../components/Input'
import Select from '../../../components/Select'
import Toggle from '../../../components/Toggle'
import { monitorService } from '../../../services/monitor'
import { ApiError } from '../../../lib/api'

const addMonitorSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  url: z.string().url('Invalid URL'),
  method: z.enum(['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD']),
  frequency: z.number().min(30, 'Frequency is required'),
  is_active: z.boolean(),
})

type AddMonitorForm = z.infer<typeof addMonitorSchema>

export default function AddMonitor({
  onClose,
  onCreated,
}: {
  onClose: () => void
  onCreated?: () => void
}) {
  const [submitError, setSubmitError] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<AddMonitorForm>({
    resolver: zodResolver(addMonitorSchema),
    defaultValues: {
      name: '',
      url: '',
      method: 'GET',
      frequency: 60,
      is_active: true,
    },
  })

  async function onSubmit(data: AddMonitorForm) {
    setSubmitError(null)

    try {
      await monitorService.create(data)
      onCreated?.()
      onClose()
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message)
        return
      }

      setSubmitError('Something went wrong. Please try again.')
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative bg-[#0c1422] border border-[#1e3558] rounded-2xl w-full max-w-md shadow-2xl">
        <div className="flex items-center justify-between px-5 py-4 border-b border-[#1a2d4a]">
          <h2 className="text-sm font-semibold text-slate-100">Add Monitor</h2>
          <button
            type="button"
            onClick={onClose}
            className="p-1 rounded text-[#3a5070] hover:text-slate-300 hover:bg-[#172240] transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
        <div className="p-5">
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
            {submitError && (
              <p className="text-xs text-red-400 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2">
                {submitError}
              </p>
            )}
            <Input
              id="name"
              label="Name"
              type="text"
              placeholder="Production API"
              error={errors.name?.message}
              {...register('name')}
            />
            <Input
              id="url"
              label="URL"
              type="text"
              placeholder="https://api.example.com/health"
              error={errors.url?.message}
              {...register('url')}
            />
            <div className="grid grid-cols-2 gap-3">
              <Select
                id="method"
                label="Method"
                error={errors.method?.message}
                {...register('method')}
              >
                <option value="GET">GET</option>
                <option value="POST">POST</option>
                <option value="PUT">PUT</option>
                <option value="PATCH">PATCH</option>
                <option value="DELETE">DELETE</option>
                <option value="HEAD">HEAD</option>
              </Select>
              <Select
                id="frequency"
                label="Frequency"
                error={errors.frequency?.message}
                {...register('frequency', { valueAsNumber: true })}
              >
                <option value="30">Every 30s</option>
                <option value="60">Every 1m</option>
                <option value="300">Every 5m</option>
                <option value="600">Every 10m</option>
                <option value="1800">Every 30m</option>
                <option value="3600">Every 1h</option>
              </Select>
            </div>
            <Toggle
              id="is_active"
              label="Active"
              error={errors.is_active?.message}
              {...register('is_active')}
            />
            <div className="flex gap-3 pt-1">
              <Button type="button" variant="secondary" onClick={onClose} className="flex-1">
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting} className="flex-1">
                {isSubmitting ? 'Creating monitor…' : 'Create Monitor'}
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}
