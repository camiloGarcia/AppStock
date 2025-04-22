import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

type Stock = {
  ticker: string
  company: string
  brokerage: string
  action: string
  target_from: string
  target_to: string
  rating_from: string
  rating_to: string
  time: string
}

export const useStockStore = defineStore('stock', () => {
  const API_BASE = import.meta.env.VITE_API_BASE_URL

  const stocks = ref<Stock[]>([])
  const loading = ref(false)
  const page = ref(1)
  const limit = 10
  const total = ref(0)
  const searchTerm = ref('')
  const sortBy = ref('')
  const sortDir = ref<'asc' | 'desc'>('asc')

  const showModal = ref(false)
  const selectedStock = ref<Stock | null>(null)

  const recommendationDate = ref(new Date().toISOString().substring(0, 10))
  const recommendedStocks = ref<Stock[]>([])
  const showRecommendedModal = ref(false)

  const totalPages = computed(() => Math.ceil(total.value / limit))

  async function fetchStocks() {
    loading.value = true
    try {
      const url = new URL(`${API_BASE}/stocks`)
      url.searchParams.append('page', page.value.toString())
      url.searchParams.append('limit', limit.toString())
      if (searchTerm.value.trim()) {
        url.searchParams.append('search', searchTerm.value.trim())
      }
      if (sortBy.value) {
        url.searchParams.append('sortBy', sortBy.value)
        url.searchParams.append('sortDir', sortDir.value)
      }

      const response = await fetch(url.toString())
      const data = await response.json()
      stocks.value = data.items
      total.value = data.total
    } catch (error) {
      console.error('❌ Error fetching stocks:', error)
    } finally {
      loading.value = false
    }
  }

  async function fetchRecommendationsByDate() {
    try {
      const url = new URL(`${API_BASE}/recommendation`)
      url.searchParams.append('date', recommendationDate.value)

      const response = await fetch(url.toString())
      if (response.ok) {
        recommendedStocks.value = await response.json()
        showRecommendedModal.value = true
      }
    } catch (err) {
      console.error('❌ Error fetching recommendations:', err)
    }
  }

  function setSort(column: string) {
    if (sortBy.value === column) {
      sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
    } else {
      sortBy.value = column
      sortDir.value = 'asc'
    }
  }

  function openStockModal(stock: Stock) {
    selectedStock.value = stock
    showModal.value = true
  }
  
  function closeStockModal() {
    selectedStock.value = null
    showModal.value = false
  }
  

  return {
    stocks,
    loading,
    page,
    limit,
    total,
    searchTerm,
    sortBy,
    sortDir,
    showModal,
    selectedStock,
    recommendationDate,
    recommendedStocks,
    showRecommendedModal,
    totalPages,
    fetchStocks,
    fetchRecommendationsByDate,
    setSort,
    openStockModal, // 👈 agregar
    closeStockModal // 👈 agregar
  }
})
