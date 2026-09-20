import { useEffect } from 'react'
import { useExpandableJobs } from '@features/curriculumVitae/components/cv/ExperienceSection/hooks/useExpandableJobs'
import { useTimelineAnimation } from '@features/curriculumVitae/components/cv/ExperienceSection/hooks/useTimelineAnimation'
import { Section } from '@features/curriculumVitae/components/cv/Section'
import { useTechnologyFilter } from '@features/curriculumVitae/hooks'
import { Experience, SectionId } from '@features/curriculumVitae/types'
import { hasTechnology } from '@features/curriculumVitae/utils'
import { useTranslate } from '@shared/hooks'
import { VISIBLE_COUNT } from './constants'
import { Job } from './Job'
import { MoreButton } from './MoreButton'
import { Timeline } from './Timeline'

type Props = {
  experience: Experience[]
}

export const ExperienceSection = ({ experience }: Props) => {
  const translate = useTranslate('CurriculumVitae.Experience')
  const timelineRef = useTimelineAnimation()
  const { isExpanded, toggle, expand, buttonRef } = useExpandableJobs(timelineRef)
  const { selected } = useTechnologyFilter()

  useEffect(() => {
    if (
      selected &&
      experience.slice(VISIBLE_COUNT).some((job) => hasTechnology(job.stack, selected))
    ) {
      expand()
    }
  }, [selected, experience, expand])

  if (!experience.length) return null

  const visibleJobs = isExpanded ? experience : experience.slice(0, VISIBLE_COUNT)
  const hiddenCount = experience.length - VISIBLE_COUNT

  return (
    <Section id={SectionId.Experience} title={translate('title')} collapsible={false}>
      <Timeline ref={timelineRef}>
        {visibleJobs.map((job, index) => (
          <Job key={job.id} job={job} isExtra={index >= VISIBLE_COUNT} />
        ))}
      </Timeline>

      {hiddenCount > 0 && (
        <MoreButton
          ref={buttonRef}
          isExpanded={isExpanded}
          hiddenCount={hiddenCount}
          onClick={toggle}
        />
      )}
    </Section>
  )
}
