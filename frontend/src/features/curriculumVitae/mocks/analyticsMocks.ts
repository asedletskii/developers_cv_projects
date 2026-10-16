import {
  ClickBreakpoint,
  ClickGrid,
  DailyPoint,
  DeviceType,
  LinkStatistics,
  LinkVisit,
  LinkVisitsPage,
  ResumeLink,
  ResumeStatistics,
  SectionId,
  SectionsStatistics,
  TechnologyInterest,
  TrafficType,
} from '@features/curriculumVitae/types'
import axios, {
  AxiosAdapter,
  AxiosInstance,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from 'axios'

/** Моки ручек аналитики из docs/analytics/backend.md — пока бэкенд их не реализовал.
 *  Включаются только в dev через `VITE_MOCK_ANALYTICS=true`; все остальные запросы идут как обычно. */

// Всего визитов вместе с ботами и владельцем; людей среди них меньше
const VISITS_TOTAL = 56
const PAGE_SIZE = 20
const HOUR = 3_600_000
const DAY = 24 * HOUR

// Детерминированный «рандом» — мок одинаково выглядит после перезагрузки
const rand = (seed: number) => {
  const x = Math.sin(seed * 9301 + 49297) * 233280
  return x - Math.floor(x)
}

const pick = <T>(items: T[], seed: number) => items[Math.floor(rand(seed) * items.length)]

const eachDay = (from: string, to: string) => {
  const days: string[] = []
  const end = new Date(`${to}T12:00:00`)

  for (let day = new Date(`${from}T12:00:00`); day <= end; day.setDate(day.getDate() + 1)) {
    days.push(day.toLocaleDateString('sv-SE'))
  }

  return days
}

export const buildDaily = (from: string, to: string): DailyPoint[] =>
  eachDay(from, to).map((date, index) => {
    // Сид — номер дня, а не позиция в периоде: у предыдущего периода получаются свои числа
    const dayNumber = Math.floor(new Date(`${date}T12:00:00`).getTime() / DAY)
    const views = Math.floor(rand(dayNumber + 1) * rand(dayNumber + 7) * 8)
    const reads = Math.floor(views * 0.6)

    return { date, views, reads, uniqueVisitors: Math.min(views, 1 + (index % 3)) }
  })

export const buildStatistics = (from: string, to: string): LinkStatistics => {
  const days = buildDaily(from, to)
  const sum = (key: 'views' | 'reads') => days.reduce((total, day) => total + day[key], 0)
  const views = sum('views')

  return {
    views,
    reads: sum('reads'),
    uniqueVisitors: Math.ceil(views * 0.6),
    returningVisitors: Math.floor(views * 0.15),
    contactViews: Math.floor(views * 0.2),
    downloads: Math.floor(views * 0.1),
    firstViewAt: new Date(Date.now() - 40 * DAY).toISOString(),
    lastViewAt: new Date(Date.now() - 2 * HOUR).toISOString(),
    medianActiveSeconds: 94,
    devices: {
      desktop: Math.ceil(views * 0.55),
      mobile: Math.floor(views * 0.4),
      tablet: Math.floor(views * 0.05),
      unknown: 0,
    },
    referrers: [
      { host: 't.me', views: Math.ceil(views * 0.6) },
      { host: 'linkedin.com', views: Math.floor(views * 0.25) },
      { host: null, views: Math.floor(views * 0.15) },
    ],
    filtered: { bot: 5, owner: 3, suspect: 0 },
  }
}

const SECTIONS = [SectionId.About, SectionId.Experience, SectionId.Projects, SectionId.Contacts]

const after = (start: Date, seconds: number) =>
  new Date(start.getTime() + seconds * 1000).toISOString()

const getTraffic = (index: number) =>
  index % 7 === 3 ? TrafficType.Bot : index % 11 === 5 ? TrafficType.Owner : TrafficType.Human

const buildVisit = (index: number): LinkVisit => {
  const isRead = rand(index + 3) > 0.35
  const device = pick([DeviceType.Desktop, DeviceType.Mobile, DeviceType.Mobile], index + 5)
  const reached = Math.floor(rand(index + 11) * SECTIONS.length)
  const startedAt = new Date(Date.now() - (index * 7 + Math.floor(rand(index) * 5)) * HOUR)

  return {
    id: `mock-visit-${index}`,
    startedAt: startedAt.toISOString(),
    visitor: { ordinal: (index % 5) + 1, isReturning: index % 4 === 0 },
    device,
    browser: device === DeviceType.Mobile ? 'safari' : pick(['chrome', 'firefox'], index + 2),
    os: device === DeviceType.Mobile ? 'ios' : pick(['macos', 'windows'], index + 4),
    referrerHost: pick(['t.me', 'linkedin.com', null], index + 9),
    isRead,
    activeSeconds: isRead ? 20 + Math.floor(rand(index + 13) * 200) : 2,
    maxSection: isRead ? SECTIONS[reached] : null,
    exitSection: isRead ? SECTIONS[Math.max(reached - (index % 2), 0)] : null,
    events:
      isRead && index % 3 === 0
        ? [
            { type: 'contact_view', at: after(startedAt, 41), contact: 'telegram' },
            { type: 'contact_copy', at: after(startedAt, 44), contact: 'telegram' },
            { type: 'contact_click', at: after(startedAt, 52), contact: 'telegram' },
            { type: 'tech_filter', at: after(startedAt, 78), technology: 'Go' },
            {
              type: 'external_click',
              at: after(startedAt, 95),
              kind: 'project',
              target: 'github.com',
            },
          ]
        : isRead && index % 5 === 0
          ? [{ type: 'print', at: after(startedAt, 63) }]
          : [],
    traffic: getTraffic(index),
  }
}

export const buildVisitsPage = (cursor?: string, includeFiltered = false): LinkVisitsPage => {
  const offset = Number(cursor) || 0
  const all = Array.from({ length: VISITS_TOTAL }, (_, index) => buildVisit(index))
  // Боты и владелец по умолчанию в ленту не попадают — как у настоящего бэкенда
  const visits = includeFiltered ? all : all.filter(({ traffic }) => traffic === TrafficType.Human)
  const items = visits.slice(offset, offset + PAGE_SIZE)

  return {
    items,
    nextCursor: offset + PAGE_SIZE < visits.length ? String(offset + PAGE_SIZE) : null,
  }
}

const SECTION_ORDER = [
  SectionId.About,
  SectionId.Experience,
  SectionId.Skills,
  SectionId.Education,
  SectionId.Projects,
  SectionId.Contacts,
]

export const buildSections = (from: string, to: string): SectionsStatistics => {
  const { reads } = buildStatistics(from, to)

  return {
    reads,
    sections: SECTION_ORDER.map((section, index) => {
      const reached = Math.max(Math.round(reads * (1 - index * 0.14)), 0)
      return { section, reached, exits: Math.min(Math.round(reached * 0.25), reached) }
    }),
  }
}

export const buildTechnologies = (): TechnologyInterest[] =>
  ['Go', 'React', 'PostgreSQL', 'TypeScript', 'Docker'].map((name, index) => ({
    name,
    selections: 9 - index * 2 + (index === 4 ? 1 : 0),
    visitors: Math.max(5 - index, 1),
  }))

const CLICK_GRID = 20
const HOTSPOTS = [
  { x: 0.25, y: 0.3, weight: 1 },
  { x: 0.7, y: 0.55, weight: 0.6 },
  { x: 0.5, y: 0.85, weight: 0.35 },
]

export const buildClicks = (
  section: string,
  breakpoint: ClickBreakpoint,
  includeStale = false,
): ClickGrid => {
  // Чем уже экран, тем плотнее клики прижаты к центру — как у настоящих мобильных кликов
  const spread = breakpoint === 'xs' ? 0.1 : 0.16
  const seed = section.length * 7 + (breakpoint === 'xs' ? 1 : breakpoint === 'md' ? 2 : 3)
  const cells: ClickGrid['cells'] = []

  for (let x = 0; x < CLICK_GRID; x++) {
    for (let y = 0; y < CLICK_GRID; y++) {
      const heat = HOTSPOTS.reduce((sum, spot, index) => {
        const dx = (x + 0.5) / CLICK_GRID - ((spot.x + rand(seed + index) * 0.2) % 1)
        const dy = (y + 0.5) / CLICK_GRID - spot.y
        return sum + spot.weight * Math.exp(-(dx * dx + dy * dy) / (2 * spread * spread))
      }, 0)
      const count = Math.round(heat * 6)
      if (count > 0) cells.push({ x, y, count })
    }
  }

  return {
    section,
    grid: CLICK_GRID,
    cells,
    total: cells.reduce((sum, cell) => sum + cell.count, 0),
    // Часть кликов пришла на прошлую версию секции: по умолчанию в сетку не входит
    staleTotal: includeStale ? 0 : 12,
  }
}

type Handler = (
  config: InternalAxiosRequestConfig,
  params: Record<string, string>,
) => Promise<unknown> | unknown

const ROUTES: [RegExp, Handler][] = [
  [/^\/lk\/link\/[^/]+\/statistics\/daily$/, (_, { from, to }) => ({ days: buildDaily(from, to) })],
  [/^\/lk\/link\/[^/]+\/statistics\/sections$/, (_, { from, to }) => buildSections(from, to)],
  [/^\/lk\/link\/[^/]+\/statistics\/technologies$/, () => ({ technologies: buildTechnologies() })],
  [
    /^\/lk\/link\/[^/]+\/clicks$/,
    (_, { section, breakpoint, includeStale }) =>
      buildClicks(section, breakpoint as ClickBreakpoint, String(includeStale) === 'true'),
  ],
  [/^\/lk\/link\/[^/]+\/statistics$/, (_, { from, to }) => buildStatistics(from, to)],
  [
    /^\/lk\/link\/[^/]+\/views$/,
    (_, { cursor, includeFiltered }) => buildVisitsPage(cursor, String(includeFiltered) === 'true'),
  ],
]

export const installAnalyticsMocks = (api: AxiosInstance) => {
  const realAdapter = axios.getAdapter(api.defaults.adapter)

  const respond = (config: InternalAxiosRequestConfig, data: unknown): AxiosResponse => ({
    data,
    status: 200,
    statusText: 'OK',
    headers: {},
    config,
  })

  const adapter: AxiosAdapter = async (config) => {
    const path = config.url ?? ''
    const params = (config.params ?? {}) as Record<string, string>

    // Разбивка по ссылкам строится из настоящих ссылок резюме, чтобы клик по строке вёл в рабочий отчёт
    const resumeMatch = path.match(/^\/lk\/resume\/([^/]+)\/statistics$/)
    if (config.method === 'get' && resumeMatch) {
      // Без бэкенда (он не запущен) берём одну выдуманную ссылку, чтобы мок работал автономно
      const links = await realAdapter({
        ...config,
        url: '/lk/links',
        params: { resumeId: resumeMatch[1] },
      })
        .then(({ data: raw }) => (typeof raw === 'string' ? JSON.parse(raw) : raw) as ResumeLink[])
        .catch(() => [{ id: 'mock-link', name: 'Яндекс, рекрутер Анна' } as ResumeLink])
      const statistics: ResumeStatistics = {
        ...buildStatistics(params.from, params.to),
        links: links.map((link, index) => ({
          id: link.id,
          name: link.name,
          views: 12 - index * 4,
          reads: 7 - index * 2,
          lastViewAt: new Date(Date.now() - (index + 1) * 5 * HOUR).toISOString(),
        })),
      }
      return respond(config, statistics)
    }

    // Приём батча событий: бэкенд ничего не отдаёт, только 204
    if (config.method === 'post' && /^\/cv\/[^/]+\/events$/.test(path)) {
      return { ...respond(config, undefined), status: 204, statusText: 'No Content' }
    }

    const route = config.method === 'get' && ROUTES.find(([pattern]) => pattern.test(path))
    if (route) return respond(config, await route[1](config, params))

    // Бэкенд пока не выдаёт `X-View-Id` — подставляем свой, чтобы трекер на публичной странице работал
    if (config.method === 'get' && /^\/cv\/[^/]+$/.test(path)) {
      const response = await realAdapter(config)
      response.headers = { ...response.headers, 'x-view-id': crypto.randomUUID() }
      return response
    }

    // Действующий бэкенд читает тело строго (неизвестное поле — 400), а про `viewId` и `initiator` он
    // ещё не знает. Пока так, вырезаем их, а скачивание (`initiator: download`) подменяем выдуманными
    // контактами: настоящей ручки для пачки ещё нет.
    if (config.method === 'post' && /^\/cv\/[^/]+\/contact-views$/.test(path)) {
      const { viewId, initiator, ...body } = JSON.parse(String(config.data ?? '{}'))

      if (initiator === 'download') {
        return respond(config, {
          contacts: [
            { contact: 'telegram', value: 'mock_telegram' },
            { contact: 'email', value: 'mock@example.com' },
            { contact: 'phone', value: '+70000000000' },
          ],
        })
      }

      return realAdapter({ ...config, data: JSON.stringify(body) })
    }

    return realAdapter(config)
  }

  api.defaults.adapter = adapter
}
