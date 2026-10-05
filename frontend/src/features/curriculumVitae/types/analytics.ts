import { AnalyticsEventType, DeviceType, TrafficType } from './enums'

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
  kind?: string
  target?: string
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
  grid: number
  cells: { x: number; y: number; count: number }[]
  total: number
  staleTotal: number
}

export type AnalyticsEvent =
  | { type: AnalyticsEventType.SectionView; section: string }
  | { type: AnalyticsEventType.Download }
  | { type: AnalyticsEventType.TechFilter; technology: string; action: 'select' | 'clear' }
  | { type: AnalyticsEventType.ExternalClick; kind: string; target: string }
  | { type: AnalyticsEventType.Print }
  | { type: AnalyticsEventType.ContactCopy; contact: string }
  | { type: AnalyticsEventType.ContactClick; contact: string }

export type AnalyticsClick = {
  section: string
  x: number
  y: number
  vw: number
  target?: string
}

export type AnalyticsBatch = {
  viewId: string
  seq: number
  final: boolean
  state: {
    activeMs: number
    maxSection?: string
    exitSection?: string
  }
  events: (AnalyticsEvent & { id: string; at: number })[]
  clicks: (AnalyticsClick & { at: number })[]
}
