import { ref, computed, onMounted, watch } from 'vue'

export default function useStockTable() {
  interface Stock {
    ticker: string
    company: string
    brokerage: string
    rating_from: string
    rating_to: string
    target_from: number
    target_to: number
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

  const totalPages = computed(() => Math.ceil(total.value / limit))

  function setSort(column: string) {
    if (sortBy.value === column) {
      sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
    } else {
      sortBy.value = column
      sortDir.value = 'asc'
    }
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

  watch([page, searchTerm, sortBy, sortDir], fetchStocks)
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
    fetchStocks
  }
}
