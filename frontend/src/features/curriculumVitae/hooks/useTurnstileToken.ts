import { useContext } from 'react'
import { TurnstileContext } from '@features/curriculumVitae/components/cv/TurnstileProvider'

/** `getToken` общего виджета Turnstile страницы (см. `TurnstileProvider`) */
export const useTurnstileToken = () => useContext(TurnstileContext).getToken
