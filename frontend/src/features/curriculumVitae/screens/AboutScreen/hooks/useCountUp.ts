import { useLayoutEffect, useRef } from 'react'
import { gsap, prefersReducedMotion } from '@shared/animations'

export const useCountUp = (to: number) => {
  const ref = useRef<HTMLSpanElement | null>(null)

  useLayoutEffect(() => {
    const element = ref.current
    if (!element || prefersReducedMotion()) return

    const counter = { value: 0 }
    element.textContent = '0'

    const ctx = gsap.context(() => {
      gsap.to(counter, {
        value: to,
        duration: 1.6,
        ease: 'power2.out',
        snap: { value: 1 },
        onUpdate: () => {
          element.textContent = String(counter.value)
        },
        scrollTrigger: { trigger: element, start: 'top 90%', once: true },
      })
    })

    return () => ctx.revert()
  }, [to])

  return ref
}
