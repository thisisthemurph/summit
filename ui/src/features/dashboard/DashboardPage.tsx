import { useAuth } from "../../hooks/useAuth.tsx";

function DashboardPage() {
  const { authenticatedUser } = useAuth();
  return (
    <>
      <h2>Dashboard</h2>
      <p>{authenticatedUser?.email ?? "anon"}</p>
    </>
  )
}

export default DashboardPage;
