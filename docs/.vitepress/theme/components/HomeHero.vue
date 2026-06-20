<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

const terminalLines = [
  { prompt: '$', cmd: 'eggl doctor', delay: 0 },
  { prompt: '✓', cmd: ' go, git, kubectl — all good', delay: 800, success: true },
  { prompt: '$', cmd: 'eggl env toggle', delay: 1400 },
  { prompt: '→', cmd: ' switched to homelab', delay: 2200, accent: true },
  { prompt: '$', cmd: 'eggl pf longhorn', delay: 2800 },
  { prompt: '→', cmd: ' forwarding localhost:8080 → svc/longhorn-frontend', delay: 3600, muted: true },
]

const visibleLines = ref<number>(0)
const typedChars = ref(0)
const currentLine = ref(0)
const mouseX = ref(50)
const mouseY = ref(50)

let interval: ReturnType<typeof setInterval> | null = null
let typeInterval: ReturnType<typeof setInterval> | null = null

const commands = [
  { name: 'dedash', desc: 'Em-dash → hyphen', icon: '—', href: '/commands/dedash' },
  { name: 'eol', desc: 'Normalize line endings', icon: '↵', href: '/commands/eol' },
  { name: 'env', desc: 'Kube + Tailscale profiles', icon: '◎', href: '/commands/env' },
  { name: 'pf', desc: 'Port-forward services', icon: '⇄', href: '/commands/pf' },
  { name: 'kill', desc: 'Free stuck ports', icon: '⊘', href: '/commands/kill' },
  { name: 'doctor', desc: 'Environment checks', icon: '♥', href: '/commands/doctor' },
]

const features = [
  {
    title: 'Workflow-first',
    desc: 'Commands shaped around real dev tasks — not another kitchen-sink CLI.',
    icon: '⚡',
  },
  {
    title: 'Zero phone-home',
    desc: 'Runs locally. Shells out to kubectl and tailscale on your machine only.',
    icon: '🔒',
  },
  {
    title: 'Beautiful TTY output',
    desc: 'Styled with lipgloss when your terminal supports it. Plain text otherwise.',
    icon: '✦',
  },
]

function onMouseMove(e: MouseEvent) {
  const el = e.currentTarget as HTMLElement
  const rect = el.getBoundingClientRect()
  mouseX.value = ((e.clientX - rect.left) / rect.width) * 100
  mouseY.value = ((e.clientY - rect.top) / rect.height) * 100
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
      }, 3000)
      return
    }

    visibleLines.value = currentLine.value + 1
    const line = terminalLines[currentLine.value]
    const fullText = line.cmd
    typedChars.value = 0

    if (typeInterval) clearInterval(typeInterval)
    typeInterval = setInterval(() => {
      typedChars.value++
      if (typedChars.value >= fullText.length) {
        if (typeInterval) clearInterval(typeInterval)
        currentLine.value++
        setTimeout(showNextLine, line.delay > 0 ? 400 : 300)
      }
    }, 28)
  }

  setTimeout(showNextLine, 600)
}

