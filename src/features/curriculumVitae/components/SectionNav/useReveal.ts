import { RefObject, useLayoutEffect } from 'react'
import { gsap, prefersReducedMotion, ScrollTrigger } from '@shared/animations'

/**
 * Навигация «материализуется» (появляется с масштабом), когда первая секция входит в экран,
 * и так же исчезает при возврате к hero.
 */
export const useReveal = (pillRef: RefObject<HTMLElement | null>, firstSectionId?: string) => {
  useLayoutEffect(() => {
    const pill = pillRef.current
    const firstSection = firstSectionId ? document.getElementById(firstSectionId) : null
    if (!pill || !firstSection) return

    const reduced = prefersReducedMotion()

    gsap.set(pill, { autoAlpha: 0, y: reduced ? 0 : -16, scale: reduced ? 1 : 0.96 })

    // Обратимая анимация: при смене направления скролла идёт от текущего состояния, а не с нуля
    const reveal = gsap.timeline({ paused: true }).to(pill, {
      autoAlpha: 1,
      y: 0,
      scale: 1,
      duration: reduced ? 0.2 : 0.5,
      ease: 'power3.out',
    })

    const trigger = ScrollTrigger.create({
      trigger: firstSection,
      start: 'top 65%',
      onEnter: () => reveal.play(),
      onLeaveBack: () => reveal.reverse(),
    })

    return () => {
      trigger.kill()
      reveal.kill()
    }
  }, [pillRef, firstSectionId])
}
