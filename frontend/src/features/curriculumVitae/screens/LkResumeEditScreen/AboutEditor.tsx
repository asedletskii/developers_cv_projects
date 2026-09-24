import { LanguageGrade, LanguageType, Person, WorkFormat } from '@features/curriculumVitae/types'
import { AddRounded, CloseRounded } from '@mui/icons-material'
import { Button, Checkbox, FormControlLabel, IconButton, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFieldArray, useFormContext } from 'react-hook-form'
import { DateField } from './DateField'
import { SelectField } from './SelectField'

/** Секция «Обо мне»: личные данные, формат работы, языки. */
export const AboutEditor = () => {
  const translate = useTranslate('Lk.Edit.About')
  const translateEdit = useTranslate('Lk.Edit')
  const translateAbout = useTranslate('CurriculumVitae.About')
  const { control, register } = useFormContext<Person>()
  const { fields, append, remove } = useFieldArray({ control, name: 'languages' })

  return (
    <Stack spacing={2}>
      <TextField label={translate('photoPath')} {...register('photoPath')} required />

      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
        <TextField
          label={translate('name')}
          {...register('name')}
          required
          fullWidth
          slotProps={{ inputLabel: { shrink: true }, htmlInput: { style: { fontSize: '1.5rem' } } }}
        />
        <TextField
          label={translate('surname')}
          {...register('surname')}
          required
          fullWidth
          slotProps={{ inputLabel: { shrink: true }, htmlInput: { style: { fontSize: '1.5rem' } } }}
        />
        <TextField label={translate('patronymic')} {...register('patronymic')} fullWidth />
      </Stack>

      <TextField
        label={translate('position')}
        {...register('position')}
        required
        sx={{ '& .MuiInput-input': { color: 'primary.main', fontSize: '1.1rem' } }}
      />

      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
        <DateField
          name="dateOfBirth"
          label={translate('dateOfBirth')}
          format="YYYY-MM-DD"
          required
        />

        <TextField label={translate('location')} {...register('location')} required fullWidth />
      </Stack>

      <FormControlLabel
        control={<Checkbox {...register('relocation')} />}
        label={translate('relocation')}
      />

      <SelectField
        name="workFormat"
        multiple
        label={translate('workFormat')}
        options={Object.values(WorkFormat).map((format) => ({
          value: format,
          label: translateAbout(`WorkFormat.${format}`),
        }))}
      />

      <Stack direction="row" flexWrap="wrap" alignItems="flex-end" gap={2}>
        {fields.map((field, index) => (
          <Stack key={field.id} direction="row" alignItems="flex-end" spacing={0.5}>
            <SelectField
              name={`languages.${index}.language`}
              label={translate('language')}
              options={Object.values(LanguageType).map((language) => ({
                value: language,
                label: translateAbout(`Languages.${language}`),
              }))}
              sx={{ width: 120 }}
            />

            <SelectField
              name={`languages.${index}.grade`}
              label={translate('languageGrade')}
              options={Object.values(LanguageGrade).map((grade) => ({
                value: grade,
                label: grade,
              }))}
              sx={{ width: 72 }}
            />

            <IconButton
              size="small"
              aria-label={translateEdit('remove')}
              onClick={() => remove(index)}
            >
              <CloseRounded fontSize="small" />
            </IconButton>
          </Stack>
        ))}

        <Button
          size="small"
          startIcon={<AddRounded />}
          onClick={() => append({ language: LanguageType.Ru })}
          sx={{ textTransform: 'none', mb: 0.5 }}
        >
          {translate('addLanguage')}
        </Button>
      </Stack>
    </Stack>
  )
}
