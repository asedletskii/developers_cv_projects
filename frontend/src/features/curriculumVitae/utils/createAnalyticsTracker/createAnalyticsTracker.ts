import { AnalyticsClick, AnalyticsEvent, AnalyticsEventType } from '@features/curriculumVitae/types'

/** Лимиты из docs/analytics/frontend.md: бэкенд всё равно отсекает лишнее, но не шлём заведомо отклоняемое */
export const FLUSH_INTERVAL_MS = 15_000
export const MAX_EVENTS_PER_BATCH = 100
export const MAX_CLICKS_PER_VISIT = 200
export const MAX_BODY_BYTES = 32 * 1024
export const CLICK_THROTTLE_MS = 100
/** Некоторые браузеры стреляют `beforeprint` дважды за одну печать */
export const PRINT_DEDUPE_MS = 2000

export type TrackedEvent = AnalyticsEvent
export type TrackedClick = AnalyticsClick

type Options = {
  /** Полный адрес `POST /cv/:linkId/events` */
  endpoint: string
  viewId: string
  /** Порядок секций на странице: по нему определяем «самую дальнюю» */
  sectionOrder: string[]
  /** Секция, занимающая большую часть экрана в момент отправки — «точка ухода» */
  getExitSection?: () => string | undefined
  now?: () => number
  sendBeacon?: (url: string, data: string) => boolean
  fetchImpl?: typeof fetch
  newId?: () => string
}

type QueuedEvent = TrackedEvent & { id: string; at: number }

/** Копит события посещения и отправляет их батчами. Ничего не знает про DOM и React —
 *  страница подключает его через `useAnalyticsTracker`, поэтому его удобно тестировать. */
export const createAnalyticsTracker = ({
  endpoint,
  viewId,
  sectionOrder,
  getExitSection,
  now = Date.now,
  sendBeacon = typeof navigator !== 'undefined' && navigator.sendBeacon
    ? (url, data) => navigator.sendBeacon(url, data)
    : undefined,
  fetchImpl = typeof fetch !== 'undefined' ? fetch : undefined,
  newId = () => crypto.randomUUID(),
}: Options) => {
  let events: QueuedEvent[] = []
  let clicks: (TrackedClick & { at: number })[] = []
  let clicksTotal = 0
  let lastClickAt = -Infinity

  let seq = 0
  let timer: ReturnType<typeof setInterval> | undefined

  let accumulatedMs = 0
  let visibleSince: number | null = null
  let lastSentActiveMs = 0

  let lastPrintAt = -Infinity

  const seenSections = new Set<string>()
  let maxSectionIndex = -1

  const getActiveMs = () =>
    Math.round(accumulatedMs + (visibleSince === null ? 0 : now() - visibleSince))

  const track = (event: TrackedEvent) => {
    if (event.type === AnalyticsEventType.SectionView) {
      // Один раз за визит на секцию; «самая дальняя» — по порядку секций на странице
      if (seenSections.has(event.section)) return
      seenSections.add(event.section)
      maxSectionIndex = Math.max(maxSectionIndex, sectionOrder.indexOf(event.section))
    }

    if (event.type === AnalyticsEventType.Print) {
      if (now() - lastPrintAt < PRINT_DEDUPE_MS) return
      lastPrintAt = now()
    }

    events.push({ ...event, id: newId(), at: now() })
  }

  const trackClick = (click: TrackedClick) => {
    const at = now()
    if (clicksTotal >= MAX_CLICKS_PER_VISIT || at - lastClickAt < CLICK_THROTTLE_MS) return

    lastClickAt = at
    clicksTotal += 1
    clicks.push({ ...click, at })
  }

  /** Вкладка стала видимой/скрытой: активное время тикает, только пока она видна */
  const setVisible = (isVisible: boolean) => {
    if (isVisible && visibleSince === null) visibleSince = now()

    if (!isVisible && visibleSince !== null) {
      accumulatedMs += now() - visibleSince
      visibleSince = null
    }
  }

  const post = (body: string) => {
    if (sendBeacon?.(endpoint, body)) return

    // text/plain, как у sendBeacon: без preflight; keepalive — чтобы запрос пережил закрытие вкладки
    fetchImpl?.(endpoint, {
      method: 'POST',
      body,
      keepalive: true,
      credentials: 'include',
      headers: { 'Content-Type': 'text/plain' },
    }).catch(() => undefined)
  }

  const buildBody = (batchEvents: QueuedEvent[], batchClicks: typeof clicks, isFinal: boolean) => {
    seq += 1
    const maxSection = sectionOrder[maxSectionIndex]

    return JSON.stringify({
      viewId,
      seq,
      final: isFinal,
      state: { activeMs: getActiveMs(), maxSection, exitSection: getExitSection?.() },
      events: batchEvents,
      clicks: batchClicks,
    })
  }

  /** Отправляет накопленное. Без новых данных ничего не шлёт, кроме финального сброса при уходе. */
  const flush = ({ final = false } = {}) => {
    const activeMs = getActiveMs()
    const hasNothingNew = !events.length && !clicks.length && activeMs === lastSentActiveMs
    if (hasNothingNew && !final) return

    lastSentActiveMs = activeMs

    do {
      let batchEvents = events.slice(0, MAX_EVENTS_PER_BATCH)
      let batchClicks = clicks
      let body = buildBody(batchEvents, batchClicks, final && events.length <= batchEvents.length)

      // Слишком большое тело сервер отклонит целиком, поэтому урезаем, а не теряем весь батч
      while (body.length > MAX_BODY_BYTES && (batchEvents.length > 1 || batchClicks.length > 1)) {
        batchEvents = batchEvents.slice(0, Math.ceil(batchEvents.length / 2))
        batchClicks = batchClicks.slice(0, Math.ceil(batchClicks.length / 2))
        seq -= 1
        body = buildBody(batchEvents, batchClicks, false)
      }

      events = events.slice(batchEvents.length)
      clicks = clicks.slice(batchClicks.length)
      post(body)
    } while (events.length > 0)
  }

  const start = () => {
    timer ??= setInterval(() => flush(), FLUSH_INTERVAL_MS)
  }

  const stop = () => {
    clearInterval(timer)
    timer = undefined
  }

  return { track, trackClick, setVisible, flush, start, stop }
}

export type AnalyticsTracker = ReturnType<typeof createAnalyticsTracker>
