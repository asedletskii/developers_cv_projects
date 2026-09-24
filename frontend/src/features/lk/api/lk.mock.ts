import { curriculumVitaeApiMock } from '@features/curriculumVitae/api'
import { Person } from '@features/curriculumVitae/types'
import { LkProfile, LkResumeSummary } from '@features/lk/types'

const delay = (ms = 300) => new Promise((resolve) => setTimeout(resolve, ms))

/** Ошибка той же формы, что у axios при 404 — включая `isAxiosError`, иначе parseLkError её не узнает */
const notFoundError = (message: string) =>
  Object.assign(new Error(message), { isAxiosError: true, response: { status: 404 } })

// userId = тот же id, что у резюме в curriculumVitae.mock.ts ('1' и '2') — отдельной сущности
// «пользователь» пока нет, поэтому не городим свои id, а используем те же самые.
const profiles: Record<string, LkProfile> = {
  '1': {
    userId: '1',
    name: 'Антон',
    surname: 'Седлецкий',
    patronymic: 'Александрович',
    position: 'Frontend',
    location: 'Новосибирск',
    photoPath: '/mock/avatar.svg',
  },

  '2': {
    userId: '2',
    name: 'Дмитрий',
    surname: 'Костерин',
    position: 'Backend',
    location: 'Томск',
    photoPath: '/mock/avatar.svg',
  },
}

const resumes: Record<string, LkResumeSummary[]> = {
  '1': [{ id: '1', title: 'Frontend-разработчик' }],
  // Два резюме у одного пользователя — чтобы список в ЛК было на чём проверить
  '2': [
    { id: '2', title: 'Backend-разработчик' },
    { id: '1', title: 'Frontend-разработчик (резерв)' },
  ],
}

class LkApiMock {
  baseUrl = '/lk'

  async getLkProfile(userId: string): Promise<LkProfile> {
    await delay()

    const profile = profiles[userId]
    if (!profile) throw notFoundError('LK profile not found')

    return profile
  }

  async getLkResumes(userId: string): Promise<LkResumeSummary[]> {
    await delay()

    const list = resumes[userId]
    if (!list) throw notFoundError('LK resumes not found')

    // Клон, а не сама запись стора: Redux (Immer) замораживает то, что попадает в state,
    // а стор мока должен оставаться мутируемым для updateLkResumeTitle
    return list.map((resume) => ({ ...resume }))
  }

  async updateLkResumeTitle(
    userId: string,
    resumeId: string,
    title: string,
  ): Promise<LkResumeSummary> {
    await delay()

    const resume = resumes[userId]?.find((r) => r.id === resumeId)
    if (!resume) throw notFoundError('LK resume not found')

    resume.title = title
    return { ...resume }
  }

  /**
   * Создаёт резюме и сразу привязывает его к пользователю. Сама сущность резюме (`Person`) живёт
   * в сторе curriculumVitae — переиспользуем его мок, а не заводим для резюме второй источник правды.
   */
  async createLkResume(userId: string, person: Person): Promise<Person> {
    await delay()

    if (!profiles[userId]) throw notFoundError('LK profile not found')

    const created = await curriculumVitaeApiMock.createPersonData(person)

    resumes[userId] = [
      ...(resumes[userId] ?? []),
      { id: created.id, title: created.position || created.name },
    ]

    return created
  }
}

export const lkApiMock = new LkApiMock()
