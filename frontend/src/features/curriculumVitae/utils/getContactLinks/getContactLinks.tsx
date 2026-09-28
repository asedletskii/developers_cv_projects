import { ReactNode } from 'react'
import { ContactKey, Contacts } from '@features/curriculumVitae/types'
import { EmailOutlined, PhoneOutlined, Telegram } from '@mui/icons-material'

export type ContactLink = {
  key: ContactKey
  icon: ReactNode

  /** Скрытый контакт: бэк не отдал `value` в публичном ответе, его нужно запросить отдельно
   *  (`POST /cv/:id/contact-views`) — см. `useContactButton` */
  isHidden: boolean

  /** Есть сразу у открытых контактов; у скрытых появляется после раскрытия по клику */
  value?: string
}

const ICONS: Record<ContactKey, ReactNode> = {
  [ContactKey.Telegram]: <Telegram />,
  [ContactKey.Email]: <EmailOutlined />,
  [ContactKey.Phone]: <PhoneOutlined />,
}

/** Как показываем значение контакта: у телеграма — с `@`, остальные как есть. */
export const getContactDisplayValue = (key: ContactKey, value: string) =>
  key === ContactKey.Telegram ? `@${value}` : value

/** Href по значению контакта — годится и для открытого значения, и для раскрытого по клику. */
export const getContactHref = (key: ContactKey, value: string) => {
  switch (key) {
    case ContactKey.Telegram:
      return `https://t.me/${value}`
    case ContactKey.Email:
      return `mailto:${value}`
    case ContactKey.Phone:
      return `tel:${value}`
  }
}

/** Ссылки на те способы связи, которые заполнены или скрыты. У скрытых `value` не приходит —
 *  бэк прячет его в публичном ответе, пока не запросили через `POST /cv/:id/contact-views`. */
export const getContactLinks = (contacts: Contacts, hidden: ContactKey[] = []): ContactLink[] =>
  ([ContactKey.Telegram, ContactKey.Email, ContactKey.Phone] as const)
    .map((key) => {
      const value = contacts[key]
      const isHidden = hidden.includes(key)

      if (!value && !isHidden) return null

      return { key, icon: ICONS[key], isHidden, value: value || undefined }
    })
    .filter((link) => link !== null)
