import {Outlet} from "react-router-dom";
import Header from "../components/Header.tsx";

function Root() {
  return (
    <>
      <Header />
      <main className="p-4">
        <Outlet />
      </main>
    </>
  )
}

export default Root;
