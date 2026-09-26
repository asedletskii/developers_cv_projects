import { useMemo } from 'react'
import { Resume } from '@features/curriculumVitae/types'
import { ChecklistTab, getResumeChecklist } from '@features/curriculumVitae/utils'

/** Чек-лист заполненности резюме: элементы по порядку, число выполненных и процент. */
export const useResumeChecklist = (view: Resume) =>
  useMemo(() => {
    const items = Object.entries(getResumeChecklist(view)) as [ChecklistTab, boolean][]
    const done = items.filter(([, isDone]) => isDone).length

    return { items, percent: Math.round((done / items.length) * 100) }
  }, [view])
