import { useState, type FormEvent } from 'react'
import { Navigate, useLocation, useNavigate } from 'react-router-dom'

import { useAuth } from '../auth/AuthContext'

export function LoginPage() {
    const { login, isAuthenticated } = useAuth()
    const navigate = useNavigate()
    const location = useLocation()

    const [organizationId, setOrganizationId] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)

    if (isAuthenticated) {
        return <Navigate to="/dashboard" replace />
    }

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()

        setError('')
        setIsSubmitting(true)

        try {
            await login({
                organization_id: organizationId.trim(),
                email: email.trim(),
                password,
            })

            const from = location.state?.from?.pathname ?? '/dashboard'
            navigate(from, { replace: true })
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Login failed')
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <main>
            <h1>Login</h1>

            <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="organizationId">Organization ID</label>
                    <input
                        id="organizationId"
                        type="text"
                        value={organizationId}
                        onChange={(event) => setOrganizationId(event.target.value)}
                        required
                    />
                </div>

                <div>
                    <label htmlFor="email">Email</label>
                    <input
                        id="email"
                        type="email"
                        value={email}
                        onChange={(event) => setEmail(event.target.value)}
                        required
                    />
                </div>

                <div>
                    <label htmlFor="password">Password</label>
                    <input
                        id="password"
                        type="password"
                        value={password}
                        onChange={(event) => setPassword(event.target.value)}
                        required
                    />
                </div>

                {error && <p role="alert">{error}</p>}

                <button type="submit" disabled={isSubmitting}>
                    {isSubmitting ? 'Signing in...' : 'Sign in'}
                </button>
            </form>
        </main>
    )
}