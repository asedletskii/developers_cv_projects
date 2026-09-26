import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, Person } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { createAsyncThunk } from '@reduxjs/toolkit'

type GetPersonDataArg = { cvId: string }

/** Публичная проекция резюме (`GET /cv/:id`) — её же показывает предпросмотр в ЛК. */
export const getPersonData = createAsyncThunk<
  Person,
  GetPersonDataArg,
  { rejectValue: CurriculumVitaeError }
>('getPersonData', async ({ cvId }, { rejectWithValue, signal }) => {
  try {
    return await curriculumVitaeApi.getPublicResume(cvId, signal)
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})
