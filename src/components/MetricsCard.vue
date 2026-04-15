<script setup lang="ts">
import { computed } from 'vue'
import type { Metrics } from '../types/dashboard'

const props = defineProps<{
  metrics: Metrics | null
}>()

const formattedUpdatedAt = computed(() => {
  if (!props.metrics?.updatedAt) {
    return '-'
  }
  const timestamp = new Date(props.metrics.updatedAt)
  if (Number.isNaN(timestamp.getTime())) {
    return props.metrics.updatedAt
  }
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(timestamp)
})
</script>

<template>
  <article class="card" aria-label="System metrics card">
    <div class="card__head">
      <h2>CPU / RAM</h2>
      <span class="chip chip--muted">System</span>
    </div>
    <div class="stats" role="list" aria-label="CPU and RAM metrics">
      <div class="stat" role="listitem">
        <span class="stat__label">CPU</span>
        <strong class="stat__value">{{ metrics?.cpuPercent ?? 0 }}%</strong>
        <div class="progress" aria-hidden="true">
          <div class="progress__fill" :style="{ width: `${metrics?.cpuPercent ?? 0}%` }"></div>
        </div>
      </div>
      <div class="stat" role="listitem">
        <span class="stat__label">RAM</span>
        <strong class="stat__value">{{ metrics?.ramPercent ?? 0 }}%</strong>
        <div class="progress" aria-hidden="true">
          <div class="progress__fill progress__fill--violet" :style="{ width: `${metrics?.ramPercent ?? 0}%` }"></div>
        </div>
      </div>
    </div>
    <p class="meta">Updated {{ formattedUpdatedAt }}</p>
  </article>
</template>

<style scoped>
.card {
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 22px;
  padding: 18px;
  background:
    linear-gradient(160deg, rgba(255, 255, 255, 0.12), rgba(255, 255, 255, 0.04)),
    rgba(8, 12, 24, 0.62);
  backdrop-filter: blur(24px) saturate(135%);
  box-shadow:
    0 24px 45px -30px rgba(0, 0, 0, 0.95),
    inset 0 1px 0 rgba(255, 255, 255, 0.16),
    inset 0 -1px 0 rgba(255, 255, 255, 0.04);
  position: relative;
  overflow: hidden;
  animation: rise-in 440ms cubic-bezier(0.2, 0.8, 0.2, 1);
  width: 100%;
  flex: 1 1 360px;
  max-width: 380px;
  min-height: 280px;
  display: flex;
  flex-direction: column;
}

.card::before {
  content: '';
  position: absolute;
  inset: 0 0 auto 0;
  height: 48%;
  background: linear-gradient(
    172deg,
    rgba(255, 255, 255, 0.24),
    rgba(255, 255, 255, 0.08) 40%,
    rgba(255, 255, 255, 0.01) 75%
  );
  pointer-events: none;
}

.card__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.card h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.01em;
}

.chip {
  border: 1px solid rgba(147, 197, 253, 0.3);
  background: rgba(96, 165, 250, 0.16);
  color: #dbeafe;
  font-size: 10px;
  line-height: 1;
  padding: 7px 11px;
  border-radius: 999px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-weight: 560;
}

.chip--muted {
  background: rgba(148, 163, 184, 0.16);
  border-color: rgba(148, 163, 184, 0.25);
  color: rgba(226, 232, 240, 0.75);
}

.stats {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 6px;
}

.stat {
  border: 1px solid rgba(255, 255, 255, 0.11);
  border-radius: 16px;
  padding: 12px;
  background: rgba(7, 13, 28, 0.42);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.07);
  width: 100%;
  box-sizing: border-box;
}

.stat__label {
  display: block;
  font-size: 12px;
  color: rgba(203, 213, 225, 0.75);
  margin-bottom: 6px;
}

.stat__value {
  font-size: 30px;
  font-weight: 650;
  letter-spacing: -0.02em;
}

.progress {
  margin-top: 10px;
  width: 100%;
  height: 6px;
  border-radius: 999px;
  background: rgba(148, 163, 184, 0.2);
  overflow: hidden;
}

.progress__fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #22d3ee, #60a5fa);
  box-shadow: 0 0 16px rgba(56, 189, 248, 0.65);
  transition: width 0.4s ease;
}

.progress__fill--violet {
  background: linear-gradient(90deg, #a78bfa, #818cf8);
  box-shadow: 0 0 16px rgba(139, 92, 246, 0.6);
}

.meta {
  margin-top: 12px;
  color: rgba(203, 213, 225, 0.66);
  font-size: 13px;
}

:global(.page--desktop) .stats {
  gap: 12px;
}

@media (max-width: 720px) {
  .card {
    border-radius: 18px;
    flex: 1 1 100%;
    max-width: none;
    padding: 14px;
  }

  .card h2 {
    font-size: 17px;
  }

  .chip {
    padding: 6px 9px;
    font-size: 9px;
  }

  .stat {
    padding: 10px;
  }

  .stat__value {
    font-size: 26px;
  }

  .stats {
    gap: 10px;
  }
}
</style>
