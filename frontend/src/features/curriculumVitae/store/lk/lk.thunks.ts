import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import {
  CurriculumVitaeError,
  LkLocalizedProfile,
  LkProfile,
  LkResumeSummary,
  ResumeLanguage,
  ResumeStatus,
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

type UpdateLkProfileArg = {
  userId: string
  ru: LkLocalizedProfile
  en: LkLocalizedProfile
  photoPath?: string
}

export const updateLkProfile = createAsyncThunk<LkProfile, UpdateLkProfileArg, ThunkConfig>(
  'updateLkProfile',
  async ({ userId, ru, en, photoPath }, { rejectWithValue }) => {
    try {
      await curriculumVitaeApi.updateProfile(userId, { ru, en, photoPath })
      return { userId, ru, en, photoPath }
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)

type UpdateLkResumeTitleArg = {
  userId: string
  resumeId: string
  version: number
  title: string
}

export const updateLkResumeTitle = createAsyncThunk<
  Pick<LkResumeSummary, 'id' | 'title' | 'version'>,
  UpdateLkResumeTitleArg,
  ThunkConfig
>('updateLkResumeTitle', async ({ userId, resumeId, version, title }, { rejectWithValue }) => {
  try {
    await curriculumVitaeApi.updateTitle(userId, resumeId, version, title)
    return { id: resumeId, title, version: version + 1 }
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})

type ChangeLkResumeStatusArg = {
  userId: string
  resumeId: string
  version: number
  status: ResumeStatus
}

export const changeLkResumeStatus = createAsyncThunk<
  Pick<LkResumeSummary, 'id' | 'status' | 'version'>,
  ChangeLkResumeStatusArg,
  ThunkConfig
>('changeLkResumeStatus', async ({ userId, resumeId, version, status }, { rejectWithValue }) => {
  try {
    await curriculumVitaeApi.changeStatus(userId, resumeId, version, status)
    return { id: resumeId, status, version: version + 1 }
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})

type CreateLkResumeArg = {
  userId: string
  title: string
  language: ResumeLanguage
}

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
