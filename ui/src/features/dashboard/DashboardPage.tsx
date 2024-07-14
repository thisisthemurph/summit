import { useAuth } from "../../hooks/useAuth.tsx";
import PageHeader from "../../shared/components/PageHeader.tsx";
import Container from "../../shared/components/Container.tsx";

function DashboardPage() {
  const { authenticatedUser } = useAuth();
  return (
    <>
      <PageHeader title={authenticatedUser?.firstName} subtitle="Welcome to your dashboard!" />
      <Container>
        <p>This is the rest of the page content.</p>
        <p>This is some more content.</p>
      </Container>
    </>
  )
}

export default DashboardPage;
