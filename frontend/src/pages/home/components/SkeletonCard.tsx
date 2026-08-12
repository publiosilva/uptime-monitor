export default function SkeletonCard() {
  return (
    <div className="bg-[#0c1422] border border-[#1a2d4a] rounded-xl p-4 flex flex-col gap-3 animate-pulse">
      <div className="flex items-center gap-2">
        <div className="w-2 h-2 rounded-full bg-[#1e3050]" />
        <div className="h-3.5 w-28 bg-[#1e3050] rounded-md" />
        <div className="ml-auto h-5 w-14 bg-[#1e3050] rounded-full" />
      </div>
      <div className="h-3 w-full bg-[#1a2840] rounded-md" />
      <div className="flex justify-end">
        <div className="h-5 w-24 bg-[#1a2840] rounded-full" />
      </div>
      <div className="h-3 w-16 bg-[#172035] rounded-md" />
    </div>
  )
}
