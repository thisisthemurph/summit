import {useAuth} from "../../hooks/useAuth.tsx";
import {Navigate, Outlet} from "react-router-dom";

/**
 * A private route for the onboarding pages.
 * The user must be authenticated, but they do not have to have their profile completed.
 */
const OnboardingPrivateRoute = () => {
  const { isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  return <Outlet />;
};

export default OnboardingPrivateRoute;