import { useCallback, useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { EMPTY_PROFILE } from '@features/curriculumVitae/constants'
import { selectResume } from '@features/curriculumVitae/store/curriculumVitae'
import { AnalyticsEventType, ContactKey } from '@features/curriculumVitae/types'
import {
  buildResumeHtml,
  getContactLinks,
  isHardSkill,
  isSoftSkill,
  ResumeHtmlLabels,
  sanitizeFileName,
} from '@features/curriculumVitae/utils'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { useAppSelector } from '@shared/store'
import { ToastType } from '@shared/types'
import { getPersonFIO, showToast } from '@shared/utils'
import { useParams } from 'react-router'
import { useAnalytics } from './useAnalytics'
import { useTurnstileToken } from './useTurnstileToken'

export const useDownloadResume = () => {
  const translate = useTranslate('CurriculumVitae')
  const resume = useAppSelector(selectResume)
  const { cvId } = useParams()
  const { viewId, track } = useAnalytics()
  const getToken = useTurnstileToken()
  const [isDownloading, setIsDownloading] = useState(false)

  const download = useCallback(async () => {
    if (isDownloading || !cvId) return

    setIsDownloading(true)
    try {
      const links = getContactLinks(resume.contacts ?? {}, resume.contacts?.hiddenContacts)
      const contacts: Partial<Record<ContactKey, string>> = {}
      links.forEach(({ key, value }) => {
        if (value) contacts[key] = value
      })

      // Запрашиваем только то, чего ещё нет: открытые и уже раскрытые значения приходят из данных
      const missing = links.filter(({ value }) => !value).map(({ key }) => key)
      let unresolvedContacts: ContactKey[] = []

      if (missing.length) {
        try {
          const token = await getToken()
          const revealed = await curriculumVitaeApi.revealContactsForDownload(cvId, viewId, token)
          revealed.forEach(({ contact, value }) => {
            contacts[contact] = value
          })
          unresolvedContacts = missing.filter((key) => !contacts[key])
        } catch {
          unresolvedContacts = missing
          showToast({ messageKey: 'downloadWithoutContacts', type: ToastType.Warning })
        }
      }

      // Язык файла — язык резюме, а не интерфейса посетителя
      const lng = resume.language
      const t = (key: Parameters<typeof translate>[0]) => translate(key, { lng })

      const labels: ResumeHtmlLabels = {
        lang: lng,
        sections: {
          technologies: t('HtmlExport.technologies'),
          experience: t('Nav.experience'),
          education: t('Nav.education'),
          courses: t('Nav.courses'),
          projects: t('Nav.projects'),
          activities: t('Nav.activities'),
          skills: t('Nav.skills'),
          languages: t('HtmlExport.languages'),
          contacts: t('Nav.contacts'),
        },
        present: t('Experience.present'),
        tasks: t('Experience.tasks'),
        achievements: t('Experience.achievements'),
        hardSkills: t('Skills.hard'),
        softSkills: t('Skills.soft'),
        hoursSuffix: t('HtmlExport.hoursSuffix'),
        contactsHidden: t('HtmlExport.contactsHidden'),
        generated: translate('HtmlExport.generated', {
          lng,
          date: new Intl.DateTimeFormat(lng, { dateStyle: 'long' }).format(new Date()),
        }),
        actualVersion: t('HtmlExport.actualVersion'),
        skillNames: Object.fromEntries([
          ...(resume.skills?.hardSkills ?? [])
            .filter(isHardSkill)
            .map((code) => [code, translate(`Skills.HardSkills.${code}`, { lng })] as const),
          ...(resume.skills?.softSkills ?? [])
            .filter(isSoftSkill)
            .map((code) => [code, translate(`Skills.SoftSkills.${code}`, { lng })] as const),
        ]),
        languageNames: Object.fromEntries(
          resume.languages.map(
            ({ language }) =>
              [language, translate(`About.Languages.${language}`, { lng })] as const,
          ),
        ),
      }

      const fullName = getPersonFIO(resume.profile ?? EMPTY_PROFILE)
      // Адрес актуальной версии есть только у публичной страницы, у предпросмотра в кабинете его нет
      const { origin, pathname } = window.location
      const sourceUrl = pathname.startsWith(ROUTES.lk) ? undefined : `${origin}${pathname}`

      const html = buildResumeHtml({
        resume,
        fullName,
        contacts,
        unresolvedContacts,
        labels,
        sourceUrl,
      })

      const url = URL.createObjectURL(new Blob([html], { type: 'text/html;charset=utf-8' }))
      const anchor = document.createElement('a')
      anchor.href = url
      anchor.download = `${sanitizeFileName(fullName)} — ${t('HtmlExport.fileSuffix')}.html`
      document.body.appendChild(anchor)
      anchor.click()
      anchor.remove()
      window.setTimeout(() => URL.revokeObjectURL(url), 1000)

      track({ type: AnalyticsEventType.Download })
    } catch {
      showToast({ messageKey: 'downloadFailed', type: ToastType.Error })
    } finally {
      setIsDownloading(false)
    }
  }, [isDownloading, cvId, resume, viewId, getToken, track, translate])

  return { download, isDownloading }
}
