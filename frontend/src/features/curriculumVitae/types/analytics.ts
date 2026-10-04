import { DeviceType, TrafficType } from './enums'

/** Типы повторяют JSON из docs/analytics/backend.md (раздел «Кабинет») один в один. */

export type LinkStatistics = {
  views: number
  reads: number
  uniqueVisitors: number
  returningVisitors: number
  contactViews: number
  downloads: number
  firstViewAt: string | null
  lastViewAt: string | null
  medianActiveSeconds: number
  devices: Record<DeviceType, number>
  referrers: { host: string | null; views: number }[]
  filtered: { bot: number; owner: number; suspect: number }
}

/** Статистика по всему резюме: то же + разбивка по ссылкам */
export type ResumeStatistics = LinkStatistics & {
  links: {
    id: string
    name: string
    views: number
    reads: number
    lastViewAt: string | null
  }[]
}

export type DailyPoint = {
  date: string
  views: number
  reads: number
  uniqueVisitors: number
}

export type VisitEvent = {
  type: string
  at: string
  contact?: string
  technology?: string
}

export type LinkVisit = {
  id: string
  startedAt: string
  visitor: { ordinal: number; isReturning: boolean }
  device: DeviceType
  browser: string
  os: string
  referrerHost: string | null
  isRead: boolean
  activeSeconds: number
  maxSection: string | null
  exitSection: string | null
  events: VisitEvent[]
  traffic: TrafficType
}

export type LinkVisitsPage = {
  items: LinkVisit[]
  nextCursor: string | null
}

export type AnalyticsPeriod = 7 | 30 | 90

export type SectionsStatistics = {
  reads: number
  sections: { section: string; reached: number; exits: number }[]
}

export type TechnologyInterest = { name: string; selections: number; visitors: number }

export type ClickBreakpoint = 'xs' | 'md' | 'lg'

export type ClickGrid = {
  section: string
  /** Сетка `grid × grid`, клетка — доля секции по ширине и высоте */
  grid: number
  cells: { x: number; y: number; count: number }[]
  total: number
  /** Клики по прежним версиям секции — в сетку не вошли (см. `includeStale`) */
  staleTotal: number
}
