<template>
  <div class="p-4">
    <h2 class="text-2xl font-semibold mb-4">📈 Stocks</h2>

    <!-- Buscador -->
    <div class="mb-4">
      <input
        type="text"
        v-model="searchTerm"
        placeholder="Search by Ticker, Company or Brokerage..."
        class="w-full px-4 py-2 border rounded-md shadow-sm"
      />
    </div>

    <!-- Input Fecha -->
    <div class="mb-4">
      <label for="recommendation-date" class="block mb-1 text-sm font-medium text-gray-700">
        Recommendation Date:
      </label>
      <input
        id="recommendation-date"
        type="date"
        v-model="recommendationDate"
        class="px-3 py-2 border rounded w-full max-w-xs shadow-sm"
      />
    </div>

    <!-- Botón Mostrar Recomendaciones -->
    <button
      @click="fetchRecommendationsByDate"
      class="mb-4 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition"
    >
      Show Recommendations for Selected Date
    </button>

    <!-- Cargando -->
    <div v-if="loading" class="text-gray-600">Loading...</div>

    <!-- Tabla -->
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
          <td class="px-4 py-2">{{ stock.target_from }}</td>
          <td class="px-4 py-2">{{ stock.target_to }}</td>
          <td class="px-4 py-2">{{ stock.rating_from }}</td>
          <td class="px-4 py-2">{{ stock.rating_to }}</td>
          <!-- <td class="px-4 py-2">{{ stock.time }}</td> -->
          <td class="px-4 py-2">
  {{
    new Intl.DateTimeFormat("es-CO", {
      dateStyle: "medium",
      timeStyle: "medium",
      timeZone: "America/Bogota"
    }).format(new Date(stock.time))
  }}
</td>

        </tr>
      </tbody>
    </table>

    <!-- Paginación -->
    <div class="mt-4 flex justify-center gap-2">
      <button
        class="px-3 py-1 bg-gray-200 rounded hover:bg-gray-300"
        :disabled="page === 1"
        @click="page--"
      >
        ← Prev
      </button>
      <span class="px-3 py-1 text-sm text-gray-700">Page {{ page }} of {{ totalPages }}</span>
      <button
        class="px-3 py-1 bg-gray-200 rounded hover:bg-gray-300"
        :disabled="page === totalPages"
        @click="page++"
      >
        Next →
      </button>
    </div>

    <!-- Modal Detalle Stock -->
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
          <!-- <p><strong>Time:</strong> {{ selectedStock?.time }}</p> -->
          <p><strong>Time:</strong>
  {{
    selectedStock?.time
      ? new Intl.DateTimeFormat("es-CO", {
          dateStyle: "medium",
          timeStyle: "medium",
          timeZone: "America/Bogota"
        }).format(new Date(selectedStock.time))
      : ''
  }}
</p>

          <button
            @click="closeStockModal"
            class="mt-4 px-4 py-2 bg-gray-300 rounded hover:bg-gray-400"
          >
            Close
          </button>
        </div>
      </div>
    </transition>

    <!-- Modal Recomendaciones -->
    <transition name="fade">
  <div
    v-if="showRecommendedModal"
    class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50"
  >
  
   <!-- ❌ Botón fuera del contenedor blanco, arriba a la derecha -->
   <button
      @click="showRecommendedModal = false"
      class="absolute top-4 right-4 text-white hover:text-red-500 transition z-50"
      aria-label="Close"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none"
        viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>

  <div class="bg-white p-6 rounded-lg shadow-lg w-[90%] max-w-xl max-h-[80vh] overflow-y-auto relative">


      <h3 class="text-xl font-semibold mb-4 text-green-600">
        📋 Recommendations for {{ recommendationDate }}
      </h3>

      <div v-if="recommendedStocks && recommendedStocks.length > 0">

        <ul class="space-y-2">
          <li v-for="stock in recommendedStocks" :key="stock.ticker" class="border-b pb-2">
            <p><strong>Ticker:</strong> {{ stock.ticker }}</p>
            <p><strong>Company:</strong> {{ stock.company }}</p>
            <p><strong>Brokerage:</strong> {{ stock.brokerage }}</p>
            <p><strong>Action:</strong> {{ stock.action }}</p>
            <p><strong>Target:</strong> {{ stock.target_from }} → {{ stock.target_to }}</p>
            <p><strong>Rating:</strong> {{ stock.rating_from }} → {{ stock.rating_to }}</p>
            <!-- <p><strong>Time:</strong> {{ new Date(stock.time).toLocaleString() }}</p> -->
            <p><strong>Time:</strong> {{
  new Intl.DateTimeFormat("es-CO", {
    dateStyle: "medium",
    timeStyle: "medium",
    timeZone: "America/Bogota"
  }).format(new Date(stock.time))
}}</p>

          </li>
        </ul>
      </div>

      <div v-else class="text-yellow-600 font-semibold">
        ⚠️ No recommendations found for the selected date.
      </div>

      <button
        @click="showRecommendedModal = false"
        class="mt-6 px-4 py-2 bg-gray-300 rounded hover:bg-gray-400"
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
  closeStockModal,
  recommendationDate,
  recommendedStocks,
  fetchRecommendationsByDate,
  showRecommendedModal
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
