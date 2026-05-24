export interface User {
  id: number
  username: string
  email: string
  role_id: number
  created_at: string
  updated_at: string
}

export interface Room {
  id: number
  name: string
  min_bet: number
  max_players: number
  description: {
    String: string
    Valid: boolean
  }
  created_at: string
  updated_at: string
}

export interface Fish {
  id: number
  name: string
  health: number
  reward_multiplier: number
  speed: number
  asset_path: string
  created_at: string
  updated_at: string
}

export interface ApiSuccess<T> {
  data: T
  error: null
}

export interface ApiError {
  data: null
  error: {
    code: string
    message: string
  }
}

export type ApiResponse<T> = ApiSuccess<T> | ApiError

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginResponse {
  access_token: string
  access_token_expires_at: number
  user: User
}

export interface RefreshTokenResponse {
  access_token: string
  access_token_expires_at: number
}

export interface RegisterResponse {
  id: number
  username: string
  role_id: number
}
