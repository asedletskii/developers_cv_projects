import { ReactNode } from 'react'
import { AboutSection } from '@features/curriculumVitae/components/AboutSection'
import { ActivitiesSection } from '@features/curriculumVitae/components/ActivitiesSection'
import { ContactsSection } from '@features/curriculumVitae/components/ContactsSection'
import { CoursesSection } from '@features/curriculumVitae/components/CoursesSection'
import { EducationSection } from '@features/curriculumVitae/components/EducationSection'
import { ExperienceSection } from '@features/curriculumVitae/components/ExperienceSection'
import { ProjectsSection } from '@features/curriculumVitae/components/ProjectsSection'
import { SkillsSection } from '@features/curriculumVitae/components/SkillsSection'
import { TechStackSlider } from '@features/curriculumVitae/components/TechStackSlider'
import { Person, SectionId, SectionSettings } from '@features/curriculumVitae/types'
import { withPinnedAbout } from '@features/curriculumVitae/utils/withPinnedAbout'

type SectionDefinition = {
  /** Есть ли что показывать: пустые секции не рендерятся, даже если включены */
  hasContent(person: Person): boolean

  render(person: Person): ReactNode

  /** Пункт в плавающей навигации */
  inNav: boolean

  /** Может ли пользователь свернуть секцию (и учитывать ли `collapsed` из API) */
  isCollapsible: boolean

  /** Занимает всю ширину экрана, а не контейнер */
  isFullWidth?: boolean
}

export const SECTION_DEFINITIONS: Record<SectionId, SectionDefinition> = {
  [SectionId.About]: {
    isCollapsible: false,
    inNav: true,
    hasContent: () => true,
    render: (person) => <AboutSection {...person} />,
  },

  [SectionId.TechStack]: {
    isCollapsible: false,
    inNav: false,
    isFullWidth: true,
    hasContent: (person) => !!person.technologies.length,
    render: (person) => <TechStackSlider technologies={person.technologies} />,
  },

  [SectionId.Experience]: {
    isCollapsible: false,
    inNav: true,
    hasContent: (person) => !!person.experience.length,
    render: (person) => <ExperienceSection experience={person.experience} />,
  },

  [SectionId.Skills]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (person) => !!person.hardSkills.length || !!person.softSkills.length,
    render: (person) => (
      <SkillsSection hardSkills={person.hardSkills} softSkills={person.softSkills} />
    ),
  },

  [SectionId.Education]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (person) => !!person.education.length,
    render: (person) => <EducationSection education={person.education} />,
  },

  [SectionId.Courses]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (person) => !!person.courses.length,
    render: (person) => <CoursesSection courses={person.courses} />,
  },

  [SectionId.Projects]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (person) => !!person.projects.length,
    render: (person) => <ProjectsSection projects={person.projects} />,
  },

  [SectionId.Activities]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (person) => !!person.activities.length,
    render: (person) => <ActivitiesSection activities={person.activities} />,
  },

  [SectionId.Contacts]: {
    isCollapsible: false,
    inNav: true,
    hasContent: (person) => Object.values(person.contacts).some(Boolean),
    render: (person) => (
      <ContactsSection
        contacts={person.contacts}
        text={person.contactsText}
        hidden={person.hiddenContacts}
      />
    ),
  },
}

/** Порядок по умолчанию, если API не прислал настройки секций */
export const DEFAULT_SECTIONS: SectionSettings[] = (
  Object.keys(SECTION_DEFINITIONS) as SectionId[]
).map((id) => ({ id, enabled: true }))

/** Секции в порядке из настроек: включённые, известные и с данными */
export const getVisibleSections = (person: Person): SectionSettings[] =>
  withPinnedAbout(person.sections?.length ? person.sections : DEFAULT_SECTIONS).filter(
    ({ id, enabled }) => enabled && SECTION_DEFINITIONS[id]?.hasContent(person),
  )
