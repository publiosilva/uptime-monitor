import { Navigate } from 'react-router-dom'
import { authService } from '../services/auth'

export default function GuestRoute({ children }: { children: React.ReactNode }) {
    if (authService.isAuthenticated()) {
        return <Navigate to="/" replace />
    }

    return children
}
