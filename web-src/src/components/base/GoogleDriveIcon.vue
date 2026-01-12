<template>
  <svg 
    xmlns="http://www.w3.org/2000/svg" 
    viewBox="0 0 24 24" 
    :width="size" 
    :height="size"
    :class="className"
  >
    <defs>
      <linearGradient id="gd-gradient-blue" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#4285F4"/>
        <stop offset="100%" style="stop-color:#1a73e8"/>
      </linearGradient>
      <linearGradient id="gd-gradient-green" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#34A853"/>
        <stop offset="100%" style="stop-color:#188038"/>
      </linearGradient>
      <linearGradient id="gd-gradient-yellow" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#FBBC04"/>
        <stop offset="100%" style="stop-color:#ea8600"/>
      </linearGradient>
      <linearGradient id="gd-gradient-red" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#EA4335"/>
        <stop offset="100%" style="stop-color:#c5221f"/>
      </linearGradient>
    </defs>
    
    <!-- Google Drive Triangle Icon -->
    <path 
      d="M7.71 2L2 13.5l2.29 3.96L12.29 2H7.71z" 
      fill="url(#gd-gradient-blue)"
    />
    <path 
      d="M12.29 2L7.71 17.46l2.29 3.96L22 6.96 12.29 2z" 
      fill="url(#gd-gradient-green)"
    />
    <path 
      d="M22 6.96L16.29 17.46l2.29 3.96L24 10.92 22 6.96z" 
      fill="url(#gd-gradient-yellow)"
    />
    <path 
      d="M10 17.46L4.29 17.46l-2.29 3.96L7.71 21.42 10 17.46z" 
      fill="url(#gd-gradient-red)"
    />
  </svg>
</template>

<script setup lang="ts">
defineProps<{
  size?: number
  className?: string
}>()
</script>
