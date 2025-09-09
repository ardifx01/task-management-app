<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import LoadingSpinner from '../components/LoadingSpinner.vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const projectId = route.params.id;

const tasks = ref([]);
const newTaskTitle = ref('');
const newTaskDescription = ref('');
const newTaskPriority = ref('low');
const isLoading = ref(true);

const editingTask = ref(null); // <-- Variabel untuk menyimpan task yang sedang diedit
const editedTitle = ref('');
const editedDescription = ref('');
const editedStatus = ref('');
const editedPriority = ref('');

const fetchTasks = async () => {
    isLoading.value = true;
    try {
        const response = await axios.get(`http://localhost:3000/projects/${projectId}/tasks`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        });
        tasks.value = Array.isArray(response.data) ? response.data : [];
    } catch (err) {
        console.error('Failed to fetch tasks:', err);
        if (err.response && err.response.status === 401) {
            authStore.clearToken();
            router.push('/login');
        }
    } finally {
        isLoading.value = false;
    }
};

const createTask = async () => {
    if (!newTaskTitle.value) return;
    try {
        await axios.post(
            `http://localhost:3000/projects/${projectId}/tasks`,
            {
                title: newTaskTitle.value,
                description: newTaskDescription.value,
                priority: newTaskPriority.value,
            },
            {
                headers: { Authorization: `Bearer ${authStore.token}` },
            }
        );
        newTaskTitle.value = '';
        newTaskDescription.value = '';
        newTaskPriority.value = 'low';
        fetchTasks();
    } catch (err) {
        console.error('Failed to create task:', err);
    }
};

// Fungsi untuk memulai edit
const startEdit = (task) => {
    editingTask.value = task;
    editedTitle.value = task.title;
    editedDescription.value = task.description;
    editedStatus.value = task.status;
    editedPriority.value = task.priority;
};

// Fungsi untuk membatalkan edit
const cancelEdit = () => {
    editingTask.value = null;
};

// Fungsi untuk menyimpan perubahan
const saveEdit = async () => {
    if (!editingTask.value) return;
    try {
        await axios.put(
            `http://localhost:3000/tasks/${editingTask.value.id}`,
            {
                title: editedTitle.value,
                description: editedDescription.value,
                status: editedStatus.value,
                priority: editedPriority.value,
            },
            {
                headers: { Authorization: `Bearer ${authStore.token}` },
            }
        );
        editingTask.value = null;
        fetchTasks();
    } catch (err) {
        console.error('Failed to update task:', err);
    }
};

const updateTaskStatus = async (taskId, newStatus) => {
    try {
        await axios.put(
            `http://localhost:3000/tasks/${taskId}`,
            { status: newStatus },
            {
                headers: { Authorization: `Bearer ${authStore.token}` },
            }
        );
        fetchTasks();
    } catch (err) {
        console.error('Failed to update task:', err);
    }
};

const deleteTask = async (taskId) => {
    try {
        await axios.delete(`http://localhost:3000/tasks/${taskId}`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        });
        fetchTasks();
    } catch (err) {
        console.error('Failed to delete task:', err);
    }
};

onMounted(() => {
    fetchTasks();
});
</script>

<template>
  <div class="tasks-container">
    <div class="header">
      <h1>Tasks</h1>
      <button @click="router.push('/projects')">Back to Projects</button>
    </div>

    <div class="task-form">
      <input v-model="newTaskTitle" placeholder="Task Title" />
      <textarea v-model="newTaskDescription" placeholder="Task Description"></textarea>
      <label>
        Priority:
        <select v-model="newTaskPriority">
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
        </select>
      </label>
      <button @click="createTask">Add Task</button>
    </div>

    <LoadingSpinner v-if="isLoading" />

    <div v-else-if="tasks.length" class="task-list">
      <ul>
        <li v-for="task in tasks" :key="task.id" class="task-item">
          <div v-if="editingTask && editingTask.id === task.id" class="edit-form">
            <input v-model="editedTitle" placeholder="Title" />
            <textarea v-model="editedDescription" placeholder="Description"></textarea>
            <div class="edit-options">
              <label>
                Status:
                <select v-model="editedStatus">
                  <option value="todo">Todo</option>
                  <option value="in-progress">In-Progress</option>
                  <option value="done">Done</option>
                </select>
              </label>
              <label>
                Priority:
                <select v-model="editedPriority">
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                </select>
              </label>
            </div>
            <div class="edit-buttons">
              <button @click="saveEdit">Save</button>
              <button @click="cancelEdit">Cancel</button>
            </div>
          </div>
          <div v-else class="task-info">
            <h3>{{ task.title }}</h3>
            <p>{{ task.description }}</p>
           <div class="badges">
            <span :class="['priority', task.priority]">{{ task.priority }}</span>
            <span :class="['status', task.status]">{{ task.status }}</span>
            </div>

            <div class="progress-bar">
          <div :class="['progress-fill', task.status]"></div>
              </div>
          </div>
          <div v-if="!editingTask || editingTask.id !== task.id" class="task-actions">
            <button v-if="task.status === 'todo'" @click="updateTaskStatus(task.id, 'in-progress')">Start</button>
            <button v-if="task.status === 'in-progress'" @click="updateTaskStatus(task.id, 'done')">Complete</button>
            <button @click="startEdit(task)">Edit</button>
            <button @click="deleteTask(task.id)">Delete</button>
          </div>
        </li>
      </ul>
    </div>
    <p v-else>This project doesn't have any tasks yet.</p>
  </div>
</template>

