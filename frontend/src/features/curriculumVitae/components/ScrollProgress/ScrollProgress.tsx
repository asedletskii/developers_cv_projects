import { useLayoutEffect, useRef } from 'react'
import { Box } from '@mui/material'
import { gsap, prefersReducedMotion } from '@shared/animations'

/** Тонкая полоска прогресса прокрутки страницы вверху экрана. */
export const ScrollProgress = () => {
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

  return (
    <Box
      aria-hidden
      data-print-hidden
      sx={{ position: 'fixed', top: 0, left: 0, right: 0, height: 3, zIndex: 1300 }}
    >
      <Box
        ref={barRef}
        sx={{
          height: '100%',
          transformOrigin: 'left center',
          background: ({ palette }) =>
            `linear-gradient(90deg, ${palette.primary.main}, ${palette.secondary.main})`,
        }}
      />
    </Box>
  )
}
