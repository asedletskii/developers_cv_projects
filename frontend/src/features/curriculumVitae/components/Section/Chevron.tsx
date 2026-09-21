import { ExpandMoreRounded } from '@mui/icons-material'
import { DURATION, EASE } from './constants'

type Props = {
  isCollapsed: boolean
}

/** Стрелка-индикатор аккордеона: вниз у свёрнутой секции, вверх у развёрнутой. */
export const Chevron = ({ isCollapsed }: Props) => (
  <ExpandMoreRounded
    data-print-hidden
    sx={{
      color: 'gray.light',
      transform: isCollapsed ? 'none' : 'rotate(180deg)',
      transition: `transform ${DURATION} ${EASE}`,
      '@media (prefers-reduced-motion: reduce)': { transition: 'none' },
    }}
  />
)
