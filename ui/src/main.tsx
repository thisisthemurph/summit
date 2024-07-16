import React from "react"
import ReactDOM from "react-dom/client"
import { createBrowserRouter, RouterProvider } from "react-router-dom"

import { AuthProvider } from "./hooks/useAuth.tsx";
import OnboardingPrivateRoute from "./shared/components/OnboardingPrivateRoute.tsx";
import PrivateRoute from "./shared/components/PrivateRoute.tsx";

import Root from "./shared/routes/root"
import ErrorPage from "./shared/routes/ErrorPage.tsx";
import LoginPage from "./features/auth/LoginPage.tsx";
import HomePage from "./features/home/HomePage.tsx";
import SignUpPage from "./features/auth/SignUpPage.tsx";
import DashboardPage from "./features/dashboard/DashboardPage.tsx";
import ProfileSetupPage from "./features/onboarding/ProfileSetupPage.tsx";
import AuthConfirmationPage from "./features/auth/AuthConfirmationPage.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "/",
        element: <HomePage />
      },
      {
        path: "/login",
        element: <LoginPage />,
      },
      {
        path: "/signup",
        element: <SignUpPage />,
      },
      {
        // https://uizsdwzazgcdodwmubcq.supabase.co/auth/v1/verify?token=3dd3309a89b8b33660e06c21ea5335f4bf4204a14269967ba991d2b6&type=signup&redirect_to=http://localhost:3000
        path: "/auth/callback",
        element: <AuthConfirmationPage />
      },
      {
        path: "/dashboard",
        element: <PrivateRoute />,
        children: [
          {
            path: "/dashboard",
            element: <DashboardPage />,
          }
        ],
      },
      {
        path: "/onboarding",
        element: <OnboardingPrivateRoute />,
        children: [
          {
            path: "/onboarding/profile",
            element: <ProfileSetupPage />,
          }
        ]
      }
    ]
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  </React.StrictMode>,
)
