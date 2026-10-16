import { RefObject, useCallback, useLayoutEffect, useRef, useState } from 'react'
import { gsap, prefersReducedMotion, ScrollTrigger } from '@shared/animations'

export const useExpandableJobs = (timelineRef: RefObject<HTMLDivElement | null>) => {
  const [isExpanded, setIsExpanded] = useState(false)
  const buttonRef = useRef<HTMLButtonElement | null>(null)

  const expand = useCallback(() => setIsExpanded(true), [])

  const toggle = () => {
    if (!isExpanded) {
      setIsExpanded(true)
      return
    }

    const collapse = () => {
      setIsExpanded(false)
      ScrollTrigger.refresh(true)
      buttonRef.current?.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
    }

    const extras = timelineRef.current?.querySelectorAll('.timeline-extra')

    if (!extras?.length || prefersReducedMotion()) {
      collapse()
      return
    }

    gsap.to(extras, {
      y: -16,
      opacity: 0,
      duration: 0.3,
      stagger: { each: 0.06, from: 'end' },
      ease: 'power2.in',
      onComplete: collapse,
    })
  }

  useLayoutEffect(() => {
    if (!isExpanded) return

    const extras = timelineRef.current?.querySelectorAll('.timeline-extra')

    if (!extras?.length || prefersReducedMotion()) {
      ScrollTrigger.refresh(true)
      return
    }

    const tween = gsap.from(extras, {
      y: 32,
      opacity: 0,
      duration: 0.6,
      stagger: 0.1,
      ease: 'power3.out',
      onComplete: () => ScrollTrigger.refresh(true),
    })

    return () => {
      tween.kill()
    }
  }, [isExpanded, timelineRef])

  return { isExpanded, toggle, expand, buttonRef }
}
