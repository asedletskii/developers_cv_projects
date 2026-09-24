import { lazy } from 'react'
import { Navigate, Route } from 'react-router'

const CurriculumVitaeScreen = lazy(
  () => import('./screens/CurriculumVitaeScreen/CurriculumVitaeScreen'),
)

// Смонтировано под общим `<Route path="cv">` из PageRouter — сюда добавляются только
// дочерние маршруты, без своей обёртки, чтобы не спорить с другими фичами за префикс /cv
export const curriculumVitaeRoutes = {
  pages: (
    <>
      {/* Без id показывать нечего — уходим на главную */}
      <Route index element={<Navigate to="/" replace />} />
      <Route path=":cvId" element={<CurriculumVitaeScreen />} />
    </>
  ),
}
