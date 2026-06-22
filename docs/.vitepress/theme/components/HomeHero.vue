<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from 'vue'
import ParticleField from './ParticleField.vue'

const terminalLines = [
  { prompt: '$', cmd: 'eggl doctor', delay: 0 },
  { prompt: '✓', cmd: ' go, git, kubectl — all good', delay: 800, success: true },
  { prompt: '$', cmd: 'eggl env toggle', delay: 1400 },
  { prompt: '→', cmd: ' switched to homelab', delay: 2200, accent: true },
  { prompt: '$', cmd: 'eggl pf longhorn', delay: 2800 },
  { prompt: '→', cmd: ' forwarding localhost:8080 → svc/longhorn-frontend', delay: 3600, muted: true },
]

const visibleLines = ref(0)
const typedChars = ref(0)
const currentLine = ref(0)
const mouseX = ref(50)
const mouseY = ref(50)
const scrollY = ref(0)
const copied = ref(false)
const terminalTilt = ref({ x: 0, y: 0 })
const heroVisible = ref(false)

const revealRefs = ref<HTMLElement[]>([])
let loopInterval: ReturnType<typeof setInterval> | null = null
let typeInterval: ReturnType<typeof setInterval> | null = null
let observer: IntersectionObserver | null = null

const floatingCommands = [
  'dedash', 'eol', 'env', 'pf', 'kill', 'doctor', 'empty', 'version',
]

const stats = [
  { value: '10+', label: 'commands' },
  { value: '3', label: 'platforms' },
  { value: '0', label: 'phone-home' },
  { value: '∞', label: 'workflow wins' },
]

const commands = [
  { name: 'dedash', desc: 'Em-dash → hyphen', icon: '—', href: '/commands/dedash', color: '#f97316' },
  { name: 'eol', desc: 'Normalize line endings', icon: '↵', href: '/commands/eol', color: '#22c55e' },
  { name: 'env', desc: 'Kube + Tailscale profiles', icon: '◎', href: '/commands/env', color: '#8b5cf6' },
  { name: 'pf', desc: 'Port-forward services', icon: '⇄', href: '/commands/pf', color: '#3b82f6' },
  { name: 'kill', desc: 'Free stuck ports', icon: '⊘', href: '/commands/kill', color: '#ef4444' },
  { name: 'doctor', desc: 'Environment checks', icon: '♥', href: '/commands/doctor', color: '#ec4899' },
]

const features = [
  { title: 'Workflow-first', desc: 'Commands shaped around real dev tasks — not another kitchen-sink CLI.', icon: '⚡' },
  { title: 'Zero phone-home', desc: 'Runs locally. Shells out to kubectl and tailscale on your machine only.', icon: '🔒' },
  { title: 'Beautiful TTY output', desc: 'Styled with lipgloss when your terminal supports it. Plain text otherwise.', icon: '✦' },
]

const titleWords = ['Your', 'daily', 'dev', 'helper', 'CLI']

const parallaxStyle = computed(() => ({
  '--mx': mouseX.value + '%',
  '--my': mouseY.value + '%',
  '--scroll': scrollY.value + 'px',
}))

const terminalStyle = computed(() => ({
  transform: `perspective(1000px) rotateX(${terminalTilt.value.x}deg) rotateY(${terminalTilt.value.y}deg)`,
}))

function onMouseMove(e: MouseEvent) {
  const el = e.currentTarget as HTMLElement
  const rect = el.getBoundingClientRect()
  mouseX.value = ((e.clientX - rect.left) / rect.width) * 100
  mouseY.value = ((e.clientY - rect.top) / rect.height) * 100
}

function onTerminalMove(e: MouseEvent) {
  const el = e.currentTarget as HTMLElement
  const rect = el.getBoundingClientRect()
  const x = (e.clientX - rect.left) / rect.width - 0.5
  const y = (e.clientY - rect.top) / rect.height - 0.5
  terminalTilt.value = { x: -y * 10, y: x * 12 }
}

function onTerminalLeave() {
  terminalTilt.value = { x: 0, y: 0 }
}

