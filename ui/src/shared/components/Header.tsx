import {Link} from "react-router-dom";
import {useAuth} from "../../hooks/useAuth.tsx";
import axiosInstance from "../requests/axiosInstance.ts";

function Header() {
  const {isAuthenticated, logoutUser} = useAuth();

  return (
    <section className="flex justify-between p-4 bg-base-100">
      <h1 className="text-3xl">Summit</h1>
      <Nav isAuthenticated={isAuthenticated} logoutUser={logoutUser} />
    </section>
  )
}

type NavProps = {
  isAuthenticated: boolean;
  logoutUser: () => void;
}

function Nav({isAuthenticated, logoutUser}: NavProps) {
  const handleLogOut = async () => {
    try {
      await axiosInstance.post("/logout");
      logoutUser();
    } catch (e) {
      alert("There was an issue logging you out");
    }
  }

  return (
    <nav>
      <ul>
        <li>
          <Link to="/">Home Page</Link>
        </li>
        {!isAuthenticated && (
          <>
            <li>
              <Link to="/login">Login</Link>
            </li>
            <li>
            <Link to="/signup">Signup</Link>
            </li>
          </>
        )}
        {isAuthenticated && (
          <>
            <li>
              <Link to="/dashboard">Dashboard</Link>
            </li>
            <li>
              <button onMouseUp={handleLogOut}>Log out</button>
            </li>
          </>
        )}
      </ul>
    </nav>
  );
}

export default Header;