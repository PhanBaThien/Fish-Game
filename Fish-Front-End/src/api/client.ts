import axios, { AxiosError, type AxiosResponse } from 'axios'

const BASE_URL = '/api/v1'

export const apiClient = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Attach Bearer token from localStorage on each request
apiClient.interceptors.request.use((config) => {
  try {
    const raw = localStorage.getItem('fish-game-auth')
    if (raw) {
      const parsed = JSON.parse(raw)
      const token: string | undefined = parsed?.state?.token
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
  } catch {
    // ignore parse errors
  }
  return config
})

// Normalize error responses into thrown Error objects
apiClient.interceptors.response.use(
  (res: AxiosResponse) => res,
  (error: AxiosError<{ error: { code: string; message: string } | null }>) => {
    const serverMessage = error.response?.data?.error?.message
    const message = serverMessage ?? error.message ?? 'Unknown error'
    return Promise.reject(new Error(message))
  },
)

export function extractData<T>(res: AxiosResponse<{ data: T; error: null }>): T {
  return res.data.data
}
