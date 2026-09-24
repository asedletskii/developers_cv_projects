import { ActivityType, Person } from '@features/curriculumVitae/types'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFieldArray, useFormContext } from 'react-hook-form'
import { ArrayItemCard } from './ArrayItemCard'
import { DateField } from './DateField'
import { SelectField } from './SelectField'

/** Секция «Дополнительная активность»: open source, выступления, волонтёрство и т.д. */
export const ActivitiesEditor = () => {
  const translate = useTranslate('Lk.Edit.Activities')
  const translateActivities = useTranslate('CurriculumVitae.Activities')
  const { control, register } = useFormContext<Person>()
  const { fields, append, remove, move } = useFieldArray({ control, name: 'activities' })

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
          <SelectField
            name={`activities.${index}.type`}
            label={translate('type')}
            options={Object.values(ActivityType).map((type) => ({
              value: type,
              label: translateActivities(`Types.${type}`),
            }))}
          />

          <TextField
            label={translate('title')}
            {...register(`activities.${index}.title`)}
            required
            fullWidth
            slotProps={{ inputLabel: { shrink: true }, htmlInput: { style: { fontSize: '1.15rem', fontWeight: 700 } } }}
          />
          <TextField
            label={translate('description')}
            {...register(`activities.${index}.description`)}
            multiline
            minRows={2}
          />

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <DateField
              name={`activities.${index}.date`}
              label={translate('date')}
              format="YYYY-MM"
              views={['year', 'month']}
            />
            <TextField
              label={translate('link')}
              {...register(`activities.${index}.link`)}
              fullWidth
            />
          </Stack>
        </ArrayItemCard>
      ))}

      <Button
        size="small"
        startIcon={<AddRounded />}
        onClick={() =>
          append({ id: crypto.randomUUID(), type: ActivityType.OpenSource, title: '' })
        }
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addActivity')}
      </Button>
    </Stack>
  )
}
