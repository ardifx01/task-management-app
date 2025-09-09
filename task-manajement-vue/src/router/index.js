import { createRouter, createWebHistory } from 'vue-router'
import Register from '@/views/Register.vue'
import LandingPage from '@/views/LandingPage.vue'
import Login from '@/views/Login.vue'
import Projects from '@/views/Projects.vue'
import Tasks from '@/views/Task.vue'
import { useAuthStore } from '@/stores/auth'
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: LandingPage,
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
    },
    {
      path: '/register',
      name: 'register',
      component: Register,
    },
    {
      path: '/projects',
      name: 'projects',
      component: Projects,
      meta: { requiresAuth: true },
    },
    {
      path: '/projects/:id/tasks', 
      name: 'tasks',
      component: Tasks,
      meta: { requiresAuth: true },
    }
  ],
})
// Navigation Guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const isAuthenticated = authStore.isAuthenticated;
  
  const publicRoutes = ['home', 'login', 'register'];
  const requiresAuth = to.meta.requiresAuth;

  if (requiresAuth && !isAuthenticated) {
    next('/login');
  } else if (!requiresAuth && isAuthenticated) {
    if (publicRoutes.includes(to.name)) {
      next('/projects');
    } else {
      next();
    }
  } else {
    next();
  }
});

export default router
