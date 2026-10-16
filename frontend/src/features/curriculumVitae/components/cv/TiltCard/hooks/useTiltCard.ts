import { useLayoutEffect, useRef } from 'react'
import { canHover, createSpring, gsap, prefersReducedMotion } from '@shared/animations'

export const useTiltCard = () => {
  const cardRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const card = cardRef.current
    if (!card || prefersReducedMotion() || !canHover()) return

    gsap.set(card, { transformPerspective: 900 })

    const rotateX = createSpring(0, (value) => gsap.set(card, { rotationX: value }), {
      response: 0.45,
    })
    const rotateY = createSpring(0, (value) => gsap.set(card, { rotationY: value }), {
      response: 0.45,
    })
    const scale = createSpring(1, (value) => gsap.set(card, { scale: value }), {
      response: 0.4,
      precision: 0.0005,
    })

    const onMove = (event: PointerEvent) => {
      const rect = card.getBoundingClientRect()
      const px = (event.clientX - rect.left) / rect.width
      const py = (event.clientY - rect.top) / rect.height

      rotateY.to((px - 0.5) * 12)
      rotateX.to((0.5 - py) * 12)

      card.style.setProperty('--mx', `${px * 100}%`)
      card.style.setProperty('--my', `${py * 100}%`)
    }

    const onEnter = () => scale.to(1.02)

    const onLeave = () => {
      rotateX.to(0)
      rotateY.to(0)
      scale.to(1)
    }

    card.addEventListener('pointermove', onMove)
    card.addEventListener('pointerenter', onEnter)
    card.addEventListener('pointerleave', onLeave)

    return () => {
      card.removeEventListener('pointermove', onMove)
      card.removeEventListener('pointerenter', onEnter)
      card.removeEventListener('pointerleave', onLeave)
      ;[rotateX, rotateY, scale].forEach((spring) => spring.kill())
    }
  }, [])

  return cardRef
}
