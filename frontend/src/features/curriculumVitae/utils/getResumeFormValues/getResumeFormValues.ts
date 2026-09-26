import { EMPTY_CONTACTS, EMPTY_PROFILE, EMPTY_SKILLS } from '@features/curriculumVitae/constants'
import { DEFAULT_SECTIONS } from '@features/curriculumVitae/sections'
import { Resume } from '@features/curriculumVitae/types'
import { withPinnedAbout } from '@features/curriculumVitae/utils/withPinnedAbout'

/**
 * Значения формы редактора из резюме владельца. Форма у всех вкладок одна и та же (`Resume`), но
 * каждая вкладка правит и отправляет только свою часть — см. `saveResumeSection`. Отсутствующие
 * секции заполняем пустыми значениями, чтобы поля были управляемыми.
 */
export const getResumeFormValues = (view: Resume): Resume => ({
  ...view,
  overview: {
    ...view.overview,
    // Только что созданное резюме приходит без секций — предлагаем стандартный набор
    sections: withPinnedAbout(
      view.overview.sections.length ? view.overview.sections : DEFAULT_SECTIONS,
    ),
  },
  profile: view.profile ?? EMPTY_PROFILE,
  contacts: view.contacts ?? EMPTY_CONTACTS,
  skills: view.skills ?? EMPTY_SKILLS,
})
