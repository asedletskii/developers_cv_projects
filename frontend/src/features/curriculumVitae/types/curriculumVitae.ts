import {
  ActivityType,
  AvailabilityStatus,
  ContactKey,
  EducationLevel,
  HardSkill,
  LanguageGrade,
  LanguageType,
  SectionId,
  SoftSkill,
  WorkFormat,
} from './enums'

export type Person = {
  id: string

  /** Имя */
  name: string

  /** Фамилия */
  surname: string

  /** Отчество */
  patronymic?: string

  /** Вакансия */
  position: string

  /** Доступность для работодателей */
  availability?: Availability

  /** Когда резюме обновлялось в последний раз (ISO-дата) */
  updatedAt?: string

  /** Ссылка на файл резюме (если нет, кнопки скачивания скрыты) */
  resumeLink?: string

  /** Акцентные цвета резюме (если не заданы, используются цвета по умолчанию) */
  theme?: ThemeSettings

  /** Настройки секций: порядок массива = порядок на странице */
  sections: SectionSettings[]

  /** Контакты, которые показываются только по клику (защита от сборщиков почты и телефонов) */
  hiddenContacts?: ContactKey[]

  /** Тексты блока контактов (если не заданы, берутся тексты по умолчанию) */
  contactsText?: ContactsText

  /** Приветствие на главном экране */
  greeting: Greeting

  /** Фото */
  photoPath: string

  /** Дата рождения */
  dateOfBirth: string

  /** Город проживания */
  location: string

  /** Готовность к релокации */
  relocation: boolean

  /** Формат работы */
  workFormat: WorkFormat[]

  /** Контакты */
  contacts: Contacts

  /** Знание языков */
  languages: {
    /** Язык */
    language: LanguageType

    /** Уровень знания языка */
    grade?: LanguageGrade
  }[]

  /** Опыт работы */
  experience: Experience[]

  /** Образование */
  education: Education[]

  /** Курсы и сертификаты */
  courses: Course[]

  /** Дополнительная активность */
  activities: Activity[]

  /** Хард скилы */
  hardSkills: HardSkill[]

  /** Софт скилы */
  softSkills: SoftSkill[]

  /** Используемые технологии/основной стэк */
  technologies: Technology[]

  /** Проекты */
  projects: Pick<Project, 'id' | 'name' | 'description' | 'stack' | 'link'>[]
}

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

  /** Роль в проекте */
  role?: string

  /** Скриншоты проекта */
  screenshots?: string[]

  /** Ссылка на проект (где посмотреть)*/
  link?: string

  /** Ссылка на гит */
  gitLink?: string

  /** Стэк проекта */
  stack: Technology[]

  /** Участники проекта */
  contributors: SimplifiedPerson[]
}

type SimplifiedPerson = Pick<Person, 'name' | 'surname' | 'contacts'> & {
  role: string
}
