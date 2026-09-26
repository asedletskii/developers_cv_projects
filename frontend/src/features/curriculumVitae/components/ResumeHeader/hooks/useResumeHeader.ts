import { useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { generatePath, useNavigate } from 'react-router'

type Params = {
  userId: string
  resumeId: string

  /** Название изменено — экран перечитывает резюме */
  onRenamed: () => void
  onError: (error: CurriculumVitaeError) => void
}

/** Действия над резюме в шапке: меню, переименование и удаление с переходом к списку. */
export const useResumeHeader = ({ userId, resumeId, onRenamed, onError }: Params) => {
  const translate = useTranslate('Lk.Edit')
  const navigate = useNavigate()

  const [menuAnchor, setMenuAnchor] = useState<HTMLElement | null>(null)
  const [isRenaming, setIsRenaming] = useState(false)

  const closeMenu = () => setMenuAnchor(null)

  const startRename = () => {
    closeMenu()
    setIsRenaming(true)
  }

  const finishRename = () => {
    setIsRenaming(false)
    onRenamed()
  }

  const remove = async () => {
    closeMenu()
    if (!window.confirm(translate('deleteConfirm'))) return

    try {
      await curriculumVitaeApi.deleteResume(userId, resumeId)
      navigate(generatePath(ROUTES.lkProfile, { userId }))
    } catch (e) {
      onError(parseApiError(e))
    }
  }

  return { menuAnchor, setMenuAnchor, closeMenu, isRenaming, startRename, finishRename, remove }
}
