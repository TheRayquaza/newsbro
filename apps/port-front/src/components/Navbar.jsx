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
    { path: "/", label: "Dashboard", icon: Home },
    { path: "/feeds", label: "Feeds", icon: Rss },
    { path: "/history", label: "My History", icon: History },
    { path: "/search", label: "Search", icon: Search, highlight: true },
  ];

  return (
    <nav className="z-50 bg-slate-900/95 backdrop-blur-xl border-b border-blue-500/20 shadow-lg shadow-blue-500/5">
      <div className="px-4 sm:px-6 lg:px-8 w-full">
        <div className="flex justify-between items-center h-16 max-w-7xl mx-auto">
          <div
            className="flex items-center gap-3 cursor-pointer group" 
            onClick={() => navigate("/")}
          >
            <div className="relative">
              <FileText className="w-8 h-8 text-blue-400 group-hover:text-cyan-400 transition-colors" />
              <div className="absolute -top-1 -right-1 w-2 h-2 bg-cyan-400 rounded-full animate-pulse"></div>
            </div>
            <div className="flex flex-col">
              <span className="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-cyan-400">
                NewsBro
              </span>
              <span className="text-[10px] text-slate-500 -mt-1">Stay Informed</span>
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
                  className={`relative flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all ${
                    active
                      ? "text-blue-400 bg-blue-500/10"
                      : "text-slate-300 hover:text-blue-400 hover:bg-slate-800/50"
                  } ${item.highlight ? "group" : ""}`}
                >
                  <Icon className={`w-4 h-4 ${item.highlight && !active ? "group-hover:animate-pulse" : ""}`} />
                  {item.label}
                  {item.highlight && (
                    <Sparkles className="w-3 h-3 text-cyan-400 animate-pulse" />
                  )}
                  {active && (
                    <div className="absolute bottom-0 left-1/2 -translate-x-1/2 w-8 h-0.5 bg-gradient-to-r from-blue-400 to-cyan-400 rounded-full"></div>
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
                className="flex items-center gap-3 px-4 py-2 rounded-lg bg-slate-800/50 border border-blue-500/20 hover:border-blue-500/40 transition-all group"
              >
                <div className="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center text-white font-bold text-sm">
                  {user?.username?.[0]?.toUpperCase() || "U"}
                </div>
                <div className="flex flex-col items-start">
                  <span className="text-sm font-medium text-slate-200">{user?.username}</span>
                  <span className="text-xs text-slate-500">{user?.email}</span>
                </div>
                <ChevronDown className={`w-4 h-4 text-slate-400 transition-transform ${profileDropdown ? "rotate-180" : ""}`} />
              </button>

              {/* Dropdown */}
              {profileDropdown && (
                <div className="absolute right-0 mt-2 w-56 bg-slate-900 border border-blue-500/20 rounded-xl shadow-xl shadow-black/50 overflow-hidden animate-in fade-in slide-in-from-top-2">
                  <div className="p-4 border-b border-blue-500/20">
                    <p className="text-sm font-medium text-slate-200">{user?.username}</p>
                    <p className="text-xs text-slate-500 truncate">{user?.email}</p>
                  </div>
                  <div className="py-2">
                    <button
                      onClick={() => {
                        setProfileDropdown(false);
                        navigate("/settings");
                      }}
                      className="w-full flex items-center gap-3 px-4 py-2 text-slate-300 hover:text-blue-400 hover:bg-slate-800/50 transition"
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
            className="lg:hidden p-2 rounded-lg text-slate-300 hover:text-blue-400 hover:bg-slate-800/50 transition"
          >
            {menuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      {menuOpen && (
        <div className="lg:hidden bg-slate-900/98 backdrop-blur-xl border-t border-blue-500/20 animate-in slide-in-from-top">
          <div className="px-4 py-4 space-y-2">
            {/* User Info */}
            <div className="flex items-center gap-3 p-4 mb-4 bg-slate-800/50 rounded-xl border border-blue-500/20">
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center text-white font-bold">
                {user?.username?.[0]?.toUpperCase() || "U"}
              </div>
              <div className="flex flex-col">
                <span className="text-sm font-medium text-slate-200">{user?.username}</span>
                <span className="text-xs text-slate-500 truncate">{user?.email}</span>
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
                  className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg font-medium transition-all ${
                    active
                      ? "text-blue-400 bg-blue-500/10 border border-blue-500/30"
                      : "text-slate-300 hover:text-blue-400 hover:bg-slate-800/50"
                  }`}
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
              className="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-slate-300 hover:text-blue-400 hover:bg-slate-800/50 transition font-medium"
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
              className="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-red-400 hover:bg-red-500/10 transition font-medium border-t border-blue-500/20 mt-4 pt-4"
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
