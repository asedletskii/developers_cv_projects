import { Resume } from '@features/curriculumVitae/types'
import { buildSectionRequest } from './buildSectionRequest'

const baseView = {
  id: 'resume-1',
  title: 'Резюме',
  language: 'ru',
  version: 4,
  overview: { greeting: { first: '', second: '', secondHighlight: '', third: '' }, sections: [] },
  languages: [],
  technologies: [],
  experiences: [],
  educations: [],
  courses: [],
  activities: [],
  projects: [],
} as unknown as Resume

describe('buildSectionRequest', () => {
  it('каждый запрос несёт текущую версию резюме', () => {
    const { body } = buildSectionRequest(baseView, 'skills', baseView)
    expect((body as { version: number }).version).toBe(4)
  })

  it('overview: пустая тема не отправляется, доступность — только при заданном статусе', () => {
    const values = {
      ...baseView,
      overview: {
        greeting: { first: 'Привет', second: '', secondHighlight: '', third: '' },
        theme: {},
        sections: [{ id: 'about', enabled: true }],
      },
    } as unknown as Resume

    const { section, body } = buildSectionRequest(baseView, 'overview', values)

    expect(section).toBe('overview')
    expect(body).toMatchObject({
      theme: undefined,
      availability: undefined,
      greeting: { first: 'Привет' },
    })
  })

  it('profile: пустой отчество не отправляется', () => {
    const values = {
      ...baseView,
      profile: { name: 'Антон', surname: 'Седлецкий', patronymic: '', position: '', photoPath: '' },
    } as unknown as Resume

    const { section, body } = buildSectionRequest(baseView, 'profile', values)

    expect(section).toBe('profile')
    expect(body).toMatchObject({ name: 'Антон', surname: 'Седлецкий', patronymic: undefined })
  })

  it('contacts: пустой контакт не отправляется, а contactsText передаётся как есть', () => {
    const values = {
      ...baseView,
      contacts: {
        phone: '',
        telegram: 'a',
        contactsText: { title: '' },
        hiddenContacts: ['phone'],
      },
    } as unknown as Resume

    const { section, body } = buildSectionRequest(baseView, 'contacts', values)

    expect(section).toBe('contacts')
    expect(body).toMatchObject({
      phone: undefined,
      telegram: 'a',
      hiddenContacts: ['phone'],
      contactsText: { title: '' },
    })
  })

  it('languages: известный id сохраняется, чужой/клиентский — отбрасывается, позиция — по индексу', () => {
    const view = { ...baseView, languages: [{ id: 'srv-1', language: 'en' }] } as unknown as Resume
    const values = {
      ...baseView,
      languages: [
        { id: 'srv-1', language: 'en', grade: 'c1' },
        { id: 'client-temp-id', language: 'de' },
      ],
    } as unknown as Resume

    const { section, body } = buildSectionRequest(view, 'languages', values)

    expect(section).toBe('languages')
    expect((body as { items: unknown[] }).items).toEqual([
      { id: 'srv-1', language: 'en', grade: 'c1', position: 0 },
      { id: undefined, language: 'de', grade: undefined, position: 1 },
    ])
  })

  it('technologies: та же фильтрация id, что и у языков', () => {
    const view = { ...baseView, technologies: [{ id: 'srv-1', name: 'Go' }] } as unknown as Resume
    const values = {
      ...baseView,
      technologies: [
        { id: 'srv-1', name: 'Go' },
        { id: 'client-temp-id', name: 'Rust' },
      ],
    } as unknown as Resume

    const { items } = buildSectionRequest(view, 'technologies', values).body as {
      items: { id?: string }[]
    }

    expect(items.map(({ id }) => id)).toEqual(['srv-1', undefined])
  })

  it('experience: стек уходит под ключом stack, а не technologies', () => {
    const view = { ...baseView, experiences: [{ id: 'srv-1' }] } as unknown as Resume
    const values = {
      ...baseView,
      experiences: [
        {
          id: 'srv-1',
          company: 'Acme',
          position: 'Dev',
          startDate: '2020-01',
          achievements: ['', 'Сделал'],
          tasks: ['', 'Ревьюил код'],
          stack: [{ name: 'Go' }],
        },
      ],
    } as unknown as Resume

    const { section, body } = buildSectionRequest(view, 'experience', values)

    expect(section).toBe('experiences')
    expect(body).toMatchObject({
      items: [
        expect.objectContaining({
          achievements: ['Сделал'],
          tasks: ['Ревьюил код'],
          stack: [{ name: 'Go', grade: undefined }],
        }),
      ],
    })
  })

  it('projects: сохраняет уже известные role/gitLink/screenshots/contributors из формы, пустой телефон участника — пустая строка', () => {
    const view = { ...baseView, projects: [{ id: 'srv-1' }] } as unknown as Resume
    const values = {
      ...baseView,
      projects: [
        {
          id: 'srv-1',
          name: 'Pet',
          role: 'Автор',
          screenshots: ['a.png'],
          gitLink: 'https://github.com/x',
          stack: [{ name: 'TS' }],
          contributors: [{ name: 'A', surname: 'B', role: 'dev', contacts: {} }],
        },
      ],
    } as unknown as Resume

    const { section, body } = buildSectionRequest(view, 'projects', values)

    expect(section).toBe('projects')
    expect(body).toMatchObject({
      items: [
        expect.objectContaining({
          role: 'Автор',
          screenshots: ['a.png'],
          gitLink: 'https://github.com/x',
          contributors: [expect.objectContaining({ contacts: { phone: '' } })],
        }),
      ],
    })
  })
})
