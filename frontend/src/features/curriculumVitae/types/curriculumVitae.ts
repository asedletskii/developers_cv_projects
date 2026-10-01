import { ActivityType, AvailabilityStatus, EducationLevel, SectionId } from './enums'

/** Тексты блока контактов: не пришло поле — текст по умолчанию, пустая строка — элемент скрыт */
export type ContactsText = {
  /** Заголовок, например «Давайте работать вместе» */
  title?: string

  /** Подпись под заголовком */
  subtitle?: string
}

export type Availability = {
  status: AvailabilityStatus

  /** С какой даты можно выйти на работу, YYYY-MM-DD */
  from?: string

  /** Пояснение, например «Ищу роль Senior Frontend, удалённо» */
  note?: string
}

/** Акцентные цвета в формате #RGB или #RRGGBB; остальные оттенки считаются на клиенте */
export type ThemeSettings = {
  primary?: string
  secondary?: string
}

export type SectionSettings = {
  id: SectionId

  /** Показывать ли секцию */
  enabled: boolean

  /** Свёрнута ли секция при открытии страницы (по умолчанию развёрнута) */
  collapsed?: boolean
}

/** Приветствие: три строки, во второй выделяется часть текста */
export type Greeting = {
  /** Первая строка, например «Привет!» */
  first: string

  /** Вторая строка, например «Меня зовут» */
  second: string

  /** Выделенная часть второй строки, например имя */
  secondHighlight: string

  /** Третья строка */
  third: string
}

export type Technology = {
  /** Название */
  name: string

  /** Уровень владения */
  grade?: string
}

export type Contacts = {
  /** Номер телефона */
  phone?: string

  /** Ник в телеграм без '@' */
  telegram?: string

  /** Электронная почта */
  email?: string
}

export type Experience = {
  id: string

  /** Компания */
  company: string

  /** Должность */
  position: string

  /** Город или формат (например, «Удалённо») */
  location?: string

  /** Дата начала, YYYY-MM */
  startDate: string

  /** Дата окончания, YYYY-MM (пусто — по настоящее время) */
  endDate?: string

  /** Краткое описание */
  description?: string

  /** Достижения */
  achievements: string[]

  /** Задачи — чем занимался на позиции */
  tasks: string[]

  /** Используемый стэк */
  stack: Technology[]
}

export type Course = {
  id: string

  /** Название курса */
  name: string

  /** Организация или платформа */
  organization: string

  /** Дата окончания, YYYY-MM */
  endDate: string

  /** Длительность в часах */
  hours?: number

  /** Ссылка на сертификат */
  certificateLink?: string
}

export type Activity = {
  id: string

  /** Вид активности */
  type: ActivityType

  /** Название */
  title: string

  /** Описание */
  description?: string

  /** Дата, YYYY-MM */
  date?: string

  /** Ссылка */
  link?: string
}

export type Education = {
  id: string

  /** Город */
  location: string

  /** Учебное заведение */
  university: string

  /** Уровень образования */
  grade: EducationLevel

  /** Наименование программы */
  programName: string

  /** Наименование профиля */
  profileName?: string

  /** Факультет */
  faculty: string

  /** Курс (при отсутствии года окончания) */
  course?: number

  /** Год окончания */
  endDate?: string
}

export type Project = {
  id: string

  /** Название проекта */
  name: string

  /** Описание проекта */
  description?: string

  /** Роль в проекте (в публичном ответе нет) */
  role?: string

  /** Скриншоты проекта (в публичном ответе всегда пусто) */
  screenshots: string[]

  /** Ссылка на проект (где посмотреть)*/
  link?: string

  /** Ссылка на гит (в публичном ответе нет) */
  gitLink?: string

  /** Стэк проекта */
  stack: Technology[]

  /** Участники проекта (в публичном ответе всегда пусто) */
  contributors: SimplifiedPerson[]
}

type SimplifiedPerson = {
  name: string
  surname: string
  contacts: Contacts
  role: string
}
