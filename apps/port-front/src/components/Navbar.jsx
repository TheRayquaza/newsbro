import { 
  FileText, 
  Home, 
  Settings, 
  User, 
  LogOut, 
  Menu, 
  X, 
  Rss,
  Search,
  ChevronDown,
  History,
  Sparkles
} from "lucide-react";
import { useNavigate, useLocation } from "react-router-dom";
import { useState } from "react";

export const Navbar = ({ menuOpen, setMenuOpen, handleLogout, user }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const [profileDropdown, setProfileDropdown] = useState(false);

  const isActive = (path) => location.pathname === path;

  const navItems = [
    { path: "/", label: "Home", icon: Home },
    { path: "/feeds", label: "Feeds", icon: Rss },
    { path: "/history", label: "My History", icon: History },
    { path: "/search", label: "Search", icon: Search, highlight: true },
  ];

  return (
    <nav style={{
      backgroundColor: 'var(--nav-bg)',
      borderBottom: '1px solid var(--nav-border)',
      boxShadow: '0 10px 15px -3px var(--nav-shadow)',
    }} className="z-50 backdrop-blur-xl">
      <div className="px-4 sm:px-6 lg:px-8 w-full">
        <div className="flex justify-between items-center h-16 max-w-7xl mx-auto">
          <div
            className="flex items-center gap-3 cursor-pointer group" 
            onClick={() => navigate("/")}
          >
            <div className="relative">
              <FileText style={{ color: 'var(--nav-active-text)' }} className="w-8 h-8 group-hover:opacity-80 transition-opacity" />
              <div className="absolute -top-1 -right-1 w-2 h-2 bg-cyan-400 rounded-full animate-pulse"></div>
            </div>
            <div className="flex flex-col">
              <span style={{
                background: `linear-gradient(to right, var(--nav-gradient-from), var(--nav-gradient-to))`,
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                backgroundClip: 'text',
              }} className="text-2xl font-bold">
                NewsBro
              </span>
              <span style={{ color: 'var(--nav-text-muted)' }} className="text-[10px] -mt-1">Stay Informed</span>
            </div>
          </div>

          {/* Desktop Menu */}
          <div className="hidden lg:flex items-center gap-2">
            {navItems.map((item) => {
              const Icon = item.icon;
              const active = isActive(item.path);
              return (
                <button
                  key={item.path}
                  onClick={() => navigate(item.path)}
                  style={{
                    color: active ? 'var(--nav-active-text)' : 'var(--nav-text)',
                    backgroundColor: active ? 'var(--nav-active-bg)' : 'transparent',
                  }}
                  onMouseEnter={(e) => {
                    if (!active) e.currentTarget.style.backgroundColor = 'var(--nav-hover-bg)';
                  }}
                  onMouseLeave={(e) => {
                    if (!active) e.currentTarget.style.backgroundColor = 'transparent';
                  }}
                  className={`relative flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all ${item.highlight ? "group" : ""}`}
                >
                  <Icon className={`w-4 h-4 ${item.highlight && !active ? "group-hover:animate-pulse" : ""}`} />
                  {item.label}
                  {item.highlight && (
                    <Sparkles className="w-3 h-3 text-cyan-400 animate-pulse" />
                  )}
                  {active && (
                    <div style={{
                      background: `linear-gradient(to right, var(--nav-gradient-from), var(--nav-gradient-to))`,
                    }} className="absolute bottom-0 left-1/2 -translate-x-1/2 w-8 h-0.5 rounded-full"></div>
                  )}
                </button>
              );
            })}
          </div>

          {/* Desktop User Menu */}
          <div className="hidden lg:flex items-center gap-3">
            <div className="relative">
              <button
                onClick={() => setProfileDropdown(!profileDropdown)}
                style={{
                  backgroundColor: 'var(--nav-profile-bg)',
                  borderColor: 'var(--nav-profile-border)',
                }}
                className="flex items-center gap-3 px-4 py-2 rounded-lg border hover:opacity-90 transition-all group"
              >
                <div style={{
                  background: `linear-gradient(to bottom right, var(--nav-gradient-from), var(--nav-gradient-to))`,
                }} className="w-8 h-8 rounded-full flex items-center justify-center text-white font-bold text-sm">
                  {user?.username?.[0]?.toUpperCase() || "U"}
                </div>
                <div className="flex flex-col items-start">
                  <span style={{ color: 'var(--nav-text)' }} className="text-sm font-medium">{user?.username}</span>
                  <span style={{ color: 'var(--nav-text-muted)' }} className="text-xs">{user?.email}</span>
                </div>
                <ChevronDown style={{ color: 'var(--nav-text-muted)' }} className={`w-4 h-4 transition-transform ${profileDropdown ? "rotate-180" : ""}`} />
              </button>

              {/* Dropdown */}
              {profileDropdown && (
                <div style={{
                  backgroundColor: 'var(--nav-dropdown-bg)',
                  borderColor: 'var(--nav-dropdown-border)',
                }} className="absolute right-0 mt-2 w-56 border rounded-xl shadow-xl shadow-black/50 overflow-hidden animate-in fade-in slide-in-from-top-2">
                  <div style={{ borderColor: 'var(--nav-dropdown-border)' }} className="p-4 border-b">
                    <p style={{ color: 'var(--nav-text)' }} className="text-sm font-medium">{user?.username}</p>
                    <p style={{ color: 'var(--nav-text-muted)' }} className="text-xs truncate">{user?.email}</p>
                  </div>
                  <div className="py-2">
                    <button
                      onClick={() => {
                        setProfileDropdown(false);
                        navigate("/settings");
                      }}
                      style={{ color: 'var(--nav-text)' }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.backgroundColor = 'var(--nav-hover-bg)';
                        e.currentTarget.style.color = 'var(--nav-active-text)';
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.backgroundColor = 'transparent';
                        e.currentTarget.style.color = 'var(--nav-text)';
                      }}
                      className="w-full flex items-center gap-3 px-4 py-2 transition"
                    >
                      <Settings className="w-4 h-4" />
                      <span>Settings</span>
                    </button>
                    <button
                      onClick={() => {
                        setProfileDropdown(false);
                        handleLogout();
                      }}
                      className="w-full flex items-center gap-3 px-4 py-2 text-red-400 hover:bg-red-500/10 transition"
                    >
                      <LogOut className="w-4 h-4" />
                      <span>Logout</span>
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Mobile Menu Button */}
          <button
            onClick={() => setMenuOpen(!menuOpen)}
            style={{ color: 'var(--nav-text)' }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = 'var(--nav-hover-bg)';
              e.currentTarget.style.color = 'var(--nav-active-text)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = 'transparent';
              e.currentTarget.style.color = 'var(--nav-text)';
            }}
            className="lg:hidden p-2 rounded-lg transition"
          >
            {menuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      {menuOpen && (
        <div style={{
          backgroundColor: 'var(--nav-bg)',
          borderTop: '1px solid var(--nav-border)',
        }} className="lg:hidden backdrop-blur-xl animate-in slide-in-from-top">
          <div className="px-4 py-4 space-y-2">
            {/* User Info */}
            <div style={{
              backgroundColor: 'var(--nav-profile-bg)',
              borderColor: 'var(--nav-profile-border)',
            }} className="flex items-center gap-3 p-4 mb-4 rounded-xl border">
              <div style={{
                background: `linear-gradient(to bottom right, var(--nav-gradient-from), var(--nav-gradient-to))`,
              }} className="w-12 h-12 rounded-full flex items-center justify-center text-white font-bold">
                {user?.username?.[0]?.toUpperCase() || "U"}
              </div>
              <div className="flex flex-col">
                <span style={{ color: 'var(--nav-text)' }} className="text-sm font-medium">{user?.username}</span>
                <span style={{ color: 'var(--nav-text-muted)' }} className="text-xs truncate">{user?.email}</span>
              </div>
            </div>

            {/* Navigation Items */}
            {navItems.map((item) => {
              const Icon = item.icon;
              const active = isActive(item.path);
              return (
                <button
                  key={item.path}
                  onClick={() => {
                    setMenuOpen(false);
                    navigate(item.path);
                  }}
                  style={{
                    color: active ? 'var(--nav-active-text)' : 'var(--nav-text)',
                    backgroundColor: active ? 'var(--nav-active-bg)' : 'transparent',
                    borderColor: active ? 'var(--nav-profile-border)' : 'transparent',
                  }}
                  className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg font-medium transition-all ${active ? 'border' : ''}`}
                >
                  <Icon className="w-5 h-5" />
                  <span>{item.label}</span>
                  {item.highlight && (
                    <Sparkles className="w-4 h-4 text-cyan-400 ml-auto animate-pulse" />
                  )}
                </button>
              );
            })}

            {/* Settings */}
            <button
              onClick={() => {
                setMenuOpen(false);
                navigate("/settings");
              }}
              style={{ color: 'var(--nav-text)' }}
              onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = 'var(--nav-hover-bg)';
                e.currentTarget.style.color = 'var(--nav-active-text)';
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = 'transparent';
                e.currentTarget.style.color = 'var(--nav-text)';
              }}
              className="w-full flex items-center gap-3 px-4 py-3 rounded-lg transition font-medium"
            >
              <Settings className="w-5 h-5" />
              <span>Settings</span>
            </button>

            {/* Logout */}
            <button
              onClick={() => {
                setMenuOpen(false);
                handleLogout();
              }}
              style={{ borderColor: 'var(--nav-border)' }}
              className="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-red-400 hover:bg-red-500/10 transition font-medium border-t mt-4 pt-4"
            >
              <LogOut className="w-5 h-5" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      )}
    </nav>
  );
};