function onScroll() {
  scrollY.value = window.scrollY
}

async function copyInstall() {
  try {
    await navigator.clipboard.writeText('brew install eggl-cli')
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch { /* ignore */ }
}

function startTerminalAnimation() {
  visibleLines.value = 0
  currentLine.value = 0
  typedChars.value = 0

  const showNextLine = () => {
    if (currentLine.value >= terminalLines.length) {
      setTimeout(() => {
        currentLine.value = 0
        visibleLines.value = 0
        typedChars.value = 0
        showNextLine()
      }, 2800)
      return
    }

    visibleLines.value = currentLine.value + 1
    const line = terminalLines[currentLine.value]
    typedChars.value = 0

    if (typeInterval) clearInterval(typeInterval)
    typeInterval = setInterval(() => {
      typedChars.value++
      if (typedChars.value >= line.cmd.length) {
        if (typeInterval) clearInterval(typeInterval)
        currentLine.value++
        setTimeout(showNextLine, 350)
      }
    }, 26)
  }

  setTimeout(showNextLine, 800)
}

function getDisplayedCmd(index: number): string {
  if (index < currentLine.value) return terminalLines[index].cmd
  if (index === currentLine.value) return terminalLines[index].cmd.slice(0, typedChars.value)
  return ''
}

function setRevealRef(el: unknown) {
  if (el instanceof HTMLElement) revealRefs.value.push(el)
}

onMounted(() => {
  heroVisible.value = true
  startTerminalAnimation()
  loopInterval = setInterval(startTerminalAnimation, 14000)
  window.addEventListener('scroll', onScroll, { passive: true })

  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) entry.target.classList.add('is-visible')
      })
    },
    { threshold: 0.12, rootMargin: '0px 0px -40px 0px' },
  )

  document.querySelectorAll('.reveal').forEach((el) => observer?.observe(el))
})

onUnmounted(() => {
  if (loopInterval) clearInterval(loopInterval)
  if (typeInterval) clearInterval(typeInterval)
  observer?.disconnect()
  window.removeEventListener('scroll', onScroll)
})
</script>

