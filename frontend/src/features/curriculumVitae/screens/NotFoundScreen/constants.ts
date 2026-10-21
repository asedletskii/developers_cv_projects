import { keyframes } from '@mui/material/styles'

/** цифры «404» по очереди выпрыгивают снизу */
export const POP = keyframes`
  0% { opacity: 0; transform: translateY(48px) scale(0.8); }
  70% { opacity: 1; transform: translateY(-6px) scale(1.03); }
  100% { opacity: 1; transform: none; }
`

/** точка обходит «нолик» по кругу, как спутник */
export const ORBIT = keyframes`
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
`

/** оборванная ссылка внутри нуля покачивается */
export const WOBBLE = keyframes`
  0%, 100% { transform: rotate(-10deg) translateY(0); }
  50% { transform: rotate(10deg) translateY(-6px); }
`

/** кольцо нуля медленно «дышит» */
export const BREATHE = keyframes`
  0%, 100% { transform: scale(1); opacity: 0.9; }
  50% { transform: scale(1.04); opacity: 1; }
`

/** текст и кнопки появляются после цифр */
export const RISE = keyframes`
  from { opacity: 0; transform: translateY(16px); }
  to { opacity: 1; transform: none; }
`

export const STILL = { '@media (prefers-reduced-motion: reduce)': { animation: 'none' } } as const
