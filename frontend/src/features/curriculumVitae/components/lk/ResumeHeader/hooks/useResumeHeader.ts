import { useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, ResumeStatus } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { useNavigate } from 'react-router'

type Params = {
  resumeId: string
  version?: number
  onRenamed: () => void
  onError: (error: CurriculumVitaeError) => void
}

export const useResumeHeader = ({ resumeId, version, onRenamed, onError }: Params) => {
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

  const changeStatus = async (status: Extract<ResumeStatus, 'opened' | 'closed'>) => {
    closeMenu()

    try {
      await curriculumVitaeApi.changeStatus(resumeId, version ?? 0, status)
      onRenamed()
    } catch (e) {
      onError(parseApiError(e))
    }
  }

  const remove = async () => {
    closeMenu()
    if (!window.confirm(translate('deleteConfirm'))) return

    try {
      await curriculumVitaeApi.deleteResume(resumeId)
      navigate(ROUTES.lk)
    } catch (e) {
      onError(parseApiError(e))
    }
  }

  return {
    menuAnchor,
    setMenuAnchor,
    closeMenu,
    isRenaming,
    startRename,
    finishRename,
    changeStatus,
    remove,
  }
}
