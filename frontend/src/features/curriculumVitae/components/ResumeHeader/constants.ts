import { glassCard } from '@features/curriculumVitae/styles'

export const PREVIEW_BUTTON_SX = {
  ...glassCard,
  borderRadius: 999,
  px: { xs: 1.5, sm: 2 },
  minWidth: 0,
  textTransform: 'none',
  '& .MuiButton-startIcon': { mr: { xs: 0, sm: 1 }, ml: { xs: 0, sm: -0.5 } },
} as const
