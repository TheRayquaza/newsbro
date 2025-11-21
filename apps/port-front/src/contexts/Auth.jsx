import React, { createContext, useState, useEffect, useCallback } from "react";
import Cookies from "js-cookie";
import api from "../api/api";

const AuthContext = createContext(null);

const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const refreshAuth = useCallback(async () => {
    try {
      const profile = await api.getProfile();
      setUser(profile);
    } catch {
      setUser(null);
    }
  }, []);

  useEffect(() => {
    setLoading(true);
    refreshAuth().finally(() => setLoading(false));
  }, [refreshAuth]);

  const logout = () => {
    Cookies.remove("auth_token");
    Cookies.remove("refresh_token");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, logout, loading, refreshAuth }}>
      {children}
    </AuthContext.Provider>
  );
};

export { AuthProvider, AuthContext };
export default AuthProvider;