onMounted(() => {
  startTerminalAnimation()
  interval = setInterval(startTerminalAnimation, 12000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
  if (typeInterval) clearInterval(typeInterval)
})

function getDisplayedCmd(index: number): string {
  if (index < currentLine.value) return terminalLines[index].cmd
  if (index === currentLine.value) return terminalLines[index].cmd.slice(0, typedChars.value)
  return ''
}
</script>

<template>
  <div class="landing" @mousemove="onMouseMove">
    <div class="mesh" :style="{ '--mx': mouseX + '%', '--my': mouseY + '%' }">
      <div class="orb orb-1" />
      <div class="orb orb-2" />
      <div class="orb orb-3" />
      <div class="grid-overlay" />
    </div>

    <nav class="nav">
      <a href="/" class="brand">
        <img src="/logo.svg" alt="eggl" width="36" height="36" />
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
        <div class="badge animate-in" style="--delay: 0.1s">
          <span class="pulse-dot" />
          Go + Cobra · MIT License
        </div>

        <h1 class="title animate-in" style="--delay: 0.2s">
          Your daily dev
          <span class="gradient-text">helper CLI</span>
        </h1>

        <p class="subtitle animate-in" style="--delay: 0.35s">
          File hygiene, environment switching, Kubernetes port-forwards, and more —
          one fast binary for the tasks you run every day.
        </p>

        <div class="cta-row animate-in" style="--delay: 0.5s">
          <a href="/guide/getting-started" class="btn btn-primary">
            Get started
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none"><path d="M3 8h10M9 4l4 4-4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </a>
          <a href="/commands/" class="btn btn-ghost">Browse commands</a>
        </div>

        <div class="install-snippet animate-in" style="--delay: 0.65s">
          <code>brew install eggl-cli</code>
          <button class="copy-hint" title="Copy">⎘</button>
        </div>
      </div>

      <div class="hero-terminal animate-in" style="--delay: 0.45s">
        <div class="terminal-window">
          <div class="terminal-bar">
            <span class="dot red" /><span class="dot yellow" /><span class="dot green" />
            <span class="terminal-title">eggl — zsh</span>
          </div>
          <div class="terminal-body">
            <div
              v-for="(line, i) in terminalLines"
              :key="i"
              class="terminal-line"
              :class="{
                visible: i < visibleLines,
                success: line.success,
                accent: line.accent,
                muted: line.muted,
              }"
            >
              <span class="prompt">{{ line.prompt }}</span>
              <span class="cmd">{{ getDisplayedCmd(i) }}</span>
              <span v-if="i === currentLine && typedChars < line.cmd.length" class="cursor">▊</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="commands-section">
      <h2 class="section-title animate-on-scroll">Command palette</h2>
      <p class="section-desc animate-on-scroll">Every command documented with flags, examples, and edge cases.</p>
      <div class="command-grid">
        <a
          v-for="(cmd, i) in commands"
          :key="cmd.name"
          :href="cmd.href"
          class="command-card animate-on-scroll"
          :style="{ '--card-delay': i * 0.08 + 's' }"
        >
          <span class="cmd-icon">{{ cmd.icon }}</span>
          <span class="cmd-name">eggl {{ cmd.name }}</span>
          <span class="cmd-desc">{{ cmd.desc }}</span>
          <span class="card-arrow">→</span>
        </a>
      </div>
    </section>

    <section class="features-section">
      <div v-for="(feat, i) in features" :key="feat.title" class="feature-card animate-on-scroll" :style="{ '--card-delay': i * 0.1 + 's' }">
        <span class="feature-icon">{{ feat.icon }}</span>
        <h3>{{ feat.title }}</h3>
        <p>{{ feat.desc }}</p>
      </div>
    </section>

    <section class="cta-section animate-on-scroll">
      <h2>Ready to simplify your workflow?</h2>
      <p>Install in seconds. No config required to get started.</p>
      <div class="cta-install">
        <pre><code>brew tap Robert27/tap && brew install eggl-cli</code></pre>
      </div>
      <a href="/guide/installation" class="btn btn-primary">Installation guide →</a>
    </section>
  </div>
</template>

<style scoped>
.landing {
  --accent: #f97316;
  --accent-glow: rgba(249, 115, 22, 0.45);
  --success: #22c55e;
  --bg: #09090b;
  --surface: #18181b;
  --border: rgba(255, 255, 255, 0.08);
  --text: #e5e7eb;
  --muted: #6b7280;
  min-height: 100vh;
  background: var(--bg);
  color: var(--text);
  font-family: 'Outfit', system-ui, sans-serif;
  overflow-x: hidden;
  position: relative;
}

.mesh {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.5;
  animation: float 18s ease-in-out infinite;
}

.orb-1 {
  width: 50vw;
  height: 50vw;
  max-width: 600px;
  max-height: 600px;
  background: radial-gradient(circle, var(--accent-glow), transparent 70%);
  top: -10%;
  left: calc(var(--mx, 50%) - 25%);
  transition: left 0.4s ease-out;
}

