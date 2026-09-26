import { Person } from '@features/curriculumVitae/types'
import { blankCourse } from '@features/curriculumVitae/utils'
import { useFormList } from '@features/curriculumVitae/hooks'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { ArrayItemCard } from '@features/curriculumVitae/components/form/ArrayItemCard'
import { DateField } from '@features/curriculumVitae/components/form/DateField'

/** Секция «Курсы»: список курсов и сертификатов. */
export const CoursesEditor = () => {
  const translate = useTranslate('Lk.Edit.Courses')
  const { control, register, fields, addItem, itemProps } = useFormList('courses', blankCourse)

  return (
    <Stack spacing={1.5}>
      {fields.map((field, index) => (
        <ArrayItemCard key={field.id} {...itemProps(index)}>
          <TextField
            label={translate('name')}
            {...register(`courses.${index}.name`)}
            required
            fullWidth
          />
          <TextField
            fullWidth
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
              fullWidth
              type="number"
              label={translate('hours')}
              {...register(`courses.${index}.hours`, { valueAsNumber: true })}
            />
          </Stack>

          <TextField
            fullWidth
            label={translate('certificateLink')}
            {...register(`courses.${index}.certificateLink`)}
          />
        </ArrayItemCard>
      ))}

      <Button
        size="small"
        startIcon={<AddRounded />}
        onClick={addItem}
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addCourse')}
      </Button>
    </Stack>
  )
}
