import {
  ActivityType,
  AvailabilityStatus,
  ContactKey,
  ContactsText,
  EducationLevel,
  Greeting,
  HardSkill,
  LanguageGrade,
  LanguageType,
  Person,
  Project,
  SectionId,
  SoftSkill,
  Technology,
  TechnologyGrade,
  WorkFormat,
} from '@features/curriculumVitae/types'
import { Languages } from '@shared/types'
import dayjs from 'dayjs'

const delay = (ms = 300) => new Promise((resolve) => setTimeout(resolve, ms))

const stack = (...names: string[]): Technology[] => names.map((name) => ({ name }))

const personStore: Person = {
  id: '1',
  name: 'Антон',
  surname: 'Седлецкий',
  patronymic: 'Александрович',
  photoPath: '/mock/avatar.svg',
  resumeLink: '/mock/resume.pdf',
  hiddenContacts: [ContactKey.Email, ContactKey.Phone],
  availability: { status: AvailabilityStatus.Open, from: '2026-10-01' },
  dateOfBirth: '1998-04-12',
  position: 'Frontend',
  sections: [
    { id: SectionId.About, enabled: true },
    { id: SectionId.TechStack, enabled: true },
    { id: SectionId.Experience, enabled: true },
    { id: SectionId.Skills, enabled: true },
    { id: SectionId.Education, enabled: true },
    { id: SectionId.Courses, enabled: true, collapsed: true },
    { id: SectionId.Projects, enabled: true },
    { id: SectionId.Activities, enabled: true, collapsed: true },
    { id: SectionId.Contacts, enabled: true },
  ],
  greeting: {
    first: 'Привет!',
    second: 'Меня зовут',
    secondHighlight: 'Седлецкий Антон',
    third: 'Я Frontend-разработчик из Сибири',
  },
  location: 'Новосибирск',
  relocation: true,

  workFormat: [WorkFormat.FullTime, WorkFormat.Remote, WorkFormat.Hybrid],

  contacts: {
    phone: '+79999999999',
    telegram: 'a_sedletskii',
    email: 'anton@example.com',
  },

  languages: [
    { language: LanguageType.Ru, grade: LanguageGrade.C2 },
    { language: LanguageType.En, grade: LanguageGrade.B2 },
    { language: LanguageType.De, grade: LanguageGrade.A2 },
  ],

  experience: [
    {
      id: 'exp-1',
      company: 'Северный ветер',
      position: 'Senior Frontend разработчик',
      location: 'Удалённо',
      startDate: '2023-03',
      description:
        'Разрабатываю клиентские приложения для финтех-продуктов и веду фронтенд-команду.',
      achievements: [
        'Перевёл монолитный SPA на модульную архитектуру, время сборки сократилось в 3 раза',
        'Внедрил дизайн-систему на MUI, которую используют 4 продуктовые команды',
        'Довёл Lighthouse-показатели ключевых страниц с 62 до 95',
      ],
      stack: stack('React', 'TypeScript', 'Redux Toolkit', 'MUI', 'Vite', 'GSAP'),
    },
    {
      id: 'exp-2',
      company: 'Tech Labs',
      position: 'Frontend разработчик',
      location: 'Новосибирск',
      startDate: '2021-02',
      endDate: '2023-02',
      description: 'Делал личный кабинет и административную панель для B2B-платформы.',
      achievements: [
        'Написал библиотеку компонентов и переиспользуемых хуков на 40+ элементов',
        'Настроил CI: линтинг, тесты и предпросмотр каждой ветки',
        'Покрыл критичные сценарии интеграционными тестами',
      ],
      stack: stack('React', 'TypeScript', 'Jest', 'Storybook', 'GitHub Actions'),
    },
    {
      id: 'exp-3',
      company: 'Freelance',
      position: 'Web-разработчик',
      location: 'Удалённо',
      startDate: '2019-06',
      endDate: '2021-01',
      description: 'Лендинги и небольшие интернет-магазины для малого бизнеса.',
      achievements: ['Выполнил более 20 проектов от вёрстки до деплоя'],
      stack: stack('JavaScript', 'HTML', 'SCSS', 'Node.js'),
    },
    {
      id: 'exp-4',
      company: 'Студия «Мозаика»',
      position: 'Junior Frontend разработчик',
      location: 'Новосибирск',
      startDate: '2018-07',
      endDate: '2019-05',
      description: 'Вёрстка и доработка корпоративных сайтов в небольшой веб-студии.',
      achievements: ['Перевёл верстку студии на адаптивную сетку и препроцессор'],
      stack: stack('JavaScript', 'jQuery', 'SCSS'),
    },
    {
      id: 'exp-5',
      company: 'ИП Кузнецов',
      position: 'Стажёр-разработчик',
      location: 'Новосибирск',
      startDate: '2017-09',
      endDate: '2018-06',
      achievements: ['Собрал первый коммерческий сайт и настроил хостинг'],
      stack: stack('HTML', 'CSS', 'PHP'),
    },
  ],

  education: [
    {
      id: 'edu-1',
      location: 'Новосибирск',
      university: 'НГТУ',
      grade: EducationLevel.Master,
      programName: 'Программная инженерия',
      profileName: 'Разработка веб-приложений',
      faculty: 'Автоматики и вычислительной техники',
      course: 2,
    },
    {
      id: 'edu-2',
      location: 'Москва',
      university: 'МГТУ им. Баумана',
      grade: EducationLevel.Bachelor,
      programName: 'Информатика и вычислительная техника',
      profileName: 'Разработка ПО',
      faculty: 'Информатика и системы управления',
      endDate: '2021',
    },
  ],

  courses: [
    {
      id: 'course-1',
      name: 'Advanced React и производительность',
      organization: 'Frontend Masters',
      endDate: '2024-05',
      hours: 40,
      certificateLink: 'https://example.com/certificate/react',
    },
    {
      id: 'course-2',
      name: 'Архитектура фронтенд-приложений',
      organization: 'Яндекс Практикум',
      endDate: '2023-11',
      hours: 72,
      certificateLink: 'https://example.com/certificate/architecture',
    },
    {
      id: 'course-3',
      name: 'Node.js: разработка серверных приложений',
      organization: 'Stepik',
      endDate: '2022-03',
      hours: 36,
    },
  ],

  activities: [
    {
      id: 'act-1',
      type: ActivityType.OpenSource,
      title: 'Контрибьютор в open-source',
      description: 'Исправления и улучшения в популярных UI-библиотеках, 15+ принятых pull request',
      date: '2024-01',
      link: 'https://github.com/example',
    },
    {
      id: 'act-2',
      type: ActivityType.Hackathon,
      title: 'Хакатон «Цифровой прорыв»',
      description: 'Финалист в команде из 4 человек: прототип сервиса за 48 часов',
      date: '2023-09',
    },
    {
      id: 'act-3',
      type: ActivityType.Speaker,
      title: 'Доклад «Анимации на GSAP без боли»',
      description: 'Выступление на митапе фронтенд-разработчиков Новосибирска',
      date: '2024-03',
      link: 'https://example.com/talk',
    },
    {
      id: 'act-4',
      type: ActivityType.Mentoring,
      title: 'Менторство джунов',
      description: 'Помогаю начинающим разработчикам с ревью кода и карьерным планом',
    },
    {
      id: 'act-5',
      type: ActivityType.Community,
      title: 'Организатор митапов',
      description: 'Помогаю проводить встречи локального IT-сообщества',
      date: '2022-10',
    },
  ],

  hardSkills: [
    HardSkill.Architecture,
    HardSkill.Testing,
    HardSkill.Cicd,
    HardSkill.Performance,
    HardSkill.Accessibility,
    HardSkill.CodeReview,
    HardSkill.UiDesign,
    HardSkill.Api,
  ],

  softSkills: [
    SoftSkill.Communication,
    SoftSkill.Teamwork,
    SoftSkill.ProblemSolving,
    SoftSkill.TimeManagement,
    SoftSkill.Mentoring,
    SoftSkill.Adaptability,
    SoftSkill.Ownership,
  ],

  technologies: [
    { name: 'TypeScript', grade: TechnologyGrade.Senior },
    { name: 'React', grade: TechnologyGrade.Senior },
    { name: 'Redux Toolkit', grade: TechnologyGrade.Middle },
    { name: 'MUI', grade: TechnologyGrade.Senior },
    { name: 'GSAP', grade: TechnologyGrade.Middle },
    { name: 'Vite', grade: TechnologyGrade.Middle },
    { name: 'Node.js', grade: TechnologyGrade.Middle },
    { name: 'PostgreSQL', grade: TechnologyGrade.Junior },
    { name: 'Docker', grade: TechnologyGrade.Junior },
  ],

  projects: [
    {
      id: 'p1',
      name: 'CV Platform',
      description: 'Платформа для отображения резюме разработчиков',
      stack: stack('React', 'TypeScript', 'Vite', 'GSAP'),
      link: 'https://example.com',
    },
    {
      id: 'p2',
      name: 'Finance Tracker',
      description: 'Приложение для учёта личных финансов с аналитикой по категориям',
      stack: stack('React', 'Redux', 'Node.js', 'PostgreSQL'),
      link: 'https://example.com/finance',
    },
    {
      id: 'p3',
      name: 'Pulse Dashboard',
      description: 'Дашборд метрик в реальном времени: графики, алерты и фильтры',
      stack: stack('React', 'TypeScript', 'WebSocket', 'D3'),
      link: 'https://example.com/pulse',
    },
    {
      id: 'p4',
      name: 'Design Kit',
      description: 'Библиотека UI-компонентов с темами и документацией в Storybook',
      stack: stack('React', 'MUI', 'Storybook'),
    },
  ],
}