.orb-2 {
  width: 40vw;
  height: 40vw;
  max-width: 500px;
  max-height: 500px;
  background: radial-gradient(circle, rgba(34, 197, 94, 0.2), transparent 70%);
  bottom: 10%;
  right: -5%;
  animation-delay: -6s;
}

.orb-3 {
  width: 30vw;
  height: 30vw;
  max-width: 400px;
  max-height: 400px;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.15), transparent 70%);
  top: 40%;
  left: 60%;
  animation-delay: -12s;
}

.grid-overlay {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
  background-size: 64px 64px;
  mask-image: radial-gradient(ellipse 80% 60% at 50% 30%, black, transparent);
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -40px) scale(1.05); }
  66% { transform: translate(-20px, 20px) scale(0.95); }
}

.nav {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem 2rem;
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-weight: 700;
  font-size: 1.25rem;
  color: var(--text);
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 2rem;
}

.nav-links a {
  color: var(--muted);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-links a:hover {
  color: var(--accent);
}

.hero {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 2rem 6rem;
}

.badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.35rem 0.9rem;
  border-radius: 999px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.04);
  font-size: 0.8rem;
  color: var(--muted);
  margin-bottom: 1.5rem;
}

.pulse-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--success);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.5); }
  50% { opacity: 0.8; box-shadow: 0 0 0 6px rgba(34, 197, 94, 0); }
}

.title {
  font-size: clamp(2.5rem, 5vw, 4rem);
  font-weight: 800;
  line-height: 1.1;
  letter-spacing: -0.03em;
  margin: 0 0 1.25rem;
}

.gradient-text {
  display: block;
  background: linear-gradient(135deg, var(--accent) 0%, #fbbf24 50%, var(--accent) 100%);
  background-size: 200% auto;
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: shimmer 4s linear infinite;
}

@keyframes shimmer {
  to { background-position: 200% center; }
}

.subtitle {
  font-size: 1.15rem;
  line-height: 1.7;
  color: var(--muted);
  max-width: 32rem;
  margin: 0 0 2rem;
}

.cta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.85rem 1.5rem;
  border-radius: 12px;
  font-weight: 600;
  font-size: 0.95rem;
  text-decoration: none;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.btn-primary {
  background: linear-gradient(135deg, var(--accent), #ea580c);
  color: #fff;
  box-shadow: 0 4px 24px var(--accent-glow);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 32px var(--accent-glow);
}

.btn-ghost {
  border: 1px solid var(--border);
  color: var(--text);
  background: rgba(255, 255, 255, 0.03);
}

.btn-ghost:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.install-snippet {
  display: inline-flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.6rem 1rem;
  border-radius: 10px;
  background: var(--surface);
  border: 1px solid var(--border);
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.85rem;
  color: var(--muted);
}

.copy-hint {
  background: none;
  border: none;
  color: var(--muted);
  cursor: pointer;
  opacity: 0.6;
}

.hero-terminal {
  perspective: 1000px;
}

.terminal-window {
  border-radius: 16px;
  border: 1px solid var(--border);
  background: rgba(15, 15, 18, 0.85);
  backdrop-filter: blur(20px);
  box-shadow:
    0 24px 80px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(255, 255, 255, 0.05) inset;
  transform: rotateY(-4deg) rotateX(2deg);
  transition: transform 0.4s ease;
  animation: terminal-float 6s ease-in-out infinite;
}

.terminal-window:hover {
  transform: rotateY(0) rotateX(0);
}

@keyframes terminal-float {
  0%, 100% { transform: rotateY(-4deg) rotateX(2deg) translateY(0); }
  50% { transform: rotateY(-2deg) rotateX(1deg) translateY(-8px); }
}

.terminal-bar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border);
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.dot.red { background: #ef4444; }
.dot.yellow { background: #eab308; }
.dot.green { background: #22c55e; }

.terminal-title {
  margin-left: 0.5rem;
  font-size: 0.75rem;
  color: var(--muted);
  font-family: 'JetBrains Mono', monospace;
}

.terminal-body {
  padding: 1.25rem;
  min-height: 220px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.8rem;
  line-height: 1.8;
}

.terminal-line {
  opacity: 0;
  transform: translateX(-8px);
  transition: opacity 0.3s, transform 0.3s;
}

.terminal-line.visible {
  opacity: 1;
  transform: translateX(0);
}

.terminal-line .prompt {
  color: var(--accent);
  margin-right: 0.5rem;
}

.terminal-line.success .cmd { color: var(--success); }
.terminal-line.accent .cmd { color: var(--accent); }
.terminal-line.muted .cmd { color: var(--muted); }

.cursor {
  color: var(--accent);
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  50% { opacity: 0; }
}

.commands-section,
.features-section {
  position: relative;
  z-index: 1;
  max-width: 1200px;
  margin: 0 auto;
  padding: 4rem 2rem;
}

.section-title {
  font-size: 2rem;
  font-weight: 700;
  text-align: center;
  margin: 0 0 0.5rem;
}

.section-desc {
  text-align: center;
  color: var(--muted);
  margin: 0 0 3rem;
}

.command-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}

.command-card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  padding: 1.5rem;
  border-radius: 16px;
  border: 1px solid var(--border);
  background: rgba(24, 24, 27, 0.6);
  backdrop-filter: blur(8px);
  text-decoration: none;
  color: inherit;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.command-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, var(--accent-glow), transparent 60%);
  opacity: 0;
  transition: opacity 0.3s;
}

.command-card:hover {
  border-color: rgba(249, 115, 22, 0.4);
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.3);
}

