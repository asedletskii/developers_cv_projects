import { Theme } from '@mui/material/styles'

/** Градиент акцентных цветов темы (основной → дополнительный) для фона, линий и текста. */
export const getBrandGradient = ({ palette }: Theme, angle = 90) =>
  `linear-gradient(${angle}deg, ${palette.primary.main}, ${palette.secondary.main})`
