import React, { createContext, useState, useEffect, useCallback } from "react";
import Cookies from "js-cookie";
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

  const refreshAuth = useCallback(() => {
    const storedToken = Cookies.get("auth_token");
    if (storedToken) {
      const decodedUser = decodeToken(storedToken);
      if (decodedUser) {
        setToken(storedToken);
        setUser(decodedUser);
        return;
      }
    }
    setToken(null);
    setUser(null);
  }, []);

  useEffect(() => {
    refreshAuth();
    setLoading(false);
  }, [refreshAuth]);

  const logout = () => {
    setUser(null);
    setToken(null);

    Cookies.remove("auth_token");
  };

  return (
    <AuthContext.Provider
      value={{ user, token, logout, loading, refreshAuth }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export { AuthProvider, AuthContext };
export default AuthProvider;