.command-card:hover::before {
  opacity: 1;
}

.cmd-icon {
  font-size: 1.5rem;
  position: relative;
}

.cmd-name {
  font-family: 'JetBrains Mono', monospace;
  font-weight: 600;
  color: var(--accent);
  position: relative;
}

.cmd-desc {
  font-size: 0.9rem;
  color: var(--muted);
  position: relative;
}

.card-arrow {
  position: absolute;
  top: 1.5rem;
  right: 1.5rem;
  opacity: 0;
  transform: translateX(-8px);
  transition: all 0.3s;
  color: var(--accent);
}

.command-card:hover .card-arrow {
  opacity: 1;
  transform: translateX(0);
}

.features-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 1.5rem;
}

.feature-card {
  padding: 2rem;
  border-radius: 16px;
  border: 1px solid var(--border);
  background: rgba(24, 24, 27, 0.4);
  text-align: center;
}

.feature-icon {
  font-size: 2rem;
  display: block;
  margin-bottom: 1rem;
}

.feature-card h3 {
  margin: 0 0 0.5rem;
  font-size: 1.1rem;
}

.feature-card p {
  margin: 0;
  color: var(--muted);
  font-size: 0.95rem;
  line-height: 1.6;
}

.cta-section {
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 5rem 2rem;
  border-radius: 24px;
  margin: 2rem auto 4rem;
  max-width: 720px;
  width: calc(100% - 4rem);
  border: 1px solid var(--border);
  background: linear-gradient(180deg, rgba(249, 115, 22, 0.08), transparent);
}

.cta-section h2 {
  font-size: 2rem;
  margin: 0 0 0.5rem;
}

.cta-section p {
  color: var(--muted);
  margin: 0 0 1.5rem;
}

.cta-install {
  margin-bottom: 1.5rem;
}

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
  font-size: 0.85rem;
  color: var(--text);
}

.animate-in {
  opacity: 0;
  transform: translateY(24px);
  animation: fadeUp 0.8s cubic-bezier(0.4, 0, 0.2, 1) forwards;
  animation-delay: var(--delay, 0s);
}

@keyframes fadeUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-on-scroll {
  opacity: 0;
  transform: translateY(20px);
  animation: fadeUp 0.6s ease forwards;
  animation-delay: var(--card-delay, 0s);
}

@media (max-width: 900px) {
  .hero {
    grid-template-columns: 1fr;
    padding-bottom: 3rem;
  }

  .hero-terminal {
    order: -1;
  }

  .terminal-window {
    transform: none;
    animation: none;
  }

  .nav-links {
    gap: 1rem;
    font-size: 0.9rem;
  }
}
</style>
