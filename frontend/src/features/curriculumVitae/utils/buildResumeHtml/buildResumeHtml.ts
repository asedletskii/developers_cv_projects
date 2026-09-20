import { ContactKey, Resume } from '@features/curriculumVitae/types'
import { getContactHref } from '@features/curriculumVitae/utils/getContactLinks'
import { safeExternalUrl } from '@features/curriculumVitae/utils/safeExternalUrl'

export type ResumeHtmlLabels = {
  lang: string
  sections: {
    technologies: string
    experience: string
    education: string
    courses: string
    projects: string
    activities: string
    skills: string
    languages: string
    contacts: string
  }
  present: string
  tasks: string
  achievements: string
  hardSkills: string
  softSkills: string
  hoursSuffix: string
  contactsHidden: string
  generated: string
  actualVersion: string
  skillNames: Record<string, string>
  languageNames: Record<string, string>
}

export type BuildResumeHtmlInput = {
  resume: Resume
  fullName: string
  contacts: Partial<Record<ContactKey, string>>
  unresolvedContacts?: ContactKey[]
  labels: ResumeHtmlLabels
  sourceUrl?: string
}

const ESCAPES: Record<string, string> = {
  '&': '&amp;',
  '<': '&lt;',
  '>': '&gt;',
  '"': '&quot;',
  "'": '&#39;',
}

