import { Person } from '@features/curriculumVitae/types'
import { lkApiMock } from '@features/lk/api'
import { LkError, LkProfile, LkResumeSummary } from '@features/lk/types'
import { parseLkError } from '@features/lk/utils'
import { createAsyncThunk } from '@reduxjs/toolkit'

// Бэкенда для ЛК пока нет (см. docs/lk-api.md) — работаем на моке, но по тому же контракту,
// что будет у реального API, чтобы замена на lkApi.getLkProfile/getLkResumes была однострочной.
// Мок отменяемость не поддерживает — signal подключим вместе с реальным API.

export const getLkProfile = createAsyncThunk<LkProfile, string, { rejectValue: LkError }>(
  'getLkProfile',
  async (userId, { rejectWithValue }) => {
    try {
      return await lkApiMock.getLkProfile(userId)
    } catch (e) {
      return rejectWithValue(parseLkError(e))
    }
  },
)

export const getLkResumes = createAsyncThunk<LkResumeSummary[], string, { rejectValue: LkError }>(
  'getLkResumes',
  async (userId, { rejectWithValue }) => {
    try {
      return await lkApiMock.getLkResumes(userId)
    } catch (e) {
      return rejectWithValue(parseLkError(e))
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
  { rejectValue: LkError }
>('updateLkResumeTitle', async ({ userId, resumeId, title }, { rejectWithValue }) => {
  try {
    return await lkApiMock.updateLkResumeTitle(userId, resumeId, title)
  } catch (e) {
    return rejectWithValue(parseLkError(e))
  }
})

type CreateLkResumeArg = {
  userId: string
  person: Person
}

export const createLkResume = createAsyncThunk<Person, CreateLkResumeArg, { rejectValue: LkError }>(
  'createLkResume',
  async ({ userId, person }, { rejectWithValue }) => {
    try {
      return await lkApiMock.createLkResume(userId, person)
    } catch (e) {
      return rejectWithValue(parseLkError(e))
    }
  },
)
