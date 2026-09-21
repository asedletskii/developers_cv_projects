import { useLayoutEffect, useMemo, useRef } from 'react'
import { Box, useTheme } from '@mui/material'
import { gsap, prefersReducedMotion } from '@shared/animations'
import { createRandom } from '@shared/utils'

type Props = {
  count?: number
  tone?: 'accent' | 'white'

  /** Зерно раскладки: одно и то же зерно всегда даёт одинаковое расположение звёзд */
  seed?: number
}

export const SparkleField = ({ count = 36, tone = 'accent', seed = 1 }: Props) => {
  const theme = useTheme()
  const rootRef = useRef<HTMLDivElement | null>(null)

  const color = tone === 'white' ? '#FFFFFF' : theme.palette.primary.light

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

  return (
    <Box
      ref={rootRef}
      aria-hidden
      sx={{ position: 'absolute', inset: 0, overflow: 'hidden', pointerEvents: 'none' }}
    >
      {sparkles.map((sparkle, index) => (
        <Box
          key={index}
          className="sparkle"
          sx={{
            position: 'absolute',
            left: `${sparkle.x}%`,
            top: `${sparkle.y}%`,
            width: sparkle.size,
            height: sparkle.size,
            borderRadius: '50%',
            opacity: 0.6,
            bgcolor: color,
            boxShadow: `0 0 ${sparkle.size * 3}px ${color}`,
          }}
        />
      ))}
    </Box>
  )
}
