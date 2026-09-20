import { EMPTY_CONTACTS, EMPTY_PROFILE, EMPTY_SKILLS } from '@features/curriculumVitae/constants'
import { Resume } from '@features/curriculumVitae/types'
import { getResumeFormValues } from './getResumeFormValues'

const baseView = {
  id: '1',
  title: 'Резюме',
  language: 'ru',
  languages: [],
  technologies: [],
  experiences: [],
  educations: [],
  courses: [],
  activities: [],
  projects: [],
} as unknown as Resume

describe('getResumeFormValues', () => {
  it('подставляет пустые значения для отсутствующих секций', () => {
    const view = { ...baseView, overview: { greeting: {}, sections: [] } } as unknown as Resume

    const values = getResumeFormValues(view)

    expect(values.profile).toEqual(EMPTY_PROFILE)
    expect(values.contacts).toEqual(EMPTY_CONTACTS)
    expect(values.skills).toEqual(EMPTY_SKILLS)
  })

  it('не трогает секции, которые уже пришли', () => {
    const profile = { ...EMPTY_PROFILE, name: 'Антон' }
    const view = {
      ...baseView,
      overview: { greeting: {}, sections: [] },
      profile,
    } as unknown as Resume

    expect(getResumeFormValues(view).profile).toBe(profile)
  })

  it('у нового резюме без секций подставляет секции по умолчанию с «Обо мне» первой', () => {
    const view = { ...baseView, overview: { greeting: {}, sections: [] } } as unknown as Resume

    const { sections } = getResumeFormValues(view).overview

    expect(sections.length).toBeGreaterThan(1)
    expect(sections[0].id).toBe('about')
  })

  it('сохраняет уже заданный порядок секций, но переставляет «Обо мне» в начало', () => {
    const view = {
      ...baseView,
      overview: {
        greeting: {},
        sections: [
          { id: 'contacts', enabled: true },
          { id: 'about', enabled: true },
        ],
      },
    } as unknown as Resume

    const { sections } = getResumeFormValues(view).overview

    expect(sections.map(({ id }) => id)).toEqual(['about', 'contacts'])
  })
})
