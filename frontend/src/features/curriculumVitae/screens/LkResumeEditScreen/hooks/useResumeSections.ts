import { useRef } from 'react'
import { ResumeEditorTab } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'
import { useSearchParams } from 'react-router'
import { TOP_SECTIONS, TopSection } from '../constants'

/**
 * Активный раздел страницы резюме (`?section=`) и переходы между разделами. Пока в форме
 * редактирования есть несохранённые правки (`isDirtyRef`), уход из раздела требует подтверждения.
 */
export const useResumeSections = () => {
  const translate = useTranslate('Lk.Edit')
  const [searchParams, setSearchParams] = useSearchParams()
  const isDirtyRef = useRef(false)

  const section = TOP_SECTIONS.find((id) => id === searchParams.get('section')) ?? TOP_SECTIONS[0]

  const changeSection = (id: TopSection) => {
    // Форма редактирования пересоздаётся — правки без сохранения пропали бы
    if (isDirtyRef.current && !window.confirm(translate('discardConfirm'))) return

    isDirtyRef.current = false
    setSearchParams({ section: id }, { replace: true })
  }

  const openEditTab = (tab: ResumeEditorTab) =>
    setSearchParams({ section: 'edit', tab }, { replace: true })

  return { section, changeSection, openEditTab, isDirtyRef }
}
