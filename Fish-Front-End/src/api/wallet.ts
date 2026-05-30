import { apiClient, extractData } from './client'
import type { Wallet, GameSession, TransactionListResponse } from '../types'

export const walletApi = {
  getWallet: async (): Promise<Wallet> => {
    const res = await apiClient.get<{ data: Wallet; error: null }>('/wallet')
    return extractData(res)
  },

  getTransactions: async (limit = 20, offset = 0): Promise<TransactionListResponse> => {
    const res = await apiClient.get<{ data: TransactionListResponse; error: null }>(
      `/wallet/transactions?limit=${limit}&offset=${offset}`,
    )
    return extractData(res)
  },

  deposit: async (amount: number, description?: string): Promise<Wallet> => {
    const res = await apiClient.post<{ data: Wallet; error: null }>('/wallet/deposit', {
      amount,
      description: description ?? null,
    })
    return extractData(res)
  },

  withdraw: async (amount: number, description?: string): Promise<Wallet> => {
    const res = await apiClient.post<{ data: Wallet; error: null }>('/wallet/withdraw', {
      amount,
      description: description ?? null,
    })
    return extractData(res)
  },

  startSession: async (roomId: number): Promise<GameSession> => {
    const res = await apiClient.post<{ data: GameSession; error: null }>('/wallet/session/start', {
      room_id: roomId,
    })
    return extractData(res)
  },

  endSession: async (params: {
    sessionId: number
    shotsFired: number
    fishKilled: number
    totalSpend: number
    totalEarn: number
  }): Promise<{ session: GameSession; wallet: Wallet }> => {
    const res = await apiClient.post<{
      data: { session: GameSession; wallet: Wallet }
      error: null
    }>('/wallet/session/end', {
      session_id:  params.sessionId,
      shots_fired: params.shotsFired,
      fish_killed: params.fishKilled,
      total_spend: params.totalSpend,
      total_earn:  params.totalEarn,
    })
    return extractData(res)
  },
}
