import { FileText, Home, Settings, User, LogOut, Menu, X } from "lucide-react";

export const Navbar = ({ menuOpen, setMenuOpen, handleLogout, user }) => {
  return (
    <nav className="bg-slate-900 border-b border-blue-500/20 w-full">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center gap-2">
            <FileText className="w-8 h-8 text-blue-400" />
            <span className="text-2xl font-bold text-blue-400">NewsBro</span>
          </div>

          <div className="hidden md:flex items-center gap-6">
            <button className="flex items-center gap-2 text-slate-300 hover:text-blue-400 transition">
              <Home className="w-5 h-5" />
              Dashboard
            </button>
            <button className="flex items-center gap-2 text-slate-300 hover:text-blue-400 transition">
              <Settings className="w-5 h-5" />
              Settings
            </button>
            <div className="flex items-center gap-3 text-slate-300">
              <User className="w-5 h-5" />
              <span>{user?.email}</span>
            </div>
            <button
              onClick={handleLogout}
              className="flex items-center gap-2 bg-red-500/10 text-red-400 px-4 py-2 rounded-lg hover:bg-red-500/20 transition"
            >
              <LogOut className="w-4 h-4" />
              Logout
            </button>
          </div>

          <button
            onClick={() => setMenuOpen(!menuOpen)}
            className="md:hidden text-slate-300"
          >
            {menuOpen ? (
              <X className="w-6 h-6" />
            ) : (
              <Menu className="w-6 h-6" />
            )}
          </button>
        </div>
      </div>

      {menuOpen && (
        <div className="md:hidden bg-slate-900/95 border-t border-blue-500/20 p-4 space-y-3">
          <button className="w-full flex items-center gap-2 text-slate-300 hover:text-blue-400 transition p-2">
            <Home className="w-5 h-5" />
            Dashboard
          </button>
          <button className="w-full flex items-center gap-2 text-slate-300 hover:text-blue-400 transition p-2">
            <Settings className="w-5 h-5" />
            Settings
          </button>
          <div className="flex items-center gap-3 text-slate-300 p-2">
            <User className="w-5 h-5" />
            <span>{user?.email}</span>
          </div>
          <button
            onClick={handleLogout}
            className="w-full flex items-center gap-2 bg-red-500/10 text-red-400 px-4 py-2 rounded-lg hover:bg-red-500/20 transition"
          >
            <LogOut className="w-4 h-4" />
            Logout
          </button>
        </div>
      )}
    </nav>
  );
};
