<template>
    <div class="p-4">
      <h2 class="text-2xl font-semibold mb-4">📈 Stocks</h2>
  
      <div v-if="loading" class="text-gray-600">Loading...</div>
  
      <table v-else class="min-w-full bg-white border border-gray-300 shadow-sm rounded-lg overflow-hidden">
        <thead class="bg-gray-100">
          <tr>
            <th class="px-4 py-2 text-left">Ticker</th>
            <th class="px-4 py-2 text-left">Company</th>
            <th class="px-4 py-2 text-left">Brokerage</th>
            <th class="px-4 py-2 text-left">Target</th>
            <th class="px-4 py-2 text-left">Rating</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="stock in stocks" :key="stock.ticker" class="hover:bg-gray-50">
            <td class="px-4 py-2 font-semibold">{{ stock.ticker }}</td>
            <td class="px-4 py-2">{{ stock.company }}</td>
            <td class="px-4 py-2">{{ stock.brokerage }}</td>
            <td class="px-4 py-2">
              {{ stock.target_from }} → {{ stock.target_to }}
            </td>
            <td class="px-4 py-2">
              {{ stock.rating_from }} → {{ stock.rating_to }}
            </td>
          </tr>
        </tbody>
      </table>

      <div class="mt-4 flex justify-center gap-2">
  <button
    class="px-3 py-1 bg-gray-200 rounded hover:bg-gray-300"
    :disabled="page === 1"
    @click="page--"
  >
    ← Prev
  </button>

  <span class="px-3 py-1 text-sm text-gray-700">
    Page {{ page }} of {{ totalPages }}
  </span>

  <button
    class="px-3 py-1 bg-gray-200 rounded hover:bg-gray-300"
    :disabled="page === totalPages"
    @click="page++"
  >
    Next →
  </button>
</div>

    </div>
  </template>
  
  <script setup lang="ts">

import { ref, computed, onMounted, watch } from 'vue'

const stocks = ref<Stock[]>([])
const loading = ref(true)
const page = ref(1)
const limit = 10
const total = ref(0)
const totalPages = computed(() => Math.ceil(total.value / limit))
const API_BASE = import.meta.env.VITE_API_BASE_URL

async function fetchStocks() {
  loading.value = true
  try {
    const response = await fetch(`${API_BASE}/stocks?page=${page.value}&limit=${limit}`)
    const data = await response.json()
    stocks.value = data.items
    total.value = data.total
  } catch (error) {
    console.error('❌ Error fetching stocks:', error)
  } finally {
    loading.value = false
  }
}

onMounted(fetchStocks)
watch(page, fetchStocks)

  </script>
  