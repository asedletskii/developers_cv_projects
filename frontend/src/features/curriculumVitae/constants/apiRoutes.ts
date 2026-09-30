export const API_ROUTES = {
  publicResume: '/cv/:resumeId',
  lkProfile: '/lk/profile',
  lkResumes: '/lk/resumes',
  lkCreateResume: '/lk/resume',
  lkResume: '/lk/resume/:resumeId',
  lkResumeStatus: '/lk/resume/:resumeId/status',
  lkSection: '/lk/resume/:resumeId/:section',
  contactViews: '/cv/:resumeId/contact-views',
  lkLinksList: '/lk/links',
  lkLinks: '/lk/link',
  lkLink: '/lk/link/:linkId',
  lkLinkStatus: '/lk/link/:linkId/status',
} as const
