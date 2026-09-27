import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { EMPTY_CONTACTS, EMPTY_PROFILE, EMPTY_SKILLS } from '@features/curriculumVitae/constants'
import { LanguageType, ListSection, Resume, ResumeEditorTab } from '@features/curriculumVitae/types'
import { ListOps, syncById, syncByIndex } from '@features/curriculumVitae/utils/syncList'
import { toOptionalNumber, toOptionalString } from '@features/curriculumVitae/utils/toOptional'
import { withPinnedAbout } from '@features/curriculumVitae/utils/withPinnedAbout'

/** Отправляет на бэк ровно одну секцию из значений формы. */
export const saveResumeSection = async (
  userId: string,
  view: Resume,
  section: ResumeEditorTab,
  values: Resume,
): Promise<void> => {
  const { id: resumeId } = view
  const { overview, languages, technologies, experiences, educations, courses, activities } = values
  const { profile = EMPTY_PROFILE, contacts = EMPTY_CONTACTS, skills = EMPTY_SKILLS } = values
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
          theme: overview.theme?.primary || overview.theme?.secondary ? overview.theme : undefined,
          greeting: overview.greeting,
          availability: overview.availability?.status
            ? {
                status: overview.availability.status,
                from: toOptionalString(overview.availability.from),
                note: toOptionalString(overview.availability.note),
              }
            : undefined,
          resumeLink: toOptionalString(overview.resumeLink),
          sections: withPinnedAbout(overview.sections).map(({ id, enabled, collapsed }) => ({
            id,
            enabled,
            collapsed: !!collapsed,
          })),
        },
        true,
      )

    case 'profile':
      return cv.saveSingleton(
        userId,
        resumeId,
        'profile',
        {
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
        !!view.profile,
      )

    case 'contacts':
      return cv.saveSingleton(
        userId,
        resumeId,
        'contacts',
        {
          phone: toOptionalString(contacts.phone),
          telegram: toOptionalString(contacts.telegram),
          email: toOptionalString(contacts.email),
          hiddenContacts: contacts.hiddenContacts,
          // Пустая строка в contactsText значима («элемент скрыт»), поэтому передаём как есть
          contactsText: contacts.contactsText,
        },
        !!view.contacts,
      )

    case 'skills':
      return cv.saveSingleton(
        userId,
        resumeId,
        'skills',
        {
          hardSkills: skills.hardSkills,
          softSkills: skills.softSkills,
        },
        !!view.skills,
      )

    case 'languages': {
      // Язык в резюме уникален, поэтому записи сопоставляем по коду языка, а не по индексу: иначе
      // перестановка двух языков местами упёрлась бы в `409` на первом же PUT. Порядок задаём отдельно.
      const byLanguage = new Map(view.languages.map((item) => [item.language, item.id as string]))
      const wanted = new Set(languages.map(({ language }) => language))
      const body = ({ grade, language }: (typeof languages)[number]) => ({
        language,
        grade: toOptionalString(grade),
      })

      for (const [language, id] of byLanguage) {
        if (!wanted.has(language as LanguageType))
          await cv.deleteItem(userId, resumeId, 'languages', id)
      }

      const ids: string[] = []
      for (const item of languages) {
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
        view.technologies as { id: string }[],
        technologies,
        list('technologies', ({ name }) => ({ name })),
      )

    case 'experience':
      return syncById(
        view.experiences,
        experiences,
        list('experiences', (e) => ({
          company: e.company,
          position: e.position,
          location: toOptionalString(e.location),
          startDate: e.startDate,
          endDate: toOptionalString(e.endDate),
          description: toOptionalString(e.description),
          achievements: e.achievements.filter(Boolean),
          stack: e.stack.map(({ name }) => ({ name })),
        })),
      )

    case 'education':
      return syncById(
        view.educations,
        educations,
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
        courses,
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
        activities,
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
        values.projects,
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
