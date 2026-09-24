import { KeyboardEvent, useState } from 'react'
import { Person } from '@features/curriculumVitae/types'
import { Chip, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFieldArray, useFormContext } from 'react-hook-form'

type Props = {
  name: 'technologies' | `experience.${number}.stack` | `projects.${number}.stack`
}

/** Список технологий (просто названия, как чипы на самой странице резюме) — переиспользуется
 *  в стэке резюме, опыте работы и проектах: везде одна и та же форма `Technology[]`. */
export const TechnologyListField = ({ name }: Props) => {
  const translate = useTranslate('Lk.Edit.TechStack')
  const { control } = useFormContext<Person>()
  const { fields, append, remove } = useFieldArray({ control, name })

  const [draft, setDraft] = useState('')

  const addFromDraft = () => {
    const value = draft.trim()
    if (!value) return

    append({ name: value })
    setDraft('')
  }

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key !== 'Enter') return
    e.preventDefault()
    addFromDraft()
  }

  return (
    <Stack direction="row" flexWrap="wrap" alignItems="center" gap={1}>
      {fields.map((field, index) => (
        <Chip key={field.id} label={field.name} onDelete={() => remove(index)} />
      ))}

      <TextField
        variant="standard"
        placeholder={translate('addTechnology')}
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        onKeyDown={handleKeyDown}
        onBlur={addFromDraft}
        sx={{ width: 160 }}
      />
    </Stack>
  )
}
