import axios from 'axios'
import type { AlertItem, Metrics, ServiceStatus } from '../types/dashboard'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api',
})

async function getJson<T>(url: string): Promise<T> {
  const response = await apiClient.get<T>(url)
  return response.data
}

export function fetchMetrics() {
  return getJson<Metrics>('/metrics')
}

export function fetchServices() {
  return getJson<ServiceStatus[]>('/services')
}

export function fetchAlerts() {
  return getJson<AlertItem[]>('/alerts')
}
