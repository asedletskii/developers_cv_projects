import { useCallback, useLayoutEffect, useRef, useState } from 'react'
import { ScrollTrigger } from '@shared/animations'
import { READING_LINE, UNLOCK_EVENTS } from '../constants'

/**
 * Определяет текущую секцию по «линии чтения» (30% высоты экрана): активна та, что её пересекает.
 * В зазорах между секциями подсветка не меняется, у нижнего края страницы активна последняя.
 * После клика по пункту навигации выбор закреплён, пока идёт доскролл или пользователь не вмешался.
 */
export const useActiveSection = (ids: string[]) => {
  const [activeId, setActiveId] = useState<string | null>(null)
  const isLockedRef = useRef(false)

  useLayoutEffect(() => {
    const sections = ids
      .map((id) => document.getElementById(id))
      .filter((section): section is HTMLElement => !!section)

    if (!sections.length) return

    const update = () => {
      if (isLockedRef.current) return

      const line = window.innerHeight * READING_LINE
      const isAtBottom =
        window.innerHeight + window.scrollY >= document.documentElement.scrollHeight - 8

      if (isAtBottom) {
        setActiveId(sections[sections.length - 1].id)
        return
      }

      const current = sections.find((section) => {
        const { top, bottom } = section.getBoundingClientRect()
        return top <= line && bottom > line
      })

      if (current) {
        setActiveId(current.id)
      } else if (sections[0].getBoundingClientRect().top > line) {
        setActiveId(null)
      }
    }

    const trigger = ScrollTrigger.create({
      start: 0,
      end: 'max',
      onUpdate: update,
      onRefresh: update,
    })

    update()

    return () => trigger.kill()
  }, [ids])

  const select = useCallback((id: string) => {
    setActiveId(id)
    isLockedRef.current = true

    let timer = 0

    const unlock = () => {
      isLockedRef.current = false
      window.clearTimeout(timer)
      UNLOCK_EVENTS.forEach((event) => window.removeEventListener(event, unlock))
    }

    UNLOCK_EVENTS.forEach((event) => window.addEventListener(event, unlock, { passive: true }))
    timer = window.setTimeout(unlock, 2000)
  }, [])

  return { activeId, select }
}
