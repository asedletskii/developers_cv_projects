export const API_ROUTES = {
  publicResume: '/cv/:resumeId',
  lkProfile: '/lk/profile',
  lkResumes: '/lk/resumes',
  lkCreateResume: '/lk/resume',
  lkResume: '/lk/resume/:resumeId',
  lkSection: '/lk/resume/:resumeId/:section',
  lkItems: '/lk/resume/:resumeId/:section',
  lkItem: '/lk/resume/:resumeId/:section/:itemId',
  lkItemsOrder: '/lk/resume/:resumeId/:section/order',
} as const
