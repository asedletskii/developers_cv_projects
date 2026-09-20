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
      return await curriculumVitaeApi.getProfile(signal)
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)

export const getLkResumes = createAsyncThunk<LkResumeSummary[], string, ThunkConfig>(
  'getLkResumes',
  async (userId, { rejectWithValue, signal }) => {
    try {
      return await curriculumVitaeApi.getResumes(signal)
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
      await curriculumVitaeApi.updateProfile({ ru, en, photoPath })
      return { userId, ru, en, photoPath }
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)

type UpdateLkResumeTitleArg = {
  resumeId: string
  version: number
  title: string
}

export const updateLkResumeTitle = createAsyncThunk<
  Pick<LkResumeSummary, 'id' | 'title' | 'version'>,
  UpdateLkResumeTitleArg,
  ThunkConfig
>('updateLkResumeTitle', async ({ resumeId, version, title }, { rejectWithValue }) => {
  try {
    await curriculumVitaeApi.updateTitle(resumeId, version, title)
    return { id: resumeId, title, version: version + 1 }
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})

type ChangeLkResumeStatusArg = {
  resumeId: string
  version: number
  status: ResumeStatus
}

export const changeLkResumeStatus = createAsyncThunk<
  Pick<LkResumeSummary, 'id' | 'status' | 'version'>,
  ChangeLkResumeStatusArg,
  ThunkConfig
>('changeLkResumeStatus', async ({ resumeId, version, status }, { rejectWithValue }) => {
  try {
    await curriculumVitaeApi.changeStatus(resumeId, version, status)
    return { id: resumeId, status, version: version + 1 }
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})

type CreateLkResumeArg = {
  title: string
  language: ResumeLanguage
}

export const createLkResume = createAsyncThunk<string, CreateLkResumeArg, ThunkConfig>(
  'createLkResume',
  async ({ title, language }, { rejectWithValue }) => {
    try {
      return await curriculumVitaeApi.createResume(title, language)
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)

type DeleteLkResumeArg = {
  resumeId: string
}

export const deleteLkResume = createAsyncThunk<string, DeleteLkResumeArg, ThunkConfig>(
  'deleteLkResume',
  async ({ resumeId }, { rejectWithValue }) => {
    try {
      await curriculumVitaeApi.deleteResume(resumeId)
      return resumeId
    } catch (e) {
      return rejectWithValue(parseApiError(e))
    }
  },
)
