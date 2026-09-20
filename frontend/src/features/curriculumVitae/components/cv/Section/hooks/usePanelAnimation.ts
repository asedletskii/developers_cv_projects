import { TransitionEvent, useEffect, useRef, useState } from 'react'
import { prefersReducedMotion, ScrollTrigger } from '@shared/animations'

export const usePanelAnimation = (isCollapsed: boolean) => {
  const [isAnimating, setIsAnimating] = useState(false)
  const previousCollapsedRef = useRef(isCollapsed)

  useEffect(() => {
    if (previousCollapsedRef.current === isCollapsed) return
    previousCollapsedRef.current = isCollapsed

    if (prefersReducedMotion()) {
      ScrollTrigger.refresh(true)
    } else {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setIsAnimating(true)
    }
  }, [isCollapsed])

  const handleTransitionEnd = (event: TransitionEvent) => {
    if (event.propertyName !== 'grid-template-rows') return

    setIsAnimating(false)
    ScrollTrigger.refresh(true)
  }

  return { isAnimating, handleTransitionEnd }
}
