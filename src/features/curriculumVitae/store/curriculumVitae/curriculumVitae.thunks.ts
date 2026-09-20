import { curriculumVitaeApiMock } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, Person } from '@features/curriculumVitae/types'
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
    return await curriculumVitaeApiMock.getPersonData(cvId, lang)
  } catch (e: any) {
    // showToast({ messageKey: 'cannotGetPersonData', type: 'error' })

    return rejectWithValue(e?.response?.status === 404 ? 'notFound' : 'cannotGetPersonData')
  }
})
