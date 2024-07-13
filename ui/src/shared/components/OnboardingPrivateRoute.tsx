import {useAuth} from "../../hooks/useAuth.tsx";
import {Navigate, Outlet} from "react-router-dom";

const OnboardingPrivateRoute = () => {
  const { isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  return <Outlet />;
};

export default OnboardingPrivateRoute;