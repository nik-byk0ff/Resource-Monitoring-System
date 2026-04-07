<template>
  <nav class="navbar" v-if="isAuthenticated">
    <div class="nav-brand">SYSTEM_MONITOR</div>
    <div>
      <span style="margin-right: 1rem;">Role: {{ userRole }}</span>
      <button class="btn" @click="logout" style="padding: 0.5rem 1rem; font-size: 1rem;">LOGOUT</button>
    </div>
  </nav>
  <main class="container">
    <router-view />
  </main>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const isAuthenticated = computed(() => {
  return !!localStorage.getItem('token')
})

const userRole = computed(() => {
  return localStorage.getItem('role') || 'unknown'
})

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('role')
  router.push('/login')
}
</script>