<template>
  <div class="landing" :style="parallaxStyle" @mousemove="onMouseMove">
  <!-- Aurora + particles -->
  <div class="mesh">
    <div class="aurora aurora-1" />
    <div class="aurora aurora-2" />
    <div class="aurora aurora-3" />
    <div class="orb orb-1" />
    <div class="orb orb-2" />
    <div class="orb orb-3" />
    <div class="grid-overlay" />
    <div class="noise" />
    <ParticleField :mouse-x="mouseX" :mouse-y="mouseY" />
    <div class="spotlight" />
  </div>

  <!-- Floating command chips -->
  <div class="float-layer" aria-hidden="true">
    <span
      v-for="(cmd, i) in floatingCommands"
      :key="cmd"
      class="float-chip"
      :style="{ '--i': i, '--drift': (i % 2 ? 1 : -1) }"
    >eggl {{ cmd }}</span>
  </div>

  <nav class="nav" :class="{ visible: heroVisible }">
    <a href="/" class="brand">
      <span class="brand-glow" />
      <img src="/logo.svg" alt="eggl" width="36" height="36" class="brand-logo" />
      <span>eggl</span>
    </a>
    <div class="nav-links">
      <a href="/guide/getting-started">Guide</a>
      <a href="/commands/">Commands</a>
      <a href="https://github.com/Robert27/eggl-cli" target="_blank" rel="noopener">GitHub</a>
    </div>
  </nav>

  <section class="hero">
    <div class="hero-content">
      <div class="badge animate-in" style="--delay: 0.15s">
        <span class="pulse-dot" />
        <span class="badge-shimmer" />
        Go + Cobra · MIT License
      </div>

      <h1 class="title">
        <span
          v-for="(word, i) in titleWords"
          :key="word"
          class="title-word animate-in"
          :class="{ accent: word === 'helper' || word === 'CLI' }"
          :style="{ '--delay': 0.2 + i * 0.07 + 's' }"
        >{{ word }}</span>
      </h1>

      <p class="subtitle animate-in" style="--delay: 0.65s">
        File hygiene, environment switching, Kubernetes port-forwards, and more —
        <span class="highlight">one fast binary</span> for the tasks you run every day.
      </p>

      <div class="stats-row animate-in" style="--delay: 0.75s">
        <div v-for="stat in stats" :key="stat.label" class="stat">
          <span class="stat-value">{{ stat.value }}</span>
          <span class="stat-label">{{ stat.label }}</span>
        </div>
      </div>

      <div class="cta-row animate-in" style="--delay: 0.85s">
        <a href="/guide/getting-started" class="btn btn-primary">
          <span class="btn-shine" />
          Get started
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none"><path d="M3 8h10M9 4l4 4-4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </a>
        <a href="/commands/" class="btn btn-ghost">Browse commands</a>
      </div>

      <div class="install-snippet animate-in" style="--delay: 0.95s">
        <span class="prompt-char">$</span>
        <code>brew install eggl-cli</code>
        <button class="copy-btn" :class="{ copied }" title="Copy" @click="copyInstall">
          {{ copied ? '✓' : '⎘' }}
        </button>
      </div>
    </div>

    <div class="hero-terminal animate-in" style="--delay: 0.5s" @mousemove="onTerminalMove" @mouseleave="onTerminalLeave">
      <div class="terminal-glow" />
      <div class="terminal-window" :style="terminalStyle">
        <div class="terminal-scanline" />
        <div class="terminal-bar">
          <span class="dot red" /><span class="dot yellow" /><span class="dot green" />
          <span class="terminal-title">eggl — zsh</span>
          <span class="terminal-live"><span class="live-dot" /> live</span>
        </div>
        <div class="terminal-body">
          <div
            v-for="(line, i) in terminalLines"
            :key="i"
            class="terminal-line"
            :class="{ visible: i < visibleLines, success: line.success, accent: line.accent, muted: line.muted }"
          >
            <span class="prompt">{{ line.prompt }}</span>
            <span class="cmd">{{ getDisplayedCmd(i) }}</span>
            <span v-if="i === currentLine && typedChars < line.cmd.length" class="cursor">▊</span>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Marquee -->
  <div class="marquee-wrap reveal" :ref="setRevealRef">
    <div class="marquee">
      <span v-for="n in 2" :key="n" class="marquee-track">
        <span v-for="cmd in commands" :key="`${n}-${cmd.name}`" class="marquee-item">
          eggl {{ cmd.name }}
        </span>
      </span>
    </div>
  </div>

  <section class="commands-section">
    <div class="section-header reveal" :ref="setRevealRef">
      <span class="section-tag">Reference</span>
      <h2 class="section-title">Command palette</h2>
      <p class="section-desc">Every command documented with flags, examples, and edge cases.</p>
    </div>
    <div class="command-grid">
      <a
        v-for="(cmd, i) in commands"
        :key="cmd.name"
        :href="cmd.href"
        class="command-card reveal"
        :ref="setRevealRef"
        :style="{ '--card-delay': i * 0.08 + 's', '--card-color': cmd.color }"
      >
        <span class="card-glow" />
        <span class="cmd-icon">{{ cmd.icon }}</span>
        <span class="cmd-name">eggl {{ cmd.name }}</span>
        <span class="cmd-desc">{{ cmd.desc }}</span>
        <span class="card-arrow">→</span>
      </a>
    </div>
  </section>

  <section class="features-section">
    <div
      v-for="(feat, i) in features"
      :key="feat.title"
      class="feature-card reveal"
      :ref="setRevealRef"
      :style="{ '--card-delay': i * 0.12 + 's' }"
    >
      <span class="feature-ring" />
      <span class="feature-icon">{{ feat.icon }}</span>
      <h3>{{ feat.title }}</h3>
      <p>{{ feat.desc }}</p>
    </div>
  </section>

  <section class="cta-section reveal" :ref="setRevealRef">
    <div class="cta-border" />
    <div class="cta-inner">
      <h2>Ready to simplify your workflow?</h2>
      <p>Install in seconds. No config required to get started.</p>
      <div class="cta-install">
        <pre><code>brew tap Robert27/tap && brew install eggl-cli</code></pre>
      </div>
      <a href="/guide/installation" class="btn btn-primary btn-lg">
        <span class="btn-shine" />
        Installation guide →
      </a>
    </div>
  </section>
  </div>
