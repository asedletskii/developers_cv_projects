import { useLayoutEffect, useRef } from 'react'
import { gsap, prefersReducedMotion } from '@shared/animations'

export const useStepIconPop = (index: number) => {
  const iconRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    if (!iconRef.current || prefersReducedMotion()) return

    gsap.fromTo(
      iconRef.current,
      { scale: 0.4, rotate: -20, opacity: 0 },
      { scale: 1, rotate: 0, opacity: 1, duration: 0.5, ease: 'back.out(2)' },
    )
  }, [index])

  return iconRef
}
