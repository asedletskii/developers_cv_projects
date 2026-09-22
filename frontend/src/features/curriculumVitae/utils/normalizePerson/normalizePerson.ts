import { Person } from '@features/curriculumVitae/types'

/**
 * Массивы, которые фронт считает обязательными (`Person['experience']` и т.д.), Go на бэкенде
 * сериализует как `null`, если срез нигде явно не проинициализирован как `[]T{}` — на эти поля
 * не навешан `omitempty`, но нулевой (`nil`) срез всё равно уходит в JSON как `null`, а не `[]`.
 * Приводим такие поля к пустому массиву сразу на границе API, чтобы остальной код мог спокойно
 * звать `.length`/`.map` без опциональной цепочки на каждом шагу.
 */
const arr = <T>(value: T[] | null | undefined): T[] => value ?? []

export const normalizePerson = (person: Person): Person => ({
  ...person,
  sections: arr(person.sections),
  workFormat: arr(person.workFormat),
  languages: arr(person.languages),
  experience: arr(person.experience).map((item) => ({
    ...item,
    achievements: arr(item.achievements),
    stack: arr(item.stack),
  })),
  education: arr(person.education),
  courses: arr(person.courses),
  activities: arr(person.activities),
  hardSkills: arr(person.hardSkills),
  softSkills: arr(person.softSkills),
  technologies: arr(person.technologies),
  projects: arr(person.projects).map((item) => ({ ...item, stack: arr(item.stack) })),
})
