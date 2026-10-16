import { MouseEvent, useState } from 'react'
import { useResumeLinks } from '@features/curriculumVitae/hooks'
import { ResumeLink } from '@features/curriculumVitae/types'
import { getResumeLinkUrl } from '@features/curriculumVitae/utils'
import { useCopyToClipboard, useTranslate } from '@shared/hooks'

type Dialog = { mode: 'create' } | { mode: 'rename'; link: ResumeLink }
type Menu = { anchor: HTMLElement; link: ResumeLink }

export const useLinksTab = (linksState: ReturnType<typeof useResumeLinks>) => {
  const translate = useTranslate('Lk.Links')
  const { links, isLoading, hasLoadError, create, rename, remove, changeStatus } = linksState
  const { copiedKey, copy } = useCopyToClipboard()

  const [dialog, setDialog] = useState<Dialog | null>(null)
  const [menu, setMenu] = useState<Menu | null>(null)

  const openCreate = () => setDialog({ mode: 'create' })
  const closeDialog = () => setDialog(null)

  const saveName = (name: string) => {
    if (dialog?.mode === 'create') create(name)
    if (dialog?.mode === 'rename') rename(dialog.link.id, name)
  }

  const openMenu = (link: ResumeLink) => (e: MouseEvent<HTMLElement>) =>
    setMenu({ anchor: e.currentTarget, link })
  const closeMenu = () => setMenu(null)

  const renameFromMenu = () => {
    if (menu) setDialog({ mode: 'rename', link: menu.link })
    setMenu(null)
  }

  const removeFromMenu = () => {
    if (menu && window.confirm(translate('deleteConfirm', { name: menu.link.name }))) {
      remove(menu.link.id)
    }
    setMenu(null)
  }

  const toggleStatusFromMenu = () => {
    if (menu) {
      const { id, version, status } = menu.link
      changeStatus(id, version, status === 'active' ? 'inactive' : 'active')
    }
    setMenu(null)
  }

  return {
    links,
    isLoading,
    hasLoadError,
    getLinkUrl: (link: ResumeLink) => getResumeLinkUrl(link.id),
    copiedKey,
    copy,
    dialog,
    openCreate,
    closeDialog,
    saveName,
    menu,
    openMenu,
    closeMenu,
    renameFromMenu,
    removeFromMenu,
    toggleStatusFromMenu,
  }
}
