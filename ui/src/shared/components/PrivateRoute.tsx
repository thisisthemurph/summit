import {useAuth} from "../../hooks/useAuth.tsx";
import {Navigate, Outlet} from "react-router-dom";

const PrivateRoute = () => {
  const { isAuthenticated } = useAuth()
  return isAuthenticated
    ? <Outlet />
    : <Navigate to="/login" />
};

export default PrivateRoute;