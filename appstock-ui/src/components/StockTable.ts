import { ref, computed, onMounted, watch } from 'vue'

export default function useStockTable() {
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

  const API_BASE = import.meta.env.VITE_API_BASE_URL

  // Estado principal
  const stocks = ref<Stock[]>([])
  const loading = ref(true)
  const page = ref(1)
  const limit = 10
  const total = ref(0)
  const totalPages = computed(() => Math.ceil(total.value / limit))

  // Filtros y ordenamiento
  const searchTerm = ref('')
  const sortBy = ref('')
  const sortDir = ref<'asc' | 'desc'>('asc')

  // Modal de detalle
  const showModal = ref(false)
  const selectedStock = ref<Stock | null>(null)

  // Modal de recomendación por fecha
  const recommendationDate = ref(new Date().toISOString().substring(0, 10))
  const recommendedStocks = ref<Stock[]>([])
  const showRecommendedModal = ref(false)

  // Modal (única recomendación - no usada en este flujo)
  const recommendedStock = ref<Stock | null>(null)
  const showRecommendation = ref(false)

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

  async function fetchRecommendation() {
    try {
      const response = await fetch(`${API_BASE}/recommendation`)
      if (response.ok) {
        recommendedStock.value = await response.json()
      }
    } catch (err) {
      console.error('❌ Error fetching recommendation:', err)
    }
  }

  function toggleRecommendation() {
    showRecommendation.value = !showRecommendation.value
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

  // Reacción a cambios de filtros
  watch(searchTerm, () => {
    page.value = 1
    fetchStocks()
  })

  watch([page, sortBy, sortDir], fetchStocks)

  onMounted(() => {
    fetchStocks()
    //  fetchRecommendation()
  })

  return {
    stocks,
    loading,
    page,
    totalPages,
    searchTerm,
    sortBy,
    sortDir,
    setSort,
    fetchStocks,
    showModal,
    selectedStock,
    openStockModal,
    closeStockModal,
    recommendedStock,
    fetchRecommendation,
    showRecommendation,
    toggleRecommendation,
    recommendationDate,
    fetchRecommendationsByDate,
    recommendedStocks,
    showRecommendedModal,
  }
}
