import { PropsWithChildren } from 'react'
import { SnackbarProvider } from 'notistack'
import { Host } from './Host'

export const ToastProvider = ({ children }: PropsWithChildren) => (
  <SnackbarProvider maxSnack={3} anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }} dense>
    {children}

    <Host />
  </SnackbarProvider>
)
