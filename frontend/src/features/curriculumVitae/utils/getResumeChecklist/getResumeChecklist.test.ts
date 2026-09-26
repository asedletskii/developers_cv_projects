import { Resume } from '@features/curriculumVitae/types'
import { getResumeChecklist } from './getResumeChecklist'

const emptyView = {
  languages: [],
  technologies: [],
  experiences: [],
  educations: [],
  courses: [],
  projects: [],
  activities: [],
} as unknown as Resume

describe('getResumeChecklist', () => {
  it('у пустого резюме ничего не заполнено', () => {
    expect(Object.values(getResumeChecklist(emptyView)).some(Boolean)).toBe(false)
  })

  it('считает секцию заполненной, когда в ней есть данные', () => {
    const view = {
      ...emptyView,
      profile: { name: 'Антон' },
      experiences: [{}],
      skills: { hardSkills: [], softSkills: ['teamwork'] },
      contacts: { telegram: 'anton' },
    } as unknown as Resume

    const checklist = getResumeChecklist(view)

    expect(checklist.profile).toBe(true)
    expect(checklist.experience).toBe(true)
    expect(checklist.skills).toBe(true)
    expect(checklist.contacts).toBe(true)
    expect(checklist.courses).toBe(false)
  })
})
