import { apiClient, extractData } from './client'
import type { Room } from '../types'

export const roomsApi = {
  list: async (): Promise<Room[]> => {
    const res = await apiClient.get<{ data: Room[]; error: null }>('/rooms')
    return extractData(res)
  },

  get: async (id: number): Promise<Room> => {
    const res = await apiClient.get<{ data: Room; error: null }>(`/rooms/${id}`)
    return extractData(res)
  },

  create: async (data: Partial<Room>): Promise<Room> => {
    const res = await apiClient.post<{ data: Room; error: null }>('/rooms', data)
    return extractData(res)
  },

  update: async (id: number, data: Partial<Room>): Promise<Room> => {
    const res = await apiClient.put<{ data: Room; error: null }>(`/rooms/${id}`, data)
    return extractData(res)
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/rooms/${id}`)
  },
}
