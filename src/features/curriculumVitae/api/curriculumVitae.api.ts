import { Person, Project } from '@features/curriculumVitae/types'
import { api } from '@shared/api'
import { Languages } from '@shared/types'

class CurriculumVitaeApi {
  baseUrl = '/cv'

  async getPersonData(cvId: string, lang: Languages): Promise<Person> {
    const url = `${this.baseUrl}/${encodeURIComponent(cvId)}`

    const response = await api.get<Person>(url, { params: { lang } })
    return response.data
  }

  async getProject(projectId: string): Promise<Project> {
    const url = `${this.baseUrl}/projectId=${projectId}`

    const response = await api.get<Project>(url)
    return response.data
  }
}

export const curriculumVitaeApi = new CurriculumVitaeApi()
