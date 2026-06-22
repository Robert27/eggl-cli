<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

const props = defineProps<{
  mouseX?: number
  mouseY?: number
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
let raf = 0
let particles: Particle[] = []

interface Particle {
  x: number
  y: number
  vx: number
  vy: number
  size: number
  alpha: number
}

function initParticles(w: number, h: number, count: number) {
  particles = Array.from({ length: count }, () => ({
    x: Math.random() * w,
    y: Math.random() * h,
    vx: (Math.random() - 0.5) * 0.35,
    vy: (Math.random() - 0.5) * 0.35,
    size: Math.random() * 1.8 + 0.4,
    alpha: Math.random() * 0.45 + 0.15,
  }))
}

function draw() {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const w = canvas.width
  const h = canvas.height
  const mx = ((props.mouseX ?? 50) / 100) * w
  const my = ((props.mouseY ?? 50) / 100) * h

  ctx.clearRect(0, 0, w, h)

  for (const p of particles) {
    p.x += p.vx
    p.y += p.vy
    if (p.x < 0) p.x = w
    if (p.x > w) p.x = 0
    if (p.y < 0) p.y = h
    if (p.y > h) p.y = 0

    const dx = mx - p.x
    const dy = my - p.y
    const dist = Math.sqrt(dx * dx + dy * dy)
    if (dist < 140) {
      p.x -= dx * 0.012
      p.y -= dy * 0.012
    }

    ctx.beginPath()
    ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2)
    ctx.fillStyle = `rgba(249, 115, 22, ${p.alpha})`
    ctx.fill()
  }

  for (let i = 0; i < particles.length; i++) {
    for (let j = i + 1; j < particles.length; j++) {
      const a = particles[i]
      const b = particles[j]
      const dx = a.x - b.x
      const dy = a.y - b.y
      const dist = Math.sqrt(dx * dx + dy * dy)
      if (dist < 110) {
        ctx.beginPath()
        ctx.moveTo(a.x, a.y)
        ctx.lineTo(b.x, b.y)
        ctx.strokeStyle = `rgba(249, 115, 22, ${0.12 * (1 - dist / 110)})`
        ctx.lineWidth = 0.6
        ctx.stroke()
      }
    }
  }

  raf = requestAnimationFrame(draw)
}

function resize() {
  const canvas = canvasRef.value
  if (!canvas) return
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  const rect = canvas.parentElement?.getBoundingClientRect()
  if (!rect) return
  const ctx = canvas.getContext('2d')
  canvas.width = rect.width * dpr
  canvas.height = rect.height * dpr
  canvas.style.width = `${rect.width}px`
  canvas.style.height = `${rect.height}px`
  if (ctx) {
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  }
  const count = window.innerWidth < 768 ? 40 : 80
  initParticles(rect.width, rect.height, count)
}

onMounted(() => {
  resize()
  draw()
  window.addEventListener('resize', resize)
})

onUnmounted(() => {
  cancelAnimationFrame(raf)
  window.removeEventListener('resize', resize)
})
</script>

<template>
  <canvas ref="canvasRef" class="particle-canvas" aria-hidden="true" />
</template>

<style scoped>
.particle-canvas {
  position: absolute;
  inset: 0;
  pointer-events: none;
}
</style>
