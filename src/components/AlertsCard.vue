<script setup lang="ts">
import { computed } from 'vue'
import type { AlertItem } from '../types/dashboard'

const props = defineProps<{
  alerts: AlertItem[]
  lang: 'en' | 'ru'
}>()

const hasCriticalAlerts = computed(() => props.alerts.some((item) => item.severity === 'critical'))
</script>

<template>
  <article class="card card--alerts" :class="{ 'card--critical': hasCriticalAlerts }" aria-label="Alerts card">
    <div class="card__head">
      <h2>{{ lang === 'ru' ? 'Алерты' : 'Alerts' }}</h2>
      <span v-if="hasCriticalAlerts" class="chip chip--danger">
        {{ lang === 'ru' ? 'есть critical' : 'critical active' }}
      </span>
      <span v-else class="chip chip--ok">{{ lang === 'ru' ? 'все в норме' : 'all good' }}</span>
    </div>

    <p v-if="alerts.length === 0" class="empty-state">
      {{ lang === 'ru' ? 'Активных алертов нет. Система стабильна.' : 'No active alerts. System looks stable.' }}
    </p>
    <ul v-else class="list" aria-label="Alert list">
      <li v-for="alert in alerts" :key="alert.id" class="alert">
        <div class="alert__head">
          <strong>{{ alert.title }}</strong>
          <small :class="alert.severity === 'critical' ? 'alert__severity alert__severity--critical' : 'alert__severity'">
            {{ alert.severity }}
          </small>
        </div>
        <p>{{ alert.message }}</p>
      </li>
    </ul>
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
  min-height: auto;
  display: flex;
  flex-direction: column;
  transition:
    transform 0.26s ease,
    border-color 0.26s ease,
    box-shadow 0.26s ease;
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

.card--critical {
  border-color: rgba(253, 164, 175, 0.52);
  background:
    linear-gradient(165deg, rgba(127, 29, 29, 0.48), rgba(12, 18, 36, 0.6)),
    rgba(8, 12, 24, 0.68);
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

.chip--danger {
  background: rgba(244, 63, 94, 0.18);
  border-color: rgba(251, 113, 133, 0.5);
  color: #fecdd3;
}

.chip--ok {
  background: rgba(74, 222, 128, 0.16);
  border-color: rgba(74, 222, 128, 0.36);
  color: #bbf7d0;
}

.empty-state {
  margin: 0;
  color: rgba(226, 232, 240, 0.82);
}

.list {
  list-style: none;
  margin: 0;
  padding: 0;
  overflow: auto;
}

.alert {
  border: 1px solid rgba(251, 113, 133, 0.28);
  background: rgba(15, 23, 42, 0.35);
  border-radius: 14px;
  padding: 11px;
  margin-bottom: 10px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

.alert__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.alert p {
  margin: 8px 0 0;
  color: rgba(226, 232, 240, 0.82);
}

.alert__severity {
  color: #fde68a;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-size: 10px;
  font-weight: 560;
}

.alert__severity--critical {
  color: #fda4af;
}

@media (hover: hover) {
  .card:hover {
    transform: translateY(-3px);
    box-shadow:
      0 28px 50px -28px rgba(0, 0, 0, 0.95),
      0 0 0 1px rgba(251, 191, 36, 0.14),
      0 0 18px rgba(251, 191, 36, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.22),
      inset 0 -1px 0 rgba(255, 255, 255, 0.05);
  }
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

  .alert__head {
    flex-wrap: wrap;
    row-gap: 4px;
  }
}
</style>
