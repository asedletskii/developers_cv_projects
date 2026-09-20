import { useEffect, useState } from 'react'
import { prefersReducedMotion } from '@shared/animations'
import { STEP_DURATION_MS } from '../constants'

export const useStepCarousel = (count: number) => {
  const [index, setIndex] = useState(0)
  const [isPaused, setIsPaused] = useState(false)

  useEffect(() => {
    if (isPaused || prefersReducedMotion()) return

    const timer = window.setTimeout(() => setIndex((i) => (i + 1) % count), STEP_DURATION_MS)
    return () => window.clearTimeout(timer)
  }, [index, isPaused, count])

  return {
    index,
    next: () => setIndex((i) => (i + 1) % count),
    prev: () => setIndex((i) => (i - 1 + count) % count),
    select: setIndex,
    pause: () => setIsPaused(true),
    resume: () => setIsPaused(false),
  }
}
