import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, Resume } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { createAsyncThunk } from '@reduxjs/toolkit'

type GetPublicResumeArg = { cvId: string }

/** Публичная проекция резюме (`GET /cv/:id`) — её же показывает предпросмотр в ЛК. */
export const getPublicResume = createAsyncThunk<
  Resume,
  GetPublicResumeArg,
  { rejectValue: CurriculumVitaeError }
>('getPublicResume', async ({ cvId }, { rejectWithValue, signal }) => {
  try {
    return await curriculumVitaeApi.getPublicResume(cvId, signal)
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})
