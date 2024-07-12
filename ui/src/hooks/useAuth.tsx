import React, {createContext, useContext, useEffect, useState} from "react";
import { AuthenticatedUser } from "../shared/types/responseTypes.ts";

interface AuthContextProps {
  isAuthenticated: boolean;
  authenticatedUser: AuthenticatedUser | null;
  loginUser: (email: string, password: string) => Promise<void>;
  logoutUser: () => void;
}

type AuthProviderProps = {
  children: React.ReactNode;
}

const AuthContext = createContext<AuthContextProps | undefined>(undefined);
const IsAuthenticatedStorageKey = "isAuthenticated";
const AuthenticatedUserStorageKey = "authenticatedUser";

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => {
    const state = localStorage.getItem(IsAuthenticatedStorageKey);
    return state === "true";
  });

  const [authenticatedUser, setAuthenticatedUser] = useState<AuthenticatedUser | null>(() => {
    const state = localStorage.getItem(AuthenticatedUserStorageKey);
    return state ? JSON.parse(state) : null;
  });

  useEffect(() => {
    localStorage.setItem(IsAuthenticatedStorageKey, isAuthenticated.toString())
  }, [isAuthenticated])

  useEffect(() => {
    if (authenticatedUser) {
      localStorage.setItem(AuthenticatedUserStorageKey, JSON.stringify(authenticatedUser));
    } else {
      localStorage.removeItem(AuthenticatedUserStorageKey);
    }
  }, [authenticatedUser]);

  const loginUser = async (email: string, password: string) => {
    const baseUrl = import.meta.env.VITE_API_BASE_URL;
    const result = await fetch(`${baseUrl}/login`, {
      method: "POST",
      body: JSON.stringify({ email, password }),
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      }
    });

    if (result.status !== 200) {
      setIsAuthenticated(false);
      setAuthenticatedUser(null);
      alert("There was an error logging you in!");
      return;
    }

    const user = await result.json();
    setIsAuthenticated(true);
    setAuthenticatedUser(user);
  }

  const logoutUser = () => {
    setIsAuthenticated(false);
    setAuthenticatedUser(null);
  }

  return (
    <AuthContext.Provider value={{isAuthenticated, authenticatedUser, loginUser, logoutUser}}>
      {children}
    </AuthContext.Provider>
  )
};

export const useAuth = (): AuthContextProps => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within a AuthProvider");
  }
  return context;
}
