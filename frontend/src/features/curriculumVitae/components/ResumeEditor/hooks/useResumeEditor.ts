import { useMemo } from 'react'
import { ResumeEditorTab, ResumeView } from '@features/curriculumVitae/types'
import { viewToPerson } from '@features/curriculumVitae/utils'
import { useTranslate } from '@shared/hooks'
import { useSearchParams } from 'react-router'
import { EDITOR_TABS } from '../constants'

type Params = {
  view: ResumeView

  /** Есть ли несохранённые правки — при уходе с вкладки спрашиваем подтверждение */
  isDirtyRef: { current: boolean }
}

/** Состояние редактора: значения формы из резюме и активная вкладка (`?tab=`). */
export const useResumeEditor = ({ view, isDirtyRef }: Params) => {
  const translate = useTranslate('Lk.Edit')
  const [searchParams, setSearchParams] = useSearchParams()

  const person = useMemo(() => viewToPerson(view), [view])
  const activeTab = EDITOR_TABS.find(({ id }) => id === searchParams.get('tab')) ?? EDITOR_TABS[0]

  const changeTab = (id: ResumeEditorTab) => {
    // Форма вкладки пересоздаётся — правки этой вкладки без сохранения пропали бы
    if (isDirtyRef.current && !window.confirm(translate('discardConfirm'))) return

    setSearchParams(
      (params) => {
        params.set('tab', id)
        return params
      },
      { replace: true },
    )
  }

  const setDirty = (isDirty: boolean) => {
    isDirtyRef.current = isDirty
  }

  return { person, activeTab, changeTab, setDirty }
}
