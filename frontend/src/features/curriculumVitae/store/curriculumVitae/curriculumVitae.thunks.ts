import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, Person } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { createAsyncThunk } from '@reduxjs/toolkit'

type GetPersonDataArg = {
  cvId: string

  /** Владелец — тогда резюме берётся из личного кабинета (`/lk/resume/:id`), иначе публичное */
  userId?: string
}

export const getPersonData = createAsyncThunk<
  Person,
  GetPersonDataArg,
  { rejectValue: CurriculumVitaeError }
>('getPersonData', async ({ cvId, userId }, { rejectWithValue, signal }) => {
  try {
    return userId
      ? await curriculumVitaeApi.getResume(userId, cvId, signal)
      : await curriculumVitaeApi.getPublicResume(cvId, signal)
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})

type UpdatePersonDataArg = {
  userId: string
  cvId: string
  person: Person
}

export const updatePersonData = createAsyncThunk<
  Person,
  UpdatePersonDataArg,
  { rejectValue: CurriculumVitaeError }
>('updatePersonData', async ({ userId, cvId, person }, { rejectWithValue }) => {
  try {
    return await curriculumVitaeApi.updateResume(userId, cvId, person)
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})
