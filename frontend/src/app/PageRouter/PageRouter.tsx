// import { UndefinedScreen } from '@shared/screens/UndefinedScreen'
import { curriculumVitaeRoutes } from '@features/curriculumVitae/routes'
import { Routes, useLocation } from 'react-router'

export const PageRouter = () => {
  const location = useLocation()
  const state = location.state as { backgroundLocation?: Location }

  return (
    <Routes location={state?.backgroundLocation || location}>
      {curriculumVitaeRoutes.pages}

      {/* <Route path="*" element={<UndefinedScreen />} /> */}
    </Routes>
  )
}
