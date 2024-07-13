import {useAuth} from "../../hooks/useAuth.tsx";
import {Navigate, Outlet} from "react-router-dom";

const PrivateRoute = () => {
  const { isAuthenticated , authenticatedUser} = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  if (authenticatedUser && !authenticatedUser.name) {
    return <Navigate to="/onboarding/profile" />;
  }

  return <Outlet />;
};

export default PrivateRoute;