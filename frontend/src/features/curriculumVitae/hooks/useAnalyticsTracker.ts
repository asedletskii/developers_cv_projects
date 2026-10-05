import { useEffect, useMemo } from 'react'
import { API_ROUTES } from '@features/curriculumVitae/constants'
import { AnalyticsEventType, SectionId } from '@features/curriculumVitae/types'
import {
  createAnalyticsTracker,
  isSectionVisibleEnough,
  pickMostVisibleSection,
} from '@features/curriculumVitae/utils'
import { api } from '@shared/api'
import { generatePath } from 'react-router'

const SECTION_IDS = new Set<string>(Object.values(SectionId))
const SECTION_VIEW_DELAY_MS = 1000
const FULL_SCALE = 10_000
const OBSERVER_THRESHOLDS = Array.from({ length: 11 }, (_, i) => i / 10)

const getSections = () =>
  [...document.querySelectorAll<HTMLElement>('section[id]')].filter(({ id }) => SECTION_IDS.has(id))

const toFullScale = (value: number, size: number) =>
  Math.min(Math.max(Math.round((value / size) * FULL_SCALE), 0), FULL_SCALE)

export const useAnalyticsTracker = (
  linkId: string | undefined,
  viewId: string | null,
  sectionOrder: string[],
) => {
  const orderKey = sectionOrder.join(',')

  const tracker = useMemo(
    () =>
      linkId && viewId
        ? createAnalyticsTracker({
            endpoint: `${api.defaults.baseURL ?? ''}${generatePath(API_ROUTES.cvEvents, { linkId })}`,
            viewId,
            sectionOrder: orderKey ? orderKey.split(',') : [],
            getExitSection: () =>
              pickMostVisibleSection(
                getSections().map((element) => {
                  const { top, bottom } = element.getBoundingClientRect()
                  return { id: element.id, top, bottom }
                }),
                window.innerHeight,
              ),
          })
        : null,
    [linkId, viewId, orderKey],
  )

  // Время на странице и отправка: таймер, скрытие вкладки, уход со страницы
  useEffect(() => {
    if (!tracker) return

    const syncVisibility = () => {
      const isVisible = document.visibilityState === 'visible'
      tracker.setVisible(isVisible)
      // Мобильный Safari часто ничего не шлёт при закрытии вкладки — сбрасываем при скрытии
      if (!isVisible) tracker.flush()
    }
    const onPageHide = () => tracker.flush({ final: true })
    // Печать и «сохранить как PDF»: после диалога сразу отправляем — вкладку могут закрыть
    const onBeforePrint = () => tracker.track({ type: AnalyticsEventType.Print })
    const onAfterPrint = () => tracker.flush()
    const onPageShow = (event: PageTransitionEvent) => {
      // Возврат из bfcache — тот же визит, таймер снова тикает
      if (event.persisted) syncVisibility()
    }

    syncVisibility()
    tracker.start()
    document.addEventListener('visibilitychange', syncVisibility)
    window.addEventListener('pagehide', onPageHide)
    window.addEventListener('pageshow', onPageShow)
    window.addEventListener('beforeprint', onBeforePrint)
    window.addEventListener('afterprint', onAfterPrint)

    return () => {
      document.removeEventListener('visibilitychange', syncVisibility)
      window.removeEventListener('pagehide', onPageHide)
      window.removeEventListener('pageshow', onPageShow)
      window.removeEventListener('beforeprint', onBeforePrint)
      window.removeEventListener('afterprint', onAfterPrint)
      tracker.flush({ final: true })
      tracker.stop()
    }
  }, [tracker])

  // Дочитывание: секция видна достаточно и дольше секунды — одно событие `section_view`
  useEffect(() => {
    if (!tracker) return

    const timers = new Map<string, number>()
    const seen = new Set<string>()

    const observer = new IntersectionObserver(
      (entries) => {
        for (const { target, intersectionRatio, intersectionRect } of entries) {
          const { id } = target

          const isEnough = isSectionVisibleEnough(
            intersectionRatio,
            intersectionRect.height,
            window.innerHeight,
          )

          if (!isEnough) {
            window.clearTimeout(timers.get(id))
            timers.delete(id)
          } else if (!seen.has(id) && !timers.has(id)) {
            timers.set(
              id,
              window.setTimeout(() => {
                timers.delete(id)
                if (document.visibilityState !== 'visible') return

                seen.add(id)
                tracker.track({ type: AnalyticsEventType.SectionView, section: id })
              }, SECTION_VIEW_DELAY_MS),
            )
          }
        }
      },
      { threshold: OBSERVER_THRESHOLDS },
    )

    getSections().forEach((element) => observer.observe(element))

    return () => {
      observer.disconnect()
      timers.forEach((timer) => window.clearTimeout(timer))
    }
  }, [tracker])

  // Click-map: один делегированный обработчик; координаты — доля секции, а не пиксели
  useEffect(() => {
    if (!tracker) return

    const onPointerDown = (event: PointerEvent) => {
      if (!(event.target instanceof Element)) return

      const section = event.target.closest<HTMLElement>('section[id]')
      if (!section || !SECTION_IDS.has(section.id)) return

      const rect = section.getBoundingClientRect()
      if (!rect.width || !rect.height) return

      tracker.trackClick({
        section: section.id,
        x: toFullScale(event.clientX - rect.left, rect.width),
        y: toFullScale(event.clientY - rect.top, rect.height),
        vw: window.innerWidth,
        target: event.target.closest<HTMLElement>('[data-analytics-id]')?.dataset.analyticsId,
      })
    }

    document.addEventListener('pointerdown', onPointerDown, { capture: true })
    return () => document.removeEventListener('pointerdown', onPointerDown, { capture: true })
  }, [tracker])

  return tracker
}
