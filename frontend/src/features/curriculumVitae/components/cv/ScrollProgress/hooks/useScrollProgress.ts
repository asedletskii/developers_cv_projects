import { useLayoutEffect, useRef } from 'react'
import { gsap, prefersReducedMotion } from '@shared/animations'

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
