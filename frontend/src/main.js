import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/brutalism.css'

const app = createApp(App)

app.use(router)

app.mount('#app')
