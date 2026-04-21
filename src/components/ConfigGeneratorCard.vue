<script setup lang="ts">
import { computed, ref } from 'vue'

type ConfigKind = 'nginx' | 'docker-compose' | 'alertmanager'

const configKind = ref<ConfigKind>('nginx')
const appName = ref('opsdash')
const domain = ref('example.com')
const upstreamHost = ref('backend')
const upstreamPort = ref('8080')
const imageTag = ref('latest')
const webhookUrl = ref('https://hooks.slack.com/services/replace/me')
const copied = ref(false)
const props = defineProps<{
  lang: 'en' | 'ru'
}>()

const configPreview = computed(() => {
  if (configKind.value === 'nginx') {
    return `server {
  listen 80;
  server_name ${domain.value};

  location / {
    proxy_pass http://${upstreamHost.value}:${upstreamPort.value};
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }
}`
  }

  if (configKind.value === 'docker-compose') {
    return `services:
  ${appName.value}-web:
    image: ghcr.io/your-org/${appName.value}-web:${imageTag.value}
    ports:
      - "80:80"
    restart: unless-stopped

  ${appName.value}-api:
    image: ghcr.io/your-org/${appName.value}-api:${imageTag.value}
    environment:
      - PORT=${upstreamPort.value}
    restart: unless-stopped`
  }

  return `route:
  receiver: 'slack-notifications'
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 30s
  repeat_interval: 2h

receivers:
  - name: 'slack-notifications'
    webhook_configs:
      - url: '${webhookUrl.value}'`
})

async function copyPreview() {
  try {
    await navigator.clipboard.writeText(configPreview.value)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 1400)
  } catch {
    copied.value = false
  }
}
</script>

<template>
  <article class="card card--generator" aria-label="Config generator card">
    <div class="card__head">
      <h2>{{ lang === 'ru' ? 'Генератор конфигов' : 'Config Generator' }}</h2>
      <span class="chip chip--muted">{{ lang === 'ru' ? 'DevOps-инструмент' : 'DevOps tool' }}</span>
    </div>

    <div class="intro">
      <p>
        {{
          lang === 'ru'
            ? 'Этот блок помогает быстро собрать рабочий конфиг без ручного написания с нуля.'
            : 'This tool helps you quickly generate a working config without writing it manually from scratch.'
        }}
      </p>
      <p>
        {{
          lang === 'ru'
            ? 'Выберите тип файла, заполните поля и скопируйте готовый шаблон. Это безопасно: вы меняете только текст предпросмотра.'
            : 'Choose a config type, fill the fields, and copy the generated template. It is safe: you only modify preview text.'
        }}
      </p>
    </div>

    <div class="form-grid">
      <label class="field">
        <span class="field__label">
          {{ lang === 'ru' ? 'Тип конфига' : 'Config type' }}
          <span
            class="help"
            tabindex="0"
            :data-tip="
              lang === 'ru'
                ? 'Выбери, какой файл нужен: nginx, docker-compose или alertmanager.'
                : 'Pick which file you need: nginx, docker-compose, or alertmanager.'
            "
            >?</span
          >
        </span>
        <select v-model="configKind">
          <option value="nginx">nginx</option>
          <option value="docker-compose">docker-compose</option>
          <option value="alertmanager">alertmanager.yml</option>
        </select>
      </label>

      <label class="field">
        <span class="field__label">
          {{ lang === 'ru' ? 'Имя приложения' : 'App name' }}
          <span
            class="help"
            tabindex="0"
            :data-tip="
              lang === 'ru'
                ? 'Простое имя проекта. Попадет в имена сервисов и Docker-образов.'
                : 'Simple project name. It will be used in service and Docker image names.'
            "
            >?</span
          >
        </span>
        <input v-model="appName" type="text" placeholder="opsdash" />
      </label>

      <label class="field">
        <span class="field__label">
          {{ lang === 'ru' ? 'Домен' : 'Domain' }}
          <span
            class="help"
            tabindex="0"
            :data-tip="
              lang === 'ru'
                ? 'Адрес сайта, например mysite.ru. Нужен для nginx.'
                : 'Website address, for example mysite.com. Used in nginx config.'
            "
            >?</span
          >
        </span>
        <input v-model="domain" type="text" placeholder="example.com" />
      </label>

      <label class="field">
        <span class="field__label">
          {{ lang === 'ru' ? 'Upstream хост' : 'Upstream host' }}
          <span
            class="help"
            tabindex="0"
            :data-tip="
              lang === 'ru'
                ? 'Имя сервиса или IP, куда отправлять запросы. Обычно backend.'
                : 'Service name or IP to forward requests to. Usually backend.'
            "
            >?</span
          >
        </span>
        <input v-model="upstreamHost" type="text" placeholder="backend" />
      </label>

      <label class="field">
        <span class="field__label">
          {{ lang === 'ru' ? 'Upstream порт' : 'Upstream port' }}
          <span
            class="help"
            tabindex="0"
            :data-tip="
              lang === 'ru'
                ? 'Порт сервиса, обычно 8080 или 3000.'
                : 'Service port, usually 8080 or 3000.'
            "
            >?</span
          >
        </span>
        <input v-model="upstreamPort" type="text" placeholder="8080" />
      </label>

      <label class="field">
        <span class="field__label">
          {{ lang === 'ru' ? 'Тег образа' : 'Image tag' }}
          <span
            class="help"
            tabindex="0"
            :data-tip="
              lang === 'ru'
                ? 'Версия Docker-образа. Пример: latest или v1.2.0.'
                : 'Docker image version. Example: latest or v1.2.0.'
            "
            >?</span
          >
        </span>
        <input v-model="imageTag" type="text" placeholder="latest" />
      </label>

      <label class="field field--wide">
        <span class="field__label">
          {{ lang === 'ru' ? 'Webhook URL (alertmanager)' : 'Webhook URL (alertmanager)' }}
          <span
            class="help"
            tabindex="0"
            :data-tip="
              lang === 'ru'
                ? 'Ссылка для уведомлений. Сюда Alertmanager шлет алерты (например в Slack).'
                : 'Notification URL. Alertmanager sends alerts here (for example Slack).'
            "
            >?</span
          >
        </span>
        <input v-model="webhookUrl" type="text" placeholder="https://hooks.slack.com/services/..." />
      </label>
    </div>

    <div class="preview">
      <div class="preview__head">
        <strong>{{ lang === 'ru' ? 'Предпросмотр' : 'Live preview' }}</strong>
        <button type="button" class="copy-btn" @click="copyPreview">
          {{ copied ? (lang === 'ru' ? 'Скопировано' : 'Copied') : lang === 'ru' ? 'Копировать' : 'Copy' }}
        </button>
      </div>
      <pre>{{ configPreview }}</pre>
    </div>
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
  flex: 1 1 760px;
  max-width: 784px;
  min-height: 280px;
  display: flex;
  flex-direction: column;
}

