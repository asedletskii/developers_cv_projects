import { lazy } from 'react'
import { ROUTES } from '@shared/constants'
import { Navigate, Route } from 'react-router'

const CurriculumVitaeScreen = lazy(
  () => import('./screens/CurriculumVitaeScreen/CurriculumVitaeScreen'),
)
const LkProfileScreen = lazy(() => import('./screens/LkProfileScreen/LkProfileScreen'))
const LkResumesScreen = lazy(() => import('./screens/LkResumesScreen/LkResumesScreen'))
const LkResumeEditScreen = lazy(() => import('./screens/LkResumeEditScreen/LkResumeEditScreen'))
const LkResumeScreen = lazy(() => import('./screens/LkResumeScreen/LkResumeScreen'))

// Пути абсолютные, из ROUTES. Статический сегмент "lk" react-router ставит выше динамического
// (:cvId), поэтому /cv/lk/... не спорит с /cv/:cvId.
export const curriculumVitaeRoutes = {
  pages: (
    <>
      {/* Без id показывать нечего — уходим на главную */}
      <Route path={ROUTES.cv} element={<Navigate to={ROUTES.home} replace />} />
      <Route path={ROUTES.cvResume} element={<CurriculumVitaeScreen />} />

      <Route path={ROUTES.lk} element={<Navigate to={ROUTES.home} replace />} />

      {/* Слева всегда карточка пользователя, справа — список резюме или редактор */}
      <Route path={ROUTES.lkProfile} element={<LkProfileScreen />}>
        <Route index element={<LkResumesScreen />} />
        <Route path={ROUTES.lkResume} element={<LkResumeEditScreen />} />
      </Route>

      {/* Просмотр как есть — отдельной страницей без каркаса ЛК */}
      <Route path={ROUTES.lkResumePreview} element={<LkResumeScreen />} />
    </>
  ),
}