</template>

<style scoped>
.landing {
  --accent: #f97316;
  --accent-glow: rgba(249, 115, 22, 0.5);
  --success: #22c55e;
  --bg: #050507;
  --surface: #121216;
  --border: rgba(255, 255, 255, 0.07);
  --text: #f4f4f5;
  --muted: #71717a;
  min-height: 100vh;
  background: var(--bg);
  color: var(--text);
  font-family: 'Outfit', system-ui, sans-serif;
  overflow-x: hidden;
  position: relative;
}

/* ── Background ── */
.mesh {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}

.aurora {
  position: absolute;
  width: 120%;
  height: 60%;
  filter: blur(60px);
  opacity: 0.35;
  mix-blend-mode: screen;
  animation: aurora-shift 12s ease-in-out infinite;
}

.aurora-1 {
  top: -20%;
  left: -10%;
  background: linear-gradient(120deg, transparent, var(--accent-glow), transparent);
  animation-delay: 0s;
}

.aurora-2 {
  top: 30%;
  right: -20%;
  background: linear-gradient(200deg, transparent, rgba(139, 92, 246, 0.35), transparent);
  animation-delay: -4s;
}

.aurora-3 {
  bottom: -10%;
  left: 20%;
  background: linear-gradient(60deg, transparent, rgba(34, 197, 94, 0.2), transparent);
  animation-delay: -8s;
}

@keyframes aurora-shift {
  0%, 100% { transform: translate(0, 0) rotate(0deg) scale(1); }
  33% { transform: translate(3%, 2%) rotate(2deg) scale(1.05); }
  66% { transform: translate(-2%, -1%) rotate(-1deg) scale(0.98); }
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  animation: float 20s ease-in-out infinite;
}

.orb-1 {
  width: 55vw; height: 55vw; max-width: 650px; max-height: 650px;
  background: radial-gradient(circle, var(--accent-glow), transparent 65%);
  top: -15%;
  left: calc(var(--mx, 50%) - 28%);
  transition: left 0.6s cubic-bezier(0.23, 1, 0.32, 1);
}

.orb-2 {
  width: 45vw; height: 45vw; max-width: 520px; max-height: 520px;
  background: radial-gradient(circle, rgba(34, 197, 94, 0.22), transparent 65%);
  bottom: 5%; right: -8%;
  animation-delay: -7s;
}

.orb-3 {
  width: 35vw; height: 35vw; max-width: 420px; max-height: 420px;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.18), transparent 65%);
  top: 45%; left: 55%;
  animation-delay: -14s;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(40px, -30px) scale(1.06); }
  50% { transform: translate(-25px, 25px) scale(0.94); }
  75% { transform: translate(15px, 10px) scale(1.02); }
}

.grid-overlay {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255,255,255,0.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255,255,255,0.025) 1px, transparent 1px);
  background-size: 72px 72px;
  mask-image: radial-gradient(ellipse 90% 70% at 50% 25%, black 20%, transparent 75%);
  animation: grid-drift 30s linear infinite;
}

@keyframes grid-drift {
  to { background-position: 72px 72px; }
}

.noise {
  position: absolute;
  inset: 0;
  opacity: 0.04;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
}

.spotlight {
  position: absolute;
  width: 600px; height: 600px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(249,115,22,0.08), transparent 70%);
  left: calc(var(--mx, 50%) - 300px);
  top: calc(var(--my, 50%) - 300px);
  transition: left 0.5s ease-out, top 0.5s ease-out;
}

/* ── Floating chips ── */
.float-layer {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}