.card--generator {
  overflow: visible;
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

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.intro {
  margin: 0 0 12px;
  border: 1px solid rgba(148, 163, 184, 0.2);
  border-radius: 12px;
  background: rgba(15, 23, 42, 0.42);
  padding: 10px 12px;
}

.intro p {
  margin: 0;
  font-size: 12px;
  line-height: 1.45;
  color: rgba(226, 232, 240, 0.82);
}

.intro p + p {
  margin-top: 6px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-height: 70px;
}

.field__label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: rgba(203, 213, 225, 0.76);
}

.help {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 999px;
  border: 1px solid rgba(148, 163, 184, 0.45);
  font-size: 10px;
  line-height: 1;
  color: rgba(148, 163, 184, 0.92);
  cursor: help;
  position: relative;
}

.help::after {
  content: attr(data-tip);
  position: absolute;
  left: 0;
  bottom: calc(100% + 8px);
  transform: none;
  width: 230px;
  max-width: 72vw;
  padding: 8px 10px;
  border-radius: 10px;
  border: 1px solid rgba(148, 163, 184, 0.32);
  background: rgba(2, 6, 23, 0.96);
  color: rgba(226, 232, 240, 0.95);
  font-size: 11px;
  line-height: 1.35;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.18s ease;
  z-index: 20;
}

.field:nth-child(even) .help::after {
  left: auto;
  right: 0;
  transform: none;
}

.field--wide .help::after {
  left: 0;
  right: auto;
  transform: none;
}

.help:hover::after,
.help:focus-visible::after {
  opacity: 1;
}

.field input,
.field select {
  border: 1px solid rgba(148, 163, 184, 0.25);
  border-radius: 10px;
  background: rgba(15, 23, 42, 0.48);
  color: #e2e8f0;
  padding: 8px 10px;
  height: 38px;
  box-sizing: border-box;
  outline: none;
}

.field input:focus,
.field select:focus {
  border-color: rgba(125, 211, 252, 0.75);
}

.field--wide {
  grid-column: 1 / -1;
}

.preview {
  margin-top: 12px;
  border: 1px solid rgba(148, 163, 184, 0.22);
  border-radius: 12px;
  background: rgba(15, 23, 42, 0.5);
  padding: 10px;
}

.preview__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.copy-btn {
  border: 1px solid rgba(125, 211, 252, 0.45);
  border-radius: 999px;
  background: rgba(56, 189, 248, 0.16);
  color: #bae6fd;
  padding: 5px 10px;
  font-size: 11px;
  cursor: pointer;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  overflow-x: auto;
  font-size: 12px;
  line-height: 1.45;
  color: rgba(226, 232, 240, 0.92);
}

@media (max-width: 720px) {
  .card {
    border-radius: 18px;
    flex: 1 1 100%;
    max-width: none;
    padding: 14px;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }

  .help::after,
  .field:nth-child(even) .help::after,
  .field--wide .help::after {
    left: 0;
    right: auto;
    transform: none;
    width: min(230px, calc(100vw - 64px));
  }
}
</style>
