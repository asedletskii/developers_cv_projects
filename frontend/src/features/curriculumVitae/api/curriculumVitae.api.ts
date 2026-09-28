import { API_ROUTES } from '@features/curriculumVitae/constants'
import {
  ContactKey,
  ListSection,
  LkLocalizedProfile,
  LkProfile,
  LkResumeSummary,
  Resume,
  ResumeLanguage,
  ResumeStatus,
  SingletonSection,
} from '@features/curriculumVitae/types'
import { api } from '@shared/api'
import { generatePath } from 'react-router'

const ownerHeaders = (userId: string) => ({ 'X-Mock-User-ID': userId })
class CurriculumVitaeApi {
  async getPublicResume(resumeId: string, signal?: AbortSignal): Promise<Resume> {
    const response = await api.get<Resume>(generatePath(API_ROUTES.publicResume, { resumeId }), {
      signal,
    })
    return response.data
  }

  async getProfile(userId: string, signal?: AbortSignal): Promise<LkProfile> {
    const response = await api.get<LkProfile>(API_ROUTES.lkProfile, {
      signal,
      headers: ownerHeaders(userId),
    })
    return response.data
  }

  async updateProfile(
    userId: string,
    body: { ru: LkLocalizedProfile; en: LkLocalizedProfile; photoPath?: string },
  ): Promise<void> {
    await api.put(API_ROUTES.lkProfile, body, { headers: ownerHeaders(userId) })
  }

  async getResumes(userId: string, signal?: AbortSignal): Promise<LkResumeSummary[]> {
    const response = await api.get<LkResumeSummary[]>(API_ROUTES.lkResumes, {
      signal,
      headers: ownerHeaders(userId),
    })
    return response.data
  }

  async getResume(userId: string, resumeId: string, signal?: AbortSignal): Promise<Resume> {
    const response = await api.get<Resume>(generatePath(API_ROUTES.lkResume, { resumeId }), {
      signal,
      headers: ownerHeaders(userId),
    })
    return response.data
  }

  async createResume(userId: string, title: string, language: ResumeLanguage): Promise<string> {
    const response = await api.post<{ id: string }>(
      API_ROUTES.lkCreateResume,
      { title, language },
      { headers: ownerHeaders(userId) },
    )
    return response.data.id
  }

  async updateTitle(
    userId: string,
    resumeId: string,
    version: number,
    title: string,
  ): Promise<void> {
    await api.put(
      generatePath(API_ROUTES.lkResume, { resumeId }),
      { version, title },
      { headers: ownerHeaders(userId) },
    )
  }

  async changeStatus(
    userId: string,
    resumeId: string,
    version: number,
    status: ResumeStatus,
  ): Promise<void> {
    await api.put(
      generatePath(API_ROUTES.lkResumeStatus, { resumeId }),
      { version, status },
      { headers: ownerHeaders(userId) },
    )
  }

  async deleteResume(userId: string, resumeId: string): Promise<void> {
    await api.delete(generatePath(API_ROUTES.lkResume, { resumeId }), {
      headers: ownerHeaders(userId),
    })
  }

  async replaceSection(
    userId: string,
    resumeId: string,
    section: SingletonSection | ListSection,
    body: unknown,
  ): Promise<void> {
    await api.put(generatePath(API_ROUTES.lkSection, { resumeId, section }), body, {
      headers: ownerHeaders(userId),
    })
  }

  async revealContact(resumeId: string, contact: ContactKey): Promise<string> {
    const response = await api.post<{ contact: ContactKey; value: string }>(
      generatePath(API_ROUTES.contactViews, { resumeId }),
      { contact },
    )
    return response.data.value
  }
}

export const curriculumVitaeApi = new CurriculumVitaeApi()