.float-chip {
  position: absolute;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.7rem;
  color: rgba(249, 115, 22, 0.2);
  border: 1px solid rgba(249, 115, 22, 0.12);
  padding: 0.3rem 0.7rem;
  border-radius: 999px;
  white-space: nowrap;
  animation: chip-float 18s ease-in-out infinite;
  animation-delay: calc(var(--i) * -2.2s);
  left: calc(5% + var(--i) * 11%);
  top: calc(10% + (var(--i) * 7%) % 80%);
}

@keyframes chip-float {
  0%, 100% { transform: translate(0, 0) rotate(0deg); opacity: 0.3; }
  25% { transform: translate(calc(var(--drift) * 30px), -20px) rotate(2deg); opacity: 0.6; }
  50% { transform: translate(calc(var(--drift) * -15px), 15px) rotate(-1deg); opacity: 0.25; }
  75% { transform: translate(calc(var(--drift) * 20px), -10px) rotate(1deg); opacity: 0.5; }
}

/* ── Nav ── */
.nav {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem 2rem;
  opacity: 0;
  transform: translateY(-16px);
  transition: opacity 0.8s ease, transform 0.8s ease;
}

.nav.visible {
  opacity: 1;
  transform: translateY(0);
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-weight: 700;
  font-size: 1.25rem;
  color: var(--text);
  text-decoration: none;
  position: relative;
}

.brand-logo {
  position: relative;
  z-index: 1;
  animation: logo-pulse 4s ease-in-out infinite;
}

@keyframes logo-pulse {
  0%, 100% { filter: drop-shadow(0 0 0 transparent); }
  50% { filter: drop-shadow(0 0 12px var(--accent-glow)); }
}

.brand-glow {
  position: absolute;
  left: 0;
  width: 36px; height: 36px;
  border-radius: 10px;
  background: var(--accent-glow);
  filter: blur(16px);
  opacity: 0.5;
}

.nav-links { display: flex; gap: 2rem; }

.nav-links a {
  color: var(--muted);
  text-decoration: none;
  font-weight: 500;
  position: relative;
  transition: color 0.25s;
}

.nav-links a::after {
  content: '';
  position: absolute;
  bottom: -4px; left: 0;
  width: 0; height: 2px;
  background: var(--accent);
  transition: width 0.3s ease;
}

.nav-links a:hover { color: var(--accent); }
.nav-links a:hover::after { width: 100%; }

/* ── Hero ── */
.hero {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 2rem 5rem;
}

.badge {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.4rem 1rem;
  border-radius: 999px;
  border: 1px solid var(--border);
  background: rgba(255,255,255,0.03);
  font-size: 0.8rem;
  color: var(--muted);
  margin-bottom: 1.5rem;
  overflow: hidden;
}

.badge-shimmer {
  position: absolute;
  inset: 0;
  background: linear-gradient(105deg, transparent 40%, rgba(255,255,255,0.06) 50%, transparent 60%);
  animation: shimmer-sweep 3s ease-in-out infinite;
}

@keyframes shimmer-sweep {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

.pulse-dot {
  width: 8px; height: 8px;
  border-radius: 50%;
  background: var(--success);
  animation: pulse 2s ease-in-out infinite;
  position: relative;
  z-index: 1;
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(34,197,94,0.5); }
  50% { box-shadow: 0 0 0 8px rgba(34,197,94,0); }
}

.title {
  font-size: clamp(2.8rem, 6vw, 4.5rem);
  font-weight: 800;
  line-height: 1.05;
  letter-spacing: -0.04em;
  margin: 0 0 1.25rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.2em;
}

.title-word {
  display: inline-block;
}

.title-word.accent {
  background: linear-gradient(135deg, var(--accent) 0%, #fbbf24 40%, #f97316 70%, #ea580c 100%);
  background-size: 300% auto;
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: shimmer 5s linear infinite, fadeUp 0.8s cubic-bezier(0.4,0,0.2,1) forwards;
}

@keyframes shimmer {
  to { background-position: 300% center; }
}

.subtitle {
  font-size: 1.15rem;
  line-height: 1.75;
  color: var(--muted);
  max-width: 32rem;
  margin: 0 0 1.75rem;
}

.highlight {
  color: var(--text);
  background: linear-gradient(180deg, transparent 60%, rgba(249,115,22,0.2) 60%);
}

.stats-row {
  display: flex;
  gap: 2rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}

.stat {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 800;
  font-family: 'JetBrains Mono', monospace;
  color: var(--accent);
  line-height: 1;
}

.stat-label {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--muted);
}

