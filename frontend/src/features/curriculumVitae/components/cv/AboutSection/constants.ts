import { keyframes } from '@mui/material/styles'

/** портрет плавно «дышит» вверх-вниз */
export const FLOAT = keyframes`
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
`

/** пунктирное кольцо вокруг портрета медленно вращается */
export const SPIN = keyframes`
  to { transform: rotate(360deg); }
`
