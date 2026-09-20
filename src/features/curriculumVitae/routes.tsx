import { lazy } from 'react'
import { Navigate, Route } from 'react-router'

const CurriculumVitaeScreen = lazy(
  () => import('./screens/CurriculumVitaeScreen/CurriculumVitaeScreen'),
)

export const curriculumVitaeRoutes = {
  pages: (
    <Route path="cv">
      {/* Без id показывать нечего — уходим на главную */}
      <Route index element={<Navigate to="/" replace />} />
      <Route path=":cvId" element={<CurriculumVitaeScreen />} />
    </Route>
  ),
}