.cta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.btn {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.85rem 1.5rem;
  border-radius: 12px;
  font-weight: 600;
  font-size: 0.95rem;
  text-decoration: none;
  overflow: hidden;
  transition: transform 0.3s cubic-bezier(0.34,1.56,0.64,1), box-shadow 0.3s;
}

.btn-lg { padding: 1rem 2rem; font-size: 1rem; }

.btn-primary {
  background: linear-gradient(135deg, var(--accent), #ea580c);
  color: #fff;
  box-shadow: 0 4px 28px var(--accent-glow), 0 0 0 1px rgba(255,255,255,0.1) inset;
}

.btn-primary:hover {
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 12px 40px var(--accent-glow);
}

.btn-shine {
  position: absolute;
  inset: 0;
  background: linear-gradient(105deg, transparent 40%, rgba(255,255,255,0.2) 50%, transparent 60%);
  transform: translateX(-100%);
  animation: shimmer-sweep 2.5s ease-in-out infinite;
}

.btn-ghost {
  border: 1px solid var(--border);
  color: var(--text);
  background: rgba(255,255,255,0.02);
  backdrop-filter: blur(8px);
}

.btn-ghost:hover {
  border-color: rgba(249,115,22,0.5);
  color: var(--accent);
  transform: translateY(-2px);
}

.install-snippet {
  display: inline-flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.65rem 1rem;
  border-radius: 12px;
  background: var(--surface);
  border: 1px solid var(--border);
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.85rem;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.install-snippet:hover {
  border-color: rgba(249,115,22,0.3);
  box-shadow: 0 0 24px rgba(249,115,22,0.1);
}

.prompt-char { color: var(--accent); font-weight: 600; }
.install-snippet code { color: var(--text); }

.copy-btn {
  background: rgba(255,255,255,0.05);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--muted);
  cursor: pointer;
  padding: 0.2rem 0.45rem;
  font-size: 0.8rem;
  transition: all 0.2s;
}

.copy-btn:hover, .copy-btn.copied {
  color: var(--success);
  border-color: rgba(34,197,94,0.4);
}

/* ── Terminal ── */
.hero-terminal {
  position: relative;
  perspective: 1000px;
}

.terminal-glow {
  position: absolute;
  inset: -20%;
  background: radial-gradient(ellipse, var(--accent-glow), transparent 70%);
  opacity: 0.35;
  filter: blur(40px);
  animation: glow-pulse 4s ease-in-out infinite;
}

@keyframes glow-pulse {
  0%, 100% { opacity: 0.25; transform: scale(1); }
  50% { opacity: 0.45; transform: scale(1.05); }
}

.terminal-window {
  position: relative;
  border-radius: 18px;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(8, 8, 12, 0.92);
  backdrop-filter: blur(24px);
  box-shadow:
    0 32px 100px rgba(0,0,0,0.6),
    0 0 0 1px rgba(255,255,255,0.06) inset,
    0 0 60px rgba(249,115,22,0.08);
  transition: transform 0.15s ease-out;
  overflow: hidden;
}

.terminal-scanline {
  position: absolute;
  inset: 0;
  background: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0,0,0,0.03) 2px,
    rgba(0,0,0,0.03) 4px
  );
  pointer-events: none;
  z-index: 2;
}

.terminal-bar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0.85rem 1.1rem;
  border-bottom: 1px solid var(--border);
  position: relative;
  z-index: 1;
}

