<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import LoadingSpinner from '../components/LoadingSpinner.vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const projects = ref([]);
const newProjectName = ref('');
const newProjectDescription = ref('');
const authStore = useAuthStore();
const router = useRouter();
const isLoading = ref(true);

const editingProject = ref(null); // Variable to hold the project being edited
const editedName = ref('');
const editedDescription = ref('');

const fetchProjects = async () => {
  isLoading.value = true;
  try {
    const response = await axios.get('http://localhost:3000/projects', {
      headers: {
        Authorization: `Bearer ${authStore.token}`,
      },
    });
    projects.value = response.data;
  } catch (err) {
    console.error('Failed to fetch projects:', err);
    if (err.response && err.response.status === 401) {
      authStore.clearToken();
      router.push('/login');
    }
  } finally {
    isLoading.value = false;
  }
};

const createProject = async () => {
  if (!newProjectName.value) return;
  try {
    await axios.post(
      'http://localhost:3000/projects',
      {
        name: newProjectName.value,
        description: newProjectDescription.value,
      },
      {
        headers: {
          Authorization: `Bearer ${authStore.token}`,
        },
      }
    );
    newProjectName.value = '';
    newProjectDescription.value = '';
    fetchProjects();
  } catch (err) {
    console.error('Failed to create project:', err);
  }
};

const deleteProject = async (projectId) => {
  try {
    await axios.delete(`http://localhost:3000/projects/${projectId}`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`,
      },
    });
    fetchProjects();
  } catch (err) {
    console.error('Failed to delete project:', err);
  }
};

// Function to start editing a project
const startEdit = (project) => {
  editingProject.value = project;
  editedName.value = project.name;
  editedDescription.value = project.description;
};

// Function to cancel editing
const cancelEdit = () => {
  editingProject.value = null;
};

// Function to save and update the project
const saveEdit = async () => {
  if (!editingProject.value) return;
  try {
    await axios.put(
      `http://localhost:3000/projects/${editingProject.value.id}`,
      {
        name: editedName.value,
        description: editedDescription.value,
      },
      {
        headers: {
          Authorization: `Bearer ${authStore.token}`,
        },
      }
    );
    editingProject.value = null; // Clear editing state
    fetchProjects(); // Reload projects list
  } catch (err) {
    console.error('Failed to update project:', err);
  }
};

const logout = () => {
  authStore.clearToken();
  router.push('/login');
};

onMounted(() => {
  fetchProjects();
});
</script>

<template>
  <div class="projects-container">
    <div class="header">
      <h1>My Projects</h1>
      <button @click="logout">Logout</button>
    </div>

    <div class="project-form">
      <input v-model="newProjectName" placeholder="Project Name" />
      <textarea v-model="newProjectDescription" placeholder="Project Description"></textarea>
      <button @click="createProject">Create Project</button>
    </div>

    <LoadingSpinner v-if="isLoading" />

    <div v-else-if="projects.length" class="project-list">
      <ul>
        <li v-for="project in projects" :key="project.id" class="project-item">
          <div v-if="editingProject && editingProject.id === project.id" class="edit-form">
            <input v-model="editedName" placeholder="Project Name" />
            <textarea v-model="editedDescription" placeholder="Project Description"></textarea>
            <div class="edit-buttons">
              <button @click="saveEdit">Save</button>
              <button @click="cancelEdit">Cancel</button>
            </div>
          </div>
          <div v-else class="project-info">
            <h3>{{ project.name }}</h3>
            <p>{{ project.description }}</p>
          </div>
          
          <div v-if="!editingProject || editingProject.id !== project.id" class="project-actions">
            <button @click="router.push(`/projects/${project.id}/tasks`)">View Tasks</button>
            <button @click="startEdit(project)">Edit</button>
            <button @click="deleteProject(project.id)">Delete</button>
          </div>
        </li>
      </ul>
    </div>
    <p v-else>You don't have any projects yet.</p>
  </div>
</template>

