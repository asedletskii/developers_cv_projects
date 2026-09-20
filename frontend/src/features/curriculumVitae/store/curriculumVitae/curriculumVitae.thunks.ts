import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, Resume } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { createAsyncThunk } from '@reduxjs/toolkit'

type GetPublicResumeArg = {
  cvId: string
  userId?: string
}
type GetPublicResumeResult = { resume: Resume; viewId: string | null }

export const getPublicResume = createAsyncThunk<
  GetPublicResumeResult,
  GetPublicResumeArg,
  { rejectValue: CurriculumVitaeError }
>('getPublicResume', async ({ cvId, userId }, { rejectWithValue, signal }) => {
  try {
    // Предпросмотр владельца визитом не считается: viewId нет, значит, трекер выключен
    return userId
      ? { resume: await curriculumVitaeApi.getResume(cvId, signal), viewId: null }
      : await curriculumVitaeApi.getPublicResume(cvId, signal)
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})
