export const ROUTES = {
  home: '/',
  main: '/main',
  about: '/about',

  cv: '/cv',
  cvResume: '/cv/:cvId',

  lk: '/cv/lk',
  lkProfile: '/cv/lk/:userId',
  lkResume: '/cv/lk/:userId/:cvId',
  lkResumePreview: '/cv/lk/:userId/:cvId/preview',
} as const
