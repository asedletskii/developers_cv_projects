import { DateField } from '@features/curriculumVitae/components/form/DateField'
import { SelectField } from '@features/curriculumVitae/components/form/SelectField'
import { Resume } from '@features/curriculumVitae/types'
import { Checkbox, FormControlLabel, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFormContext } from 'react-hook-form'
import { useProfileOptions } from './hooks'

/** Вкладка «Профиль»: личные данные и формат работы. */
export const AboutEditor = () => {
  const translate = useTranslate('Lk.Edit.About')
  const { workFormatOptions } = useProfileOptions()
  const { register } = useFormContext<Resume>()

  return (
    <Stack spacing={2}>
      <TextField
        fullWidth
        label={translate('photoPath')}
        {...register('profile.photoPath')}
        required
      />

      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
        <TextField label={translate('name')} {...register('profile.name')} required fullWidth />
        <TextField
          label={translate('surname')}
          {...register('profile.surname')}
          required
          fullWidth
        />
        <TextField label={translate('patronymic')} {...register('profile.patronymic')} fullWidth />
      </Stack>

      <TextField
        fullWidth
        label={translate('position')}
        {...register('profile.position')}
        required
      />

      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
        <DateField
          name="profile.dateOfBirth"
          label={translate('dateOfBirth')}
          format="YYYY-MM-DD"
          required
        />

        <TextField
          label={translate('location')}
          {...register('profile.location')}
          required
          fullWidth
        />
      </Stack>

      <FormControlLabel
        control={<Checkbox {...register('profile.relocation')} />}
        label={translate('relocation')}
      />

      <SelectField
        name="profile.workFormat"
        multiple
        label={translate('workFormat')}
        options={workFormatOptions}
      />
    </Stack>
  )
}
