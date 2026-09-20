import { curriculumVitaeApi, curriculumVitaeApiMock } from '@features/curriculumVitae/api'
import { Project } from '@features/curriculumVitae/types'
import { createAsyncThunk } from '@reduxjs/toolkit'

export const getProject = createAsyncThunk<Project, string, { rejectValue: 'cannotGetProject' }>(
  'getProject',
  async (payload, { rejectWithValue }) => {
    try {
      const response = await curriculumVitaeApiMock.getProject(payload)

      return response
    } catch (e: any) {
      // showToast({ messageKey: 'cannotGetProject', type: 'error' })

      return rejectWithValue('cannotGetProject')
    }
  },
)
