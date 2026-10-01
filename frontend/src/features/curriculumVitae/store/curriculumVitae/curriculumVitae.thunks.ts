import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, Resume } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { createAsyncThunk } from '@reduxjs/toolkit'

type GetPublicResumeArg = {
  cvId: string
  userId?: string
}
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
