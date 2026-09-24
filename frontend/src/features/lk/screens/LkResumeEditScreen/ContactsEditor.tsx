import { ContactKey, Person } from '@features/curriculumVitae/types'
import { Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFormContext } from 'react-hook-form'
import { SelectField } from './SelectField'

/** Секция «Контакты»: сами контакты, какие скрыты по умолчанию, тексты блока. */
export const ContactsEditor = () => {
  const translate = useTranslate('Lk.Edit.Contacts')
  const { register } = useFormContext<Person>()

  return (
    <Stack spacing={2}>
      <TextField label={translate('phone')} {...register('contacts.phone')} />
      <TextField label={translate('telegram')} {...register('contacts.telegram')} />
      <TextField label={translate('email')} {...register('contacts.email')} />

      <SelectField
        name="hiddenContacts"
        multiple
        label={translate('hiddenContacts')}
        options={Object.values(ContactKey).map((key) => ({ value: key, label: translate(key) }))}
      />

      <TextField label={translate('contactsTitle')} {...register('contactsText.title')} />
      <TextField
        label={translate('contactsSubtitle')}
        {...register('contactsText.subtitle')}
      />
    </Stack>
  )
}
