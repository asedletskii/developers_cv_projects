import { Person } from '@features/curriculumVitae/types'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Controller, useFieldArray, useFormContext } from 'react-hook-form'
import { ArrayItemCard } from './ArrayItemCard'
import { DateField } from './DateField'
import { TechnologyListField } from './TechnologyListField'

/** Секция «Опыт работы»: список мест работы, каждое — со своим стэком технологий. */
export const ExperienceEditor = () => {
  const translate = useTranslate('Lk.Edit.Experience')
  const { control, register } = useFormContext<Person>()
  const { fields, append, remove, move } = useFieldArray({ control, name: 'experience' })

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
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <TextField
              label={translate('position')}
              {...register(`experience.${index}.position`)}
              required
              fullWidth
              slotProps={{ inputLabel: { shrink: true }, htmlInput: { style: { fontSize: '1.15rem', fontWeight: 700 } } }}
            />
            <TextField
              label={translate('company')}
              {...register(`experience.${index}.company`)}
              required
              fullWidth
            />
          </Stack>

          <TextField label={translate('location')} {...register(`experience.${index}.location`)} />

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <DateField
              name={`experience.${index}.startDate`}
              label={translate('startDate')}
              format="YYYY-MM"
              views={['year', 'month']}
              required
            />
            <DateField
              name={`experience.${index}.endDate`}
              label={translate('endDate')}
              format="YYYY-MM"
              views={['year', 'month']}
            />
          </Stack>

          <TextField
            label={translate('description')}
            {...register(`experience.${index}.description`)}
            multiline
            minRows={2}
          />

          <Controller
            control={control}
            name={`experience.${index}.achievements`}
            render={({ field }) => (
              <TextField
                label={translate('achievements')}
                multiline
                minRows={2}
                value={(field.value ?? []).join('\n')}
                onChange={(e) => field.onChange(e.target.value.split('\n'))}
                onBlur={field.onBlur}
              />
            )}
          />

          <TechnologyListField name={`experience.${index}.stack`} />
        </ArrayItemCard>
      ))}

      <Button
        size="small"
        startIcon={<AddRounded />}
        onClick={() =>
          append({
            id: crypto.randomUUID(),
            company: '',
            position: '',
            startDate: '',
            achievements: [],
            stack: [],
          })
        }
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addExperience')}
      </Button>
    </Stack>
  )
}
