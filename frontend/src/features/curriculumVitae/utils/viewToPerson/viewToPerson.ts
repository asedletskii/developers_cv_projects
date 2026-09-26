import { DEFAULT_SECTIONS } from '@features/curriculumVitae/sections'
import { getCurriculumVitaeInitialState } from '@features/curriculumVitae/store/curriculumVitae/curriculumVitae.state'
import {
  ContactKey,
  HardSkill,
  Person,
  ResumeView,
  SoftSkill,
  WorkFormat,
} from '@features/curriculumVitae/types'
import { withPinnedAbout } from '@features/curriculumVitae/utils/withPinnedAbout'

/**
 * Ответ `GET /lk/resume/:id` → форма. Форма у всех вкладок одна и та же (`Person`), но каждая
 * вкладка правит и отправляет только свою часть — см. `saveSection`.
 */
export const viewToPerson = (view: ResumeView): Person => {
  const { overview, profile, contacts, skills } = view

  return {
    ...getCurriculumVitaeInitialState().person,
    id: view.id,

    theme: overview.theme,
    greeting: overview.greeting,
    availability: overview.availability,
    resumeLink: overview.resumeLink,
    // Только что созданное резюме приходит без секций — предлагаем стандартный набор
    sections: withPinnedAbout(overview.sections?.length ? overview.sections : DEFAULT_SECTIONS),

    name: profile?.name ?? '',
    surname: profile?.surname ?? '',
    patronymic: profile?.patronymic,
    position: profile?.position ?? '',
    photoPath: profile?.photoPath ?? '',
    dateOfBirth: profile?.dateOfBirth ?? '',
    location: profile?.location ?? '',
    relocation: profile?.relocation ?? false,
    workFormat: (profile?.workFormat ?? []) as WorkFormat[],

    contacts: { phone: contacts?.phone, telegram: contacts?.telegram, email: contacts?.email },
    hiddenContacts: (contacts?.hiddenContacts ?? []) as ContactKey[],
    contactsText: contacts?.contactsText,

    hardSkills: (skills?.hardSkills ?? []) as HardSkill[],
    softSkills: (skills?.softSkills ?? []) as SoftSkill[],

    languages: (view.languages ?? []).map(({ language, grade }) => ({
      language,
      grade,
    })) as Person['languages'],
    technologies: (view.technologies ?? []).map(({ name }) => ({ name })),
    experience: (view.experiences ?? []).map(({ technologies, ...rest }) => ({
      ...rest,
      achievements: rest.achievements ?? [],
      stack: technologies ?? [],
    })),
    education: (view.educations ?? []) as Person['education'],
    courses: view.courses ?? [],
    activities: (view.activities ?? []) as Person['activities'],
    projects: (view.projects ?? []).map(({ id, name, description, link, stack }) => ({
      id,
      name,
      description,
      link,
      stack: stack ?? [],
    })),
  }
}
