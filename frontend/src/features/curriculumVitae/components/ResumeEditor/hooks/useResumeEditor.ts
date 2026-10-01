import { useMemo } from 'react'
import { Resume, ResumeEditorTab } from '@features/curriculumVitae/types'
import { getResumeFormValues } from '@features/curriculumVitae/utils'
import { useTranslate } from '@shared/hooks'
import { useSearchParams } from 'react-router'
import { EDITOR_TABS } from '../constants'

type Params = {
  view: Resume
  isDirtyRef: { current: boolean }
}

export const useResumeEditor = ({ view, isDirtyRef }: Params) => {
  const translate = useTranslate('Lk.Edit')
  const [searchParams, setSearchParams] = useSearchParams()

  const values = useMemo(() => getResumeFormValues(view), [view])
  const activeTab = EDITOR_TABS.find(({ id }) => id === searchParams.get('tab')) ?? EDITOR_TABS[0]

  const changeTab = (id: ResumeEditorTab) => {
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

  return { values, activeTab, changeTab, setDirty }
}
