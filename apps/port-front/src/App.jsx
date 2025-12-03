import "./App.css";
import { useContext, useEffect, useState } from "react";
import { BrowserRouter, Route, Routes, useNavigate, Outlet } from "react-router-dom";
import { AuthContext, AuthProvider } from "./contexts/Auth";
import { Navbar } from "./components/Navbar";
import { Footer } from "./components/Footer";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import HomePage from "./pages/HomePage";
import SettingsPage from "./pages/SettingsPage";
import NotFoundPage from "./pages/NotFoundPage";
import DeepSearchPage from "./pages/SearchPage";
import HistoryPage from "./pages/HistoryPage";
import FeedPage from "./pages/FeedPage";

const AuthRoute = () => {
  const { user, loading } = useContext(AuthContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (!loading && !user) {
      navigate("/login");
    }
  }, [user, loading, navigate]);

  if (loading) {
    return (
      <div className="min-h-screen min-w-screen flex items-center justify-center">
        <div className="relative">
          <div className="w-24 h-24 border-4 border-cyan-200/20 border-t-cyan-500 rounded-full animate-spin"></div>
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-16 h-16 border-4 border-blue-200/20 border-b-blue-400 rounded-full animate-spin" style={{ animationDirection: 'reverse', animationDuration: '0.8s' }}></div>
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-4 h-4 bg-gradient-to-r from-cyan-400 to-blue-400 rounded-full animate-pulse"></div>
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-32 h-32 bg-cyan-500/20 rounded-full blur-2xl animate-pulse"></div>
        </div>

        <p className="absolute top-full mt-8 left-1/2 -translate-x-1/2 text-transparent bg-clip-text bg-gradient-to-r from-cyan-400 to-blue-400 text-xl font-semibold tracking-wider animate-pulse whitespace-nowrap">
          Loading
          <span className="inline-block animate-bounce" style={{ animationDelay: '0s' }}>.</span>
          <span className="inline-block animate-bounce" style={{ animationDelay: '0.2s' }}>.</span>
          <span className="inline-block animate-bounce" style={{ animationDelay: '0.4s' }}>.</span>
        </p>
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
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          <Route element={<AuthRoute />}>
            <Route element={<Layout />}>
              <Route path="/settings" element={<SettingsPage />} />

              <Route element={<AuthRoute />}>
                <Route path="/search" element={<DeepSearchPage />} />
              </Route>

              <Route element={<AuthRoute />}>
                <Route path="/history" element={<HistoryPage />} />
              </Route>

              <Route element={<AuthRoute />}>
                <Route path="/feeds" element={<FeedPage />} />
              </Route>

              {/* Home route with wildcard to capture RSS path */}
              <Route path="/" element={<HomePage />} />
              <Route path="/*" element={<HomePage />} />

              <Route path="*" element={<NotFoundPage />} />
            </Route>
          </Route>

          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}