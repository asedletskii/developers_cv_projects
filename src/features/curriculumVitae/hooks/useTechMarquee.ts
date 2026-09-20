import { useLayoutEffect, useRef } from 'react'
import { gsap, prefersReducedMotion, ScrollTrigger } from '@shared/animations'

/**
 * Бесконечная бегущая строка: трек состоит из двух одинаковых половин.
 * При быстрой прокрутке страницы лента ускоряется.
 */
export const useTechMarquee = (speed = 60) => {
  const trackRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const track = trackRef.current
    if (!track || prefersReducedMotion()) return

    const tween = gsap.to(track, {
      xPercent: -50,
      duration: track.scrollWidth / 2 / speed,
      ease: 'none',
      repeat: -1,
    })

    // Пока курсор над лентой, она замирает — иначе по бегущему слову не попасть
    let isHovered = false

    const onEnter = () => {
      isHovered = true
      gsap.to(tween, { timeScale: 0, duration: 0.4, overwrite: true })
    }

    const onLeave = () => {
      isHovered = false
      gsap.to(tween, { timeScale: 1, duration: 0.6, overwrite: true })
    }

    track.addEventListener('pointerenter', onEnter)
    track.addEventListener('pointerleave', onLeave)

    const scrollTrigger = ScrollTrigger.create({
      start: 0,
      end: 'max',
      onUpdate: (self) => {
        if (isHovered) return

        const boost = gsap.utils.clamp(0, 8, Math.abs(self.getVelocity()) / 350)

        gsap.to(tween, {
          timeScale: 1 + boost,
          duration: 0.25,
          overwrite: true,
          onComplete: () => {
            gsap.to(tween, { timeScale: 1, duration: 1.2, ease: 'power2.out' })
          },
        })
      },
    })

    return () => {
      track.removeEventListener('pointerenter', onEnter)
      track.removeEventListener('pointerleave', onLeave)
      scrollTrigger.kill()
      gsap.killTweensOf(tween)
      tween.kill()
    }
  }, [speed])

  return trackRef
}
