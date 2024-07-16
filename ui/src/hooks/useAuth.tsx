import React, {createContext, useContext, useEffect, useState} from "react";
import { AuthenticatedUser } from "../shared/types/responseTypes";
import axiosInstance from "../shared/requests/axiosInstance";
import {AxiosError} from "axios";

interface AuthContextProps {
  isAuthenticated: boolean;
  authenticatedUser: AuthenticatedUser | null;
  loginUser: (email: string, password: string) => Promise<void>;
  loginUserWithToken: (token: string) => Promise<void>;
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
    try {
      const response = await axiosInstance.post<AuthenticatedUser>("/login", {email, password});
      setIsAuthenticated(true);
      setAuthenticatedUser(response.data);
    } catch (e) {
      let message = "There was an issue logging you in";
      if (e instanceof AxiosError) {
        message = e.response?.data?.message ? e.response.data.message : message;
      }

      setIsAuthenticated(false);
      setAuthenticatedUser(null);
      alert(message);
    }
  }

  const loginUserWithToken = async (token: string) => {
    try {
      await axiosInstance.post("/login/token", {token});
      setIsAuthenticated(true);
      setAuthenticatedUser(null); // TODO: Send user back
    } catch (e) {
      let message = "There was an issue logging you in";
      if (e instanceof AxiosError) {
        message = e.response?.data?.message ? e.response.data.message : message;
      }

      setIsAuthenticated(false);
      setAuthenticatedUser(null);
      alert(message);
    }
  }

  const logoutUser = () => {
    setIsAuthenticated(false);
    setAuthenticatedUser(null);
  }

  return (
    <AuthContext.Provider value={{isAuthenticated, authenticatedUser, loginUser, logoutUser, loginUserWithToken}}>
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
