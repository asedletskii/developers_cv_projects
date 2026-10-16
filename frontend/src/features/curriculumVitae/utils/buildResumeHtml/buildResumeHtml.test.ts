import { ContactKey, Resume } from '@features/curriculumVitae/types'
import { buildResumeHtml, escapeHtml, ResumeHtmlLabels, sanitizeFileName } from './buildResumeHtml'

const labels: ResumeHtmlLabels = {
  lang: 'ru',
  sections: {
    technologies: 'Стек',
    experience: 'Опыт',
    education: 'Образование',
    courses: 'Курсы',
    projects: 'Проекты',
    activities: 'Активность',
    skills: 'Навыки',
    languages: 'Языки',
    contacts: 'Контакты',
  },
  present: 'по настоящее время',
  tasks: 'Задачи',
  achievements: 'Достижения',
  hardSkills: 'Профессиональные',
  softSkills: 'Личные',
  hoursSuffix: 'ч',
  contactsHidden: 'Контакты скрыты — смотрите на странице',
  generated: 'Сформировано 5 октября 2026',
  actualVersion: 'Актуальная версия',
  skillNames: { testing: 'Тестирование' },
  languageNames: { en: 'Английский' },
}

const resume = {
  id: 'r1',
  title: 'CV',
  language: 'ru',
  overview: { greeting: { first: '', second: '', secondHighlight: '', third: '' }, sections: [] },
  profile: {
    name: 'Антон',
    surname: 'Седлецкий',
    position: 'Frontend-разработчик',
    photoPath: '',
    dateOfBirth: '',
    location: 'Новосибирск',
    relocation: false,
    workFormat: [],
  },
  languages: [{ language: 'en', grade: 'B2' }],
  technologies: [{ name: 'TypeScript' }],
  experiences: [
    {
      id: 'e1',
      company: 'Acme',
      position: 'Разработчик',
      startDate: '2020-01',
      achievements: ['Ускорил сборку'],
      tasks: ['Писал фичи'],
      stack: [{ name: 'React' }],
    },
  ],
  educations: [],
  courses: [],
  activities: [],
  projects: [],
  skills: { hardSkills: ['testing'], softSkills: [] },
} as unknown as Resume

const build = (overrides: Partial<Parameters<typeof buildResumeHtml>[0]> = {}) =>
  buildResumeHtml({
    resume,
    fullName: 'Седлецкий Антон',
    contacts: {},
    labels,
    sourceUrl: 'https://cv.example.ru/cv/abc',
    ...overrides,
  })

describe('escapeHtml', () => {
  it('экранирует разметку и кавычки', () => {
    expect(escapeHtml(`<script>alert("x")</script> & 'y'`)).toBe(
      '&lt;script&gt;alert(&quot;x&quot;)&lt;/script&gt; &amp; &#39;y&#39;',
    )
  })
})

describe('buildResumeHtml', () => {
  it('собирает самодостаточный документ на языке резюме', () => {
    const html = build()

    expect(html).toContain('<html lang="ru">')
    expect(html).toContain('<h1>Седлецкий Антон</h1>')
    expect(html).toContain('Frontend-разработчик')
    expect(html).toContain('Разработчик — Acme')
    expect(html).toContain('2020-01 — по настоящее время')
    expect(html).toContain('Тестирование')
    expect(html).toContain('Английский (B2)')
    // ничего внешнего: ни скриптов, ни подключаемых стилей и шрифтов
    expect(html).not.toMatch(/<script|<link /i)
  })

  it('не вставляет пользовательский текст как разметку', () => {
    const html = build({
      fullName: '<img src=x onerror=alert(1)>',
      resume: {
        ...resume,
        experiences: [{ ...resume.experiences[0], description: '<b>жирный</b>' }],
      } as Resume,
    })

    expect(html).not.toContain('<img src=x')
    expect(html).not.toContain('<b>жирный</b>')
    expect(html).toContain('&lt;b&gt;жирный&lt;/b&gt;')
  })

  it('ссылки автора только с http(s): javascript: превращается в обычный текст', () => {
    const html = build({
      resume: {
        ...resume,
        projects: [{ id: 'p', name: 'P', link: 'javascript:alert(1)', stack: [] }],
        courses: [
          {
            id: 'c',
            name: 'C',
            organization: 'O',
            endDate: '2024-01',
            certificateLink: 'https://cert.example.ru/1',
          },
        ],
      } as unknown as Resume,
    })

    expect(html).not.toContain('javascript:')
    expect(html).toContain('href="https://cert.example.ru/1"')
  })

  it('контакты: значения со ссылками нужного вида', () => {
    const html = build({
      contacts: {
        [ContactKey.Telegram]: 'a_sedletskii',
        [ContactKey.Email]: 'me@example.ru',
        [ContactKey.Phone]: '+79998882233',
      },
    })

    expect(html).toContain('href="https://t.me/a_sedletskii"')
    expect(html).toContain('>@a_sedletskii<')
    expect(html).toContain('href="mailto:me@example.ru"')
    expect(html).toContain('href="tel:+79998882233"')
  })

  it('если скрытые контакты получить не удалось — пояснение вместо значений', () => {
    const html = build({ unresolvedContacts: [ContactKey.Email, ContactKey.Phone] })

    expect(html).toContain('Контакты скрыты — смотрите на странице')
    expect(html).not.toContain('mailto:')
  })

  it('пустые секции не попадают в файл', () => {
    const html = build()

    expect(html).not.toContain('<h2>Образование</h2>')
    expect(html).not.toContain('<h2>Проекты</h2>')
    expect(html).not.toContain('<h2>Контакты</h2>')
    expect(html).toContain('<h2>Опыт</h2>')
  })

  it('внизу — дата формирования и ссылка на актуальную версию', () => {
    const html = build()

    expect(html).toContain('Сформировано 5 октября 2026')
    expect(html).toContain('href="https://cv.example.ru/cv/abc"')
  })
})

describe('sanitizeFileName', () => {
  it('убирает символы, недопустимые в именах файлов', () => {
    expect(sanitizeFileName('Седлецкий Антон: резюме/CV?')).toBe('Седлецкий Антон резюме CV')
  })

  it('пустое имя заменяется запасным', () => {
    expect(sanitizeFileName('  ///  ')).toBe('resume')
  })
})
