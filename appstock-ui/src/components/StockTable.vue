<template>
  <div class="p-4">
    <h2 class="text-2xl font-semibold mb-4">📈 Stocks</h2>

    <div class="mb-4">
      <input
        type="text"
        v-model="searchTerm"
        placeholder="Search by Ticker, Company or Brokerage..."
        class="w-full px-4 py-2 border rounded-md shadow-sm"
      />
    </div>

    <div v-if="loading" class="text-gray-600">Loading...</div>
    <table v-else class="min-w-full bg-white border border-gray-300 shadow-sm rounded-lg overflow-hidden">
      <thead class="bg-gray-100">
        <tr>
          <th class="px-4 py-2 text-left cursor-pointer" @click="setSort('ticker')">
            Ticker <span v-if="sortBy === 'ticker'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
          </th>
          <th class="px-4 py-2 text-left cursor-pointer" @click="setSort('company')">
            Company <span v-if="sortBy === 'company'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
          </th>
          <th class="px-4 py-2 text-left cursor-pointer" @click="setSort('brokerage')">
            Brokerage <span v-if="sortBy === 'brokerage'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
          </th>
          <th class="px-4 py-2 text-left">Target</th>
          <th class="px-4 py-2 text-left">Rating</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="stock in stocks" :key="stock.ticker">
          <td class="px-4 py-2 font-semibold">{{ stock.ticker }}</td>
          <td class="px-4 py-2">{{ stock.company }}</td>
          <td class="px-4 py-2">{{ stock.brokerage }}</td>
          <td class="px-4 py-2">{{ stock.target_from }} → {{ stock.target_to }}</td>
          <td class="px-4 py-2">{{ stock.rating_from }} → {{ stock.rating_to }}</td>
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
import useStockTable from './StockTable.ts'
const {
  stocks,
  loading,
  page,
  totalPages,
  searchTerm,
  sortBy,
  sortDir,
  setSort
} = useStockTable()
</script>
