import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, Resume } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { createAsyncThunk } from '@reduxjs/toolkit'

type GetPublicResumeArg = {
  cvId: string

  /** Предпросмотр в ЛК: тот же экран, что у публичной страницы, но данные — владельческой
   *  ручкой (`GET /lk/resume/:id`), а не публичной. Публичная отвечает `404`, пока резюме не
   *  `opened` — владелец должен видеть черновик и до публикации. */
  userId?: string
}

/** Резюме для страницы: публичная проекция (`GET /cv/:id`) или, в предпросмотре ЛК, владельческая
 *  (`GET /lk/resume/:id`) — см. `userId` у `GetPublicResumeArg`. */
export const getPublicResume = createAsyncThunk<
  Resume,
  GetPublicResumeArg,
  { rejectValue: CurriculumVitaeError }
>('getPublicResume', async ({ cvId, userId }, { rejectWithValue, signal }) => {
  try {
    return userId
      ? await curriculumVitaeApi.getResume(userId, cvId, signal)
      : await curriculumVitaeApi.getPublicResume(cvId, signal)
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})
