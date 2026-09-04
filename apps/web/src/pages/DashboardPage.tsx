import { useNavigate } from 'react-router-dom'

import { useAuth } from '../auth/AuthContext'

export function DashboardPage() {
    const { user, logout } = useAuth()
    const navigate = useNavigate()

    async function handleLogout() {
        try {
            await logout()
            navigate('/login', { replace: true })
        } catch {
            // Keep the user on the page if logout request fails.
        }
    }

    return (
        <main>
            <h1>Dashboard</h1>

            <p>Welcome, {user?.first_name}.</p>
            <p>{user?.email}</p>

            <button type="button" onClick={handleLogout}>
                Logout
            </button>
        </main>
    )
}