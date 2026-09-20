import { SelectField } from '@features/curriculumVitae/components/editor/form/SelectField'
import { LanguageGrade } from '@features/curriculumVitae/types'
import {
  DeleteOutlineRounded,
  KeyboardArrowDownRounded,
  KeyboardArrowUpRounded,
} from '@mui/icons-material'
import { IconButton, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  index: number
  languageOptions: { value: string; label: string }[]
  canMoveUp: boolean
  canMoveDown: boolean
  onMoveUp: () => void
  onMoveDown: () => void
  onRemove: () => void
}

export const LanguageRow = ({
  index,
  languageOptions,
  canMoveUp,
  canMoveDown,
  onMoveUp,
  onMoveDown,
  onRemove,
}: Props) => {
  const translate = useTranslate('Lk.Edit.About')
  const translateEdit = useTranslate('Lk.Edit')

  return (
    <Stack direction="row" alignItems="flex-end" spacing={1}>
      <SelectField
        name={`languages.${index}.language`}
        clearable={false}
        label={translate('language')}
        options={languageOptions}
        sx={{ flex: 1, minWidth: 140 }}
      />

      <SelectField
        name={`languages.${index}.grade`}
        label={translate('languageGrade')}
        options={Object.values(LanguageGrade).map((grade) => ({ value: grade, label: grade }))}
        sx={{ width: 96 }}
      />

      <IconButton
        size="small"
        aria-label={translateEdit('moveUp')}
        disabled={!canMoveUp}
        onClick={onMoveUp}
      >
        <KeyboardArrowUpRounded />
      </IconButton>

      <IconButton
        size="small"
        aria-label={translateEdit('moveDown')}
        disabled={!canMoveDown}
        onClick={onMoveDown}
      >
        <KeyboardArrowDownRounded />
      </IconButton>

      <IconButton size="small" aria-label={translateEdit('remove')} onClick={onRemove}>
        <DeleteOutlineRounded />
      </IconButton>
    </Stack>
  )
}
