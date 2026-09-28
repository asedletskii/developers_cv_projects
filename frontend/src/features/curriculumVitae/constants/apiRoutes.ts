export const API_ROUTES = {
  publicResume: '/cv/:resumeId',
  lkProfile: '/lk/profile',
  lkResumes: '/lk/resumes',
  lkCreateResume: '/lk/resume',
  lkResume: '/lk/resume/:resumeId',
  lkResumeStatus: '/lk/resume/:resumeId/status',
  /** Любая секция (синглтон или список) сохраняется целиком одним `PUT` на этот адрес */
  lkSection: '/lk/resume/:resumeId/:section',
  /** Значение скрытого контакта — бэк отдаёт его только по явному запросу, не в самом резюме */
  contactViews: '/cv/:resumeId/contact-views',
} as const
