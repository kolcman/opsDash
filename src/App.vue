<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import AlertsCard from './components/AlertsCard.vue'
import ConfigGeneratorCard from './components/ConfigGeneratorCard.vue'
import MetricsCard from './components/MetricsCard.vue'
import ServicesCard from './components/ServicesCard.vue'
import { fetchAlerts, fetchMetrics, fetchServices } from './services/dashboardApi'
import type { AlertItem, Metrics, ServiceStatus } from './types/dashboard'
type DashboardMode = 'monitoring' | 'generator'
type Lang = 'en' | 'ru'

const metrics = ref<Metrics | null>(null)
const services = ref<ServiceStatus[]>([])
const alerts = ref<AlertItem[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const dashboardMode = ref<DashboardMode>('monitoring')
const lang = ref<Lang>('en')
let pollTimer: ReturnType<typeof setInterval> | null = null
const dashboardModes: DashboardMode[] = ['monitoring', 'generator']
const langs: Lang[] = ['en', 'ru']

async function loadDashboardData() {
  try {
    error.value = null
    const [nextMetrics, nextServices, nextAlerts] = await Promise.all([
      fetchMetrics(),
      fetchServices(),
      fetchAlerts(),
    ])
    metrics.value = nextMetrics
    services.value = nextServices
    alerts.value = nextAlerts
  } catch (requestError) {
    error.value = requestError instanceof Error ? requestError.message : 'Unknown request error'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadDashboardData()
  pollTimer = setInterval(loadDashboardData, 10000)
})

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
})
</script>

<template>
  <main class="page page--auto">
    <div class="ambient ambient--top"></div>
    <div class="ambient ambient--bottom"></div>

    <header class="page__header">
      <div>
        <p class="eyebrow">{{ lang === 'ru' ? 'Инфраструктура' : 'Infrastructure' }}</p>
        <h1>OpsDash</h1>
        <p class="subtitle">
          {{
            lang === 'ru' ? 'Инструменты DevOps' : 'DevOps tools'
          }}
        </p>
      </div>
      <div class="header-controls">
        <span class="chip chip--live">{{ lang === 'ru' ? 'Обновление каждые 10с' : 'Live every 10s' }}</span>
        <section class="lang-switch lang-switch--mobile" aria-label="Language switch mobile">
          <button
            v-for="currentLang in langs"
            :key="`mobile-${currentLang}`"
            type="button"
            class="lang-switch__button"
            :class="{ 'lang-switch__button--active': lang === currentLang }"
            @click="lang = currentLang"
          >
            {{ currentLang.toUpperCase() }}
          </button>
        </section>
      </div>
    </header>

    <section class="switch-row">
      <section class="mode-switch" aria-label="Dashboard mode">
        <button
          v-for="mode in dashboardModes"
          :key="mode"
          type="button"
          class="mode-switch__button"
          :class="{ 'mode-switch__button--active': dashboardMode === mode }"
          @click="dashboardMode = mode"
        >
          {{
            mode === 'monitoring'
              ? lang === 'ru'
                ? 'Мониторинг'
                : 'Monitoring'
              : lang === 'ru'
                ? 'Генератор конфигов'
                : 'Config Generator'
          }}
        </button>
      </section>

      <section class="lang-switch lang-switch--desktop" aria-label="Language switch">
        <button
          v-for="currentLang in langs"
          :key="currentLang"
          type="button"
          class="lang-switch__button"
          :class="{ 'lang-switch__button--active': lang === currentLang }"
          @click="lang = currentLang"
        >
          {{ currentLang.toUpperCase() }}
        </button>
      </section>
    </section>

    <section v-if="loading" class="skeleton-grid" aria-label="Loading dashboard">
      <article v-for="item in dashboardMode === 'monitoring' ? 3 : 1" :key="item" class="skeleton-card"></article>
    </section>
    <p v-else-if="error" class="state state--error">
      {{ lang === 'ru' ? 'Не удалось загрузить дашборд:' : 'Failed to load dashboard:' }} {{ error }}
    </p>

    <section v-else-if="dashboardMode === 'monitoring'" class="grid">
      <MetricsCard :metrics="metrics" :lang="lang" />
      <ServicesCard :services="services" :lang="lang" />
      <AlertsCard :alerts="alerts" :lang="lang" />
    </section>

    <section v-else class="grid grid--generator">
      <ConfigGeneratorCard :lang="lang" />
    </section>
  </main>
