import { SelectField } from '@features/curriculumVitae/components/editor/form/SelectField'
import { Resume } from '@features/curriculumVitae/types'
import { Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFormContext } from 'react-hook-form'
import { useContactOptions } from './hooks'

export const ContactsEditor = () => {
  const translate = useTranslate('Lk.Edit.Contacts')
  const { hiddenContactOptions } = useContactOptions()
  const { register } = useFormContext<Resume>()

  return (
    <Stack spacing={2}>
      <TextField fullWidth label={translate('phone')} {...register('contacts.phone')} />
      <TextField fullWidth label={translate('telegram')} {...register('contacts.telegram')} />
      <TextField fullWidth label={translate('email')} {...register('contacts.email')} />

      <SelectField
        name="contacts.hiddenContacts"
        multiple
        label={translate('hiddenContacts')}
        options={hiddenContactOptions}
      />

      <TextField
        fullWidth
        label={translate('contactsTitle')}
        {...register('contacts.contactsText.title')}
      />
      <TextField
        fullWidth
        label={translate('contactsSubtitle')}
        {...register('contacts.contactsText.subtitle')}
      />
    </Stack>
  )
}
