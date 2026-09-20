import { DateField } from '@features/curriculumVitae/components/editor/form/DateField'
import { SelectField } from '@features/curriculumVitae/components/editor/form/SelectField'
import { Resume } from '@features/curriculumVitae/types'
import { PersonRounded } from '@mui/icons-material'
import { Button, Checkbox, FormControlLabel, Menu, MenuItem, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Languages } from '@shared/types'
import { useFormContext } from 'react-hook-form'
import { useFillFromProfile, useProfileOptions } from './hooks'

export const AboutEditor = () => {
  const translate = useTranslate('Lk.Edit.About')
  const { workFormatOptions } = useProfileOptions()
  const { register } = useFormContext<Resume>()
  const fillFromProfile = useFillFromProfile()

  return (
    <Stack spacing={2}>
      {fillFromProfile.hasProfile && (
        <Button
          onClick={fillFromProfile.openMenu}
          startIcon={<PersonRounded />}
          size="small"
          sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
        >
          {translate('fillFromProfile')}
        </Button>
      )}

      <Menu
        anchorEl={fillFromProfile.anchor}
        open={!!fillFromProfile.anchor}
        onClose={fillFromProfile.closeMenu}
      >
        <MenuItem onClick={() => fillFromProfile.fill(Languages.Ru)}>RU</MenuItem>
        <MenuItem onClick={() => fillFromProfile.fill(Languages.En)}>EN</MenuItem>
      </Menu>

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
