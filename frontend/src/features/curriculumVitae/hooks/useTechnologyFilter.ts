import { useContext } from 'react'
import { TechnologyFilterContext } from '@features/curriculumVitae/components/TechnologyFilterProvider'

export const useTechnologyFilter = () => useContext(TechnologyFilterContext)
