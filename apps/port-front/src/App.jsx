import "./App.css";
import { useContext, useEffect } from "react";
import { BrowserRouter, Route, Routes, useNavigate } from "react-router-dom";
import { AuthContext, AuthProvider } from "./contexts/Auth";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import DashboardPage from "./pages/DashboardPage";
import AdminPage from "./pages/AdminPage";

const ProtectedRoute = ({ element, adminOnly = false }) => {
  const { user, loading } = useContext(AuthContext);
  let navigate = useNavigate();

  useEffect(() => {
    if (!loading && !user) {
      navigate("/login");
    }
    if (!loading && adminOnly && user?.role !== "admin") {
      navigate("/");
    }
  }, [user, loading, navigate, adminOnly]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-black">
        <div className="text-blue-400 text-xl">Loading...</div>
      </div>
    );
  }

  return user ? element : null;
};

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route
            path="/"
            element={<ProtectedRoute element={<DashboardPage />} />}
          />
          <Route
            path="/admin"
            element={<ProtectedRoute element={<AdminPage />} adminOnly />}
          />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}
