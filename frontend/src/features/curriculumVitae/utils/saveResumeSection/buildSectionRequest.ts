import { EMPTY_CONTACTS, EMPTY_PROFILE, EMPTY_SKILLS } from '@features/curriculumVitae/constants'
import {
  ListSection,
  Resume,
  ResumeEditorTab,
  SingletonSection,
} from '@features/curriculumVitae/types'
import { toOptionalNumber, toOptionalString } from '@features/curriculumVitae/utils/toOptional'
import { withPinnedAbout } from '@features/curriculumVitae/utils/withPinnedAbout'

/** Свои клиентские id есть только у новых записей списков (`blankItems`) — на бэк их отправлять
 *  нельзя, он их не знает и ответит 400 «does not belong to resume». Известный id — тот, что
 *  реально пришёл с сервера в `view`; остальные считаем новыми записями (id не шлём). */
const knownId = (knownIds: Set<string | undefined>, id?: string) =>
  id && knownIds.has(id) ? id : undefined

export type SectionRequest = {
  section: SingletonSection | ListSection
  body: object
}

/**
 * Собирает тело `PUT`-запроса ровно для одной секции из значений формы — какую секцию сохраняем
 * и что в неё положить, включая текущую версию резюме (бэк её сверяет: несовпадение значит, что
 * резюме уже сохранили в другой вкладке/окне, и отвечает `412`). Чистая функция без похода на
 * бэк — в своём файле специально, чтобы не тянуть за собой API-клиент (и через него `axios`)
 * туда, где нужна только сборка тела запроса, включая тесты. Сам запрос делает `saveResumeSection`.
 */
export const buildSectionRequest = (
  view: Resume,
  section: ResumeEditorTab,
  values: Resume,
): SectionRequest => {
  const { version } = view
  const { overview, languages, technologies, experiences, educations, courses, activities } = values
  const { profile = EMPTY_PROFILE, contacts = EMPTY_CONTACTS, skills = EMPTY_SKILLS } = values

  switch (section) {
    case 'overview':
      return {
        section: 'overview',
        body: {
          version,
          theme: overview.theme?.primary || overview.theme?.secondary ? overview.theme : undefined,
          greeting: overview.greeting,
          availability: overview.availability?.status
            ? {
                status: overview.availability.status,
                from: toOptionalString(overview.availability.from),
                note: toOptionalString(overview.availability.note),
              }
            : undefined,
          sections: withPinnedAbout(overview.sections).map(({ id, enabled, collapsed }) => ({
            id,
            enabled,
            collapsed: !!collapsed,
          })),
        },
      }

    case 'profile':
      return {
        section: 'profile',
        body: {
          version,
          name: profile.name,
          surname: profile.surname,
          patronymic: toOptionalString(profile.patronymic),
          position: profile.position,
          photoPath: profile.photoPath,
          dateOfBirth: profile.dateOfBirth,
          location: profile.location,
          relocation: profile.relocation,
          workFormat: profile.workFormat,
        },
      }

    case 'contacts':
      return {
        section: 'contacts',
        body: {
          version,
          phone: toOptionalString(contacts.phone),
          telegram: toOptionalString(contacts.telegram),
          email: toOptionalString(contacts.email),
          hiddenContacts: contacts.hiddenContacts,
          // Пустая строка в contactsText значима («элемент скрыт»), поэтому передаём как есть
          contactsText: contacts.contactsText,
        },
      }

    case 'skills':
      return {
        section: 'skills',
        body: { version, hardSkills: skills.hardSkills, softSkills: skills.softSkills },
      }

    case 'languages': {
      const knownIds = new Set(view.languages.map(({ id }) => id))

      return {
        section: 'languages',
        body: {
          version,
          items: languages.map(({ id, language, grade }, position) => ({
            id: knownId(knownIds, id),
            language,
            grade: toOptionalString(grade),
            position,
          })),
        },
      }
    }

    case 'technologies': {
      const knownIds = new Set(view.technologies.map(({ id }) => id))

      return {
        section: 'technologies',
        body: {
          version,
          items: technologies.map(({ id, name, grade }, position) => ({
            id: knownId(knownIds, id),
            name,
            grade: toOptionalString(grade),
            position,
          })),
        },
      }
    }

    case 'experience': {
      const knownIds = new Set(view.experiences.map(({ id }) => id))

      return {
        section: 'experiences',
        body: {
          version,
          items: experiences.map((e) => ({
            id: knownId(knownIds, e.id),
            company: e.company,
            position: e.position,
            location: toOptionalString(e.location),
            startDate: e.startDate,
            endDate: toOptionalString(e.endDate),
            description: toOptionalString(e.description),
            achievements: e.achievements.filter(Boolean),
            tasks: e.tasks.filter(Boolean),
            stack: e.stack.map(({ name, grade }) => ({ name, grade: toOptionalString(grade) })),
          })),
        },
      }
    }

    case 'education': {
      const knownIds = new Set(view.educations.map(({ id }) => id))

      return {
        section: 'educations',
        body: {
          version,
          items: educations.map((e) => ({
            id: knownId(knownIds, e.id),
            location: e.location,
            university: e.university,
            grade: e.grade,
            programName: e.programName,
            profileName: toOptionalString(e.profileName),
            faculty: e.faculty,
            course: toOptionalNumber(e.course),
            endDate: toOptionalString(e.endDate),
          })),
        },
      }
    }

    case 'courses': {
      const knownIds = new Set(view.courses.map(({ id }) => id))

      return {
        section: 'courses',
        body: {
          version,
          items: courses.map((c) => ({
            id: knownId(knownIds, c.id),
            name: c.name,
            organization: c.organization,
            endDate: c.endDate,
            hours: toOptionalNumber(c.hours),
            certificateLink: toOptionalString(c.certificateLink),
          })),
        },
      }
    }

    case 'activities': {
      const knownIds = new Set(view.activities.map(({ id }) => id))

      return {
        section: 'activities',
        body: {
          version,
          items: activities.map((a) => ({
            id: knownId(knownIds, a.id),
            type: a.type,
            title: a.title,
            description: toOptionalString(a.description),
            date: toOptionalString(a.date),
            link: toOptionalString(a.link),
          })),
        },
      }
    }

    case 'projects': {
      const knownIds = new Set(view.projects.map(({ id }) => id))

      return {
        section: 'projects',
        body: {
          version,
          items: values.projects.map((project) => ({
            id: knownId(knownIds, project.id),
            name: project.name,
            description: toOptionalString(project.description),
            role: toOptionalString(project.role),
            screenshots: project.screenshots,
            link: toOptionalString(project.link),
            gitLink: toOptionalString(project.gitLink),
            stack: project.stack.map(({ name, grade }) => ({
              name,
              grade: toOptionalString(grade),
            })),
            contributors: project.contributors.map((c) => ({
              ...c,
              contacts: { ...c.contacts, phone: c.contacts.phone ?? '' },
            })),
          })),
        },
      }
    }
  }
}
