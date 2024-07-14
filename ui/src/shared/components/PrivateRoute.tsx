import {useAuth} from "../../hooks/useAuth.tsx";
import {Navigate, Outlet} from "react-router-dom";
import {AuthenticatedUser} from "../types/responseTypes.ts";

/**
 * A generic private route.
 * Users must be authenticated and have completed the onboarding process.
 * If the onboarding process is incomplete, they are navigated to the onboarding/profile page.
 */
const PrivateRoute = () => {
  const { isAuthenticated , authenticatedUser} = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  if (!profileIsComplete(authenticatedUser)) {
    return <Navigate to="/onboarding/profile" />;
  }

  return <Outlet />;
};

function profileIsComplete(user: AuthenticatedUser|null): boolean {
  if (!user) {
    return false;
  }
  if (!user.firstName && !user.lastName) {
    return false;
  }
  return true;
}

export default PrivateRoute;