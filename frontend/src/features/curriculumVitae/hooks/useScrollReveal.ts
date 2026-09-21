import { useLayoutEffect, useRef } from 'react'
import { gsap, prefersReducedMotion } from '@shared/animations'

/**
 * Проявление при прокрутке:
 * `.reveal` — блок целиком, `.section-line` — прорисовка линии заголовка,
 * `.reveal-group` — каскадно проявляет вложенные `.reveal-item`.
 *
 * При reduced motion остаётся только мягкий кроссфейд без движения.
 */
export const useScrollReveal = () => {
  const containerRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const root = containerRef.current
    if (!root) return

    const reduced = prefersReducedMotion()
    const shift = reduced ? 0 : 1

    const ctx = gsap.context(() => {
      const trigger = (element: Element, start = 'top 88%') => ({
        trigger: element,
        start,
        once: true,
      })

      gsap.utils.toArray<HTMLElement>('.reveal', root).forEach((element) => {
        gsap.from(element, {
          y: 56 * shift,
          opacity: 0,
          duration: reduced ? 0.4 : 1,
          ease: 'power3.out',
          scrollTrigger: trigger(element),
        })
      })

      if (!reduced) {
        gsap.utils.toArray<HTMLElement>('.section-line', root).forEach((element) => {
          gsap.from(element, {
            scaleX: 0,
            transformOrigin: 'left center',
            duration: 1.4,
            ease: 'power3.inOut',
            scrollTrigger: trigger(element, 'top 92%'),
          })
        })
      }

      gsap.utils.toArray<HTMLElement>('.reveal-group', root).forEach((group) => {
        gsap.from(group.querySelectorAll('.reveal-item'), {
          y: 40 * shift,
          opacity: 0,
          duration: reduced ? 0.4 : 0.8,
          stagger: reduced ? 0 : 0.09,
          ease: 'power3.out',
          scrollTrigger: trigger(group, 'top 85%'),
        })
      })
    }, containerRef)

    return () => ctx.revert()
  }, [])

  return containerRef
}
