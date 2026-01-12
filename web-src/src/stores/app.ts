import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  // State
  const isAuthenticated = ref(false)
  const userName = ref('')
  
  // Check authentication status on load
  function checkAuth() {
    const authToken = localStorage.getItem('auth_token')
    const savedUser = localStorage.getItem('user_name')
    
    if (authToken) {
      isAuthenticated.value = true
      userName.value = savedUser || ''
    }
  }
  
  // Login
  function login(username: string) {
    isAuthenticated.value = true
    userName.value = username
    localStorage.setItem('auth_token', 'authenticated')
    localStorage.setItem('user_name', username)
  }
  
  // Logout
  function logout() {
    isAuthenticated.value = false
    userName.value = ''
    localStorage.removeItem('auth_token')
    localStorage.removeItem('user_name')
    
    // Mark that user just logged out (for wallpaper refresh)
    sessionStorage.setItem('just_logged_out', 'true')
    
    // Also clear any session cookies
    document.cookie = 'session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;'
  }
  
  return {
    isAuthenticated,
    userName,
    checkAuth,
    login,
    logout
  }
})
