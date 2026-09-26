import { ActivityType, Person } from '@features/curriculumVitae/types'
import { blankActivity } from '@features/curriculumVitae/utils'
import { useFormList } from '@features/curriculumVitae/hooks'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { ArrayItemCard } from '@features/curriculumVitae/components/form/ArrayItemCard'
import { DateField } from '@features/curriculumVitae/components/form/DateField'
import { SelectField } from '@features/curriculumVitae/components/form/SelectField'
import { useActivityOptions } from './hooks'

/** Секция «Дополнительная активность»: open source, выступления, волонтёрство и т.д. */
export const ActivitiesEditor = () => {
  const translate = useTranslate('Lk.Edit.Activities')
  const { typeOptions } = useActivityOptions()
  const { control, register, fields, addItem, itemProps } = useFormList('activities', blankActivity)

  return (
    <Stack spacing={1.5}>
      {fields.map((field, index) => (
        <ArrayItemCard key={field.id} {...itemProps(index)}>
          <SelectField
            name={`activities.${index}.type`}
            clearable={false}
            label={translate('type')}
            options={typeOptions}
          />

          <TextField
            label={translate('title')}
            {...register(`activities.${index}.title`)}
            required
            fullWidth
          />
          <TextField
            fullWidth
            label={translate('description')}
            {...register(`activities.${index}.description`)}
            multiline
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
        onClick={addItem}
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addActivity')}
      </Button>
    </Stack>
  )
}
