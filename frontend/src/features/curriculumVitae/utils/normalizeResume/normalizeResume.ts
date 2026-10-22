import { Resume } from '@features/curriculumVitae/types'

export const normalizeResume = (raw: Resume): Resume => ({
  ...raw,
  overview: {
    ...raw.overview,
    greeting: raw.overview?.greeting ?? { first: '', second: '', secondHighlight: '', third: '' },
    sections: raw.overview?.sections ?? [],
  },
  profile: raw.profile && { ...raw.profile, workFormat: raw.profile.workFormat ?? [] },
  contacts: raw.contacts && { ...raw.contacts, hiddenContacts: raw.contacts.hiddenContacts ?? [] },
  skills: raw.skills && {
    ...raw.skills,
    hardSkills: raw.skills.hardSkills ?? [],
    softSkills: raw.skills.softSkills ?? [],
  },
  languages: raw.languages ?? [],
  technologies: raw.technologies ?? [],
  experiences: (raw.experiences ?? []).map((job) => ({
    ...job,
    achievements: job.achievements ?? [],
    tasks: job.tasks ?? [],
    stack: job.stack ?? [],
  })),
  educations: raw.educations ?? [],
  courses: raw.courses ?? [],
  activities: raw.activities ?? [],
  projects: (raw.projects ?? []).map((project) => ({
    ...project,
    screenshots: project.screenshots ?? [],
    stack: project.stack ?? [],
    contributors: project.contributors ?? [],
  })),
})
