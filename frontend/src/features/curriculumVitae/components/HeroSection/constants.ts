import { AvailabilityStatus } from '@features/curriculumVitae/types'
import { keyframes } from '@mui/material/styles'

/** Пульсация индикатора доступности */
export const PULSE = keyframes`
  0% { box-shadow: 0 0 0 0 currentColor; opacity: 1; }
  70% { box-shadow: 0 0 0 8px transparent; opacity: 0.7; }
  100% { box-shadow: 0 0 0 0 transparent; opacity: 1; }
`

/** Цвет индикатора: зелёный — ищу, жёлтый — рассматриваю, серый — не ищу */
export const STATUS_COLORS: Record<AvailabilityStatus, string> = {
  [AvailabilityStatus.Open]: 'success.main',
  [AvailabilityStatus.Considering]: 'alert.main',
  [AvailabilityStatus.Busy]: 'gray.light',
}

/** Колёсико индикатора «листайте вниз» */
export const WHEEL = keyframes`
  0% { transform: translateY(0); opacity: 1; }
  70% { transform: translateY(14px); opacity: 0; }
  100% { transform: translateY(0); opacity: 0; }
`

/** Переливающийся градиент выделенной части приветствия */
export const SHIMMER = keyframes`
  from { background-position: 0% 50%; }
  to { background-position: 200% 50%; }
`

/** Световые пятна фона: размер, глубина параллакса, цвет, прозрачность и позиция */
export const ORBS: {
  size: number
  depth: number
  color: string
  opacity: number
  position: Partial<Record<'top' | 'bottom' | 'left' | 'right', string>>
}[] = [
  {
    size: 460,
    depth: 30,
    color: 'primary.main',
    opacity: 0.32,
    position: { top: '-10%', left: '-8%' },
  },
  {
    size: 380,
    depth: 55,
    color: 'secondary.main',
    opacity: 0.26,
    position: { top: '25%', right: '-8%' },
  },
  {
    size: 320,
    depth: 20,
    color: 'primary.light',
    opacity: 0.28,
    position: { bottom: '-8%', left: '28%' },
  },
]
