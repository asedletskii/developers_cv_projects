import { authDict } from '@features/auth/locale'
import { aboutPageDict, curriculumVitaeDict, lkDict } from '@features/curriculumVitae/locale'
import { sharedDict } from '@shared/locale/dicts'
import { mergeLocaleDicts } from '@shared/utils/mergeLocaleDicts'

export const appLocaleDict = mergeLocaleDicts(
  curriculumVitaeDict,
  lkDict,
  authDict,
  aboutPageDict,
  sharedDict,
)
