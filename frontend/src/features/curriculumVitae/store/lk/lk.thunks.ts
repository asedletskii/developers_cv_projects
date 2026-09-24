import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import {
  CurriculumVitaeError,
  LkProfile,
  LkResumeSummary,
  Person,
} from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { createAsyncThunk } from '@reduxjs/toolkit'

type ThunkConfig = { rejectValue: CurriculumVitaeError }

export const getLkProfile = createAsyncThunk<LkProfile, string, ThunkConfig>(
  'getLkProfile',
  async (userId, { rejectWithValue, signal }) => {
    try {
      return await curriculumVitaeApi.getProfile(userId, signal)
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)

export const getLkResumes = createAsyncThunk<LkResumeSummary[], string, ThunkConfig>(
  'getLkResumes',
  async (userId, { rejectWithValue, signal }) => {
    try {
      return await curriculumVitaeApi.getResumes(userId, signal)
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)

type UpdateLkResumeTitleArg = {
  userId: string
  resumeId: string
  title: string
}

export const updateLkResumeTitle = createAsyncThunk<
  LkResumeSummary,
  UpdateLkResumeTitleArg,
  ThunkConfig
>('updateLkResumeTitle', async ({ userId, resumeId, title }, { rejectWithValue }) => {
  try {
    return await curriculumVitaeApi.renameResume(userId, resumeId, title)
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})

type CreateLkResumeArg = {
  userId: string
  person: Person
}

export const createLkResume = createAsyncThunk<Person, CreateLkResumeArg, ThunkConfig>(
  'createLkResume',
  async ({ userId, person }, { rejectWithValue }) => {
    try {
      return await curriculumVitaeApi.createResume(userId, person)
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)
