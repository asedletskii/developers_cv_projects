export const ROUTES = {
  home: '/',
  main: '/main',
  about: '/about',

  cv: '/cv',
  cvResume: '/cv/:cvId',

  lk: '/cv/lk',
  lkResume: '/cv/lk/:cvId',
  lkResumePreview: '/cv/lk/:cvId/preview',

  authError: '/cv/auth/error',
} as const
