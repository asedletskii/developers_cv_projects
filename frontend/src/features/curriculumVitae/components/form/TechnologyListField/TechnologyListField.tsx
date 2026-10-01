import { useState } from 'react'
import { TECHNOLOGY_SUGGESTIONS } from '@features/curriculumVitae/constants'
import { AddRounded } from '@mui/icons-material'
import { Autocomplete, Chip, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useTechnologyListField } from './hooks'

type Props = {
  name: 'technologies' | `experiences.${number}.stack` | `projects.${number}.stack`
}

export const TechnologyListField = ({ name }: Props) => {
  const translate = useTranslate('Lk.Edit.TechStack')
  const { names, onChange, onBlur } = useTechnologyListField(name)
  const [isAdding, setIsAdding] = useState(false)

  if (isAdding) {
    return (
      <Autocomplete
        multiple
        freeSolo
        autoSelect
        filterSelectedOptions
        openOnFocus
        options={TECHNOLOGY_SUGGESTIONS}
        value={names}
        onChange={(_, next) => onChange(next)}
        onBlur={() => {
          onBlur()
          setIsAdding(false)
        }}
        renderInput={(params) => (
          <TextField
            {...params}
            variant="standard"
            autoFocus
            placeholder={translate('addTechnology')}
          />
        )}
      />
    )
  }

  return (
    <Stack direction="row" flexWrap="wrap" gap={1}>
      {names.map((item) => (
        <Chip key={item} label={item} onDelete={() => onChange(names.filter((n) => n !== item))} />
      ))}

      <Chip
        icon={<AddRounded fontSize="small" />}
        label={translate('addTechnology')}
        variant="outlined"
        onClick={() => setIsAdding(true)}
        sx={{ cursor: 'pointer' }}
      />
    </Stack>
  )
}
