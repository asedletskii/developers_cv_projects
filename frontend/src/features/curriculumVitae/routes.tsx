import { lazy } from 'react'
import { Navigate, Route } from 'react-router'

const CurriculumVitaeScreen = lazy(
  () => import('./screens/CurriculumVitaeScreen/CurriculumVitaeScreen'),
)
const LkProfileScreen = lazy(() => import('./screens/LkProfileScreen/LkProfileScreen'))
const LkResumeEditScreen = lazy(() => import('./screens/LkResumeEditScreen/LkResumeEditScreen'))
const LkResumeScreen = lazy(() => import('./screens/LkResumeScreen/LkResumeScreen'))

// Смонтировано под `<Route path="cv">` из PageRouter. Статические сегменты ("lk", "new")
// react-router ставит выше динамических (:cvId), поэтому /cv/lk/... не спорит с /cv/:cvId.
export const curriculumVitaeRoutes = {
  pages: (
    <>
      {/* Без id показывать нечего — уходим на главную */}
      <Route index element={<Navigate to="/" replace />} />
      <Route path=":cvId" element={<CurriculumVitaeScreen />} />

      <Route path="lk">
        <Route index element={<Navigate to="/" replace />} />

        <Route path=":userId">
          <Route index element={<LkProfileScreen />} />
          <Route path="new" element={<LkResumeEditScreen />} />
          {/* Клик по резюме открывает редактирование; просмотр как есть — отдельная кнопка */}
          <Route path=":cvId" element={<LkResumeEditScreen />} />
          <Route path=":cvId/preview" element={<LkResumeScreen />} />
        </Route>
      </Route>
    </>
  ),
}
