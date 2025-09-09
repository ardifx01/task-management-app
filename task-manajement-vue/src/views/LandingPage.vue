<script setup>
import { useRouter } from 'vue-router';
const router = useRouter();
</script>

<template>
  <div class="landing-container">
    <h1>Welcome to Task Manager</h1>
    <p>A simple and efficient tool to manage your projects and tasks.</p>
    <div class="cta-buttons">
      <button @click="router.push('/login')">Get Started</button>
    </div>
  </div>
</template>

<style scoped>
.landing-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  min-height: 80vh;
}
.landing-container h1 {
  font-size: 2.5em;
  color: #333;
}
.landing-container p {
  color: #666;
  margin-top: 10px;
  max-width: 600px;
}
.cta-buttons {
  margin-top: 20px;
}
</style>