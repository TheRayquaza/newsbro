import "./App.css";
import { useContext, useEffect, useState } from "react";
import { BrowserRouter, Route, Routes, useNavigate, Outlet } from "react-router-dom";
import { AuthContext, AuthProvider } from "./contexts/Auth";
import { Navbar } from "./components/Navbar";
import { Footer } from "./components/Footer";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import DashboardPage from "./pages/DashboardPage";
import AdminPage from "./pages/AdminPage";
import SettingsPage from "./pages/SettingsPage";
import NotFoundPage from "./pages/NotFoundPage";
import DeepSearchPage from "./pages/SearchPage";

const AuthRoute = ({ adminOnly = false }) => {
  const { user, loading } = useContext(AuthContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (!loading && !user) {
      navigate("/login");
    } else if (!loading && adminOnly && user?.role !== "admin") {
      navigate("/");
    }
  }, [user, loading, navigate, adminOnly]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-blue-400 text-xl">Loading...</div>
      </div>
    );
  }

  return user ? <Outlet /> : null;
};

const Layout = () => {
  const { user, logout } = useContext(AuthContext);
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <div className="flex flex-col min-h-screen main-div" style={{ width: "100vw" }}>
      <Navbar
        menuOpen={menuOpen}
        setMenuOpen={setMenuOpen}
        handleLogout={logout}
        user={user}
      />
      <Outlet />
      <Footer />
    </div>
  );
};

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          {/* Public routes */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          {/* Protected routes (user & admin) */}
          <Route element={<AuthRoute />}>
            <Route element={<Layout />}>
              <Route path="/" element={<DashboardPage />} />
              <Route path="/settings" element={<SettingsPage />} />
              
              {/* Admin-only routes */}
              <Route element={<AuthRoute adminOnly />}>
                <Route path="/admin" element={<AdminPage />} />
              </Route>

              <Route element={<AuthRoute />}>
                <Route path="/search" element={<DeepSearchPage />} />
              </Route>
              
              {/* 404 for authenticated users - with navbar/footer */}
              <Route path="*" element={<NotFoundPage />} />
            </Route>
          </Route>

          {/* 404 for public routes - without navbar/footer */}
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}