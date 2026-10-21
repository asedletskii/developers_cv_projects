import { keyframes } from '@mui/material/styles'

/** кольцо расходится от значка Telegram, как сигнал «ждём подтверждения» */
export const PULSE_RING = keyframes`
  0% { transform: scale(0.7); opacity: 0.6; }
  100% { transform: scale(1.8); opacity: 0; }
`

/** точка бежит по пунктиру от Telegram к кабинету */
export const TRAVEL = keyframes`
  0% { left: 0%; opacity: 0; }
  15%, 85% { opacity: 1; }
  100% { left: 100%; opacity: 0; }
`

/** строки макета кабинета дорисовываются слева направо */
export const DRAW = keyframes`
  from { transform: scaleX(0); }
  to { transform: scaleX(1); }
`

/** значок ошибки коротко «вздрагивает» при появлении */
export const SHAKE = keyframes`
  0%, 100% { transform: translateX(0); }
  20% { transform: translateX(-8px); }
  40% { transform: translateX(7px); }
  60% { transform: translateX(-4px); }
  80% { transform: translateX(3px); }
`

/** карточка и блоки мягко выезжают снизу */
export const RISE = keyframes`
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: none; }
`

export const STILL = { '@media (prefers-reduced-motion: reduce)': { animation: 'none' } } as const
