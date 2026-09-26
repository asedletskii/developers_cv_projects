import { KeyboardEvent, useState } from 'react'

type Params = {
  initialName?: string
  onSave: (name: string) => void
  onClose: () => void
}

/** Состояние диалога с именем ссылки: одно и для создания, и для переименования. */
export const useLinkNameDialog = ({ initialName = '', onSave, onClose }: Params) => {
  const [name, setName] = useState(initialName)

  const canSave = !!name.trim()

  const save = () => {
    if (!canSave) return

    onSave(name.trim())
    onClose()
  }

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'Enter') save()
  }

  return { name, setName, canSave, save, handleKeyDown }
}
