function cn(...classes: (string | false | null | undefined)[]) {
  return classes.filter(Boolean).join(' ')
}

export default function StatCard({
  label,
  value,
  sub,
  accent = 'slate',
}: {
  label: string
  value: string | number
  sub?: string
  accent?: 'cyan' | 'emerald' | 'red' | 'slate'
}) {
  const accentCls = {
    cyan: 'text-cyan-400',
    emerald: 'text-emerald-400',
    red: 'text-red-400',
    slate: 'text-slate-300',
  }[accent]

  return (
    <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl p-4 flex flex-col gap-1">
      <span className="text-xs text-[#4a6080] font-semibold uppercase tracking-widest">{label}</span>
      <span className={cn('text-2xl font-mono font-semibold leading-tight', accentCls)}>{value}</span>
      {sub && <span className="text-xs text-[#3a5070]">{sub}</span>}
    </div>
  )
}
