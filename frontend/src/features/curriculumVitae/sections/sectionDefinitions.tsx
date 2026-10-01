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
import { EMPTY_CONTACTS, EMPTY_PROFILE, EMPTY_SKILLS } from '@features/curriculumVitae/constants'
import { Resume, SectionId, SectionSettings } from '@features/curriculumVitae/types'
import { withPinnedAbout } from '@features/curriculumVitae/utils/withPinnedAbout'
import { DEFAULT_SECTIONS } from './defaultSections'

type SectionDefinition = {
  hasContent(resume: Resume): boolean
  render(resume: Resume): ReactNode
  inNav: boolean
  isCollapsible: boolean
  isFullWidth?: boolean
}

export const SECTION_DEFINITIONS: Record<SectionId, SectionDefinition> = {
  [SectionId.About]: {
    isCollapsible: false,
    inNav: true,
    hasContent: () => true,
    render: (resume) => (
      <AboutSection profile={resume.profile ?? EMPTY_PROFILE} languages={resume.languages} />
    ),
  },

  [SectionId.TechStack]: {
    isCollapsible: false,
    inNav: false,
    isFullWidth: true,
    hasContent: (resume) => !!resume.technologies?.length,
    render: (resume) => <TechStackSlider technologies={resume.technologies} />,
  },

  [SectionId.Experience]: {
    isCollapsible: false,
    inNav: true,
    hasContent: (resume) => !!resume.experiences?.length,
    render: (resume) => <ExperienceSection experience={resume.experiences} />,
  },

  [SectionId.Skills]: {
    isCollapsible: true,
    inNav: true,
    hasContent: ({ skills = EMPTY_SKILLS }) =>
      !!skills.hardSkills.length || !!skills.softSkills.length,
    render: ({ skills = EMPTY_SKILLS }) => (
      <SkillsSection hardSkills={skills.hardSkills} softSkills={skills.softSkills} />
    ),
  },

  [SectionId.Education]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (resume) => !!resume.educations?.length,
    render: (resume) => <EducationSection education={resume.educations} />,
  },

  [SectionId.Courses]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (resume) => !!resume.courses?.length,
    render: (resume) => <CoursesSection courses={resume.courses} />,
  },

  [SectionId.Projects]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (resume) => !!resume.projects?.length,
    render: (resume) => <ProjectsSection projects={resume.projects} />,
  },

  [SectionId.Activities]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (resume) => !!resume.activities?.length,
    render: (resume) => <ActivitiesSection activities={resume.activities} />,
  },

  [SectionId.Contacts]: {
    isCollapsible: false,
    inNav: true,
    hasContent: ({ contacts = EMPTY_CONTACTS }) =>
      !!(contacts.phone || contacts.telegram || contacts.email || contacts.hiddenContacts.length),
    render: ({ contacts = EMPTY_CONTACTS }) => (
      <ContactsSection
        contacts={contacts}
        text={contacts.contactsText}
        hidden={contacts.hiddenContacts}
      />
    ),
  },
}

export const getVisibleSections = (resume: Resume): SectionSettings[] =>
  withPinnedAbout(
    resume.overview.sections.length ? resume.overview.sections : DEFAULT_SECTIONS,
  ).filter(({ id, enabled }) => enabled && SECTION_DEFINITIONS[id]?.hasContent(resume))
