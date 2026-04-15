export interface Metrics {
  cpuPercent: number
  ramPercent: number
  updatedAt: string
}

export interface ServiceStatus {
  name: string
  status: 'up' | 'down' | 'degraded'
}

export interface AlertItem {
  id: string
  title: string
  severity: 'warning' | 'critical'
  message: string
}