</template>

<style scoped>
:global(html),
:global(body),
:global(#app) {
  margin: 0;
  min-height: 100%;
  background:
    radial-gradient(circle at 18% 8%, rgba(76, 115, 255, 0.32), transparent 38%),
    radial-gradient(circle at 80% 10%, rgba(153, 97, 255, 0.26), transparent 36%),
    radial-gradient(circle at 50% 105%, rgba(51, 171, 255, 0.24), transparent 42%),
    linear-gradient(160deg, #02030a 0%, #060913 45%, #0a1020 100%);
  background-attachment: fixed;
}

.page {
  max-width: 1440px;
  width: min(1440px, 100%);
  box-sizing: border-box;
  margin: 0 auto;
  min-height: 100dvh;
  padding: 30px 16px 42px;
  font-family:
    'SF Pro Display',
    'Inter',
    -apple-system,
    BlinkMacSystemFont,
    'Segoe UI',
    Roboto,
    sans-serif;
  color: #f8fafc;
  position: relative;
  overflow: visible;
}

.ambient {
  position: fixed;
  width: 360px;
  height: 360px;
  border-radius: 999px;
  filter: blur(86px);
  opacity: 0.34;
  pointer-events: none;
  z-index: 0;
}

.ambient--top {
  top: -110px;
  right: -90px;
  background: #5071ff;
}

.ambient--bottom {
  bottom: -170px;
  left: -140px;
  background: #7b54ff;
}

.page__header {
  margin-bottom: 18px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 12px;
  position: relative;
  z-index: 1;
}

.header-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
  position: relative;
  z-index: 1;
}

.eyebrow {
  margin: 0 0 6px;
  text-transform: uppercase;
  letter-spacing: 0.14em;
  font-size: 12px;
  color: rgba(186, 201, 255, 0.9);
}

.page__header h1 {
  margin: 0;
  line-height: 1.1;
  font-size: clamp(32px, 5vw, 48px);
  font-weight: 650;
  letter-spacing: -0.02em;
}

.subtitle {
  margin: 8px 0 0;
  color: rgba(203, 213, 225, 0.8);
}

.state--error {
  color: #fecdd3;
}

.state {
  margin-top: 16px;
  color: rgba(226, 232, 240, 0.85);
  position: relative;
  z-index: 1;
}

.mode-switch {
  display: inline-flex;
  box-sizing: border-box;
  max-width: 100%;
  align-items: center;
  gap: 6px;
  padding: 5px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.18);
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(14px);
  margin: 0;
  position: relative;
  z-index: 1;
}

.lang-switch {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  background: rgba(15, 23, 42, 0.35);
  margin: 0;
  position: relative;
  z-index: 1;
}

.lang-switch--mobile {
  display: none;
}

.lang-switch__button {
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: rgba(226, 232, 240, 0.74);
  font-size: 9px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  padding: 5px 9px;
  cursor: pointer;
}

.lang-switch__button--active {
  background: rgba(125, 211, 252, 0.22);
  color: #e0f2fe;
}

