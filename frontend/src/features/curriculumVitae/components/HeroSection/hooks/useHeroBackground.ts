import { useLayoutEffect, useRef } from 'react'
import { createSpring, gsap, prefersReducedMotion, Spring } from '@shared/animations'

/** Живой фон первого экрана: плавающие пятна, реакция на мышь и параллакс при прокрутке. */
export const useHeroBackground = () => {
  const rootRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const root = rootRef.current
    if (!root || prefersReducedMotion()) return

    let onMove: ((event: PointerEvent) => void) | undefined
    const springs: Spring[] = []

    const ctx = gsap.context(() => {
      gsap.utils.toArray<HTMLElement>('.orb', root).forEach((orb) => {
        gsap.to(orb, {
          x: 'random(-80, 80)',
          y: 'random(-80, 80)',
          scale: 'random(0.85, 1.25)',
          duration: 'random(6, 11)',
          ease: 'sine.inOut',
          repeat: -1,
          repeatRefresh: true,
        })
      })

      gsap.to('.hero-bg', {
        yPercent: 25,
        ease: 'none',
        scrollTrigger: {
          trigger: root.parentElement,
          start: 'top top',
          end: 'bottom top',
          scrub: true,
        },
      })

      const layers = gsap.utils.toArray<HTMLElement>('.orb-parallax', root).map((element) => ({
        depth: Number(element.dataset.depth),
        springX: createSpring(0, (value) => gsap.set(element, { x: value }), { response: 1.1 }),
        springY: createSpring(0, (value) => gsap.set(element, { y: value }), { response: 1.1 }),
      }))

      springs.push(...layers.flatMap(({ springX, springY }) => [springX, springY]))

      onMove = (event) => {
        const nx = event.clientX / window.innerWidth - 0.5
        const ny = event.clientY / window.innerHeight - 0.5

        layers.forEach(({ depth, springX, springY }) => {
          springX.to(-nx * depth * 3)
          springY.to(-ny * depth * 3)
        })
      }
    }, rootRef)

    if (onMove) window.addEventListener('pointermove', onMove)

    return () => {
      if (onMove) window.removeEventListener('pointermove', onMove)
      springs.forEach((spring) => spring.kill())
      ctx.revert()
    }
  }, [])

  return rootRef
}
