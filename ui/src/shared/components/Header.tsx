import {Link} from "react-router-dom";

function Header() {
  return (
    <section className="p-4 bg-base-200">
      <div>
        <h1>React Router</h1>
      </div>
      <nav>
        <ul>
          <li>
            <Link to="/">Home Page</Link>
          </li>
          <li>
            <Link to="/login">Login</Link>
          </li>
          <li>
            <Link to="/signup">Signup</Link>
          </li>
        </ul>
      </nav>
    </section>
  )
}

export default Header;