.mode-switch__button {
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: rgba(226, 232, 240, 0.8);
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  padding: 7px 11px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.mode-switch__button:hover {
  color: #f8fafc;
  background: rgba(186, 230, 253, 0.12);
}

.mode-switch__button--active {
  background: rgba(125, 211, 252, 0.22);
  color: #e0f2fe;
}

.mode-switch__button:focus-visible {
  outline: 2px solid rgba(125, 211, 252, 0.9);
  outline-offset: 2px;
}

.skeleton-grid {
  margin: 16px auto 0;
  width: min(1240px, 100%);
  display: flex;
  flex-wrap: wrap;
  align-items: stretch;
  justify-content: center;
  gap: 24px;
  position: relative;
  z-index: 1;
}

.skeleton-card {
  width: 100%;
  flex: 1 1 360px;
  max-width: 380px;
  min-height: 280px;
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.16);
  background:
    linear-gradient(130deg, rgba(148, 163, 184, 0.08), rgba(148, 163, 184, 0.2), rgba(148, 163, 184, 0.08))
      0 0 / 210% 100%,
    rgba(8, 12, 24, 0.62);
  backdrop-filter: blur(24px) saturate(135%);
  animation: shimmer 1.8s linear infinite;
}

.grid {
  margin: 16px auto 0;
  width: min(1240px, 100%);
  display: flex;
  flex-wrap: wrap;
  align-items: stretch;
  justify-content: center;
  gap: 24px;
  position: relative;
  z-index: 1;
}

.grid--generator {
  width: min(1240px, 100%);
}

.grid--generator :deep(.card--generator) {
  flex: 1 1 100%;
  width: 100%;
  max-width: none;
}

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

.card h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.01em;
}

.card__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.card--critical {
  border-color: rgba(253, 164, 175, 0.52);
  background:
    linear-gradient(165deg, rgba(127, 29, 29, 0.48), rgba(12, 18, 36, 0.6)),
    rgba(8, 12, 24, 0.68);
}

.card--alerts {
  min-height: auto;
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

.chip--live {
  background: rgba(56, 189, 248, 0.14);
  border-color: rgba(125, 211, 252, 0.36);
  color: #bae6fd;
}

.chip--muted {
  background: rgba(148, 163, 184, 0.16);
  border-color: rgba(148, 163, 184, 0.25);
  color: rgba(226, 232, 240, 0.75);
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

.stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 6px;
}

.stat {
  border: 1px solid rgba(255, 255, 255, 0.11);
  border-radius: 16px;
  padding: 12px;
  background: rgba(7, 13, 28, 0.42);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.07);
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

.list {
  list-style: none;
  margin: 0;
  padding: 0;
  overflow: auto;
}

.list__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);
  padding: 11px 0;
}

.service-name {
  letter-spacing: 0.01em;
}

.status {
  font-size: 10px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  border: 1px solid transparent;
  border-radius: 999px;
  padding: 6px 10px;
  font-weight: 560;
}

.status--up {
  color: #bbf7d0;
  border-color: rgba(74, 222, 128, 0.38);
  background: rgba(22, 163, 74, 0.2);
}

