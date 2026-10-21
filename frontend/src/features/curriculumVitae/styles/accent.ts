import { alpha, Theme } from '@mui/material/styles'

const isLight = ({ palette }: Theme) => palette.mode === 'light'

/** Плитка под значок. Жёлтый значок на светлом фоне не читается, поэтому в светлой теме плитка
 *  жёлтая, а значок тёмный; в тёмной — серая плитка с жёлтым значком */
export const iconTile = {
  bgcolor: (theme: Theme) =>
    isLight(theme) ? theme.palette.primary.main : theme.palette.action.hover,
  color: (theme: Theme) =>
    isLight(theme) ? theme.palette.primary.contrastText : theme.palette.primary.main,
}

/** Бейдж: в светлой теме жёлтая заливка с тёмным текстом, в тёмной — жёлтый контур */
export const badge = {
  bgcolor: (theme: Theme) => (isLight(theme) ? theme.palette.primary.main : 'transparent'),
  color: (theme: Theme) =>
    isLight(theme) ? theme.palette.primary.contrastText : theme.palette.primary.main,
  borderColor: (theme: Theme) =>
    isLight(theme) ? 'transparent' : alpha(theme.palette.primary.main, 0.5),
}

/** Вторичная контурная кнопка: жёлтый контур и текст на светлом фоне не читаются, поэтому нейтральные */
export const secondaryButton = {
  color: 'text.primary',
  borderColor: 'divider',
  '&:hover': { borderColor: 'text.secondary', bgcolor: 'action.hover' },
}
