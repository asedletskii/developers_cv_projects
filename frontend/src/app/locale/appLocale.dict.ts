import { curriculumVitaeDict } from '@features/curriculumVitae/locale'
import { lkDict } from '@features/lk/locale'
import { sharedDict } from '@shared/locale/dicts'
import { mergeLocaleDicts } from '@shared/utils/mergeLocaleDicts'

export const appLocaleDict = mergeLocaleDicts(curriculumVitaeDict, lkDict, sharedDict)