<style scoped>
.projects-container {
  min-height: 100vh;
  padding: 40px 20px;
  background: linear-gradient(135deg, #0f2027, #203a43, #2c5364); /* Dark gradient */
  font-family: 'Inter', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #eaeaea;
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
  color: #fff;
  font-size: 2.2rem;
  font-weight: 700;
  text-align: center;
  text-shadow: 0 2px 6px rgba(0,0,0,0.35);
}

.header button {
  background: linear-gradient(135deg, #f40303, #b1372f, #f40303);
  border: none;
  padding: 10px 18px;
  border-radius: 8px;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  transition: 0.3s ease;
}


.header button:hover {
  background: linear-gradient(135deg, #e63e26, #d9325d);
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(0,0,0,0.3);
}

/* Project Form */
.project-form {
  width: 100%;
  max-width: 600px;
  background: rgba(30, 30, 30, 0.75);
  backdrop-filter: blur(14px);
  padding: 25px 20px;
  border-radius: 14px;
  box-shadow: 0 8px 20px rgba(0,0,0,0.35);
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 35px;
}

.project-form input,
.project-form textarea {
  padding: 12px 14px;
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 10px;
  font-size: 14px;
  background: rgba(255,255,255,0.05);
  color: #fff;
  transition: 0.25s ease;
}

.project-form input::placeholder,
.project-form textarea::placeholder {
  color: rgba(255,255,255,0.5);
}

.project-form input:focus,
.project-form textarea:focus {
  border-color: #00d4ff;
  box-shadow: 0 0 8px rgba(0,212,255,0.4);
  outline: none;
  transform: scale(1.01);
}

.project-form button {
  align-self: flex-end;
  background: linear-gradient(135deg, #00d4ff, #0077ff);
  border: none;
  padding: 12px 20px;
  border-radius: 10px;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.project-form button:hover {
  background: linear-gradient(135deg, #009ecb, #005bbb);
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(0,0,0,0.3);
}

/* Project List */
.project-list {
  width: 100%;
  max-width: 900px;
}

.project-list ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.project-item {
  background: rgba(30,30,30,0.75);
  backdrop-filter: blur(14px);
  padding: 20px;
  border-radius: 14px;
  box-shadow: 0 6px 18px rgba(0,0,0,0.35);
  transition: transform 0.2s ease, box-shadow 0.3s ease;
}

.project-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 25px rgba(0,0,0,0.4);
}

.project-info h3 {
  margin: 0 0 10px 0;
  color: #fff;
  font-size: 1.3rem;
}

.project-info p {
  margin: 0;
  color: rgba(255,255,255,0.85);
  font-size: 0.95rem;
}

/* Project Actions */
.project-actions {
  margin-top: 15px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.project-actions button {
  flex: 1;
  padding: 10px;
  border-radius: 8px;
  border: none;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: 0.3s ease;
}

/* View Tasks */
.project-actions button:nth-child(1) {
  background: #00aaff;
  color: #fff;
}
.project-actions button:nth-child(1):hover {
  background: #0086cc;
}

/* Edit */
.project-actions button:nth-child(2) {
  background: #ffb703;
  color: #fff;
}
.project-actions button:nth-child(2):hover {
  background: #e09a00;
}

/* Delete */
.project-actions button:nth-child(3) {
  background: #ff4b2b;
  color: #fff;
}
.project-actions button:nth-child(3):hover {
  background: #d93a1f;
}

/* Edit Form */
.edit-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.edit-form input,
.edit-form textarea {
  padding: 10px 12px;
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 8px;
  background: rgba(255,255,255,0.05);
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
.edit-buttons button:first-child:hover {
  background: #229977;
}

.edit-buttons button:last-child {
  background: #666;
  color: #fff;
}
.edit-buttons button:last-child:hover {
  background: #444;
}

/* Empty state */
.projects-container p {
  margin-top: 20px;
  color: #ccc;
  font-size: 1rem;
  font-weight: 500;
}
</style>

