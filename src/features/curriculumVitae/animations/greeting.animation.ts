import { gsap } from '@shared/animations'

const hasTarget = (selector: string) => !!document.querySelector(selector)

/** Анимирует только те строки приветствия, которые реально есть на странице. */
export const greetingAnimation = () => {
  const tl = gsap.timeline()

  if (hasTarget('.letter')) {
    tl.from('.letter', {
      y: 120,
      opacity: 0,
      filter: 'blur(10px)',
      stagger: 0.05,
      duration: 1,
      ease: 'power4.out',
    })
  }

  if (hasTarget('.name')) {
    tl.from('.name', { y: 60, opacity: 0, duration: 0.9, ease: 'power3.out' }, '-=0.4')
  }

  if (hasTarget('.subtitle')) {
    tl.from('.subtitle', { y: 40, opacity: 0, duration: 0.8, ease: 'power2.out' }, '-=0.5')
  }

  if (hasTarget('.hero-availability')) {
    tl.from('.hero-availability', { y: 24, opacity: 0, duration: 0.8, ease: 'power3.out' }, '-=0.5')
  }

  if (hasTarget('.hero-action')) {
    tl.from('.hero-action', { y: 30, opacity: 0, duration: 0.8, ease: 'power3.out' }, '-=0.5')
  }

  tl.from('.hero-hint', { y: 20, opacity: 0, duration: 0.8, ease: 'power2.out' }, '-=0.2')

  return tl
}
