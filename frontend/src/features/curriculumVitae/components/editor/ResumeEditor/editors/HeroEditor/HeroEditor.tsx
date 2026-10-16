import { DateField } from '@features/curriculumVitae/components/editor/form/DateField'
import { SelectField } from '@features/curriculumVitae/components/editor/form/SelectField'
import { HeroPreview } from '@features/curriculumVitae/components/editor/HeroPreview'
import { glassCard } from '@features/curriculumVitae/styles'
import { Resume } from '@features/curriculumVitae/types'
import { Stack, TextField, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFormContext } from 'react-hook-form'
import { useHeroOptions } from './hooks'

export const HeroEditor = () => {
  const translate = useTranslate('Lk.Edit.Hero')
  const { availabilityOptions } = useHeroOptions()
  const { register } = useFormContext<Resume>()

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Typography variant="h6">{translate('title')}</Typography>

      <HeroPreview />

      <TextField
        fullWidth
        label={translate('greetingFirst')}
        {...register('overview.greeting.first')}
        required
      />

      <Stack direction="row" spacing={2}>
        <TextField
          fullWidth
          label={translate('greetingSecond')}
          {...register('overview.greeting.second')}
          required
        />

        <TextField
          label={translate('greetingHighlight')}
          {...register('overview.greeting.secondHighlight')}
          required
          fullWidth
        />
      </Stack>

      <TextField
        fullWidth
        label={translate('greetingThird')}
        {...register('overview.greeting.third')}
        required
        sx={{ color: 'text.secondary' }}
      />

      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
        <SelectField
          name="overview.availability.status"
          label={translate('availabilityStatus')}
          sx={{ minWidth: 200 }}
          options={availabilityOptions}
        />

        <DateField
          name="overview.availability.from"
          label={translate('availabilityFrom')}
          format="YYYY-MM-DD"
        />
      </Stack>

      <TextField
        fullWidth
        label={translate('availabilityNote')}
        {...register('overview.availability.note')}
        multiline
      />
    </Stack>
  )
}
