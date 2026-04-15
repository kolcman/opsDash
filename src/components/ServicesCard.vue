<script setup lang="ts">
import { computed } from 'vue'
import type { ServiceStatus } from '../types/dashboard'

const props = defineProps<{
  services: ServiceStatus[]
}>()

const statusClassMap: Record<ServiceStatus['status'], string> = {
  up: 'status status--up',
  down: 'status status--down',
  degraded: 'status status--degraded',
}

const totalLabel = computed(() => `${props.services.length} services`)
</script>

<template>
  <article class="card" aria-label="Services status card">
    <div class="card__head">
      <h2>Services</h2>
      <span class="chip chip--muted">{{ totalLabel }}</span>
    </div>

    <p v-if="services.length === 0" class="empty-state">No services configured yet.</p>
    <ul v-else class="list" aria-label="Service statuses">
      <li v-for="service in services" :key="service.name" class="list__row">
        <span class="service-name">{{ service.name }}</span>
        <strong :class="statusClassMap[service.status]">{{ service.status }}</strong>
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

.status--degraded {
  color: #fde68a;
  border-color: rgba(251, 191, 36, 0.55);
  background: rgba(161, 98, 7, 0.22);
}

.empty-state {
  margin: 0;
  color: rgba(203, 213, 225, 0.8);
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

  .list__row {
    padding: 10px 0;
  }

  .status {
    padding: 5px 8px;
    font-size: 9px;
  }
}
</style>