<style scoped>
.tasks-container {
  min-height: 100vh;
  padding: 40px 20px;
background: linear-gradient(135deg, #0f2027, #203a43, #2c5364);
  font-family: 'Inter', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  display: flex;
  flex-direction: column;
  align-items: center;
}

/* Header */
.header {
  position: relative;
  width: 100%;
  max-width: 900px;
  display: flex;
  justify-content: flex-end; /* fokus ke button */
  align-items: center;
  margin-bottom: 30px;
}

.header h1 {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  color: #f5f5f5;
  font-size: 2.2rem;
  font-weight: 700;
  text-align: center;
}

.header button {
  background: linear-gradient(135deg, #ff416c, #ff4b2b);
  border: none;
  padding: 10px 18px;
  border-radius: 8px;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  transition: 0.3s ease;
}

.header button:hover {
  box-shadow: 0 0 12px #ff416c;
  transform: translateY(-2px);
}

/* Task Form */
.task-form {
  width: 100%;
  color : white;
  max-width: 600px;
  background: #1e1e2f;
  padding: 25px 20px;
  border-radius: 14px;
  box-shadow: 0 8px 20px rgba(0,0,0,0.5);
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 35px;
}

.task-form input,
.task-form textarea,
.task-form select {
  padding: 12px 14px;
  border: 1px solid #333;
  border-radius: 10px;
  font-size: 14px;
  background: #2a2a40;
  color: #f5f5f5;
  transition: 0.25s ease;
}

.task-form input:focus,
.task-form textarea:focus,
.task-form select:focus {
  border-color: #00b4ff;
  box-shadow: 0 0 8px #00b4ff;
  outline: none;
}

.task-form button {
  align-self: flex-end;
  background: linear-gradient(135deg, #00b4ff, #7b5cff);
  border: none;
  padding: 12px 20px;
  border-radius: 10px;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  transition: 0.3s ease;
}

.task-form button:hover {
  box-shadow: 0 0 12px #00ffd5;
  transform: translateY(-2px);
}

/* Task List */
.task-list {
  width: 100%;
  max-width: 900px;
}

.task-list ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.task-item {
  background: #1e1e2f;
  padding: 20px;
  border-radius: 14px;
  box-shadow: 0 6px 18px rgba(0,0,0,0.5);
  transition: transform 0.2s ease, box-shadow 0.3s ease;
}

.task-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 0 20px rgba(0, 180, 255, 0.3);
}

.task-info h3 {
  margin: 0 0 8px 0;
  color: #fff;
  font-size: 1.2rem;
}

.task-info p {
  margin: 0 0 10px 0;
  color: #bbb;
  font-size: 0.95rem;
}

/* Badges wrapper */
.badges {
  display: flex;
  gap: 8px;
  margin: 8px 0;
}

/* Badges */
.priority,
.status {
  display: inline-block;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  color: #fff;
  text-transform: capitalize;
}

.priority.low { background: #2abf91; }
.priority.medium { background: #ffb703; }
.priority.high { background: #ff4b2b; }

.status.todo { background: #7b5cff; }
.status.in-progress { background: #00b4ff; }
.status.done { background: #28a745; }

/* Progress Bar */
.progress-bar {
  width: 100%;
  height: 6px;
  border-radius: 4px;
  background: #333;
  margin-top: 12px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.4s ease;
}

.progress-fill.todo { width: 10%; background: #7b5cff; }
.progress-fill.in-progress { width: 60%; background: #00b4ff; }
.progress-fill.done { width: 100%; background: #28a745; }

/* Task Actions */
.task-actions {
  margin-top: 12px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.task-actions button {
  flex: 1;
  padding: 10px;
  border-radius: 8px;
  border: none;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: 0.3s ease;
}

.task-actions button:nth-child(1) {
  background: #ffb703;
  color: #fff;
}
.task-actions button:nth-child(1):hover { box-shadow: 0 0 8px #ffb703; }

.task-actions button:nth-child(2) {
  background: #28a745;
  color: #fff;
}
.task-actions button:nth-child(2):hover { box-shadow: 0 0 8px #28a745; }

.task-actions button:nth-child(3) {
  background: #00b4ff;
  color: #fff;
}
.task-actions button:nth-child(3):hover { box-shadow: 0 0 8px #00b4ff; }

.task-actions button:nth-child(4) {
  background: #ff4b2b;
  color: #fff;
}
.task-actions button:nth-child(4):hover { box-shadow: 0 0 8px #ff4b2b; }

/* Edit Form */
.edit-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.edit-form input,
.edit-form textarea,
.edit-form select {
  padding: 10px 12px;
  border: 1px solid #333;
  border-radius: 8px;
  background: #2a2a40;
  color: #fff;
}

.edit-options {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.edit-options label {
  display: flex;
  flex-direction: column;
  font-size: 0.9rem;
  font-weight: 600;
  color: #fff;
}

.edit-buttons {
  display: flex;
  gap: 10px;
}

.edit-buttons button {
  flex: 1;
  padding: 10px;
  border-radius: 8px;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: 0.3s ease;
}

.edit-buttons button:first-child {
  background: #2abf91;
  color: #fff;
}
.edit-buttons button:first-child:hover { box-shadow: 0 0 8px #2abf91; }

.edit-buttons button:last-child {
  background: #999;
  color: #fff;
}
.edit-buttons button:last-child:hover { box-shadow: 0 0 8px #777; }

/* Empty state */
.tasks-container p {
  margin-top: 20px;
  color: #aaa;
  font-size: 1rem;
  font-weight: 500;
  text-align: center;
}
</style>
