<script setup lang="ts">
/**
 * App.vue - Root Application Component
 * Premium Dashboard with Glass Effects
 */

import { computed, ref, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import Login from '@/views/Login.vue'
import Dashboard from '@/views/Dashboard.vue'

const appStore = useAppStore()
const isInitialized = ref(false)

// Check auth synchronously before rendering
appStore.checkAuth()
isInitialized.value = true

const isAuthenticated = computed(() => appStore.isAuthenticated)
</script>

<template>
  <div class="app-root">
    <!-- Background Effects -->
    <div class="app-background">
      <div class="bg-gradient-orb bg-orb-1"></div>
      <div class="bg-gradient-orb bg-orb-2"></div>
      <div class="bg-noise"></div>
    </div>

    <!-- Main Content -->
    <Transition name="fade" mode="out-in" v-if="isInitialized">
      <Login v-if="!isAuthenticated" key="login" />
      <Dashboard v-else key="dashboard" />
    </Transition>
  </div>
</template>

<style scoped>
.app-root {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
}

/* ========== Background Effects ========== */
.app-background {
  position: fixed;
  inset: 0;
  z-index: -1;
  overflow: hidden;
}

.bg-gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: float 20s ease-in-out infinite;
}

.bg-orb-1 {
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, var(--color-accent) 0%, transparent 70%);
  top: -200px;
  right: -100px;
  animation-delay: 0s;
}

.bg-orb-2 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, var(--color-secondary) 0%, transparent 70%);
  bottom: -150px;
  left: -100px;
  animation-delay: -10s;
}

.bg-noise {
  position: absolute;
  inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
  opacity: 0.03;
  pointer-events: none;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  25% {
    transform: translate(30px, -30px) scale(1.05);
  }
  50% {
    transform: translate(-20px, 20px) scale(0.95);
  }
  75% {
    transform: translate(-30px, -10px) scale(1.02);
  }
}

/* ========== Transitions ========== */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--duration-normal) var(--ease-default);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
