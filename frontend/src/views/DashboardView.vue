<template>
  <div class="dashboard">
    <h1 style="border-bottom: 4px solid black; padding-bottom: 0.5rem;">SYSTEM DASHBOARD</h1>
    
    <div v-if="error" style="color: red; border: 3px solid red; padding: 1rem; margin-bottom: 1rem;">
      {{ error }}
    </div>

    <div v-if="loading" class="card">
      LOADING METRICS...
    </div>
    
    <div v-else>
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1rem;">
        <div v-for="metric in metrics" :key="metric.time + metric.node_name" class="card">
          <h3 style="margin-top: 0; display: flex; justify-content: space-between;">
            <span>{{ metric.node_name }}</span>
            <span style="font-size: 0.8rem; background: black; color: white; padding: 0.2rem 0.5rem;">{{ new Date(metric.time).toLocaleTimeString() }}</span>
          </h3>
          <div style="margin-top: 1rem;">
            <div style="margin-bottom: 0.5rem;">
              <strong>CPU:</strong> {{ metric.cpu_usage.toFixed(1) }}%
              <div style="height: 10px; background: #ddd; border: 2px solid black; margin-top: 0.2rem;">
                <div :style="`width: ${metric.cpu_usage}%; height: 100%; background: black;`"></div>
              </div>
            </div>
            <div style="margin-bottom: 0.5rem;">
              <strong>MEM:</strong> {{ metric.memory_usage.toFixed(1) }}%
              <div style="height: 10px; background: #ddd; border: 2px solid black; margin-top: 0.2rem;">
                <div :style="`width: ${metric.memory_usage}%; height: 100%; background: black;`"></div>
              </div>
            </div>
            <div>
              <strong>DSK:</strong> {{ metric.disk_usage.toFixed(1) }}%
              <div style="height: 10px; background: #ddd; border: 2px solid black; margin-top: 0.2rem;">
                <div :style="`width: ${metric.disk_usage}%; height: 100%; background: black;`"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-if="metrics.length === 0" class="card" style="text-align: center; font-weight: bold;">
        NO METRICS DATA AVAILABLE
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const metrics = ref([])
const error = ref('')
const loading = ref(true)
let interval

const fetchMetrics = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }
    
    const res = await axios.get(import.meta.env.VITE_API_URL || 'http://localhost:8080/api/metrics', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    metrics.value = res.data
    error.value = ''
  } catch (err) {
    if (err.response?.status === 401 || err.response?.status === 403) {
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      router.push('/login')
    } else {
      error.value = 'Failed to fetch metrics: ' + (err.response?.data || err.message)
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchMetrics()
  interval = setInterval(fetchMetrics, 5000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})
</script>
