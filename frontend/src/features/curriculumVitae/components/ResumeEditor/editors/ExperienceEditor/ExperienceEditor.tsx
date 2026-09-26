import { ArrayItemCard } from '@features/curriculumVitae/components/form/ArrayItemCard'
import { DateField } from '@features/curriculumVitae/components/form/DateField'
import { TechnologyListField } from '@features/curriculumVitae/components/form/TechnologyListField'
import { useFormList } from '@features/curriculumVitae/hooks'
import { blankExperience } from '@features/curriculumVitae/utils'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Controller } from 'react-hook-form'

/** Секция «Опыт работы»: список мест работы, каждое — со своим стэком технологий. */
export const ExperienceEditor = () => {
  const translate = useTranslate('Lk.Edit.Experience')
  const { control, register, fields, addItem, itemProps } = useFormList(
    'experience',
    blankExperience,
  )

  return (
    <Stack spacing={1.5}>
      {fields.map((field, index) => (
        <ArrayItemCard key={field.id} {...itemProps(index)}>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <TextField
              label={translate('position')}
              {...register(`experience.${index}.position`)}
              required
              fullWidth
            />
            <TextField
              label={translate('company')}
              {...register(`experience.${index}.company`)}
              required
              fullWidth
            />
          </Stack>

          <TextField
            fullWidth
            label={translate('location')}
            {...register(`experience.${index}.location`)}
          />

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
            fullWidth
            label={translate('description')}
            {...register(`experience.${index}.description`)}
            multiline
          />

          <Controller
            control={control}
            name={`experience.${index}.achievements`}
            render={({ field }) => (
              <TextField
                fullWidth
                label={translate('achievements')}
                multiline
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
        onClick={addItem}
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addExperience')}
      </Button>
    </Stack>
  )
}