.dot { width: 10px; height: 10px; border-radius: 50%; }
.dot.red { background: #ef4444; }
.dot.yellow { background: #eab308; }
.dot.green { background: #22c55e; }

.terminal-title {
  margin-left: 0.5rem;
  font-size: 0.75rem;
  color: var(--muted);
  font-family: 'JetBrains Mono', monospace;
  flex: 1;
}

.terminal-live {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.65rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--success);
  font-family: 'JetBrains Mono', monospace;
}

.live-dot {
  width: 6px; height: 6px;
  border-radius: 50%;
  background: var(--success);
  animation: pulse 1.5s ease-in-out infinite;
}

.terminal-body {
  padding: 1.25rem;
  min-height: 230px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.8rem;
  line-height: 1.85;
  position: relative;
  z-index: 1;
}

.terminal-line {
  opacity: 0;
  transform: translateX(-12px);
  transition: opacity 0.35s, transform 0.35s;
}

.terminal-line.visible { opacity: 1; transform: translateX(0); }
.terminal-line .prompt { color: var(--accent); margin-right: 0.5rem; }
.terminal-line.success .cmd { color: var(--success); }
.terminal-line.accent .cmd { color: var(--accent); }
.terminal-line.muted .cmd { color: var(--muted); }

.cursor {
  color: var(--accent);
  animation: blink 0.9s step-end infinite;
}

@keyframes blink { 50% { opacity: 0; } }

/* ── Marquee ── */
.marquee-wrap {
  position: relative;
  z-index: 1;
  overflow: hidden;
  padding: 1.5rem 0;
  border-block: 1px solid var(--border);
  mask-image: linear-gradient(90deg, transparent, black 10%, black 90%, transparent);
}

.marquee {
  display: flex;
  width: max-content;
  animation: marquee 30s linear infinite;
}

.marquee-track {
  display: flex;
  gap: 3rem;
  padding-right: 3rem;
}

.marquee-item {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.85rem;
  color: rgba(249,115,22,0.35);
  white-space: nowrap;
}

@keyframes marquee {
  to { transform: translateX(-50%); }
}

/* ── Sections ── */
.commands-section,
.features-section {
  position: relative;
  z-index: 1;
  max-width: 1200px;
  margin: 0 auto;
  padding: 5rem 2rem;
}

.section-header { text-align: center; margin-bottom: 3rem; }

.section-tag {
  display: inline-block;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.15em;
  color: var(--accent);
  border: 1px solid rgba(249,115,22,0.3);
  padding: 0.3rem 0.8rem;
  border-radius: 999px;
  margin-bottom: 1rem;
}

.section-title {
  font-size: clamp(1.8rem, 4vw, 2.5rem);
  font-weight: 800;
  margin: 0 0 0.5rem;
  letter-spacing: -0.02em;
}

.section-desc {
  color: var(--muted);
  margin: 0;
  font-size: 1.05rem;
}

.command-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.25rem;
}

.command-card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  padding: 1.6rem;
  border-radius: 18px;
  border: 1px solid var(--border);
  background: rgba(18, 18, 22, 0.7);
  backdrop-filter: blur(12px);
  text-decoration: none;
  color: inherit;
  overflow: hidden;
  transition: transform 0.4s cubic-bezier(0.34,1.56,0.64,1), border-color 0.3s, box-shadow 0.3s;
}

.card-glow {
  position: absolute;
  inset: -1px;
  border-radius: 18px;
  background: linear-gradient(135deg, var(--card-color, var(--accent)), transparent 60%);
  opacity: 0;
  transition: opacity 0.4s;
  z-index: 0;
}

.command-card::after {
  content: '';
  position: absolute;
  inset: 1px;
  border-radius: 17px;
  background: rgba(12, 12, 16, 0.9);
  z-index: 0;
}

.command-card > *:not(.card-glow) { position: relative; z-index: 1; }

.command-card:hover {
  transform: translateY(-6px) scale(1.02);
  border-color: rgba(249,115,22,0.35);
  box-shadow: 0 20px 50px rgba(0,0,0,0.4), 0 0 40px color-mix(in srgb, var(--card-color) 20%, transparent);
}

.command-card:hover .card-glow { opacity: 0.25; }

.cmd-icon { font-size: 1.6rem; }
.cmd-name {
  font-family: 'JetBrains Mono', monospace;
  font-weight: 600;
  color: var(--card-color, var(--accent));
}
.cmd-desc { font-size: 0.9rem; color: var(--muted); }

