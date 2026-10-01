import { lazy } from 'react'
import { ROUTES } from '@shared/constants'
import { Navigate, Route, useLocation, useParams } from 'react-router'

const AuthErrorScreen = lazy(() => import('@features/auth/screens/AuthErrorScreen'))
const CurriculumVitaeScreen = lazy(
  () => import('./screens/CurriculumVitaeScreen/CurriculumVitaeScreen'),
)
const LkProfileScreen = lazy(() => import('./screens/LkProfileScreen/LkProfileScreen'))
const LkResumesScreen = lazy(() => import('./screens/LkResumesScreen/LkResumesScreen'))
const LkResumeEditScreen = lazy(() => import('./screens/LkResumeEditScreen/LkResumeEditScreen'))
const LkResumeScreen = lazy(() => import('./screens/LkResumeScreen/LkResumeScreen'))

const LegacyLkRedirect = () => {
  const { '*': rest } = useParams()
  const { search } = useLocation()
  return <Navigate to={`${ROUTES.lk}${rest ? `/${rest}` : ''}${search}`} replace />
}

export const curriculumVitaeRoutes = {
  pages: (
    <>
      <Route path={ROUTES.cv} element={<Navigate to={ROUTES.home} replace />} />
      <Route path={ROUTES.cvResume} element={<CurriculumVitaeScreen />} />

      <Route path={ROUTES.authError} element={<AuthErrorScreen />} />

      <Route path="/cv/lk/:legacyUserId/*" element={<LegacyLkRedirect />} />

      <Route path={ROUTES.lk} element={<LkProfileScreen />}>
        <Route index element={<LkResumesScreen />} />
        <Route path={ROUTES.lkResume} element={<LkResumeEditScreen />} />
      </Route>

      <Route path={ROUTES.lkResumePreview} element={<LkResumeScreen />} />
    </>
  ),
}
