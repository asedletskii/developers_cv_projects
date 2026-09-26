import { ResumeEditorTab, ResumeView } from '@features/curriculumVitae/types'

export type ChecklistTab = Exclude<ResumeEditorTab, 'overview'>

/** Что считается заполненным в резюме: вкладка редактора → есть ли в ней данные. */
export const getResumeChecklist = (view: ResumeView): Record<ChecklistTab, boolean> => ({
  profile: !!view.profile?.name,
  languages: view.languages.length > 0,
  technologies: view.technologies.length > 0,
  experience: view.experiences.length > 0,
  skills: !!view.skills && (view.skills.hardSkills.length > 0 || view.skills.softSkills.length > 0),
  education: view.educations.length > 0,
  courses: view.courses.length > 0,
  projects: view.projects.length > 0,
  activities: view.activities.length > 0,
  contacts: !!(view.contacts?.phone || view.contacts?.telegram || view.contacts?.email),
})
