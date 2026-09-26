import { Person } from '@features/curriculumVitae/types'
import { normalizePerson } from './normalizePerson'

const basePerson = {
  id: '1',
  name: 'Антон',
  surname: 'Седлецкий',
  position: 'Frontend',
  photoPath: '',
  dateOfBirth: '1998-04-12',
  location: 'Новосибирск',
  relocation: true,
  contacts: {},
  greeting: { first: '', second: '', secondHighlight: '', third: '' },
} as unknown as Person

describe('normalizePerson', () => {
  it('превращает null в пустой массив для всех обязательных полей верхнего уровня', () => {
    const raw = {
      ...basePerson,
      sections: null,
      workFormat: null,
      languages: null,
      experience: null,
      education: null,
      courses: null,
      activities: null,
      hardSkills: null,
      softSkills: null,
      technologies: null,
      projects: null,
    } as unknown as Person

    const result = normalizePerson(raw)

    expect(result.sections).toEqual([])
    expect(result.workFormat).toEqual([])
    expect(result.languages).toEqual([])
    expect(result.experience).toEqual([])
    expect(result.education).toEqual([])
    expect(result.courses).toEqual([])
    expect(result.activities).toEqual([])
    expect(result.hardSkills).toEqual([])
    expect(result.softSkills).toEqual([])
    expect(result.technologies).toEqual([])
    expect(result.projects).toEqual([])
  })

  it('то же самое для undefined', () => {
    const raw = { ...basePerson, technologies: undefined } as unknown as Person

    expect(normalizePerson(raw).technologies).toEqual([])
  })

  it('не трогает уже заполненные массивы', () => {
    const raw = { ...basePerson, technologies: [{ name: 'React' }] }

    expect(normalizePerson(raw).technologies).toEqual([{ name: 'React' }])
  })

  it('приводит к пустому массиву achievements и stack внутри каждого места работы', () => {
    const raw = {
      ...basePerson,
      experience: [
        {
          id: 'e1',
          company: 'A',
          position: 'Dev',
          startDate: '2020-01',
          achievements: null,
          stack: null,
        },
      ],
    } as unknown as Person

    const [job] = normalizePerson(raw).experience

    expect(job.achievements).toEqual([])
    expect(job.stack).toEqual([])
  })

  it('сохраняет остальные поля места работы без изменений', () => {
    const raw = {
      ...basePerson,
      experience: [
        {
          id: 'e1',
          company: 'Северный ветер',
          position: 'Senior Frontend',
          startDate: '2023-03',
          achievements: ['Сделал X'],
          stack: [{ name: 'React' }],
        },
      ],
    } as unknown as Person

    const [job] = normalizePerson(raw).experience

    expect(job).toMatchObject({
      id: 'e1',
      company: 'Северный ветер',
      position: 'Senior Frontend',
      achievements: ['Сделал X'],
      stack: [{ name: 'React' }],
    })
  })

  it('приводит к пустому массиву stack внутри каждого проекта', () => {
    const raw = {
      ...basePerson,
      projects: [{ id: 'p1', name: 'CV Platform', stack: null }],
    } as unknown as Person

    expect(normalizePerson(raw).projects[0].stack).toEqual([])
  })

  it('не падает на пустом человеке без опциональных/массивных полей вовсе', () => {
    expect(() => normalizePerson(basePerson)).not.toThrow()
  })
})
