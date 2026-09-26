import { ReactNode } from 'react'
import {
  DeleteOutlineRounded,
  KeyboardArrowDownRounded,
  KeyboardArrowUpRounded,
} from '@mui/icons-material'
import { IconButton, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  canMoveUp: boolean
  canMoveDown: boolean
  onMoveUp: () => void
  onMoveDown: () => void
  onRemove: () => void
  children: ReactNode
}

/** Одна повторяющаяся запись внутри секции (место работы, курс и т.д.): стрелки, удаление, форма. */
export const ArrayItemCard = ({
  canMoveUp,
  canMoveDown,
  onMoveUp,
  onMoveDown,
  onRemove,
  children,
}: Props) => {
  const translate = useTranslate('Lk.Edit')

  return (
    <Stack
      direction="row"
      spacing={1}
      sx={{ p: 2, border: 1, borderColor: 'divider', borderRadius: 2 }}
    >
      <Stack spacing={2} sx={{ flex: 1, minWidth: 0 }}>
        {children}
      </Stack>

      <Stack alignItems="center">
        <IconButton
          size="small"
          aria-label={translate('moveUp')}
          disabled={!canMoveUp}
          onClick={onMoveUp}
        >
          <KeyboardArrowUpRounded />
        </IconButton>

        <IconButton
          size="small"
          aria-label={translate('moveDown')}
          disabled={!canMoveDown}
          onClick={onMoveDown}
        >
          <KeyboardArrowDownRounded />
        </IconButton>

        <IconButton size="small" aria-label={translate('remove')} onClick={onRemove}>
          <DeleteOutlineRounded />
        </IconButton>
      </Stack>
    </Stack>
  )
}
