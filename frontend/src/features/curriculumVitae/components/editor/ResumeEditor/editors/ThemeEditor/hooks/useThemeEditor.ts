import { useState } from 'react'
import { useTheme } from '@mui/material/styles'

/** Состояние блока «Цветовая гамма»: цвета темы по умолчанию и открыт ли предпросмотр. */
export const useThemeEditor = () => {
  const { palette } = useTheme()
  const [isPreviewOpen, setIsPreviewOpen] = useState(false)

  return {
    defaultPrimary: palette.primary.main,
    defaultSecondary: palette.secondary.main,
    isPreviewOpen,
    openPreview: () => setIsPreviewOpen(true),
    closePreview: () => setIsPreviewOpen(false),
  }
}
