import { useEffect, useState } from 'react'
import { useTechnologyFilter } from '@features/curriculumVitae/hooks'
import { Resume } from '@features/curriculumVitae/types'
import { hasTechnology } from '@features/curriculumVitae/utils'

type Params = Pick<Resume, 'experiences' | 'projects'>

export const useFilterBar = ({ experiences, projects }: Params) => {
  const { selected, clear } = useTechnologyFilter()

  const [shown, setShown] = useState(selected)

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    if (selected) setShown(selected)
  }, [selected])

  const jobs = shown ? experiences.filter((job) => hasTechnology(job.stack, shown)).length : 0
  const matchedProjects = shown
    ? projects.filter((project) => hasTechnology(project.stack, shown)).length
    : 0

  return { isVisible: !!selected, shown, jobs, matchedProjects, clear }
}
