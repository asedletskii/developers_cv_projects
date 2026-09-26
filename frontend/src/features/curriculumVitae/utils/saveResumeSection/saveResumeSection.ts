import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import {
  LanguageType,
  ListSection,
  Person,
  ResumeEditorTab,
  ResumeView,
} from '@features/curriculumVitae/types'
import { ListOps, syncById, syncByIndex } from '@features/curriculumVitae/utils/syncList'
import { toOptionalNumber, toOptionalString } from '@features/curriculumVitae/utils/toOptional'
import { withPinnedAbout } from '@features/curriculumVitae/utils/withPinnedAbout'

/** Отправляет на бэк ровно одну секцию из значений формы. */
export const saveResumeSection = async (
  userId: string,
  view: ResumeView,
  section: ResumeEditorTab,
  p: Person,
): Promise<void> => {
  const { id: resumeId } = view
  const cv = curriculumVitaeApi

  const list = <F>(kind: ListSection, toBody: (item: F, index: number) => unknown): ListOps<F> => ({
    create: (item, index) => cv.createItem(userId, resumeId, kind, toBody(item, index)),
    update: (id, item, index) => cv.updateItem(userId, resumeId, kind, id, toBody(item, index)),
    remove: (id) => cv.deleteItem(userId, resumeId, kind, id),
  })

  switch (section) {
    case 'overview':
      return cv.saveSingleton(
        userId,
        resumeId,
        'overview',
        {
          presentation: {
            theme: p.theme?.primary || p.theme?.secondary ? p.theme : undefined,
            greeting: p.greeting,
            availability: p.availability?.status
              ? {
                  status: p.availability.status,
                  from: toOptionalString(p.availability.from),
                  note: toOptionalString(p.availability.note),
                }
              : undefined,
            resumeLink: toOptionalString(p.resumeLink),
            sections: withPinnedAbout(p.sections).map(({ id, enabled, collapsed }) => ({
              id,
              enabled,
              collapsed: !!collapsed,
            })),
          },
        },
        true,
      )

    case 'profile':
      return cv.saveSingleton(
        userId,
        resumeId,
        'profile',
        {
          name: p.name,
          surname: p.surname,
          patronymic: toOptionalString(p.patronymic),
          position: p.position,
          photoPath: p.photoPath,
          dateOfBirth: p.dateOfBirth,
          location: p.location,
          relocation: p.relocation,
          workFormat: p.workFormat,
        },
        !!view.profile,
      )

    case 'contacts':
      return cv.saveSingleton(
        userId,
        resumeId,
        'contacts',
        {
          phone: toOptionalString(p.contacts.phone),
          telegram: toOptionalString(p.contacts.telegram),
          email: toOptionalString(p.contacts.email),
          hiddenContacts: p.hiddenContacts ?? [],
          // Пустая строка в contactsText значима («элемент скрыт»), поэтому передаём как есть
          contactsText: p.contactsText,
        },
        !!view.contacts,
      )

    case 'skills':
      return cv.saveSingleton(
        userId,
        resumeId,
        'skills',
        {
          hardSkills: p.hardSkills,
          softSkills: p.softSkills,
        },
        !!view.skills,
      )

    case 'languages': {
      // Язык в резюме уникален, поэтому записи сопоставляем по коду языка, а не по индексу: иначе
      // перестановка двух языков местами упёрлась бы в `409` на первом же PUT. Порядок задаём отдельно.
      const byLanguage = new Map(view.languages.map((item) => [item.language, item.id]))
      const wanted = new Set(p.languages.map(({ language }) => language))
      const body = ({ grade, language }: (typeof p.languages)[number]) => ({
        language,
        grade: toOptionalString(grade),
      })

      for (const [language, id] of byLanguage) {
        if (!wanted.has(language as LanguageType))
          await cv.deleteItem(userId, resumeId, 'languages', id)
      }

      const ids: string[] = []
      for (const item of p.languages) {
        const existingId = byLanguage.get(item.language)

        if (existingId) {
          await cv.updateItem(userId, resumeId, 'languages', existingId, body(item))
          ids.push(existingId)
        } else {
          ids.push(await cv.createItem(userId, resumeId, 'languages', body(item)))
        }
      }

      if (ids.length > 1) await cv.reorderItems(userId, resumeId, 'languages', ids)
      return
    }

    case 'technologies':
      return syncByIndex(
        view.technologies,
        p.technologies,
        list('technologies', ({ name }) => ({ name })),
      )

    case 'experience':
      return syncById(
        view.experiences,
        p.experience,
        list('experiences', (e) => ({
          company: e.company,
          position: e.position,
          location: toOptionalString(e.location),
          startDate: e.startDate,
          endDate: toOptionalString(e.endDate),
          description: toOptionalString(e.description),
          achievements: e.achievements.filter(Boolean),
          technologies: e.stack.map(({ name }) => ({ name })),
        })),
      )

    case 'education':
      return syncById(
        view.educations,
        p.education,
        list('educations', (e) => ({
          location: e.location,
          university: e.university,
          grade: e.grade,
          programName: e.programName,
          profileName: toOptionalString(e.profileName),
          faculty: e.faculty,
          course: toOptionalNumber(e.course),
          endDate: toOptionalString(e.endDate),
        })),
      )

    case 'courses':
      return syncById(
        view.courses,
        p.courses,
        list('courses', (c) => ({
          name: c.name,
          organization: c.organization,
          endDate: c.endDate,
          hours: toOptionalNumber(c.hours),
          certificateLink: toOptionalString(c.certificateLink),
        })),
      )

    case 'activities':
      return syncById(
        view.activities,
        p.activities,
        list('activities', (a) => ({
          type: a.type,
          title: a.title,
          description: toOptionalString(a.description),
          date: toOptionalString(a.date),
          link: toOptionalString(a.link),
        })),
      )

    case 'projects':
      return syncById(
        view.projects,
        p.projects,
        list('projects', (project) => {
          // Форма правит только название/описание/ссылку/стэк — остальные поля берём с сервера,
          // иначе PUT затёр бы роль, скриншоты, git-ссылку и участников
          const original = view.projects.find(({ id }) => id === project.id)

          return {
            name: project.name,
            description: toOptionalString(project.description),
            link: toOptionalString(project.link),
            stack: project.stack.map(({ name }) => ({ name })),
            role: original?.role,
            screenshots: original?.screenshots ?? [],
            gitLink: original?.gitLink,
            contributors: (original?.contributors ?? []).map((c) => ({
              ...c,
              contacts: { ...c.contacts, phone: c.contacts.phone ?? '' },
            })),
          }
        }),
      )
  }
}