const secondPersonStore: Person = {
  ...personStore,
  id: '2',
  name: 'Дмитрий',
  surname: 'Костерин',
  patronymic: undefined,
  position: 'Backend',
  location: 'Томск',
  contacts: { telegram: 'KostDmitry', email: 'dmitry@example.com' },
  greeting: {
    first: 'Здравствуйте!',
    second: 'Я',
    secondHighlight: 'Дмитрий Костерин',
    third: 'Backend-разработчик: API, базы данных и всё, что работает под капотом',
  },
  theme: { primary: '#2E9E6B', secondary: '#7C5CFA' },
  hiddenContacts: [],
  availability: { status: AvailabilityStatus.Considering },
  courses: [],
  activities: [],
  sections: [
    { id: SectionId.About, enabled: true },
    { id: SectionId.Experience, enabled: true },
    { id: SectionId.TechStack, enabled: true },
    { id: SectionId.Skills, enabled: false },
    { id: SectionId.Education, enabled: true, collapsed: true },
    { id: SectionId.Projects, enabled: true },
    { id: SectionId.Contacts, enabled: true },
  ],
}

const personStores: Record<string, Person> = {
  '1': personStore,
  '2': secondPersonStore,
}

const greetings: Record<string, Record<Languages, Greeting>> = {
  '1': {
    ru: personStore.greeting,
    en: {
      first: 'Hello!',
      second: 'My name is',
      secondHighlight: 'Anton Sedletskii',
      third: "I'm a Frontend developer from Siberia",
    },
  },
  '2': {
    ru: secondPersonStore.greeting,
    en: {
      first: 'Hello!',
      second: "I'm",
      secondHighlight: 'Dmitry Kosterin',
      third: 'Backend developer: APIs, databases and everything under the hood',
    },
  },
}

