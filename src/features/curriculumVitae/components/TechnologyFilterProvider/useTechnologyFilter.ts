import { useContext } from 'react'
import { TechnologyFilterContext } from './TechnologyFilterProvider'

export const useTechnologyFilter = () => useContext(TechnologyFilterContext)
