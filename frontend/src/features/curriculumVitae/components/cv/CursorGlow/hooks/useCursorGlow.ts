import { useLayoutEffect, useRef } from 'react'
import { canHover, createSpring, gsap, prefersReducedMotion } from '@shared/animations'

export const useCursorGlow = () => {
  const glowRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const glow = glowRef.current
    if (!glow || prefersReducedMotion() || !canHover()) return

    const position = { x: window.innerWidth / 2, y: window.innerHeight / 2 }
    const apply = () => gsap.set(glow, { x: position.x, y: position.y })

    gsap.set(glow, { xPercent: -50, yPercent: -50 })
    apply()

    const springX = createSpring(
      position.x,
      (value) => {
        position.x = value
        apply()
      },
      { response: 0.6 },
    )

    const springY = createSpring(
      position.y,
      (value) => {
        position.y = value
        apply()
      },
      { response: 0.6 },
    )

    const onMove = (event: PointerEvent) => {
      springX.to(event.clientX)
      springY.to(event.clientY)
      gsap.to(glow, { opacity: 1, duration: 0.8, overwrite: 'auto' })
    }

    window.addEventListener('pointermove', onMove)

    return () => {
      window.removeEventListener('pointermove', onMove)
      springX.kill()
      springY.kill()
      gsap.killTweensOf(glow)
    }
  }, [])

  return glowRef
}
