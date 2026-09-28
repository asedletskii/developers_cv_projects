import { useState } from 'react'
import { changeLkResumeStatus, deleteLkResume } from '@features/curriculumVitae/store/lk'
import { LkResumeSummary } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'
import { useAppDispatch } from '@shared/store'
import { ToastType } from '@shared/types'
import { showToast } from '@shared/utils'

type Menu = { anchor: HTMLElement; resume: LkResumeSummary }

/** Логика списка резюме: меню действий, переименование, создание и удаление. */
export const useResumeList = (userId: string) => {
  const dispatch = useAppDispatch()
  const translate = useTranslate('Lk.Resumes')

  const [menu, setMenu] = useState<Menu | null>(null)
  const [editingResume, setEditingResume] = useState<LkResumeSummary | null>(null)
  const [isCreating, setIsCreating] = useState(false)

  const closeMenu = () => setMenu(null)

  const renameFromMenu = () => {
    if (menu) setEditingResume(menu.resume)
    closeMenu()
  }

  const removeFromMenu = () => {
    if (menu && window.confirm(translate('deleteConfirm', { title: menu.resume.title }))) {
      dispatch(deleteLkResume({ userId, resumeId: menu.resume.id }))
    }
    closeMenu()
  }

  const changeStatusFromMenu = async (status: 'opened' | 'closed') => {
    const { resume } = menu ?? {}
    closeMenu()
    if (!resume) return

    try {
      await dispatch(
        changeLkResumeStatus({ userId, resumeId: resume.id, version: resume.version, status }),
      ).unwrap()
    } catch {
      showToast({ messageKey: 'statusChangeFailed', type: ToastType.Error })
    }
  }

  return {
    menu,
    openMenu: (anchor: HTMLElement, resume: LkResumeSummary) => setMenu({ anchor, resume }),
    closeMenu,
    renameFromMenu,
    publishFromMenu: () => changeStatusFromMenu('opened'),
    unpublishFromMenu: () => changeStatusFromMenu('closed'),
    removeFromMenu,
    editingResume,
    closeEditing: () => setEditingResume(null),
    isCreating,
    openCreating: () => setIsCreating(true),
    closeCreating: () => setIsCreating(false),
  }
}
