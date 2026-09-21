import { useLayoutEffect, useRef } from 'react'
import { gsap } from '@shared/animations'

export const useGsap = <T extends HTMLElement>(animation: () => void) => {
  const ref = useRef<T | null>(null)

  useLayoutEffect(() => {
    if (!ref.current) return

    const ctx = gsap.context(() => {
      animation()
    }, ref)

    return () => ctx.revert()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return ref
}
