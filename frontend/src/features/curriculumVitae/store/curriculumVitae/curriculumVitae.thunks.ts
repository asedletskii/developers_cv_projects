import { curriculumVitaeApiMock } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, Person } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { createAsyncThunk } from '@reduxjs/toolkit'
import { Languages } from '@shared/types'
// import { showToast } from '@shared/utils'

type GetPersonDataArg = {
  cvId: string
  lang: Languages
}

export const getPersonData = createAsyncThunk<
  Person,
  GetPersonDataArg,
  { rejectValue: CurriculumVitaeError }
>('getPersonData', async ({ cvId, lang }, { rejectWithValue }) => {
  try {
    // На моках, без бэкенда (см. docs/backend-api.md). Чтобы вернуться на реальный API —
    // curriculumVitaeApiMock.getPersonData(cvId, lang) → curriculumVitaeApi.getPersonData(cvId, lang, signal)
    return await curriculumVitaeApiMock.getPersonData(cvId, lang)
  } catch (e) {
    // showToast({ messageKey: 'cannotGetPersonData', type: 'error' })

    return rejectWithValue(parseApiError(e))
  }
})

type UpdatePersonDataArg = {
  cvId: string
  person: Person
}

export const updatePersonData = createAsyncThunk<
  Person,
  UpdatePersonDataArg,
  { rejectValue: CurriculumVitaeError }
>('updatePersonData', async ({ cvId, person }, { rejectWithValue }) => {
  try {
    // На моках, без бэкенда — см. комментарий у getPersonData
    return await curriculumVitaeApiMock.updatePersonData(cvId, person)
  } catch (e) {
    return rejectWithValue(parseApiError(e))
  }
})
