import { ElementType } from 'react'
import {
  CodeOutlined,
  EmojiEventsOutlined,
  ForumOutlined,
  MicOutlined,
  SupervisorAccountOutlined,
  VolunteerActivismOutlined,
} from '@mui/icons-material'
import { ActivityType } from '@features/curriculumVitae/types'

export const ACTIVITY_ICONS: Record<ActivityType, ElementType> = {
  [ActivityType.OpenSource]: CodeOutlined,
  [ActivityType.Hackathon]: EmojiEventsOutlined,
  [ActivityType.Speaker]: MicOutlined,
  [ActivityType.Volunteering]: VolunteerActivismOutlined,
  [ActivityType.Mentoring]: SupervisorAccountOutlined,
  [ActivityType.Community]: ForumOutlined,
}
