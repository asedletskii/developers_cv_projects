import { useLayoutEffect, useRef } from 'react'
import { greetingAnimation } from '@features/curriculumVitae/animations'
import { gsap } from '@shared/animations'

export const useGreetingAnimation = () => {
  const containerRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const ctx = gsap.context(() => {
      greetingAnimation()
    }, containerRef)

    return () => ctx.revert()
  }, [])

  return containerRef
}
