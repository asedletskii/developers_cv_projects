import { keyframes } from '@mui/material/styles'

/** лупа качается из стороны в сторону, как будто ищет */
export const SWEEP = keyframes`
  0%, 100% { transform: rotate(-14deg) translateX(-4px); }
  50% { transform: rotate(14deg) translateX(4px); }
`

/** облако мигает, как плохая связь */
export const FLICKER = keyframes`
  0%, 100% { opacity: 1; }
  45% { opacity: 1; }
  50% { opacity: 0.35; }
  55% { opacity: 1; }
  70% { opacity: 0.55; }
  75% { opacity: 1; }
`
