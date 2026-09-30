import { API_ROUTES } from '@features/curriculumVitae/constants'
import {
  ContactKey,
  LinkStatus,
  ListSection,
  LkLocalizedProfile,
  LkProfile,
  LkResumeSummary,
  Resume,
  ResumeLanguage,
  ResumeLink,
  ResumeStatus,
  SingletonSection,
} from '@features/curriculumVitae/types'
import { api } from '@shared/api'
import { generatePath } from 'react-router'

class CurriculumVitaeApi {
  async getPublicResume(resumeId: string, signal?: AbortSignal): Promise<Resume> {
    const response = await api.get<Resume>(generatePath(API_ROUTES.publicResume, { resumeId }), {
      signal,
    })
    return response.data
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

  async revealContact(resumeId: string, contact: ContactKey): Promise<string> {
    const response = await api.post<{ contact: ContactKey; value: string }>(
      generatePath(API_ROUTES.contactViews, { resumeId }),
      { contact },
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
}

export const curriculumVitaeApi = new CurriculumVitaeApi()
