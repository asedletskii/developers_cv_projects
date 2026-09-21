import { useLayoutEffect, useRef } from 'react'
import { gsap, prefersReducedMotion } from '@shared/animations'

/** Линия таймлайна «дорисовывается» вместе со скроллом, точки и карточки появляются по очереди. */
export const useTimelineAnimation = () => {
  const containerRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const root = containerRef.current
    if (!root || prefersReducedMotion()) return

    const ctx = gsap.context(() => {
      gsap.fromTo(
        '.timeline-fill',
        { scaleY: 0 },
        {
          scaleY: 1,
          ease: 'none',
          scrollTrigger: { trigger: root, start: 'top 70%', end: 'bottom 70%', scrub: 0.6 },
        },
      )

      gsap.utils.toArray<HTMLElement>('.timeline-item', root).forEach((item) => {
        gsap
          .timeline({ scrollTrigger: { trigger: item, start: 'top 85%', once: true } })
          .from(item.querySelector('.timeline-dot'), {
            scale: 0,
            duration: 0.6,
            ease: 'power3.out',
          })
          .from(
            item.querySelector('.timeline-card'),
            { x: -48, opacity: 0, duration: 0.8, ease: 'power3.out' },
            '<0.1',
          )
      })
    }, containerRef)

    return () => ctx.revert()
  }, [])

  return containerRef
}
