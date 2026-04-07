<template>
  <div class="login-container">
    <div class="card" style="max-width: 400px; margin: 2rem auto;">
      <h2 style="margin-top: 0;">{{ isLogin ? 'LOGIN' : 'REGISTER' }}</h2>
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>USERNAME</label>
          <input type="text" v-model="username" class="input" required minlength="3" />
        </div>
        <div class="form-group" style="margin-top: 1rem;">
          <label>PASSWORD</label>
          <input type="password" v-model="password" class="input" required minlength="6" />
        </div>
        <div style="margin-top: 2rem;">
          <button type="submit" class="btn" style="width: 100%;">{{ isLogin ? 'LOGIN' : 'REGISTER' }}</button>
        </div>
      </form>
      <div style="margin-top: 1rem; text-align: center;">
        <button class="btn" style="background: none; border: none; text-decoration: underline;" @click="isLogin = !isLogin">
          {{ isLogin ? 'Need an account? Register' : 'Already have an account? Login' }}
        </button>
      </div>
      <div v-if="error" style="color: red; margin-top: 1rem; border: 2px solid red; padding: 0.5rem;">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const isLogin = ref(true)
const username = ref('')
const password = ref('')
const error = ref('')

const handleSubmit = async () => {
  error.value = ''
  try {
    const endpoint = isLogin.value ? '/api/auth/login' : '/api/auth/register'
    const res = await axios.post(import.meta.env.VITE_API_URL || 'http://localhost:8080' + endpoint, {
      username: username.value,
      password: password.value
    })
    
    if (isLogin.value) {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('role', res.data.role)
      router.push('/')
      // reload to update App.vue state if necessary, or just push
      setTimeout(() => { window.location.reload() }, 100)
    } else {
      isLogin.value = true
      error.value = 'Registration successful. Please login.'
    }
  } catch (err) {
    error.value = err.response?.data || 'An error occurred'
  }
}
</script>
