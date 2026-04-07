<template>
  <div class="login-container">
    <div class="card box" style="max-width: 400px; margin: 2rem auto;">
      <h2 style="margin-top: 0;">{{ isLogin ? 'LOGIN' : 'REGISTER' }}</h2>
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>USERNAME</label>
          <input type="text" v-model="username" class="input" required minlength="3" />
        </div>
        
        <div class="form-group" style="margin-top: 1rem;">
          <label>PASSWORD</label>
          <input type="password" v-model="password" class="input" required />
        </div>

        <!-- Dynamic password validations when registering -->
        <Transition name="fade">
          <ul v-if="!isLogin" class="validation-checklist">
            <li :class="{ valid: passHasMinLength }">At least 8 characters</li>
            <li :class="{ valid: passHasUpper }">At least 1 uppercase letter</li>
            <li :class="{ valid: passHasLower }">At least 1 lowercase letter</li>
            <li :class="{ valid: passHasNumber }">At least 1 number</li>
            <li :class="{ valid: passHasSpecial }">At least 1 special character</li>
          </ul>
        </Transition>

        <div style="margin-top: 2rem;">
          <button type="submit" class="btn" style="width: 100%;" :disabled="!isLogin && !isPasswordValid">
            {{ isLogin ? 'LOGIN' : 'REGISTER' }}
          </button>
        </div>
      </form>
      
      <div style="margin-top: 1rem; text-align: center;">
        <button class="btn" style="background: none; border: none; text-decoration: underline; box-shadow: none; padding: 0.5rem;" @click="toggleMode">
          {{ isLogin ? 'Need an account? Register' : 'Already have an account? Login' }}
        </button>
      </div>

      <Transition name="fade">
        <div v-if="error" style="color: red; margin-top: 1rem; border: 2px solid red; padding: 0.5rem; background: white;">
          {{ error }}
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const isLogin = ref(true)
const username = ref('')
const password = ref('')
const error = ref('')

// Password validation rules
const passHasMinLength = computed(() => password.value.length >= 8)
const passHasUpper = computed(() => /[A-Z]/.test(password.value))
const passHasLower = computed(() => /[a-z]/.test(password.value))
const passHasNumber = computed(() => /[0-9]/.test(password.value))
const passHasSpecial = computed(() => /[!@#$%^&*(),.?":{}|<>]/.test(password.value))

const isPasswordValid = computed(() => {
  return passHasMinLength.value &&
         passHasUpper.value &&
         passHasLower.value &&
         passHasNumber.value &&
         passHasSpecial.value
})

const toggleMode = () => {
  isLogin.value = !isLogin.value
  error.value = ''
  password.value = ''
}

const handleSubmit = async () => {
  error.value = ''
  
  if (!isLogin.value && !isPasswordValid.value) {
    error.value = 'Password does not meet all criteria'
    return
  }

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
      setTimeout(() => { window.location.reload() }, 100)
    } else {
      isLogin.value = true
      error.value = 'Registration successful. Please login.'
      password.value = ''
    }
  } catch (err) {
    error.value = err.response?.data || 'An error occurred'
  }
}
</script>
