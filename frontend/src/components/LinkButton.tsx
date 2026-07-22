import { Link } from 'react-router-dom'

type LinkButtonProps = {
    children: React.ReactNode
    to?: string
    onClick?: () => void
}

export default function LinkButton({ children, to, onClick }: LinkButtonProps) {
    const className = 'text-cyan-400 hover:text-cyan-300 font-semibold text-base'

    if (to) {
        return (
            <Link to={to} className={className}>
                {children}
            </Link>
        )
    }

    return (
        <button type="button" onClick={onClick} className={className}>
            {children}
        </button>
    )
}
