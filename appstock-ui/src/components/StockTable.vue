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
    <th class="px-4 py-2 text-left cursor-pointer" @click="setSort('action')">
      Action <span v-if="sortBy === 'action'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
    </th>
    <th class="px-4 py-2 text-left cursor-pointer" @click="setSort('target_from')">
      Target From <span v-if="sortBy === 'target_from'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
    </th>
    <th class="px-4 py-2 text-left cursor-pointer" @click="setSort('target_to')">
      Target To <span v-if="sortBy === 'target_to'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
    </th>
    <th class="px-4 py-2 text-left cursor-pointer" @click="setSort('rating_from')">
      Rating From <span v-if="sortBy === 'rating_from'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
    </th>
    <th class="px-4 py-2 text-left cursor-pointer" @click="setSort('rating_to')">
      Rating To <span v-if="sortBy === 'rating_to'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
    </th>
    <th class="px-4 py-2 text-left cursor-pointer" @click="setSort('time')">
      Time <span v-if="sortBy === 'time'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
    </th>
  </tr>
</thead>

  <tbody>
    <tr
      v-for="stock in stocks"
      :key="stock.ticker"
      @click="openStockModal(stock)"
      class="cursor-pointer hover:bg-gray-50"
    >
      <td class="px-4 py-2 font-semibold">{{ stock.ticker }}</td>
      <td class="px-4 py-2">{{ stock.company }}</td>
      <td class="px-4 py-2">{{ stock.brokerage }}</td>
      <td class="px-4 py-2">{{ stock.action }}</td>
      <!-- <td class="px-4 py-2">{{ stock.target_from }}</td> -->
      <td class="px-4 py-2">{{ stock.target_from }}</td>
      <!-- <td class="px-4 py-2">{{ stock.target_to }}</td> -->
      <td class="px-4 py-2">{{ stock.target_to }}</td>
      <td class="px-4 py-2">{{ stock.rating_from }}</td>
      <td class="px-4 py-2">{{ stock.rating_to }}</td>
      <td class="px-4 py-2">{{ stock.time }}</td>
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

    <!-- Modal -->
    <transition name="fade">
      <div
        v-if="showModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50"
      >
        <div class="bg-white p-6 rounded-lg shadow-lg w-[90%] max-w-xl">
          <h3 class="text-xl font-semibold mb-4">📊 Stock Details</h3>
          <p><strong>Ticker:</strong> {{ selectedStock?.ticker }}</p>
          <p><strong>Company:</strong> {{ selectedStock?.company }}</p>
          <p><strong>Brokerage:</strong> {{ selectedStock?.brokerage }}</p>
          <p><strong>Action:</strong> {{ selectedStock?.action }}</p>
          <p><strong>Target from:</strong> {{ selectedStock?.target_from }}</p>
          <p><strong>Target to:</strong> {{ selectedStock?.target_to }}</p>
          <p><strong>Rating from:</strong> {{ selectedStock?.rating_from }}</p>
          <p><strong>Rating to:</strong> {{ selectedStock?.rating_to }}</p>
          <p><strong>Time:</strong> {{ selectedStock?.time }}</p>
          <button
            @click="closeStockModal"
            class="mt-4 px-4 py-2 bg-gray-300 rounded hover:bg-gray-400"
          >
            Close
          </button>
        </div>
      </div>
    </transition>
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
  showModal,
  selectedStock,
  setSort,
  openStockModal,
  closeStockModal
} = useStockTable()
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
