import { useLayoutEffect, useRef } from 'react'
import { gsap, prefersReducedMotion } from '@shared/animations'

/** Полоска прогресса растёт от 0 до 1 по мере прокрутки страницы. */
export const useScrollProgress = () => {
  const barRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const bar = barRef.current
    if (!bar || prefersReducedMotion()) return

    const tween = gsap.fromTo(
      bar,
      { scaleX: 0 },
      {
        scaleX: 1,
        ease: 'none',
        scrollTrigger: { start: 0, end: 'max', scrub: 0.3 },
      },
    )

    return () => {
      tween.scrollTrigger?.kill()
      tween.kill()
    }
  }, [])

  return barRef
}
