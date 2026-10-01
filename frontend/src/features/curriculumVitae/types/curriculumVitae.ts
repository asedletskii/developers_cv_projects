import { ActivityType, AvailabilityStatus, EducationLevel, SectionId } from './enums'

export type ContactsText = {
  title?: string
  subtitle?: string
}

export type Availability = {
  status: AvailabilityStatus
  from?: string
  note?: string
}

export type ThemeSettings = {
  primary?: string
  secondary?: string
}

export type SectionSettings = {
  id: SectionId
  enabled: boolean
  collapsed?: boolean
}

export type Greeting = {
  first: string
  second: string
  secondHighlight: string
  third: string
}

export type Technology = {
  name: string
  grade?: string
}

export type Contacts = {
  phone?: string
  telegram?: string
  email?: string
}

export type Experience = {
  id: string
  company: string
  position: string
  location?: string
  startDate: string
  endDate?: string
  description?: string
  achievements: string[]
  tasks: string[]
  stack: Technology[]
}

export type Course = {
  id: string
  name: string
  organization: string
  endDate: string
  hours?: number
  certificateLink?: string
}

export type Activity = {
  id: string
  type: ActivityType
  title: string
  description?: string
  date?: string
  link?: string
}

export type Education = {
  id: string
  location: string
  university: string
  grade: EducationLevel
  programName: string
  profileName?: string
  faculty: string
  course?: number
  endDate?: string
}

export type Project = {
  id: string
  name: string
  description?: string
  role?: string
  screenshots: string[]
  link?: string
  gitLink?: string
  stack: Technology[]
  contributors: SimplifiedPerson[]
}

type SimplifiedPerson = {
  name: string
  surname: string
  contacts: Contacts
  role: string
}
