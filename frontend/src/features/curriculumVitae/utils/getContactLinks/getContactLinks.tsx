import { ReactNode } from 'react'
import { ContactKey, Contacts } from '@features/curriculumVitae/types'
import { EmailOutlined, PhoneOutlined, Telegram } from '@mui/icons-material'

export type ContactLink = {
  key: ContactKey
  href: string

  /** Что показываем и копируем после раскрытия: @ник, адрес почты или номер */
  value: string
  icon: ReactNode
}

/** Ссылки на те способы связи, которые заполнены. */
export const getContactLinks = (contacts: Contacts): ContactLink[] =>
  [
    contacts.telegram
      ? {
          key: ContactKey.Telegram,
          href: `https://t.me/${contacts.telegram}`,
          value: `@${contacts.telegram}`,
          icon: <Telegram />,
        }
      : null,
    contacts.email
      ? {
          key: ContactKey.Email,
          href: `mailto:${contacts.email}`,
          value: contacts.email,
          icon: <EmailOutlined />,
        }
      : null,
    contacts.phone
      ? {
          key: ContactKey.Phone,
          href: `tel:${contacts.phone}`,
          value: contacts.phone,
          icon: <PhoneOutlined />,
        }
      : null,
  ].filter((link) => link !== null)