/** Заметки о доступности на двух языках (сервер отдаёт текст на языке запроса) */
const availabilityNotes: Record<string, Record<Languages, string>> = {
  '1': {
    ru: 'Ищу роль Senior Frontend, удалённо или гибрид',
    en: 'Looking for a Senior Frontend role, remote or hybrid',
  },
}

/** У резюме «2» текстов нет — там работают значения по умолчанию */
const contactsTexts: Record<string, Record<Languages, ContactsText>> = {
  '1': {
    ru: {
      title: 'Расскажите о своём проекте',
      subtitle: 'Обсудим задачу, сроки и то, чем я могу быть полезен',
    },
    en: {
      title: 'Tell me about your project',
      subtitle: "Let's discuss the task, timeline and how I can help",
    },
  },
}

/** Ошибка той же формы, что у axios при 404 */
const notFoundError = (message: string) =>
  Object.assign(new Error(message), { response: { status: 404 } })

const projectsStore: Record<string, Project> = Object.fromEntries(
  personStore.projects.map((project) => [
    project.id,
    {
      ...project,
      role: 'Frontend developer',
      screenshots: [],
      gitLink: `https://github.com/example/${project.id}`,
      contributors: [
        {
          name: 'Антон',
          surname: 'Седлецкий',
          role: 'Frontend',
          contacts: { telegram: 'a_sedletskii', email: 'anton@example.com' },
        },
      ],
    },
  ]),
)

class CurriculumVitaeApiMock {
  baseUrl = '/cv'

  async getPersonData(cvId: string, lang: Languages = Languages.Ru): Promise<Person> {
    await delay()

    const person = personStores[cvId]

    if (!person) {
      throw notFoundError('CV not found')
    }

    return {
      ...person,
      greeting: greetings[cvId][lang],
      contactsText: contactsTexts[cvId]?.[lang],
      availability: person.availability && {
        ...person.availability,
        note: availabilityNotes[cvId]?.[lang],
      },
      updatedAt: dayjs().subtract(3, 'day').toISOString(),
    }
  }

  async getProject(projectId: string): Promise<Project> {
    await delay()

    const project = projectsStore[projectId]

    if (!project) {
      throw new Error('Project not found')
    }

    return project
  }
}

export const curriculumVitaeApiMock = new CurriculumVitaeApiMock()
