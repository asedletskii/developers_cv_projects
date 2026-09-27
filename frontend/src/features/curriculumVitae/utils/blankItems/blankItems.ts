import {
  Activity,
  ActivityType,
  Course,
  Education,
  EducationLevel,
  Experience,
  Project,
} from '@features/curriculumVitae/types'

/** Новые пустые записи списков в форме редактора. Клиентский `id` нужен только форме: бэкенд
 *  выдаёт свой при создании, а запись без известного ему id отправляется `POST`-ом. */

export const blankExperience = (): Experience => ({
  id: crypto.randomUUID(),
  company: '',
  position: '',
  startDate: '',
  achievements: [],
  stack: [],
})

export const blankEducation = (): Education => ({
  id: crypto.randomUUID(),
  location: '',
  university: '',
  grade: EducationLevel.Bachelor,
  programName: '',
  faculty: '',
})

export const blankCourse = (): Course => ({
  id: crypto.randomUUID(),
  name: '',
  organization: '',
  endDate: '',
})

export const blankActivity = (): Activity => ({
  id: crypto.randomUUID(),
  type: ActivityType.OpenSource,
  title: '',
})

export const blankProject = (): Project => ({
  id: crypto.randomUUID(),
  name: '',
  screenshots: [],
  stack: [],
  contributors: [],
})
