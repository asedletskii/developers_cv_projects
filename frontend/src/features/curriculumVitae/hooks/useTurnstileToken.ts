import { useContext } from 'react'
import { TurnstileContext } from '@features/curriculumVitae/components/cv/TurnstileProvider'

export const useTurnstileToken = () => useContext(TurnstileContext).getToken
