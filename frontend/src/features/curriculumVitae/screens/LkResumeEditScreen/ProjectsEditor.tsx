import { Person } from '@features/curriculumVitae/types'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFieldArray, useFormContext } from 'react-hook-form'
import { ArrayItemCard } from './ArrayItemCard'
import { TechnologyListField } from './TechnologyListField'

/** Секция «Проекты». */
export const ProjectsEditor = () => {
  const translate = useTranslate('Lk.Edit.Projects')
  const { control, register } = useFormContext<Person>()
  const { fields, append, remove, move } = useFieldArray({ control, name: 'projects' })

  return (
    <Stack spacing={1.5}>
      {fields.map((field, index) => (
        <ArrayItemCard
          key={field.id}
          canMoveUp={index > 0}
          canMoveDown={index < fields.length - 1}
          onMoveUp={() => move(index, index - 1)}
          onMoveDown={() => move(index, index + 1)}
          onRemove={() => remove(index)}
        >
          <TextField
            label={translate('name')}
            {...register(`projects.${index}.name`)}
            required
            fullWidth
            slotProps={{ inputLabel: { shrink: true }, htmlInput: { style: { fontSize: '1.15rem', fontWeight: 700 } } }}
          />
          <TextField
            label={translate('description')}
            {...register(`projects.${index}.description`)}
            multiline
            minRows={2}
          />
          <TextField label={translate('link')} {...register(`projects.${index}.link`)} />

          <TechnologyListField name={`projects.${index}.stack`} />
        </ArrayItemCard>
      ))}

      <Button
        size="small"
        startIcon={<AddRounded />}
        onClick={() => append({ id: crypto.randomUUID(), name: '', stack: [] })}
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addProject')}
      </Button>
    </Stack>
  )
}
