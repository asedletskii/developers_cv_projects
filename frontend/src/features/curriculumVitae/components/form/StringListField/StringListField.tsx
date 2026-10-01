import { Resume } from '@features/curriculumVitae/types'
import { AddRounded, DeleteOutlineRounded } from '@mui/icons-material'
import { Button, IconButton, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useController, useFormContext } from 'react-hook-form'

type Props = {
  name: `experiences.${number}.tasks` | `experiences.${number}.achievements`
  addLabel: string
}

export const StringListField = ({ name, addLabel }: Props) => {
  const translate = useTranslate('Lk.Edit')
  const { control } = useFormContext<Resume>()
  const { field } = useController({ control, name })
  const items = field.value ?? []

  const updateItem = (index: number, value: string) => {
    const next = [...items]
    next[index] = value
    field.onChange(next)
  }

  const removeItem = (index: number) => {
    field.onChange(items.filter((_, i) => i !== index))
  }

  return (
    <Stack spacing={1}>
      {items.map((item, index) => (
        <Stack key={index} direction="row" spacing={1} alignItems="center">
          <TextField
            fullWidth
            size="small"
            value={item}
            onChange={(e) => updateItem(index, e.target.value)}
            onBlur={field.onBlur}
          />
          <IconButton
            size="small"
            aria-label={translate('remove')}
            onClick={() => removeItem(index)}
          >
            <DeleteOutlineRounded fontSize="small" />
          </IconButton>
        </Stack>
      ))}

      <Button
        size="small"
        startIcon={<AddRounded />}
        onClick={() => field.onChange([...items, ''])}
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {addLabel}
      </Button>
    </Stack>
  )
}
