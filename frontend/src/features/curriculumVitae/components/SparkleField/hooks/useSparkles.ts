import { useLayoutEffect, useMemo, useRef } from 'react'
import { gsap, prefersReducedMotion } from '@shared/animations'
import { createRandom } from '@shared/utils'

type Params = {
  count: number

  /** Зерно раскладки: одно и то же зерно всегда даёт одинаковое расположение звёзд */
  seed: number
}

/** Раскладка звёзд по зерну и их мерцание/дрейф (GSAP) внутри контейнера, реф которого возвращается. */
export const useSparkles = ({ count, seed }: Params) => {
  const rootRef = useRef<HTMLDivElement | null>(null)

  const sparkles = useMemo(() => {
    const random = createRandom(seed)

    return Array.from({ length: count }, () => ({
      x: random() * 100,
      y: random() * 100,
      size: 2 + random() * 4,
    }))
  }, [count, seed])

  useLayoutEffect(() => {
    const root = rootRef.current
    if (!root || prefersReducedMotion()) return

    const ctx = gsap.context(() => {
      gsap.utils.toArray<HTMLElement>('.sparkle', root).forEach((element) => {
        gsap.fromTo(
          element,
          { opacity: 0, scale: 0.2 },
          {
            opacity: gsap.utils.random(0.4, 1),
            scale: gsap.utils.random(0.8, 1.7),
            duration: gsap.utils.random(1.2, 3),
            delay: Math.random() * 3,
            repeat: -1,
            yoyo: true,
            ease: 'sine.inOut',
          },
        )

        gsap.to(element, {
          x: gsap.utils.random(-24, 24),
          y: gsap.utils.random(-36, 36),
          duration: gsap.utils.random(4, 9),
          repeat: -1,
          yoyo: true,
          ease: 'sine.inOut',
        })
      })
    }, rootRef)

    return () => ctx.revert()
  }, [sparkles])

  return { rootRef, sparkles }
}
