import { ArrayItemCard } from '@features/curriculumVitae/components/form/ArrayItemCard'
import { DateField } from '@features/curriculumVitae/components/form/DateField'
import { SelectField } from '@features/curriculumVitae/components/form/SelectField'
import { useFormList } from '@features/curriculumVitae/hooks'
import { blankEducation } from '@features/curriculumVitae/utils'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useEducationOptions } from './hooks'

/** Секция «Образование»: список учебных заведений. */
export const EducationEditor = () => {
  const translate = useTranslate('Lk.Edit.Education')
  const { levelOptions } = useEducationOptions()
  const { register, fields, addItem, itemProps } = useFormList('educations', blankEducation)

  return (
    <Stack spacing={1.5}>
      {fields.map((field, index) => (
        <ArrayItemCard key={field.id} {...itemProps(index)}>
          <TextField
            label={translate('university')}
            {...register(`educations.${index}.university`)}
            required
            fullWidth
          />

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <TextField
              label={translate('location')}
              {...register(`educations.${index}.location`)}
              required
              fullWidth
            />

            <SelectField
              name={`educations.${index}.grade`}
              clearable={false}
              label={translate('grade')}
              options={levelOptions}
              sx={{ minWidth: 220 }}
            />
          </Stack>

          <TextField
            fullWidth
            label={translate('programName')}
            {...register(`educations.${index}.programName`)}
            required
          />
          <TextField
            fullWidth
            label={translate('profileName')}
            {...register(`educations.${index}.profileName`)}
          />
          <TextField
            fullWidth
            label={translate('faculty')}
            {...register(`educations.${index}.faculty`)}
            required
          />

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <TextField
              fullWidth
              type="number"
              label={translate('course')}
              {...register(`educations.${index}.course`, { valueAsNumber: true })}
            />
            <DateField
              name={`educations.${index}.endDate`}
              label={translate('endDate')}
              format="YYYY"
              views={['year']}
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
        {translate('addEducation')}
      </Button>
    </Stack>
  )
}
