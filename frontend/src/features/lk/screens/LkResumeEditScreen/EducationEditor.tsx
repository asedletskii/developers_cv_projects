import { EducationLevel, Person } from '@features/curriculumVitae/types'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFieldArray, useFormContext } from 'react-hook-form'
import { ArrayItemCard } from './ArrayItemCard'
import { DateField } from './DateField'
import { SelectField } from './SelectField'

const levels = [0, 1, 2, 3, 4, 5, 6] as const

/** Секция «Образование»: список учебных заведений. */
export const EducationEditor = () => {
  const translate = useTranslate('Lk.Edit.Education')
  const translateEducation = useTranslate('CurriculumVitae.Education')
  const { control, register } = useFormContext<Person>()
  const { fields, append, remove, move } = useFieldArray({ control, name: 'education' })

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
            label={translate('university')}
            {...register(`education.${index}.university`)}
            required
            fullWidth
            slotProps={{ inputLabel: { shrink: true }, htmlInput: { style: { fontSize: '1.15rem', fontWeight: 700 } } }}
          />

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <TextField
              label={translate('location')}
              {...register(`education.${index}.location`)}
              required
              fullWidth
            />

            <SelectField
              name={`education.${index}.grade`}
              label={translate('grade')}
              options={levels.map((level) => ({
                value: level,
                label: translateEducation(`Levels.level${level}`),
              }))}
              sx={{ minWidth: 220 }}
            />
          </Stack>

          <TextField
            label={translate('programName')}
            {...register(`education.${index}.programName`)}
            required
          />
          <TextField
            label={translate('profileName')}
            {...register(`education.${index}.profileName`)}
          />
          <TextField
            label={translate('faculty')}
            {...register(`education.${index}.faculty`)}
            required
          />

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <TextField
              type="number"
              label={translate('course')}
              {...register(`education.${index}.course`)}
            />
            <DateField
              name={`education.${index}.endDate`}
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
        onClick={() =>
          append({
            id: crypto.randomUUID(),
            location: '',
            university: '',
            grade: EducationLevel.Bachelor,
            programName: '',
            faculty: '',
          })
        }
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addEducation')}
      </Button>
    </Stack>
  )
}
