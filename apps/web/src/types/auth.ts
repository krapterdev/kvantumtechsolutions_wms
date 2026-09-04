export interface AuthUser {
    id: string
    organization_id: string
    email: string
    first_name: string
    last_name?: string
    is_active: boolean
}

export interface LoginUser {
    id: string
    organization_id: string
    email: string
    first_name: string
}

export interface LoginRequest {
    organization_id: string
    email: string
    password: string
}