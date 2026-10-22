import { Resume } from '@features/curriculumVitae/types'
import { normalizeResume } from './normalizeResume'

const rawEmpty = {
  id: 'r1',
  title: 'CV',
  language: 'ru',
  overview: { sections: null },
  skills: {},
  contacts: {},
  profile: { name: 'Антон' },
  languages: null,
  technologies: null,
  experiences: [{ id: 'e1', company: 'A', position: 'B', startDate: '2020-01' }],
  projects: [{ id: 'p1', name: 'P' }],
} as unknown as Resume

describe('normalizeResume', () => {
  it('заменяет отсутствующие и null-списки пустыми, не трогая заполненное', () => {
    const resume = normalizeResume(rawEmpty)

    expect(resume.overview.sections).toEqual([])
    expect(resume.skills).toEqual({ hardSkills: [], softSkills: [] })
    expect(resume.contacts?.hiddenContacts).toEqual([])
    expect(resume.profile?.workFormat).toEqual([])
    expect(resume.languages).toEqual([])
    expect(resume.experiences[0].tasks).toEqual([])
    expect(resume.projects[0].stack).toEqual([])
    expect(resume.profile?.name).toBe('Антон')
  })
})
