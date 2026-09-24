// import { UndefinedScreen } from '@shared/screens/UndefinedScreen'
import { curriculumVitaeRoutes } from '@features/curriculumVitae/routes'
import { Route, Routes, useLocation } from 'react-router'

export const PageRouter = () => {
  const location = useLocation()
  const state = location.state as { backgroundLocation?: Location }

  return (
    <Routes location={state?.backgroundLocation || location}>
      <Route path="/">
        <Route path="cv">{curriculumVitaeRoutes.pages}</Route>
      </Route>

      {/* <Route path="*" element={<UndefinedScreen />} /> */}
    </Routes>
  )
}
