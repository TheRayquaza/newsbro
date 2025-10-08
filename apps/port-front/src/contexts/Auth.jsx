import React, { createContext, useState, useEffect } from "react";
import Cookies from "js-cookie";
import api from "../api/api";
import { jwtDecode } from "jwt-decode";

const AuthContext = createContext(null);

const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);
  const [loading, setLoading] = useState(true);

  const decodeToken = (token) => {
    try {
      const decoded = jwtDecode(token);
      return {
        id: decoded.user_id,
        email: decoded.email,
        firstName: decoded.first_name,
        lastName: decoded.last_name,
        username: decoded.username,
        role: decoded.role,
      };
    } catch (e) {
      console.error("Failed to decode JWT", e);
      return null;
    }
  };

  useEffect(() => {
    const storedToken = Cookies.get("auth_token");
    if (storedToken) {
      const decodedUser = decodeToken(storedToken);
      if (decodedUser) {
        setToken(storedToken);
        setUser(decodedUser);
      }
    }
    setLoading(false);
  }, []);

  const login = async (email, password) => {
    const data = await api.login(email, password);
    const authToken = data.access_token || data.token;

    const decodedUser = decodeToken(authToken);

    setToken(authToken);
    setUser(decodedUser);

    return data;
  };

  const logout = () => {
    setUser(null);
    setToken(null);

    Cookies.remove("auth_token");
  };

  return (
    <AuthContext.Provider
      value={{ user, token, login, logout, loading }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export { AuthProvider, AuthContext };
export default AuthProvider;
