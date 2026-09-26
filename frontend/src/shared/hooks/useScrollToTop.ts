import { useCallback } from 'react'

/** Прокрутка страницы наверх: плавно, а при `prefers-reduced-motion` — сразу. */
export const useScrollToTop = () =>
  useCallback(() => {
    const isReduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    window.scrollTo({ top: 0, behavior: isReduced ? 'auto' : 'smooth' })
  }, [])
