<script setup>
import { RouterView } from 'vue-router';
</script>

<template>
  <div class="main-container">
    <RouterView />
  </div>
</template>

<style>
/* CSS global */
body {
  font-family: sans-serif;
  background-color: #f4f4f4;
  margin: 0;
  padding: 0;
}
.main-container {
  max-width: 800px;
  margin: 40px auto;
  padding: 20px;
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}
h1 {
  text-align: center;
  color: #333;
}
button {
  cursor: pointer;
  background-color: #3b82f6;
  color: white;
  border: none;
  padding: 10px 15px;
  border-radius: 5px;
  transition: background-color 0.3s;
}
button:hover {
  background-color: #2563eb;
}
input,
textarea {
  width: 100%;
  box-sizing: border-box;
  padding: 10px;
  margin: 8px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.project-form, .task-form {
  border: 1px dashed #ccc;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 8px;
}
.project-list, .task-list {
  list-style: none;
  padding: 0;
}
.project-item, .task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #eee;
}
</style>