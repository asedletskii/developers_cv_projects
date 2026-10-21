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

export const getContactDisplayValue = (key: ContactKey, value: string) =>
  key === ContactKey.Telegram ? `@${value}` : value

export const getContactHref = (key: ContactKey, value: string) => {
  switch (key) {
    case ContactKey.Telegram:
      return `https://t.me/${encodeURIComponent(value.replace(/^@/, ''))}`
    case ContactKey.Email:
      return `mailto:${encodeURIComponent(value.trim()).replace(/%40/g, '@')}`
    case ContactKey.Phone:
      return `tel:${value.replace(/[^\d+*#]/g, '')}`
  }
}

export const getContactLinks = (contacts: Contacts, hidden: ContactKey[] = []): ContactLink[] =>
  ([ContactKey.Telegram, ContactKey.Email, ContactKey.Phone] as const)
    .map((key) => {
      const value = contacts[key]
      const isHidden = hidden.includes(key)

      if (!value && !isHidden) return null

      return { key, icon: ICONS[key], isHidden, value: value || undefined }
    })
    .filter((link) => link !== null)
