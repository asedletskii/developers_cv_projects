import { Person } from '@features/curriculumVitae/types'
import { LkProfile, LkResumeSummary } from '@features/lk/types'
import { api } from '@shared/api'

/** Пока без авторизации: бэкенд узнаёт пользователя по этому заголовку, а не по id в пути. */
const userHeader = (userId: string) => ({ 'X-Mock-User-ID': userId })

class LkApi {
  baseUrl = '/lk'

  async getLkProfile(userId: string, signal?: AbortSignal): Promise<LkProfile> {
    const response = await api.get<LkProfile>(`${this.baseUrl}/profile`, {
      signal,
      headers: userHeader(userId),
    })
    return response.data
  }

  async getLkResumes(userId: string, signal?: AbortSignal): Promise<LkResumeSummary[]> {
    const response = await api.get<LkResumeSummary[]>(`${this.baseUrl}/resumes`, {
      signal,
      headers: userHeader(userId),
    })
    return response.data
  }

  async updateLkResumeTitle(
    userId: string,
    resumeId: string,
    title: string,
    signal?: AbortSignal,
  ): Promise<LkResumeSummary> {
    const url = `${this.baseUrl}/resumes/${encodeURIComponent(resumeId)}`

    const response = await api.patch<LkResumeSummary>(
      url,
      { title },
      { signal, headers: userHeader(userId) },
    )
    return response.data
  }

  /** `id` в теле игнорируется бэкендом — создаётся новый. */
  async createResume(userId: string, person: Person, signal?: AbortSignal): Promise<Person> {
    const response = await api.post<Person>(`${this.baseUrl}/resumes`, person, {
      signal,
      headers: userHeader(userId),
    })
    return response.data
  }
}

export const lkApi = new LkApi()
