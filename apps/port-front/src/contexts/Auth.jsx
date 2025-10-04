import React, { createContext, useState, useEffect } from "react";
import api from "../api/api";

const AuthContext = createContext(null);

const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const storedToken = localStorage.getItem("auth_token");
    const storedUser = localStorage.getItem("user_email");

    if (storedToken && storedUser) {
      setToken(storedToken);
      setUser({ email: storedUser });
    }
    setLoading(false);
  }, []);

  const login = async (email, password) => {
    const data = await api.login(email, password);
    const authToken = data.access_token || data.token;

    setToken(authToken);
    setUser({ email });
    localStorage.setItem("auth_token", authToken);
    localStorage.setItem("user_email", email);

    return data;
  };

  const register = async (userData) => {
    return await api.register(userData);
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    localStorage.removeItem("auth_token");
    localStorage.removeItem("user_email");
    api.clearCache();
  };

  return (
    <AuthContext.Provider
      value={{ user, token, login, register, logout, loading }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export { AuthProvider, AuthContext };

export default AuthProvider;
