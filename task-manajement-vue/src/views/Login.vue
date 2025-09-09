<script setup>
import { ref } from 'vue'
import axios from 'axios'
import {useRouter} from 'vue-router'
import {useAuthStore} from '../stores/auth'

const email= ref('')
const password = ref('')
const error = ref('')
const router = useRouter()
const authStore = useAuthStore()


const handleLogin = async () => {
    error.value = ''
    try{
        const response = await axios.post('http://localhost:3000/auth/login', {
            email : email.value,
            password : password.value,
        })
        authStore.setToken(response.data.token)
        router.push('/projects')
    } catch(err){
        error.value = "Login Failed, please check your credential"
        console.error(err)
    }
}
</script>

<template>
  <div class="login-container">
    <h1>Login to Task Manager</h1>
    <form @submit.prevent="handleLogin">
      <input type="email" v-model="email" placeholder="Email" required />
      <input type="password" v-model="password" placeholder="Password" required />
      <button type="submit">Login</button>
    </form>
    <p v-if="error" class="error-message">{{ error }}</p>
    <div class="switch-link">
      <p>Don't have an account? <router-link to="/register">Register</router-link></p>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #0f2027, #203a43, #2c5364); /* dark blue gradient */
  font-family: 'Inter', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  padding: 20px;
  color: #eaeaea;
}

.login-container h1 {
  color: #ffffff;
  font-size: 2.2rem;
  margin-bottom: 25px;
  font-weight: 700;
  letter-spacing: 0.5px;
  text-shadow: 0 2px 8px rgba(0,0,0,0.35);
  animation: fadeDown 0.8s ease;
}

form {
  background: rgba(30, 30, 30, 0.75);
  backdrop-filter: blur(14px);
  padding: 35px 28px;
  border-radius: 18px;
  box-shadow: 0 12px 28px rgba(0,0,0,0.45);
  width: 100%;
  max-width: 380px;
  display: flex;
  flex-direction: column;
  animation: fadeUp 0.8s ease;
}

input {
  margin: 12px 0;
  padding: 14px 16px;
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 10px;
  font-size: 15px;
  background: rgba(255,255,255,0.05);
  color: #fff;
  transition: border 0.25s ease, box-shadow 0.25s ease, transform 0.1s ease;
}

input::placeholder {
  color: rgba(255,255,255,0.5);
}

input:focus {
  border-color: #00d4ff;
  box-shadow: 0 0 12px rgba(0, 212, 255, 0.4);
  outline: none;
  transform: scale(1.02);
}

button {
  margin-top: 20px;
  padding: 14px;
  background: linear-gradient(135deg, #00d4ff, #0077ff);
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.35s ease, transform 0.15s ease, box-shadow 0.25s ease;
}

button:hover {
  background: linear-gradient(135deg, #009ecb, #005bbb);
  box-shadow: 0 6px 20px rgba(0,0,0,0.35);
  transform: translateY(-2px);
}

button:active {
  transform: scale(0.97);
}

.error-message {
  margin-top: 18px;
  color: #ff6b6b;
  background: rgba(255, 77, 79, 0.12);
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 14px;
  text-align: center;
  animation: shake 0.4s ease;
}

.switch-link {
  margin-top: 22px;
  color: #bbb;
  font-size: 14px;
}

.switch-link a {
  color: #00d4ff;
  font-weight: 600;
  text-decoration: none;
  transition: color 0.2s ease;
}

.switch-link a:hover {
  color: #00a8cc;
}

/* Animations */
@keyframes fadeUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes fadeDown {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-6px); }
  50% { transform: translateX(6px); }
  75% { transform: translateX(-4px); }
}
</style>
