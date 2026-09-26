import { Person } from '@features/curriculumVitae/types'
import { Checkbox, FormControlLabel, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFormContext } from 'react-hook-form'
import { useProfileOptions } from './hooks'
import { DateField } from '@features/curriculumVitae/components/form/DateField'
import { SelectField } from '@features/curriculumVitae/components/form/SelectField'

/** Вкладка «Профиль»: личные данные и формат работы. */
export const AboutEditor = () => {
  const translate = useTranslate('Lk.Edit.About')
  const { workFormatOptions } = useProfileOptions()
  const { register } = useFormContext<Person>()

  return (
    <Stack spacing={2}>
      <TextField fullWidth label={translate('photoPath')} {...register('photoPath')} required />

      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
        <TextField label={translate('name')} {...register('name')} required fullWidth />
        <TextField label={translate('surname')} {...register('surname')} required fullWidth />
        <TextField label={translate('patronymic')} {...register('patronymic')} fullWidth />
      </Stack>

      <TextField fullWidth label={translate('position')} {...register('position')} required />

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
        options={workFormatOptions}
      />
    </Stack>
  )
}
