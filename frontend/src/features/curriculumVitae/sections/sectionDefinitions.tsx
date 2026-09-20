import { ReactNode } from 'react'
import { AboutSection } from '@features/curriculumVitae/components/cv/AboutSection'
import {
  AboutCompact,
  AboutRich,
} from '@features/curriculumVitae/components/cv/AboutSection/AboutVariants'
import { ActivitiesSection } from '@features/curriculumVitae/components/cv/ActivitiesSection'
import {
  ActivitiesCompact,
  ActivitiesRich,
} from '@features/curriculumVitae/components/cv/ActivitiesSection/ActivitiesVariants'
import { ContactsSection } from '@features/curriculumVitae/components/cv/ContactsSection'
import {
  ContactsCompact,
  ContactsRich,
} from '@features/curriculumVitae/components/cv/ContactsSection/ContactsVariants'
import { CoursesSection } from '@features/curriculumVitae/components/cv/CoursesSection'
import {
  CoursesCompact,
  CoursesRich,
} from '@features/curriculumVitae/components/cv/CoursesSection/CoursesVariants'
import { EducationSection } from '@features/curriculumVitae/components/cv/EducationSection'
import {
  EducationCompact,
  EducationRich,
} from '@features/curriculumVitae/components/cv/EducationSection/EducationVariants'
import { ExperienceSection } from '@features/curriculumVitae/components/cv/ExperienceSection'
import {
  ExperienceCompact,
  ExperienceRich,
} from '@features/curriculumVitae/components/cv/ExperienceSection/ExperienceVariants'
import { ProjectsSection } from '@features/curriculumVitae/components/cv/ProjectsSection'
import {
  ProjectsCompact,
  ProjectsRich,
} from '@features/curriculumVitae/components/cv/ProjectsSection/ProjectsVariants'
import { SkillsSection } from '@features/curriculumVitae/components/cv/SkillsSection'
import { TechStackSlider } from '@features/curriculumVitae/components/cv/TechStackSlider'
import {
  TechStackCompact,
  TechStackRich,
} from '@features/curriculumVitae/components/cv/TechStackSlider/TechStackVariants'
import { EMPTY_CONTACTS, EMPTY_PROFILE, EMPTY_SKILLS } from '@features/curriculumVitae/constants'
import { Resume, SectionId, SectionSettings } from '@features/curriculumVitae/types'
import { withPinnedAbout } from '@features/curriculumVitae/utils/withPinnedAbout'
import { DEFAULT_SECTIONS } from './defaultSections'
import { byStyle, resolveSectionStyle, SectionStyle } from './sectionStyles'

type SectionDefinition = {
  hasContent(resume: Resume): boolean
  render(resume: Resume, style: SectionStyle): ReactNode
  inNav: boolean
  isCollapsible: boolean
  isFullWidth?: boolean
}

export const SECTION_DEFINITIONS: Record<SectionId, SectionDefinition> = {
  [SectionId.About]: {
    isCollapsible: false,
    inNav: true,
    hasContent: () => true,
    render: (resume, style) => {
      const props = { profile: resume.profile ?? EMPTY_PROFILE, languages: resume.languages }

      return byStyle(style, {
        first: <AboutSection {...props} />,
        second: <AboutRich {...props} />,
        third: <AboutCompact {...props} />,
      })
    },
  },

  [SectionId.TechStack]: {
    isCollapsible: false,
    inNav: false,
    isFullWidth: true,
    hasContent: (resume) => !!resume.technologies?.length,
    render: (resume, style) =>
      byStyle(style, {
        first: <TechStackSlider technologies={resume.technologies} />,
        second: <TechStackRich technologies={resume.technologies} />,
        third: <TechStackCompact technologies={resume.technologies} />,
      }),
  },

  [SectionId.Experience]: {
    isCollapsible: false,
    inNav: true,
    hasContent: (resume) => !!resume.experiences?.length,
    render: (resume, style) =>
      byStyle(style, {
        first: <ExperienceSection experience={resume.experiences} />,
        second: <ExperienceRich experience={resume.experiences} />,
        third: <ExperienceCompact experience={resume.experiences} />,
      }),
  },

  [SectionId.Skills]: {
    isCollapsible: true,
    inNav: true,
    hasContent: ({ skills = EMPTY_SKILLS }) =>
      !!skills.hardSkills.length || !!skills.softSkills.length,
    render: ({ skills = EMPTY_SKILLS }, style) => (
      <SkillsSection
        hardSkills={skills.hardSkills}
        softSkills={skills.softSkills}
        variant={style}
      />
    ),
  },

  [SectionId.Education]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (resume) => !!resume.educations?.length,
    render: (resume, style) =>
      byStyle(style, {
        first: <EducationSection education={resume.educations} />,
        second: <EducationRich education={resume.educations} />,
        third: <EducationCompact education={resume.educations} />,
      }),
  },

  [SectionId.Courses]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (resume) => !!resume.courses?.length,
    render: (resume, style) =>
      byStyle(style, {
        first: <CoursesSection courses={resume.courses} />,
        second: <CoursesRich courses={resume.courses} />,
        third: <CoursesCompact courses={resume.courses} />,
      }),
  },

  [SectionId.Projects]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (resume) => !!resume.projects?.length,
    render: (resume, style) =>
      byStyle(style, {
        first: <ProjectsSection projects={resume.projects} />,
        second: <ProjectsRich projects={resume.projects} />,
        third: <ProjectsCompact projects={resume.projects} />,
      }),
  },

  [SectionId.Activities]: {
    isCollapsible: true,
    inNav: true,
    hasContent: (resume) => !!resume.activities?.length,
    render: (resume, style) =>
      byStyle(style, {
        first: <ActivitiesSection activities={resume.activities} />,
        second: <ActivitiesRich activities={resume.activities} />,
        third: <ActivitiesCompact activities={resume.activities} />,
      }),
  },

  [SectionId.Contacts]: {
    isCollapsible: false,
    inNav: true,
    hasContent: ({ contacts = EMPTY_CONTACTS }) =>
      !!(contacts.phone || contacts.telegram || contacts.email || contacts.hiddenContacts.length),
    render: ({ contacts = EMPTY_CONTACTS }, style) => {
      const props = { contacts, text: contacts.contactsText, hidden: contacts.hiddenContacts }

      return byStyle(style, {
        first: <ContactsSection {...props} />,
        second: <ContactsRich {...props} />,
        third: <ContactsCompact {...props} />,
      })
    },
  },
}

export const getVisibleSections = (resume: Resume): SectionSettings[] =>
  withPinnedAbout(
    resume.overview.sections.length ? resume.overview.sections : DEFAULT_SECTIONS,
  ).filter(({ id, enabled }) => enabled && SECTION_DEFINITIONS[id]?.hasContent(resume))

/** Вид секции в этом резюме: из настроек секции, а если его нет или он незнакомый — основной */
export const getSectionStyle = (resume: Resume, id: SectionId) =>
  resolveSectionStyle(resume.overview?.sections?.find((section) => section.id === id)?.style)
