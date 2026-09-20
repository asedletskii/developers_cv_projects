import { useContext } from 'react'
import { TechnologyFilterContext } from '@features/curriculumVitae/components/cv/TechnologyFilterProvider'

export const useTechnologyFilter = () => useContext(TechnologyFilterContext)
