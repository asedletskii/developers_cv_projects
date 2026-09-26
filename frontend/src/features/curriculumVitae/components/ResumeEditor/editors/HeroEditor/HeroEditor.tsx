import { DateField } from '@features/curriculumVitae/components/form/DateField'
import { SelectField } from '@features/curriculumVitae/components/form/SelectField'
import { HeroPreview } from '@features/curriculumVitae/components/HeroPreview'
import { glassCard } from '@features/curriculumVitae/styles'
import { Person } from '@features/curriculumVitae/types'
import { Stack, TextField, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFormContext } from 'react-hook-form'
import { useHeroOptions } from './hooks'

export const HeroEditor = () => {
  const translate = useTranslate('Lk.Edit.Hero')
  const { availabilityOptions } = useHeroOptions()
  const { register } = useFormContext<Person>()

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Typography variant="h6">{translate('title')}</Typography>

      <HeroPreview />

      <TextField
        fullWidth
        label={translate('greetingFirst')}
        {...register('greeting.first')}
        required
      />

      <Stack direction="row" spacing={2}>
        <TextField
          fullWidth
          label={translate('greetingSecond')}
          {...register('greeting.second')}
          required
        />

        <TextField
          label={translate('greetingHighlight')}
          {...register('greeting.secondHighlight')}
          required
          fullWidth
        />
      </Stack>

      <TextField
        fullWidth
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
          options={availabilityOptions}
        />

        <DateField
          name="availability.from"
          label={translate('availabilityFrom')}
          format="YYYY-MM-DD"
        />
      </Stack>

      <TextField
        fullWidth
        label={translate('availabilityNote')}
        {...register('availability.note')}
        multiline
      />

      <TextField fullWidth label={translate('resumeLink')} {...register('resumeLink')} />
    </Stack>
  )
}
