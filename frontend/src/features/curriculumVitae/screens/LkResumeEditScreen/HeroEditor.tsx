import { glassCard } from '@features/curriculumVitae/styles'
import { AvailabilityStatus, Person } from '@features/curriculumVitae/types'
import { Stack, TextField, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFormContext } from 'react-hook-form'
import { DateField } from './DateField'
import { SelectField } from './SelectField'

/** Приветствие и доступность — не секция, а фиксированная шапка резюме (всегда первая). */
export const HeroEditor = () => {
  const translate = useTranslate('Lk.Edit.Hero')
  const translateAvailability = useTranslate('CurriculumVitae.Availability')
  const { register } = useFormContext<Person>()

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Typography variant="h6">{translate('title')}</Typography>

      <TextField
        label={translate('greetingFirst')}
        {...register('greeting.first')}
        required
        slotProps={{ inputLabel: { shrink: true }, htmlInput: { style: { fontSize: '1.5rem' } } }}
      />

      <Stack direction="row" spacing={2}>
        <TextField label={translate('greetingSecond')} {...register('greeting.second')} required />

        <TextField
          label={translate('greetingHighlight')}
          {...register('greeting.secondHighlight')}
          required
          fullWidth
          slotProps={{
            inputLabel: { shrink: true },
            htmlInput: { style: { fontSize: '1.25rem', fontWeight: 700 } },
          }}
          sx={{ '& .MuiInput-input': { color: 'primary.main' } }}
        />
      </Stack>

      <TextField
        label={translate('greetingThird')}
        {...register('greeting.third')}
        required
        sx={{ color: 'text.secondary' }}
      />

      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
        <SelectField
          name="availability.status"
          label={translate('availabilityStatus')}
          sx={{ minWidth: 200 }}
          options={Object.values(AvailabilityStatus).map((status) => ({
            value: status,
            label: translateAvailability(status),
          }))}
        />

        <DateField name="availability.from" label={translate('availabilityFrom')} format="YYYY-MM-DD" />
      </Stack>

      <TextField
        label={translate('availabilityNote')}
        {...register('availability.note')}
        multiline
        minRows={2}
      />

      <TextField label={translate('resumeLink')} {...register('resumeLink')} />
    </Stack>
  )
}
