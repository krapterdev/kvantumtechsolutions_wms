import {
    createContext,
    useContext,
    useEffect,
    useState,
    type ReactNode,
} from 'react'

import { getCurrentUser, login as loginRequest, logout as logoutRequest } from '../api/auth'
import type { AuthUser, LoginRequest } from '../types/auth'

interface AuthContextValue {
    user: AuthUser | null
    isLoading: boolean
    isAuthenticated: boolean
    login: (data: LoginRequest) => Promise<void>
    logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

interface AuthProviderProps {
    children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
    const [user, setUser] = useState<AuthUser | null>(null)
    const [isLoading, setIsLoading] = useState(true)

    useEffect(() => {
        async function loadCurrentUser() {
            try {
                const response = await getCurrentUser()
                setUser(response.user)
            } catch {
                setUser(null)
            } finally {
                setIsLoading(false)
            }
        }

        void loadCurrentUser()
    }, [])

    async function login(data: LoginRequest) {
        await loginRequest(data)

        const response = await getCurrentUser()
        setUser(response.user)
    }

    async function logout() {
        await logoutRequest()
        setUser(null)
    }

    const value: AuthContextValue = {
        user,
        isLoading,
        isAuthenticated: user !== null,
        login,
        logout,
    }

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
    const context = useContext(AuthContext)

    if (!context) {
        throw new Error('useAuth must be used inside AuthProvider')
    }

    return context
}