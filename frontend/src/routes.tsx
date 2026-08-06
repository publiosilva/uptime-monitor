import { Navigate, Route, Routes } from 'react-router-dom'
import AuthListener from './components/AuthListener'
import GuestRoute from './components/GuestRoute'
import ProtectedRoute from './components/ProtectedRoute'
import Home from './pages/home/Home'
import SignIn from './pages/signin/SignIn'
import SignUp from './pages/signup/SignUp'
import Details from './pages/details/Details'

export default function AppRoutes() {
    return (
        <>
            <AuthListener />
            <Routes>
            <Route
                path="/signin"
                element={(
                    <GuestRoute>
                        <SignIn />
                    </GuestRoute>
                )}
            />
            <Route
                path="/signup"
                element={(
                    <GuestRoute>
                        <SignUp />
                    </GuestRoute>
                )}
            />
            <Route
                path="/"
                element={(
                    <ProtectedRoute>
                        <Home />
                    </ProtectedRoute>
                )}
            />
            <Route
                path="/details/:id"
                element={(
                    <ProtectedRoute>
                        <Details />
                    </ProtectedRoute>
                )}
            />
            <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
        </>
    )
}
