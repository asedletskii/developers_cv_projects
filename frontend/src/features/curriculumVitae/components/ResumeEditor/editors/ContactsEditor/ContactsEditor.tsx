import { SelectField } from '@features/curriculumVitae/components/form/SelectField'
import { Person } from '@features/curriculumVitae/types'
import { Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFormContext } from 'react-hook-form'
import { useContactOptions } from './hooks'

/** Секция «Контакты»: сами контакты, какие скрыты по умолчанию, тексты блока. */
export const ContactsEditor = () => {
  const translate = useTranslate('Lk.Edit.Contacts')
  const { hiddenContactOptions } = useContactOptions()
  const { register } = useFormContext<Person>()

  return (
    <Stack spacing={2}>
      <TextField fullWidth label={translate('phone')} {...register('contacts.phone')} />
      <TextField fullWidth label={translate('telegram')} {...register('contacts.telegram')} />
      <TextField fullWidth label={translate('email')} {...register('contacts.email')} />

      <SelectField
        name="hiddenContacts"
        multiple
        label={translate('hiddenContacts')}
        options={hiddenContactOptions}
      />

      <TextField fullWidth label={translate('contactsTitle')} {...register('contactsText.title')} />
      <TextField
        fullWidth
        label={translate('contactsSubtitle')}
        {...register('contactsText.subtitle')}
      />
    </Stack>
  )
}
