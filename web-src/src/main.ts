import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import i18n from './i18n'

// Styles
import './assets/styles/variables.css'
import './assets/styles/base.css'
import './assets/styles/animations.css'

// Create app
const app = createApp(App)

// Use plugins
app.use(createPinia())
app.use(i18n)

// Mount
app.mount('#app')
