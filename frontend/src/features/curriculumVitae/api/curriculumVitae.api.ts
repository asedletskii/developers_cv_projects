import { API_ROUTES } from '@features/curriculumVitae/constants'
import {
  AnalyticsBatch,
  ClickBreakpoint,
  ClickGrid,
  ContactInitiator,
  ContactKey,
  DailyPoint,
  LinkStatistics,
  LinkStatus,
  LinkVisitsPage,
  ListSection,
  LkLocalizedProfile,
  LkProfile,
  LkResumeSummary,
  Resume,
  ResumeLanguage,
  ResumeLink,
  ResumeStatistics,
  ResumeStatus,
  RevealedContact,
  SectionsStatistics,
  SingletonSection,
  TechnologyInterest,
} from '@features/curriculumVitae/types'
import { api } from '@shared/api'
import { generatePath } from 'react-router'

type AnalyticsParams = { from: string; to: string; tz: string }

// Referer у запроса резюме — адрес самого SPA, поэтому реальный источник визита передаём отдельно
const getReferrerHost = (): string | undefined => {
  try {
    const host = document.referrer ? new URL(document.referrer).hostname : ''
    return host && host !== window.location.hostname ? host : undefined
  } catch {
    return undefined
  }
}

class CurriculumVitaeApi {
  async getPublicResume(
    resumeId: string,
    signal?: AbortSignal,
  ): Promise<{ resume: Resume; viewId: string | null }> {
    const response = await api.get<Resume>(generatePath(API_ROUTES.publicResume, { resumeId }), {
      signal,
      headers: { 'X-Referrer-Host': getReferrerHost() },
    })
    return { resume: response.data, viewId: response.headers['x-view-id'] ?? null }
  }

  async getProfile(signal?: AbortSignal): Promise<LkProfile> {
    const response = await api.get<LkProfile>(API_ROUTES.lkProfile, { signal })
    return response.data
  }

  async updateProfile(body: {
    ru: LkLocalizedProfile
    en: LkLocalizedProfile
    photoPath?: string
  }): Promise<void> {
    await api.put(API_ROUTES.lkProfile, body)
  }

  async getResumes(signal?: AbortSignal): Promise<LkResumeSummary[]> {
    const response = await api.get<LkResumeSummary[]>(API_ROUTES.lkResumes, { signal })
    return response.data
  }

  async getResume(resumeId: string, signal?: AbortSignal): Promise<Resume> {
    const response = await api.get<Resume>(generatePath(API_ROUTES.lkResume, { resumeId }), {
      signal,
    })
    return response.data
  }

  async createResume(title: string, language: ResumeLanguage): Promise<string> {
    const response = await api.post<{ id: string }>(API_ROUTES.lkCreateResume, { title, language })
    return response.data.id
  }

  async updateTitle(resumeId: string, version: number, title: string): Promise<void> {
    await api.put(generatePath(API_ROUTES.lkResume, { resumeId }), { version, title })
  }

  async changeStatus(resumeId: string, version: number, status: ResumeStatus): Promise<void> {
    await api.put(generatePath(API_ROUTES.lkResumeStatus, { resumeId }), { version, status })
  }

  async deleteResume(resumeId: string): Promise<void> {
    await api.delete(generatePath(API_ROUTES.lkResume, { resumeId }))
  }

  async replaceSection(
    resumeId: string,
    section: SingletonSection | ListSection,
    body: unknown,
  ): Promise<void> {
    await api.put(generatePath(API_ROUTES.lkSection, { resumeId, section }), body)
  }

  async revealContact(
    resumeId: string,
    contact: ContactKey,
    viewId?: string | null,
    turnstileToken?: string,
  ): Promise<string> {
    const response = await api.post<{ contact: ContactKey; value: string }>(
      generatePath(API_ROUTES.contactViews, { resumeId }),
      {
        initiator: ContactInitiator.Contact,
        contact,
        viewId: viewId ?? undefined,
        turnstileToken,
      },
    )
    return response.data.value
  }

  /** Все скрытые контакты разом — для сборки файла по кнопке «Скачать». Тот же эндпоинт, что и
   *  `revealContact`, но с `initiator: download`: один запрос и один токен Turnstile на всю пачку,
   *  а в аналитике это отдельное действие, а не серия раскрытий */
  async revealContactsForDownload(
    resumeId: string,
    viewId?: string | null,
    turnstileToken?: string,
  ): Promise<RevealedContact[]> {
    const response = await api.post<{ contacts: RevealedContact[] }>(
      generatePath(API_ROUTES.contactViews, { resumeId }),
      { initiator: ContactInitiator.Download, viewId: viewId ?? undefined, turnstileToken },
    )
    return response.data.contacts
  }

