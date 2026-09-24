import { Person } from '@features/curriculumVitae/types'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFieldArray, useFormContext } from 'react-hook-form'
import { ArrayItemCard } from './ArrayItemCard'
import { DateField } from './DateField'

/** Секция «Курсы»: список курсов и сертификатов. */
export const CoursesEditor = () => {
  const translate = useTranslate('Lk.Edit.Courses')
  const { control, register } = useFormContext<Person>()
  const { fields, append, remove, move } = useFieldArray({ control, name: 'courses' })

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
            {...register(`courses.${index}.name`)}
            required
            fullWidth
            slotProps={{ inputLabel: { shrink: true }, htmlInput: { style: { fontSize: '1.15rem', fontWeight: 700 } } }}
          />
          <TextField
            label={translate('organization')}
            {...register(`courses.${index}.organization`)}
            required
          />

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <DateField
              name={`courses.${index}.endDate`}
              label={translate('endDate')}
              format="YYYY-MM"
              views={['year', 'month']}
              required
            />
            <TextField
              type="number"
              label={translate('hours')}
              {...register(`courses.${index}.hours`)}
            />
          </Stack>

          <TextField
            label={translate('certificateLink')}
            {...register(`courses.${index}.certificateLink`)}
          />
        </ArrayItemCard>
      ))}

      <Button
        size="small"
        startIcon={<AddRounded />}
        onClick={() =>
          append({ id: crypto.randomUUID(), name: '', organization: '', endDate: '' })
        }
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addCourse')}
      </Button>
    </Stack>
  )
}
