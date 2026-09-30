import { useOutletContext } from 'react-router'

export type LkOutletContext = { userId: string }

/** userId для вложенных в `LkProfileScreen` экранов — передаётся через `<Outlet context>`. */
export const useLkOutletContext = () => useOutletContext<LkOutletContext>()