.status--down {
  color: #fecdd3;
  border-color: rgba(251, 113, 133, 0.5);
  background: rgba(190, 24, 93, 0.2);
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

@keyframes shimmer {
  from {
    background-position: 120% 0, 0 0;
  }
  to {
    background-position: -120% 0, 0 0;
  }
}

@keyframes rise-in {
  from {
    transform: translateY(12px) scale(0.99);
    opacity: 0;
  }
  to {
    transform: translateY(0) scale(1);
    opacity: 1;
  }
}

@media (max-width: 720px) {
  .page {
    padding: 18px 12px 28px;
  }

  .page__header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    margin-bottom: 12px;
  }

  .eyebrow {
    font-size: 11px;
    margin-bottom: 4px;
  }

  .page__header h1 {
    font-size: clamp(30px, 8vw, 40px);
  }

  .subtitle {
    font-size: 13px;
    margin-top: 6px;
  }

  .grid,
  .skeleton-grid {
    width: 100%;
    justify-content: stretch;
    gap: 12px;
  }

  .card,
  .skeleton-card {
    border-radius: 18px;
    aspect-ratio: auto;
    flex: 1 1 100%;
    max-width: none;
  }

  .card {
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

  .list__row {
    padding: 10px 0;
  }

  .status {
    padding: 5px 8px;
    font-size: 9px;
  }

  .alert__head {
    flex-wrap: wrap;
    row-gap: 4px;
  }

  .ambient {
    opacity: 0.28;
  }

  .stats {
    grid-template-columns: 1fr;
  }

  .mode-switch {
    width: 100%;
    justify-content: space-between;
  }

  .switch-row {
    margin-bottom: 12px;
  }

  .lang-switch--desktop {
    display: none;
  }

  .lang-switch--mobile {
    display: inline-flex;
  }

  .header-controls {
    width: 100%;
    justify-content: space-between;
  }
}

@media (min-width: 721px) and (max-width: 1024px) {
  .page {
    max-width: 920px;
    padding: 24px 20px 34px;
  }

  .page__header {
    margin-bottom: 14px;
  }

  .grid,
  .skeleton-grid {
    width: 100%;
    justify-content: center;
    gap: 14px;
  }

  .card {
    border-radius: 20px;
    padding: 16px;
    flex: 1 1 calc(50% - 14px);
    max-width: 440px;
  }

  .chip {
    font-size: 9px;
  }

  .stat__value {
    font-size: 28px;
  }
}

.page--phone {
  max-width: 430px;
  padding: 18px 12px 28px;
}

.page--phone .page__header {
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 12px;
}

.page--phone .eyebrow {
  font-size: 11px;
  margin-bottom: 4px;
}

.page--phone .page__header h1 {
  font-size: clamp(30px, 8vw, 40px);
}

.page--phone .subtitle {
  font-size: 13px;
  margin-top: 6px;
}

.page--phone .grid,
.page--phone .skeleton-grid {
  width: 100%;
  justify-content: stretch;
  gap: 12px;
}

.page--phone .card,
.page--phone .skeleton-card {
  border-radius: 18px;
  aspect-ratio: auto;
  flex: 1 1 100%;
  max-width: none;
}

.page--phone .card {
  padding: 14px;
}

.page--phone .card h2 {
  font-size: 17px;
}

.page--phone .chip {
  padding: 6px 9px;
  font-size: 9px;
}

.page--phone .stat {
  padding: 10px;
}

.page--phone .stat__value {
  font-size: 26px;
}

.page--phone .list__row {
  padding: 10px 0;
}

.page--phone .status {
  padding: 5px 8px;
  font-size: 9px;
}

.page--phone .alert__head {
  flex-wrap: wrap;
  row-gap: 4px;
}

.page--phone .ambient {
  opacity: 0.28;
}

.page--phone .stats {
  grid-template-columns: 1fr;
}

.page--phone .mode-switch {
  width: 100%;
  justify-content: space-between;
}

.page--tablet {
  max-width: 920px;
  padding: 24px 20px 34px;
}

.page--tablet .page__header {
  margin-bottom: 14px;
}

.page--tablet .grid,
.page--tablet .skeleton-grid {
  width: min(920px, 100%);
  display: flex;
  flex-wrap: wrap;
  align-items: stretch;
  justify-content: center;
  gap: 20px;
}

.page--tablet .card {
  border-radius: 20px;
  padding: 16px;
}

.page--tablet .grid > .card:not(.card--alerts),
.page--tablet .skeleton-grid > .skeleton-card {
  flex: 1 1 360px;
  max-width: 420px;
}

.page--tablet .card--alerts {
  flex: 1 1 100%;
  max-width: none;
  min-height: fit-content;
}

.page--tablet .chip {
  font-size: 9px;
}

.page--tablet .stat__value {
  font-size: 28px;
}

.page--tablet .stats {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.page--tablet .stats .stat {
  flex: 1 1 180px;
  min-width: 0;
}

.page--desktop {
  max-width: 1440px;
}

.page--auto {
  max-width: 1440px;
}

.page--desktop .grid,
.page--desktop .skeleton-grid {
  width: min(1240px, 100%);
  justify-content: center;
  gap: 24px;
}

.page--desktop .card,
.page--desktop .skeleton-card {
  flex: 1 1 360px;
  max-width: 380px;
}

.page--desktop .card--alerts {
  flex: 1 1 760px;
  max-width: 784px;
  min-height: fit-content;
}

.page--desktop .grid {
  display: grid;
  grid-template-columns: minmax(320px, 380px) minmax(640px, 1fr);
  grid-auto-rows: minmax(220px, auto);
  align-items: stretch;
}

.page--desktop .grid > .card {
  width: auto;
  max-width: none;
  min-height: 280px;
}

.page--desktop .grid > .card:not(.card--alerts) {
  grid-column: 1;
}

.page--desktop .grid > .card--alerts {
  grid-column: 2;
  grid-row: 1;
  align-self: start;
  min-height: fit-content;
}

.page--desktop .stats {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.page--desktop .stats .stat {
  width: auto;
  box-sizing: border-box;
  padding: 10px;
}

.page--desktop .stats .stat__label {
  margin-bottom: 4px;
  font-size: 11px;
}

.page--desktop .stats .stat__value {
  font-size: 26px;
  line-height: 1.1;
}

.page--desktop .stats .progress {
  margin-top: 8px;
}

@media (min-width: 721px) and (max-width: 1024px) {
  .page--auto {
    max-width: 920px;
    padding: 24px 20px 34px;
  }

  .page--auto .page__header {
    margin-bottom: 14px;
  }

  .page--auto .grid,
  .page--auto .skeleton-grid {
    width: min(920px, 100%);
    display: flex;
    flex-wrap: wrap;
    align-items: stretch;
    justify-content: center;
    gap: 20px;
  }

  .page--auto .card {
    border-radius: 20px;
    padding: 16px;
  }

  .page--auto .grid > .card:not(.card--alerts),
  .page--auto .skeleton-grid > .skeleton-card {
    flex: 1 1 300px;
    max-width: 340px;
  }

  .page--auto .grid--generator > * {
    flex: 1 1 100%;
    max-width: none;
  }

  .page--auto .grid.grid--generator > .card.card--generator {
    flex: 1 1 100%;
    width: 100%;
    max-width: none;
  }

  .page--auto .card--alerts {
    flex: 1 1 100%;
    max-width: none;
    min-height: fit-content;
  }

  .page--auto .chip {
    font-size: 9px;
  }

  .page--auto .stat__value {
    font-size: 28px;
  }

  .page--auto .stats {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
  }

  .page--auto .stats .stat {
    flex: 1 1 180px;
    min-width: 0;
  }
}

@media (min-width: 1025px) {
  .page--auto {
    max-width: 1440px;
  }

  .page--auto .grid,
  .page--auto .skeleton-grid {
    width: min(1240px, 100%);
    justify-content: center;
    gap: 24px;
  }

  .page--auto .card,
  .page--auto .skeleton-card {
    flex: 1 1 360px;
    max-width: 380px;
  }

  .page--auto .card--alerts {
    flex: 1 1 760px;
    max-width: 784px;
    min-height: fit-content;
  }

  .page--auto .grid {
    display: grid;
    grid-template-columns: minmax(320px, 380px) minmax(640px, 1fr);
    grid-auto-rows: minmax(220px, auto);
    align-items: stretch;
  }

  .page--auto .grid > .card {
    width: auto;
    max-width: none;
    min-height: 280px;
  }

  .page--auto .grid > .card:not(.card--alerts) {
    grid-column: 1;
  }

  .page--auto .grid > .card--alerts {
    grid-column: 2;
    grid-row: 1;
    align-self: start;
    min-height: fit-content;
  }

  .page--auto .grid--generator {
    display: flex;
    width: min(1240px, 100%);
  }

  .page--auto .grid--generator > * {
    flex: 1 1 100%;
    max-width: none;
  }

  .page--auto .stats {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .page--auto .stats .stat {
    width: auto;
    box-sizing: border-box;
    padding: 10px;
  }

  .page--auto .stats .stat__label {
    margin-bottom: 4px;
    font-size: 11px;
  }

  .page--auto .stats .stat__value {
    font-size: 26px;
    line-height: 1.1;
  }

  .page--auto .stats .progress {
    margin-top: 8px;
  }
}
</style>
