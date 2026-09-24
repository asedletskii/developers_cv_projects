import { LkProfile, LkResumeSummary, Person } from '@features/curriculumVitae/types'
import { normalizePerson } from '@features/curriculumVitae/utils'
import { api } from '@shared/api'

/** Пока без авторизации: бэкенд узнаёт владельца `/lk/*` по этому заголовку. */
const ownerHeaders = (userId: string) => ({ 'X-Mock-User-ID': userId })

const resumeUrl = (resumeId: string) => `/lk/resume/${encodeURIComponent(resumeId)}`

/** Пути ручек — как на бэкенде (`backend/internal/platform/http/router.go`). */
class CurriculumVitaeApi {
  /** Публичное резюме — без владельца */
  async getPublicResume(resumeId: string, signal?: AbortSignal): Promise<Person> {
    const response = await api.get<Person>(`/cv/${encodeURIComponent(resumeId)}`, { signal })
    return normalizePerson(response.data)
  }

  async getProfile(userId: string, signal?: AbortSignal): Promise<LkProfile> {
    const response = await api.get<LkProfile>('/lk/profile', {
      signal,
      headers: ownerHeaders(userId),
    })
    return response.data
  }

  async getResumes(userId: string, signal?: AbortSignal): Promise<LkResumeSummary[]> {
    const response = await api.get<LkResumeSummary[]>('/lk/resumes', {
      signal,
      headers: ownerHeaders(userId),
    })
    return response.data
  }

  /** Резюме владельца (бэк отдаёт 404, если оно чужое) */
  async getResume(userId: string, resumeId: string, signal?: AbortSignal): Promise<Person> {
    const response = await api.get<Person>(resumeUrl(resumeId), {
      signal,
      headers: ownerHeaders(userId),
    })
    return normalizePerson(response.data)
  }

  /** `id` в теле бэк игнорирует и генерирует новый */
  async createResume(userId: string, person: Person): Promise<Person> {
    const response = await api.post<Person>('/lk/resume', person, {
      headers: ownerHeaders(userId),
    })
    return normalizePerson(response.data)
  }

  async updateResume(userId: string, resumeId: string, person: Person): Promise<Person> {
    const response = await api.put<Person>(resumeUrl(resumeId), person, {
      headers: ownerHeaders(userId),
    })
    return normalizePerson(response.data)
  }

  async renameResume(userId: string, resumeId: string, title: string): Promise<LkResumeSummary> {
    const response = await api.patch<LkResumeSummary>(
      resumeUrl(resumeId),
      { title },
      { headers: ownerHeaders(userId) },
    )
    return response.data
  }
}

export const curriculumVitaeApi = new CurriculumVitaeApi()
