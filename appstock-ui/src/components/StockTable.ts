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

  const stocks = ref<Stock[]>([])
  const loading = ref(true)
  const page = ref(1)
  const limit = 10
  const total = ref(0)
  const searchTerm = ref('')
  const API_BASE = import.meta.env.VITE_API_BASE_URL

  const sortBy = ref('')
  const sortDir = ref<'asc' | 'desc'>('asc')

  const showModal = ref(false)
  const selectedStock = ref<Stock | null>(null)

  const totalPages = computed(() => Math.ceil(total.value / limit))

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
      console.log('🔎 Items:', data.items)
      total.value = data.total
    } catch (error) {
      console.error('❌ Error fetching stocks:', error)
    } finally {
      loading.value = false
    }
  }

  // watch([page, searchTerm, sortBy, sortDir], fetchStocks)
  // 👇 Reiniciar página al cambiar el término de búsqueda
  watch(searchTerm, () => {
    page.value = 1
    fetchStocks()
  })

  // 👇 Ejecutar fetch cuando cambian página u orden
  watch([page, sortBy, sortDir], fetchStocks)
  
  onMounted(fetchStocks)

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
    closeStockModal
  }
}
