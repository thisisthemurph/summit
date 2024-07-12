import React from "react"
import ReactDOM from "react-dom/client"
import { createBrowserRouter, RouterProvider } from "react-router-dom"

import Root from "./shared/routes/root"
import ErrorPage from "./shared/routes/ErrorPage.tsx";
import LoginPage from "./features/auth/LoginPage.tsx";
import HomePage from "./features/home/HomePage.tsx";
import SignUpPage from "./features/auth/SignUpPage.tsx";
import DashboardPage from "./features/dashboard/DashboardPage.tsx";
import PrivateRoute from "./shared/components/PrivateRoute.tsx";

import {AuthProvider} from "./hooks/useAuth.tsx";

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
        path: "/dashboard",
        element: <PrivateRoute />,
        children: [
          {
            path: "/dashboard",
            element: <DashboardPage />,
          }
        ],
      },
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
