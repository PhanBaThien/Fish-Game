import { apiClient, extractData } from './client'
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, User } from '../types'

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const res = await apiClient.post<{ data: LoginResponse; error: null }>('/auth/login', data)
    return extractData(res)
  },

  register: async (data: RegisterRequest): Promise<RegisterResponse> => {
    const res = await apiClient.post<{ data: RegisterResponse; error: null }>('/auth/register', data)
    return extractData(res)
  },

  me: async (): Promise<User> => {
    const res = await apiClient.get<{ data: User; error: null }>('/auth/me')
    return extractData(res)
  },
}
