import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import {
  CurriculumVitaeError,
  LkProfile,
  LkResumeSummary,
  ResumeLanguage,
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
  Pick<LkResumeSummary, 'id' | 'title'>,
  UpdateLkResumeTitleArg,
  ThunkConfig
>('updateLkResumeTitle', async ({ userId, resumeId, title }, { rejectWithValue }) => {
  try {
    await curriculumVitaeApi.updateTitle(userId, resumeId, title)
    return { id: resumeId, title }
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})

type CreateLkResumeArg = {
  userId: string
  title: string
  language: ResumeLanguage
}

/** Возвращает id созданного резюме — дальше его наполняют секциями в редакторе. */
export const createLkResume = createAsyncThunk<string, CreateLkResumeArg, ThunkConfig>(
  'createLkResume',
  async ({ userId, title, language }, { rejectWithValue }) => {
    try {
      return await curriculumVitaeApi.createResume(userId, title, language)
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)

type DeleteLkResumeArg = {
  userId: string
  resumeId: string
}

export const deleteLkResume = createAsyncThunk<string, DeleteLkResumeArg, ThunkConfig>(
  'deleteLkResume',
  async ({ userId, resumeId }, { rejectWithValue }) => {
    try {
      await curriculumVitaeApi.deleteResume(userId, resumeId)
      return resumeId
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)
