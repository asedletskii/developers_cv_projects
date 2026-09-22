import { curriculumVitaeApi } from '@features/curriculumVitae/api'
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
>('getPersonData', async ({ cvId, lang }, { rejectWithValue, signal }) => {
  try {
    // signal приходит из thunkAPI: отмена диспатча (см. .abort() в компоненте) реально обрывает
    // HTTP-запрос, а не просто перестаёт ждать его в браузере
    return await curriculumVitaeApi.getPersonData(cvId, lang, signal)
  } catch (e) {
    // showToast({ messageKey: 'cannotGetPersonData', type: 'error' })

    return rejectWithValue(parseApiError(e))
  }
})