.card-arrow {
  position: absolute;
  top: 1.5rem; right: 1.5rem;
  opacity: 0;
  transform: translateX(-10px) rotate(-45deg);
  transition: all 0.35s cubic-bezier(0.34,1.56,0.64,1);
  color: var(--card-color, var(--accent));
  z-index: 2;
}

.command-card:hover .card-arrow {
  opacity: 1;
  transform: translateX(0) rotate(0deg);
}

.features-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 1.5rem;
  padding-top: 2rem;
}

.feature-card {
  position: relative;
  padding: 2.25rem 2rem;
  border-radius: 20px;
  border: 1px solid var(--border);
  background: rgba(18, 18, 22, 0.5);
  text-align: center;
  overflow: hidden;
  transition: transform 0.4s ease, border-color 0.3s;
}

.feature-card:hover {
  transform: translateY(-4px);
  border-color: rgba(249,115,22,0.25);
}

.feature-ring {
  position: absolute;
  top: 50%; left: 50%;
  width: 120px; height: 120px;
  border: 1px solid rgba(249,115,22,0.15);
  border-radius: 50%;
  transform: translate(-50%, -50%);
  animation: ring-expand 3s ease-in-out infinite;
}

@keyframes ring-expand {
  0%, 100% { transform: translate(-50%,-50%) scale(0.8); opacity: 0.3; }
  50% { transform: translate(-50%,-50%) scale(1.2); opacity: 0.1; }
}

.feature-icon {
  font-size: 2.2rem;
  display: block;
  margin-bottom: 1rem;
  position: relative;
  animation: icon-bob 3s ease-in-out infinite;
}

@keyframes icon-bob {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}

.feature-card h3 { margin: 0 0 0.5rem; font-size: 1.15rem; }
.feature-card p { margin: 0; color: var(--muted); font-size: 0.95rem; line-height: 1.65; }

/* ── CTA ── */
.cta-section {
  position: relative;
  z-index: 1;
  margin: 2rem auto 5rem;
  max-width: 720px;
  width: calc(100% - 4rem);
  padding: 3px;
  border-radius: 26px;
  overflow: hidden;
}

.cta-border {
  position: absolute;
  inset: -50%;
  background: conic-gradient(from 0deg, var(--accent), #fbbf24, #8b5cf6, var(--accent));
  animation: spin 6s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.cta-inner {
  position: relative;
  text-align: center;
  padding: 4.5rem 2.5rem;
  border-radius: 24px;
  background: var(--bg);
}

.cta-inner h2 {
  font-size: clamp(1.6rem, 3vw, 2rem);
  margin: 0 0 0.5rem;
  font-weight: 800;
}

.cta-inner p { color: var(--muted); margin: 0 0 1.75rem; }

.cta-install { margin-bottom: 1.75rem; }

.cta-install pre {
  display: inline-block;
  margin: 0;
  padding: 1rem 1.5rem;
  border-radius: 12px;
  background: var(--surface);
  border: 1px solid var(--border);
}

.cta-install code {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.82rem;
  color: var(--text);
}

/* ── Animations ── */
.animate-in {
  opacity: 0;
  transform: translateY(28px);
  animation: fadeUp 0.9s cubic-bezier(0.34,1.56,0.64,1) forwards;
  animation-delay: var(--delay, 0s);
}

@keyframes fadeUp {
  to { opacity: 1; transform: translateY(0); }
}

.reveal {
  opacity: 0;
  transform: translateY(32px);
  transition: opacity 0.7s ease, transform 0.7s cubic-bezier(0.34,1.56,0.64,1);
  transition-delay: var(--card-delay, 0s);
}

.reveal.is-visible {
  opacity: 1;
  transform: translateY(0);
}

@media (max-width: 900px) {
  .hero { grid-template-columns: 1fr; padding-bottom: 3rem; }
  .hero-terminal { order: -1; }
  .float-chip { display: none; }
  .stats-row { gap: 1.25rem; }
  .nav-links { gap: 1rem; font-size: 0.9rem; }
}

@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
</style>
