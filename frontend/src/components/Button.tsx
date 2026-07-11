export default function Button({ children, ...props }: React.ButtonHTMLAttributes<HTMLButtonElement> & { children: React.ReactNode }) {
    return (
        <button className="inline-flex items-center justify-center gap-1.5 font-semibold rounded-lg transition-all duration-150 cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed select-none px-5 py-2.5 text-sm bg-cyan-500 text-black hover:bg-cyan-400 active:scale-[0.98] w-full mt-1" {...props}>
            {children}
        </button>
    )
}
