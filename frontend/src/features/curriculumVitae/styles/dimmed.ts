export const dimmed = (isDimmed: boolean) => ({
  filter: isDimmed ? 'grayscale(0.8) opacity(0.35)' : 'none',
  transition: 'filter .3s',
  '@media (prefers-reduced-motion: reduce)': { transition: 'none' },
})
