import {useNavigate} from "react-router-dom";
import Container from "../../shared/components/Container.tsx";
import {useEffect, useState} from "react";
import {useAuth} from "../../hooks/useAuth.tsx";

function AuthConfirmationPage() {
  const { loginUserWithToken } = useAuth();
  const navigate = useNavigate();
  const [error, setError] = useState<string>("");

  const parseErrorFromUrl = (urlString: string): string | null => {
    const url = new URL(urlString);
    const fragment = url.hash.substring(1);
    const params = new URLSearchParams(fragment);
    return params.get("error_description")
  }

  const parseTokenFromUrl = (urlString: string): string | null => {
    const url = new URL(urlString);
    const fragment = url.hash.substring(1);
    const params = new URLSearchParams(fragment);
    return params.get("access_token")
  }

  const url = window.location.href;

  useEffect(() => {
    const error = parseErrorFromUrl(url);
    if (error) {
      setError(error);
      return;
    }

    const token = parseTokenFromUrl(url);
    if (!token) {
      setError("Invalid token.");
      return;
    }

    loginUserWithToken(token).finally(() => {
      navigate("/dashboard")
    });
  }, [url]);

  if (error) return <Container><p>Error: {error}</p></Container>;

  return (
    <Container>
      <p>Please wait to be redirected...</p>
    </Container>
  )
}

export default AuthConfirmationPage;