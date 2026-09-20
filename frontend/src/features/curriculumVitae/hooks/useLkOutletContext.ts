import { useOutletContext } from 'react-router'

export type LkOutletContext = { userId: string }

export const useLkOutletContext = () => useOutletContext<LkOutletContext>()