export const escapeHtml = (value: string) => value.replace(/[&<>"']/g, (char) => ESCAPES[char])

const safeHref = (value?: string) => safeExternalUrl(value) ?? ''

const link = (href: string, text: string) =>
  href
    ? `<a href="${escapeHtml(href)}" rel="noopener noreferrer">${escapeHtml(text)}</a>`
    : escapeHtml(text)

const list = (items: string[]) =>
  items.length ? `<ul>${items.map((item) => `<li>${escapeHtml(item)}</li>`).join('')}</ul>` : ''

const chips = (items: string[]) =>
  items.length
    ? `<p class="chips">${items.map((item) => `<span>${escapeHtml(item)}</span>`).join('')}</p>`
    : ''

const section = (title: string, body: string) =>
  body.trim() ? `<section><h2>${escapeHtml(title)}</h2>${body}</section>` : ''

const meta = (parts: (string | undefined | false)[]) => {
  const text = parts.filter(Boolean).join(' · ')
  return text ? `<p class="meta">${escapeHtml(text)}</p>` : ''
}

const contactText = (key: ContactKey, value: string) =>
  key === ContactKey.Telegram ? `@${value.replace(/^@/, '')}` : value

const STYLES = `
  :root { color-scheme: light; }
  body { margin: 0; font: 16px/1.55 system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif; color: #1c1c1e; background: #fff; }
  main { max-width: 760px; margin: 0 auto; padding: 40px 20px 64px; }
  h1 { margin: 0; font-size: 2rem; line-height: 1.2; }
  h2 { margin: 2rem 0 .75rem; padding-bottom: .35rem; border-bottom: 1px solid #e3e3e6; font-size: 1.15rem; }
  h3 { margin: 1.1rem 0 .15rem; font-size: 1rem; }
  p { margin: .35rem 0; }
  ul { margin: .35rem 0; padding-left: 1.2rem; }
  a { color: #2556c4; }
  .position { margin-top: .25rem; font-size: 1.1rem; color: #4a4a4f; }
  .meta { color: #6b6b72; font-size: .9rem; }
  .label { margin-top: .6rem; font-weight: 600; font-size: .9rem; }
  .chips span { display: inline-block; margin: 0 .35rem .35rem 0; padding: .1rem .55rem; border-radius: 999px; background: #f0f0f3; font-size: .85rem; }
  footer { margin-top: 3rem; color: #6b6b72; font-size: .85rem; }
  @media print { main { padding: 0; } a { color: inherit; } }
`

export const buildResumeHtml = ({
  resume,
  fullName,
  contacts,
  unresolvedContacts = [],
  labels,
  sourceUrl,
}: BuildResumeHtmlInput) => {
  const { sections } = labels
  const profile = resume.profile

  const contactItems = ([ContactKey.Telegram, ContactKey.Email, ContactKey.Phone] as const)
    .filter((key) => contacts[key])
    .map((key) => {
      const value = contacts[key] as string
      return `<li>${link(getContactHref(key, value), contactText(key, value))}</li>`
    })
  const contactsBody = contactItems.length
    ? `<ul>${contactItems.join('')}</ul>`
    : unresolvedContacts.length
      ? `<p>${escapeHtml(labels.contactsHidden)}</p>`
      : ''

  const experience = resume.experiences
    .map(
      (job) => `
      <h3>${escapeHtml(job.position)} — ${escapeHtml(job.company)}</h3>
      ${meta([`${job.startDate} — ${job.endDate ?? labels.present}`, job.location])}
      ${job.description ? `<p>${escapeHtml(job.description)}</p>` : ''}
      ${job.tasks.length ? `<p class="label">${escapeHtml(labels.tasks)}</p>${list(job.tasks)}` : ''}
      ${job.achievements.length ? `<p class="label">${escapeHtml(labels.achievements)}</p>${list(job.achievements)}` : ''}
      ${chips(job.stack.map(({ name }) => name))}`,
    )
    .join('')

  const education = resume.educations
    .map(
      (item) => `
      <h3>${escapeHtml(item.university)}</h3>
      <p>${escapeHtml([item.programName, item.profileName].filter(Boolean).join(', '))}</p>
      ${meta([item.faculty, item.location, item.endDate])}`,
    )
    .join('')

  const courses = resume.courses
    .map(
      (course) => `
      <h3>${escapeHtml(course.name)}</h3>
      ${meta([
        course.organization,
        course.endDate,
        course.hours ? `${course.hours} ${labels.hoursSuffix}` : '',
      ])}
      ${course.certificateLink && safeHref(course.certificateLink) ? `<p>${link(safeHref(course.certificateLink), safeHref(course.certificateLink))}</p>` : ''}`,
    )
    .join('')

  const projects = resume.projects
    .map(
      (project) => `
      <h3>${escapeHtml(project.name)}</h3>
      ${meta([project.role])}
      ${project.description ? `<p>${escapeHtml(project.description)}</p>` : ''}
      ${safeHref(project.link) ? `<p>${link(safeHref(project.link), safeHref(project.link))}</p>` : ''}
      ${chips(project.stack.map(({ name }) => name))}`,
    )
    .join('')

  const activities = resume.activities
    .map(
      (activity) => `
      <h3>${escapeHtml(activity.title)}</h3>
      ${meta([activity.date])}
      ${activity.description ? `<p>${escapeHtml(activity.description)}</p>` : ''}
      ${safeHref(activity.link) ? `<p>${link(safeHref(activity.link), safeHref(activity.link))}</p>` : ''}`,
    )
    .join('')

  const hardSkills = (resume.skills?.hardSkills ?? []).map(
    (code) => labels.skillNames[code] ?? code,
  )
  const softSkills = (resume.skills?.softSkills ?? []).map(
    (code) => labels.skillNames[code] ?? code,
  )
  const skills = [
    hardSkills.length
      ? `<p class="label">${escapeHtml(labels.hardSkills)}</p>${chips(hardSkills)}`
      : '',
    softSkills.length
      ? `<p class="label">${escapeHtml(labels.softSkills)}</p>${chips(softSkills)}`
      : '',
  ].join('')

  const languages = chips(
    resume.languages.map(
      ({ language, grade }) =>
        `${labels.languageNames[language] ?? language}${grade ? ` (${grade})` : ''}`,
    ),
  )

  const footer = [
    escapeHtml(labels.generated),
    sourceUrl && safeHref(sourceUrl)
      ? `${escapeHtml(labels.actualVersion)}: ${link(safeHref(sourceUrl), safeHref(sourceUrl))}`
      : '',
  ]
    .filter(Boolean)
    .join(' · ')

  return `<!doctype html>
<html lang="${escapeHtml(labels.lang)}">
<head>
<meta charset="utf-8">
<meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'; img-src data:">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>${escapeHtml(fullName)}</title>
<style>${STYLES}</style>
</head>
<body>
<main>
<header>
<h1>${escapeHtml(fullName)}</h1>
${profile?.position ? `<p class="position">${escapeHtml(profile.position)}</p>` : ''}
${profile?.location ? `<p class="meta">${escapeHtml(profile.location)}</p>` : ''}
</header>
${section(sections.contacts, contactsBody)}
${section(sections.technologies, chips(resume.technologies.map(({ name }) => name)))}
${section(sections.experience, experience)}
${section(sections.education, education)}
${section(sections.courses, courses)}
${section(sections.projects, projects)}
${section(sections.activities, activities)}
${section(sections.skills, skills)}
${section(sections.languages, languages)}
<footer>${footer}</footer>
</main>
</body>
</html>
`
}

const RESERVED_FILE_NAMES = /^(con|prn|aux|nul|com\d|lpt\d)(\..*)?$/i

export const sanitizeFileName = (name: string, fallback = 'resume') => {
  const cleaned = [...name]
    .map((char) => (char.charCodeAt(0) < 32 || '\\/:*?"<>|'.includes(char) ? ' ' : char))
    .join('')
    .replace(/\s+/g, ' ')
    .replace(/[. ]+$/, '')
    .slice(0, 120)
    .trim()

  return !cleaned || RESERVED_FILE_NAMES.test(cleaned) ? fallback : cleaned
}
