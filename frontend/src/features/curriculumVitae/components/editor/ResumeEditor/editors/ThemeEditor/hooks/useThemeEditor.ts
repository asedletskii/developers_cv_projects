import { useState } from 'react'
import { useTheme } from '@mui/material/styles'

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
