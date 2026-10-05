import { createContext, PropsWithChildren, useMemo } from 'react'
import { useAnalyticsTracker } from '@features/curriculumVitae/hooks/useAnalyticsTracker'
import { isDoNotTrackEnabled, TrackedEvent } from '@features/curriculumVitae/utils'
import { getRuntimeConfig } from '@shared/utils'

export type AnalyticsValue = {
  track: (event: TrackedEvent) => void
  /** Идентификатор визита — его же передаём при раскрытии контакта */
  viewId: string | null
}

/** По умолчанию — «пустая» аналитика: компоненты зовут `track` всегда и не проверяют, включена ли она
 *  (предпросмотр в кабинете, ошибка записи визита на бэке, выключенная ссылка). */
export const AnalyticsContext = createContext<AnalyticsValue>({ track: () => {}, viewId: null })

type Props = PropsWithChildren<{
  linkId?: string
  viewId: string | null
  sectionOrder: string[]
}>

export const AnalyticsProvider = ({ linkId, viewId, sectionOrder, children }: Props) => {
  // Браузер просит не отслеживать — клиентский трекер не запускаем (сам факт открытия сервер всё равно
  // считает, но без времени, кликов и действий)
  const isAllowed = !(getRuntimeConfig().respectDoNotTrack && isDoNotTrackEnabled())
  const tracker = useAnalyticsTracker(linkId, isAllowed ? viewId : null, sectionOrder)

  const value = useMemo<AnalyticsValue>(
    () => ({ track: (event) => tracker?.track(event), viewId: tracker ? viewId : null }),
    [tracker, viewId],
  )

  return <AnalyticsContext.Provider value={value}>{children}</AnalyticsContext.Provider>
}
