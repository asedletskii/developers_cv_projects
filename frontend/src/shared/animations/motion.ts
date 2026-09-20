export const prefersReducedMotion = () =>
  window.matchMedia('(prefers-reduced-motion: reduce)').matches

export const canHover = () => window.matchMedia('(hover: hover)').matches
