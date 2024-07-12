import React, { createContext, useContext, useState } from "react";

interface AuthContextProps {
  isAuthenticated: boolean;
  loginUser: (email: string, password: string) => Promise<void>;
  logoutUser: () => void;
}

type AuthProviderProps = {
  children: React.ReactNode;
}

const AuthContext = createContext<AuthContextProps | undefined>(undefined);

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

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
      alert("There was an error logging you in!");
      return;
    }
    setIsAuthenticated(true);
  }

  const logoutUser = () => {
    setIsAuthenticated(false);
  }

  return (
    <AuthContext.Provider value={{isAuthenticated, loginUser, logoutUser}}>
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
