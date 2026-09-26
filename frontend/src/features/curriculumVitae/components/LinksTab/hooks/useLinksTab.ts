import { MouseEvent, useState } from 'react'
import { getResumeLinkUrl } from '@features/curriculumVitae/utils'
import { useCopyToClipboard, useTranslate } from '@shared/hooks'
import { ResumeLink, useResumeLinks } from './useResumeLinks'

type Dialog = { mode: 'create' } | { mode: 'rename'; link: ResumeLink }
type Menu = { anchor: HTMLElement; link: ResumeLink }

/** Логика вкладки «Ссылки»: список ссылок, диалог имени, меню действий и копирование. */
export const useLinksTab = (resumeId: string) => {
  const translate = useTranslate('Lk.Links')
  const { links, create, rename, remove } = useResumeLinks(resumeId)
  const { copiedKey, copy } = useCopyToClipboard()

  const [dialog, setDialog] = useState<Dialog | null>(null)
  const [menu, setMenu] = useState<Menu | null>(null)

  const publicUrl = getResumeLinkUrl(resumeId)

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

  return {
    links,
    publicUrl,
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
  }
}
