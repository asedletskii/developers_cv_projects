import { useEffect, useState } from 'react'
import { useTechnologyFilter } from '@features/curriculumVitae/hooks'
import { Person } from '@features/curriculumVitae/types'
import { hasTechnology } from '@features/curriculumVitae/utils'

type Params = Pick<Person, 'experience' | 'projects'>

export const useFilterBar = ({ experience, projects }: Params) => {
  const { selected, clear } = useTechnologyFilter()

  const [shown, setShown] = useState(selected)

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    if (selected) setShown(selected)
  }, [selected])

  const jobs = shown ? experience.filter((job) => hasTechnology(job.stack, shown)).length : 0
  const matchedProjects = shown
    ? projects.filter((project) => hasTechnology(project.stack, shown)).length
    : 0

  return { isVisible: !!selected, shown, jobs, matchedProjects, clear }
}
