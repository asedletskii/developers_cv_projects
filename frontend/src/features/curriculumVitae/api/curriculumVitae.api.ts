import {
  ListSection,
  LkProfile,
  LkResumeSummary,
  Person,
  ResumeLanguage,
  ResumeView,
  SingletonSection,
} from '@features/curriculumVitae/types'
import { normalizePerson } from '@features/curriculumVitae/utils'
import { API_ROUTES } from '@features/curriculumVitae/constants'
import { api } from '@shared/api'
import { generatePath } from 'react-router'

const ownerHeaders = (userId: string) => ({ 'X-Mock-User-ID': userId })
class CurriculumVitaeApi {
  async getPublicResume(resumeId: string, signal?: AbortSignal): Promise<Person> {
    const response = await api.get<Person>(generatePath(API_ROUTES.publicResume, { resumeId }), {
      signal,
    })
    return normalizePerson(response.data)
  }

  async getProfile(userId: string, signal?: AbortSignal): Promise<LkProfile> {
    const response = await api.get<LkProfile>(API_ROUTES.lkProfile, {
      signal,
      headers: ownerHeaders(userId),
    })
    return response.data
  }

  async getResumes(userId: string, signal?: AbortSignal): Promise<LkResumeSummary[]> {
    const response = await api.get<LkResumeSummary[]>(API_ROUTES.lkResumes, {
      signal,
      headers: ownerHeaders(userId),
    })
    return response.data
  }

  async getResume(userId: string, resumeId: string, signal?: AbortSignal): Promise<ResumeView> {
    const response = await api.get<ResumeView>(generatePath(API_ROUTES.lkResume, { resumeId }), {
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

  async updateTitle(userId: string, resumeId: string, title: string): Promise<void> {
    await api.put(
      generatePath(API_ROUTES.lkResume, { resumeId }),
      { title },
      { headers: ownerHeaders(userId) },
    )
  }

  async deleteResume(userId: string, resumeId: string): Promise<void> {
    await api.delete(generatePath(API_ROUTES.lkResume, { resumeId }), {
      headers: ownerHeaders(userId),
    })
  }

  async saveSingleton(
    userId: string,
    resumeId: string,
    section: SingletonSection,
    body: unknown,
    exists: boolean,
  ): Promise<void> {
    const url = generatePath(API_ROUTES.lkSection, { resumeId, section })
    const config = { headers: ownerHeaders(userId) }

    if (exists || section === 'overview') await api.put(url, body, config)
    else await api.post(url, body, config)
  }

  async createItem(
    userId: string,
    resumeId: string,
    section: ListSection,
    body: unknown,
  ): Promise<string> {
    const response = await api.post<{ id: string }>(
      generatePath(API_ROUTES.lkItems, { resumeId, section }),
      body,
      {
        headers: ownerHeaders(userId),
      },
    )
    return response.data.id
  }

  async reorderItems(userId: string, resumeId: string, section: ListSection, ids: string[]) {
    await api.put(
      generatePath(API_ROUTES.lkItemsOrder, { resumeId, section }),
      { ids },
      { headers: ownerHeaders(userId) },
    )
  }

  async updateItem(
    userId: string,
    resumeId: string,
    section: ListSection,
    itemId: string,
    body: unknown,
  ) {
    await api.put(generatePath(API_ROUTES.lkItem, { resumeId, section, itemId }), body, {
      headers: ownerHeaders(userId),
    })
  }

  async deleteItem(userId: string, resumeId: string, section: ListSection, itemId: string) {
    await api.delete(generatePath(API_ROUTES.lkItem, { resumeId, section, itemId }), {
      headers: ownerHeaders(userId),
    })
  }
}

export const curriculumVitaeApi = new CurriculumVitaeApi()
