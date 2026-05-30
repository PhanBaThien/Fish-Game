import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { walletApi } from '../api/wallet'
import { useWalletStore } from '../stores/walletStore'
import Navbar from '../components/Navbar'
import type { Transaction } from '../types'

const PAGE_SIZE = 10

const TX_META = {
  play:     { icon: '🎮', label: 'Chơi game',  color: 'bg-cyan-500/15'   },
  deposit:  { icon: '💰', label: 'Nạp vàng',   color: 'bg-yellow-500/15' },
  withdraw: { icon: '💸', label: 'Rút vàng',   color: 'bg-red-500/15'    },
}

function TxRow({ tx }: { tx: Transaction }) {
  const meta = TX_META[tx.type]
  const isPositive = tx.amount >= 0
  return (
    <div className="flex items-center justify-between px-4 py-3 rounded-xl bg-white/5 border border-white/5 hover:bg-white/8 transition-colors">
      <div className="flex items-center gap-3">
        <div className={`w-9 h-9 rounded-full flex items-center justify-center text-lg shrink-0 ${meta.color}`}>
          {meta.icon}
        </div>
        <div>
          <p className="text-white/85 text-sm font-medium">{meta.label}</p>
          {tx.description && (
            <p className="text-white/35 text-xs mt-0.5 truncate max-w-[200px]">{tx.description}</p>
          )}
          <p className="text-white/25 text-xs mt-0.5">
            {new Date(tx.created_at).toLocaleString('vi-VN', {
              day: '2-digit',
              month: '2-digit',
              year: 'numeric',
              hour: '2-digit',
              minute: '2-digit',
            })}
          </p>
        </div>
      </div>

      <span className={`text-base font-bold tabular-nums ${isPositive ? 'text-yellow-400' : 'text-red-400'}`}>
        {isPositive ? '+' : ''}{tx.amount.toLocaleString()}
      </span>
    </div>
  )
}

export default function WalletPage() {
  const [offset, setOffset] = useState(0)
  const { balance, setBalance } = useWalletStore()

  // fetch wallet balance
  const { isLoading: walletLoading } = useQuery({
    queryKey: ['wallet'],
    queryFn: async () => {
      const w = await walletApi.getWallet()
      setBalance(w.balance)
      return w
    },
    staleTime: 10_000,
  })

  // fetch transactions
  const {
    data: txData,
    isLoading: txLoading,
    isError: txError,
    isFetching,
  } = useQuery({
    queryKey: ['wallet-transactions', offset],
    queryFn: () => walletApi.getTransactions(PAGE_SIZE, offset),
    placeholderData: (prev) => prev,
  })

  const totalPages = txData ? Math.ceil(txData.total / PAGE_SIZE) : 0
  const currentPage = Math.floor(offset / PAGE_SIZE) + 1

  return (
    <div className="min-h-screen bg-slate-900">
      {/* Ambient */}
      <div className="fixed inset-0 pointer-events-none overflow-hidden">
        <div className="absolute top-1/3 left-1/3 w-[500px] h-[500px] bg-yellow-500/4 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 right-1/4 w-[400px] h-[400px] bg-teal-600/5 rounded-full blur-3xl" />
      </div>

      <Navbar />

      <main className="relative pt-24 pb-16 px-4 sm:px-6 lg:px-8 max-w-2xl mx-auto">
        {/* Balance card */}
        <div className="rounded-3xl border border-yellow-500/20 bg-gradient-to-br from-yellow-500/10 via-amber-500/5 to-transparent p-8 mb-8 text-center">
          <div className="text-5xl mb-3">🪙</div>
          <p className="text-white/40 text-sm uppercase tracking-widest mb-2">Số dư hiện tại</p>
          {walletLoading && balance === null ? (
            <div className="h-12 w-40 mx-auto bg-white/10 rounded-xl animate-pulse" />
          ) : (
            <p className="text-5xl font-extrabold text-yellow-400 tabular-nums">
              {(balance ?? 0).toLocaleString()}
            </p>
          )}
          <p className="text-white/25 text-xs mt-3">Vàng · Fish Game</p>
        </div>

        {/* Transaction history */}
        <div>
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-white font-semibold text-lg">Lịch sử giao dịch</h3>
            {txData && (
              <span className="text-white/30 text-xs">{txData.total} giao dịch</span>
            )}
          </div>

          {/* Loading skeleton */}
          {txLoading && !txData && (
            <div className="space-y-2">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="h-16 rounded-xl bg-white/5 animate-pulse" />
              ))}
            </div>
          )}

          {/* Error */}
          {txError && (
            <div className="text-center py-10 text-red-400/70 text-sm">
              Không thể tải lịch sử giao dịch
            </div>
          )}

          {/* Empty */}
          {!txLoading && txData && txData.transactions.length === 0 && (
            <div className="text-center py-16">
              <span className="text-5xl block mb-3">📭</span>
              <p className="text-white/35">Chưa có giao dịch nào</p>
              <p className="text-white/20 text-sm mt-1">Vào phòng bắn cá để kiếm vàng!</p>
            </div>
          )}

          {/* List */}
          {txData && txData.transactions.length > 0 && (
            <>
              <div className={`space-y-2 transition-opacity ${isFetching ? 'opacity-60' : ''}`}>
                {txData.transactions.map((tx) => (
                  <TxRow key={tx.id} tx={tx} />
                ))}
              </div>

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="flex items-center justify-center gap-2 mt-6">
                  <button
                    disabled={currentPage === 1}
                    onClick={() => setOffset((p) => Math.max(0, p - PAGE_SIZE))}
                    className="w-9 h-9 rounded-lg border border-white/10 flex items-center justify-center text-white/50 hover:text-white hover:border-white/30 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                    </svg>
                  </button>

                  <div className="flex items-center gap-1">
                    {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
                      let page: number
                      if (totalPages <= 5) {
                        page = i + 1
                      } else if (currentPage <= 3) {
                        page = i + 1
                      } else if (currentPage >= totalPages - 2) {
                        page = totalPages - 4 + i
                      } else {
                        page = currentPage - 2 + i
                      }
                      return (
                        <button
                          key={page}
                          onClick={() => setOffset((page - 1) * PAGE_SIZE)}
                          className={`w-9 h-9 rounded-lg text-sm font-medium transition-all ${
                            page === currentPage
                              ? 'bg-yellow-500/20 border border-yellow-500/40 text-yellow-400'
                              : 'border border-white/10 text-white/40 hover:text-white hover:border-white/30'
                          }`}
                        >
                          {page}
                        </button>
                      )
                    })}
                  </div>

                  <button
                    disabled={currentPage === totalPages}
                    onClick={() => setOffset((p) => p + PAGE_SIZE)}
                    className="w-9 h-9 rounded-lg border border-white/10 flex items-center justify-center text-white/50 hover:text-white hover:border-white/30 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                  </button>
                </div>
              )}
            </>
          )}
        </div>
      </main>
    </div>
  )
}
