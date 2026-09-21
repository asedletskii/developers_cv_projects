import { PropsWithChildren } from 'react'
import { SnackbarProvider } from 'notistack'
import { Host } from './Host'

/** Подключает notistack и выводит в нём тосты из стора (см. `showToast`). */
export const ToastProvider = ({ children }: PropsWithChildren) => (
  <SnackbarProvider
    maxSnack={3}
    anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
    dense
  >
    {children}

    <Host />
  </SnackbarProvider>
)
