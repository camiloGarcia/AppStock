import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia' // 👈 importa Pinia

const app = createApp(App)

app.use(router)
app.use(createPinia()) // 👈 agrega Pinia al app

app.mount('#app')
