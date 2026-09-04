import { apiRequest } from './client'
import type { AuthUser, LoginRequest, LoginUser } from '../types/auth'

export function login(data: LoginRequest) {
    return apiRequest<{ user: LoginUser }>('/api/v1/auth/login', {
        method: 'POST',
        body: JSON.stringify(data),
    })
}

export function getCurrentUser() {
    return apiRequest<{ user: AuthUser }>('/api/v1/auth/me')
}

export async function logout() {
    const response = await fetch(
        `${import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:3003'}/api/v1/auth/logout`,
        {
            method: 'POST',
            credentials: 'include',
        },
    )

    if (!response.ok) {
        throw new Error('Logout failed')
    }
}