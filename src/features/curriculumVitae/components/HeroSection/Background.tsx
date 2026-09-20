import { useLayoutEffect, useRef } from 'react'
import { Box } from '@mui/material'
import { createSpring, gsap, prefersReducedMotion, Spring } from '@shared/animations'
import { SparkleField } from '@features/curriculumVitae/components/SparkleField'

const ORBS: {
  size: number
  depth: number
  color: string
  opacity: number
  position: Partial<Record<'top' | 'bottom' | 'left' | 'right', string>>
}[] = [
  {
    size: 460,
    depth: 30,
    color: 'primary.main',
    opacity: 0.32,
    position: { top: '-10%', left: '-8%' },
  },
  {
    size: 380,
    depth: 55,
    color: 'secondary.main',
    opacity: 0.26,
    position: { top: '25%', right: '-8%' },
  },
  {
    size: 320,
    depth: 20,
    color: 'primary.light',
    opacity: 0.28,
    position: { bottom: '-8%', left: '28%' },
  },
]

/** Плавающие световые пятна: живут своей жизнью, реагируют на мышь и уезжают медленнее контента при скролле. */
export const Background = () => {
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

  return (
    <Box
      ref={rootRef}
      aria-hidden
      data-print-hidden
      sx={{ position: 'absolute', inset: 0, overflow: 'hidden', pointerEvents: 'none' }}
    >
      <Box className="hero-bg" sx={{ position: 'absolute', inset: 0 }}>
        {ORBS.map((orb, index) => (
          <Box
            key={index}
            className="orb-parallax"
            data-depth={orb.depth}
            sx={{ position: 'absolute', width: orb.size, height: orb.size, ...orb.position }}
          >
            <Box
              className="orb"
              sx={{
                width: '100%',
                height: '100%',
                borderRadius: '50%',
                bgcolor: orb.color,
                opacity: orb.opacity,
                filter: 'blur(80px)',
              }}
            />
          </Box>
        ))}

        <SparkleField count={44} />
      </Box>
    </Box>
  )
}
