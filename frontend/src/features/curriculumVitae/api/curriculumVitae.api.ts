import { API_ROUTES } from '@features/curriculumVitae/constants'
import {
  AnalyticsBatch,
  ClickBreakpoint,
  ClickGrid,
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
  SectionsStatistics,
  SingletonSection,
  TechnologyInterest,
} from '@features/curriculumVitae/types'
import { api } from '@shared/api'
import { generatePath } from 'react-router'

type AnalyticsParams = { from: string; to: string; tz: string }

class CurriculumVitaeApi {
  /** `viewId` — идентификатор визита из заголовка `X-View-Id`; без него (ошибка записи на бэке)
   *  трекер выключен */
  async getPublicResume(
    resumeId: string,
    signal?: AbortSignal,
  ): Promise<{ resume: Resume; viewId: string | null }> {
    const response = await api.get<Resume>(generatePath(API_ROUTES.publicResume, { resumeId }), {
      signal,
    })
    return { resume: response.data, viewId: response.headers['x-view-id'] ?? null }
  }

  async getProfile(userId: string, signal?: AbortSignal): Promise<LkProfile> {
    const response = await api.get<LkProfile>(API_ROUTES.lkProfile, { signal })
    return response.data
  }

  async updateProfile(
    userId: string,
    body: { ru: LkLocalizedProfile; en: LkLocalizedProfile; photoPath?: string },
  ): Promise<void> {
    await api.put(API_ROUTES.lkProfile, body)
  }

  async getResumes(userId: string, signal?: AbortSignal): Promise<LkResumeSummary[]> {
    const response = await api.get<LkResumeSummary[]>(API_ROUTES.lkResumes, { signal })
    return response.data
  }

  async getResume(userId: string, resumeId: string, signal?: AbortSignal): Promise<Resume> {
    const response = await api.get<Resume>(generatePath(API_ROUTES.lkResume, { resumeId }), {
      signal,
    })
    return response.data
  }

  async createResume(userId: string, title: string, language: ResumeLanguage): Promise<string> {
    const response = await api.post<{ id: string }>(API_ROUTES.lkCreateResume, { title, language })
    return response.data.id
  }

  async updateTitle(
    userId: string,
    resumeId: string,
    version: number,
    title: string,
  ): Promise<void> {
    await api.put(generatePath(API_ROUTES.lkResume, { resumeId }), { version, title })
  }

  async changeStatus(
    userId: string,
    resumeId: string,
    version: number,
    status: ResumeStatus,
  ): Promise<void> {
    await api.put(generatePath(API_ROUTES.lkResumeStatus, { resumeId }), { version, status })
  }

  async deleteResume(userId: string, resumeId: string): Promise<void> {
    await api.delete(generatePath(API_ROUTES.lkResume, { resumeId }))
  }

  async replaceSection(
    userId: string,
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
  ): Promise<string> {
    const response = await api.post<{ contact: ContactKey; value: string }>(
      generatePath(API_ROUTES.contactViews, { resumeId }),
      { contact, viewId: viewId ?? undefined },
    )
    return response.data.value
  }

  async getResumeLinks(
    userId: string,
    resumeId: string,
    signal?: AbortSignal,
  ): Promise<ResumeLink[]> {
    const response = await api.get<ResumeLink[]>(API_ROUTES.lkLinksList, {
      signal,
      params: { resumeId },
    })
    return response.data
  }

  async createResumeLink(userId: string, resumeId: string, name: string): Promise<string> {
    const response = await api.post<{ id: string }>(API_ROUTES.lkLinks, { resumeId, name })
    return response.data.id
  }

  async renameResumeLink(userId: string, linkId: string, name: string): Promise<void> {
    await api.put(generatePath(API_ROUTES.lkLink, { linkId }), { name })
  }

  async deleteResumeLink(userId: string, linkId: string): Promise<void> {
    await api.delete(generatePath(API_ROUTES.lkLink, { linkId }))
  }

  async changeLinkStatus(
    userId: string,
    linkId: string,
    version: number,
    status: LinkStatus,
  ): Promise<void> {
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
    return response.data
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