  async getResumeLinks(resumeId: string, signal?: AbortSignal): Promise<ResumeLink[]> {
    const response = await api.get<ResumeLink[]>(API_ROUTES.lkLinksList, {
      signal,
      params: { resumeId },
    })
    return response.data
  }

  async createResumeLink(resumeId: string, name: string): Promise<string> {
    const response = await api.post<{ id: string }>(API_ROUTES.lkLinks, { resumeId, name })
    return response.data.id
  }

  async renameResumeLink(linkId: string, name: string): Promise<void> {
    await api.put(generatePath(API_ROUTES.lkLink, { linkId }), { name })
  }

  async deleteResumeLink(linkId: string): Promise<void> {
    await api.delete(generatePath(API_ROUTES.lkLink, { linkId }))
  }

  async changeLinkStatus(linkId: string, version: number, status: LinkStatus): Promise<void> {
    await api.put(generatePath(API_ROUTES.lkLinkStatus, { linkId }), { version, status })
  }

  async sendAnalyticsBatch(linkId: string, batch: AnalyticsBatch): Promise<void> {
    await api.post(generatePath(API_ROUTES.cvEvents, { linkId }), batch)
  }

  async getLinkStatistics(
    linkId: string,
    params: AnalyticsParams,
    signal?: AbortSignal,
  ): Promise<LinkStatistics> {
    const response = await api.get<LinkStatistics>(
      generatePath(API_ROUTES.lkLinkStatistics, { linkId }),
      { params, signal },
    )
    return response.data
  }

  async getLinkDaily(
    linkId: string,
    params: AnalyticsParams,
    signal?: AbortSignal,
  ): Promise<DailyPoint[]> {
    const response = await api.get<{ days: DailyPoint[] }>(
      generatePath(API_ROUTES.lkLinkDaily, { linkId }),
      { params, signal },
    )
    return response.data.days
  }

  async getLinkVisits(
    linkId: string,
    cursor: string | null | undefined,
    includeFiltered: boolean,
    signal?: AbortSignal,
  ): Promise<LinkVisitsPage> {
    const response = await api.get<LinkVisitsPage>(
      generatePath(API_ROUTES.lkLinkVisits, { linkId }),
      {
        params: {
          cursor: cursor ?? undefined,
          limit: 20,
          includeFiltered: includeFiltered || undefined,
        },
        signal,
      },
    )
    // Go отдаёт пустой слайс как null: приводим к массиву, чтобы дальше не проверять
    return {
      ...response.data,
      items: response.data.items.map((visit) => ({ ...visit, events: visit.events ?? [] })),
    }
  }

  async getResumeStatistics(
    resumeId: string,
    params: AnalyticsParams,
    signal?: AbortSignal,
  ): Promise<ResumeStatistics> {
    const response = await api.get<ResumeStatistics>(
      generatePath(API_ROUTES.lkResumeStatistics, { resumeId }),
      { params, signal },
    )
    return response.data
  }

  async getLinkSections(
    linkId: string,
    params: AnalyticsParams,
    signal?: AbortSignal,
  ): Promise<SectionsStatistics> {
    const response = await api.get<SectionsStatistics>(
      generatePath(API_ROUTES.lkLinkSections, { linkId }),
      { params, signal },
    )
    return response.data
  }

  async getLinkTechnologies(
    linkId: string,
    params: AnalyticsParams,
    signal?: AbortSignal,
  ): Promise<TechnologyInterest[]> {
    const response = await api.get<{ technologies: TechnologyInterest[] }>(
      generatePath(API_ROUTES.lkLinkTechnologies, { linkId }),
      { params, signal },
    )
    return response.data.technologies
  }

  async getLinkClicks(
    linkId: string,
    section: string,
    breakpoint: ClickBreakpoint,
    includeStale: boolean,
    signal?: AbortSignal,
  ): Promise<ClickGrid> {
    const response = await api.get<ClickGrid>(generatePath(API_ROUTES.lkLinkClicks, { linkId }), {
      params: { section, breakpoint, includeStale: includeStale || undefined },
      signal,
    })
    return response.data
  }
}

export const curriculumVitaeApi = new CurriculumVitaeApi()
