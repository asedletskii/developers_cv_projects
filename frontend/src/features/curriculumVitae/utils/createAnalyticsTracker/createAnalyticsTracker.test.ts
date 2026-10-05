import { AnalyticsEventType } from '@features/curriculumVitae/types'
import {
  CLICK_THROTTLE_MS,
  createAnalyticsTracker,
  FLUSH_INTERVAL_MS,
  MAX_CLICKS_PER_VISIT,
  MAX_EVENTS_PER_BATCH,
  PRINT_DEDUPE_MS,
} from './createAnalyticsTracker'

const ORDER = ['about', 'experience', 'projects', 'contacts']

const setup = (overrides: Partial<Parameters<typeof createAnalyticsTracker>[0]> = {}) => {
  const sent: any[] = []
  let counter = 0
  const tracker = createAnalyticsTracker({
    endpoint: '/api/v1/cv/link/events',
    viewId: 'view-1',
    sectionOrder: ORDER,
    sendBeacon: (_, body) => (sent.push(JSON.parse(body)), true),
    newId: () => `id-${(counter += 1)}`,
    ...overrides,
  })

  return { tracker, sent }
}

describe('createAnalyticsTracker', () => {
  beforeEach(() => {
    jest.useFakeTimers()
    jest.setSystemTime(0)
  })

  afterEach(() => jest.useRealTimers())

  it('шлёт накопленные события батчем и растит seq', () => {
    const { tracker, sent } = setup()

    tracker.track({ type: AnalyticsEventType.Download })
    tracker.flush()
    tracker.track({
      type: AnalyticsEventType.TechFilter,
      technology: 'Go',
      action: 'select',
    })
    tracker.flush()

    expect(sent.map((batch) => batch.seq)).toEqual([1, 2])
    expect(sent[0].viewId).toBe('view-1')
    expect(sent[0].events).toEqual([{ id: 'id-1', type: 'download', at: 0 }])
    expect(sent[1].events[0]).toMatchObject({ type: 'tech_filter', technology: 'Go' })
  })

  it('без новых данных ничего не шлёт', () => {
    const { tracker, sent } = setup()

    tracker.flush()

    expect(sent).toHaveLength(0)
  })

  it('активное время растёт только пока вкладка видна', () => {
    const { tracker, sent } = setup()

    tracker.setVisible(true)
    jest.advanceTimersByTime(3000)
    tracker.setVisible(false)
    jest.advanceTimersByTime(10_000)
    tracker.flush()

    expect(sent[0].state.activeMs).toBe(3000)
  })

  it('отправляет, когда выросло только активное время', () => {
    const { tracker, sent } = setup()

    tracker.setVisible(true)
    jest.advanceTimersByTime(2000)
    tracker.flush()
    tracker.flush()

    expect(sent).toHaveLength(1)
  })

  it('section_view один раз за секцию, а maxSection — самая дальняя по порядку страницы', () => {
    const { tracker, sent } = setup()

    tracker.track({ type: AnalyticsEventType.SectionView, section: 'projects' })
    tracker.track({ type: AnalyticsEventType.SectionView, section: 'about' })
    tracker.track({ type: AnalyticsEventType.SectionView, section: 'projects' })
    tracker.flush()

    expect(sent[0].events).toHaveLength(2)
    expect(sent[0].state.maxSection).toBe('projects')
  })

  it('финальный сброс помечается final и уходит даже без новых данных', () => {
    const { tracker, sent } = setup()

    tracker.flush({ final: true })

    expect(sent).toHaveLength(1)
    expect(sent[0].final).toBe(true)
  })

  it('делит больше 100 событий на несколько батчей, final — только на последнем', () => {
    const { tracker, sent } = setup()

    for (let i = 0; i < MAX_EVENTS_PER_BATCH + 20; i += 1) {
      tracker.track({ type: AnalyticsEventType.Download })
    }
    tracker.flush({ final: true })

    expect(sent.map((batch) => batch.events.length)).toEqual([MAX_EVENTS_PER_BATCH, 20])
    expect(sent.map((batch) => batch.final)).toEqual([false, true])
  })

  it('уходит через fetch с keepalive, если sendBeacon недоступен или отказал', () => {
    const fetchImpl = jest.fn().mockResolvedValue({})
    const { tracker } = setup({ sendBeacon: () => false, fetchImpl })

    tracker.track({ type: AnalyticsEventType.Download })
    tracker.flush()

    expect(fetchImpl).toHaveBeenCalledWith(
      '/api/v1/cv/link/events',
      expect.objectContaining({ method: 'POST', keepalive: true }),
    )
  })

  it('клики: не чаще раза в 100 мс и не больше 200 за визит', () => {
    const { tracker, sent } = setup()
    const click = { section: 'about', x: 1, y: 1, vw: 1280 }

    tracker.trackClick(click)
    tracker.trackClick(click)
    jest.advanceTimersByTime(CLICK_THROTTLE_MS)
    tracker.trackClick(click)
    tracker.flush()

    expect(sent[0].clicks).toHaveLength(2)

    for (let i = 0; i < MAX_CLICKS_PER_VISIT + 10; i += 1) {
      jest.advanceTimersByTime(CLICK_THROTTLE_MS)
      tracker.trackClick(click)
    }
    tracker.flush()

    expect(sent[0].clicks.length + sent[1].clicks.length).toBe(MAX_CLICKS_PER_VISIT)
  })

  it('периодический флаш каждые 15 секунд после start()', () => {
    const { tracker, sent } = setup()

    tracker.start()
    tracker.track({ type: AnalyticsEventType.Download })
    jest.advanceTimersByTime(FLUSH_INTERVAL_MS)
    tracker.stop()
    tracker.track({ type: AnalyticsEventType.Download })
    jest.advanceTimersByTime(FLUSH_INTERVAL_MS * 2)

    expect(sent).toHaveLength(1)
  })

  it('двойной beforeprint считается одной печатью, а повторная через время — новой', () => {
    const { tracker, sent } = setup()

    tracker.track({ type: AnalyticsEventType.Print })
    tracker.track({ type: AnalyticsEventType.Print })
    jest.advanceTimersByTime(PRINT_DEDUPE_MS)
    tracker.track({ type: AnalyticsEventType.Print })
    tracker.flush()

    expect(sent[0].events.map((event: { type: string }) => event.type)).toEqual(['print', 'print'])
  })

  it('точка ухода берётся в момент отправки', () => {
    let exit = 'about'
    const { tracker, sent } = setup({ getExitSection: () => exit })

    tracker.track({ type: AnalyticsEventType.Download })
    exit = 'projects'
    tracker.flush()

    expect(sent[0].state.exitSection).toBe('projects')
  })
})
