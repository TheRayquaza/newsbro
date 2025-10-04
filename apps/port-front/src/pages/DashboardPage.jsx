import { useState, useContext } from "react";
import { useNavigate } from "react-router";
import { AuthContext } from "../contexts/Auth";
import { Footer } from "../components/Footer";
import { Navbar } from "../components/Navbar";

const DashboardPage = () => {
  const { user, logout } = useContext(AuthContext);
  let navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
        width: "100vw",
      }}
    >
      <div className="flex flex-col min-h-screen w-full">
        {/* Navbar */}
        <Navbar
          menuOpen={menuOpen}
          setMenuOpen={setMenuOpen}
          handleLogout={handleLogout}
          user={user}
        />

        {/* Centered main section */}
        <main className="flex-1 flex items-center justify-center w-full px-4 sm:px-6 lg:px-8">
          <div className="w-full max-w-7xl bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-blue-500/20 p-8">
            <h2 className="text-3xl font-bold text-blue-400 mb-4 text-center">
              🎉 Welcome, {user?.email}!
            </h2>

            <div className="bg-blue-500/10 border border-blue-500/30 text-blue-300 px-6 py-4 rounded-xl mb-8 text-center">
              <p className="font-semibold">
                🚧 This site is currently under development. New features coming
                soon!
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 w-full">
              <div className="bg-slate-800/40 border border-blue-500/20 rounded-xl p-6 hover:bg-slate-800/60 transition cursor-pointer w-full">
                <h3 className="text-xl font-semibold text-blue-400 mb-2">
                  Latest Articles
                </h3>
                <p className="text-slate-400">
                  Stay updated with trending news
                </p>
              </div>
              <div className="bg-slate-800/40 border border-blue-500/20 rounded-xl p-6 hover:bg-slate-800/60 transition cursor-pointer w-full">
                <h3 className="text-xl font-semibold text-blue-400 mb-2">
                  Saved Items
                </h3>
                <p className="text-slate-400">Your bookmarked content</p>
              </div>
              <div className="bg-slate-800/40 border border-blue-500/20 rounded-xl p-6 hover:bg-slate-800/60 transition cursor-pointer w-full">
                <h3 className="text-xl font-semibold text-blue-400 mb-2">
                  Preferences
                </h3>
                <p className="text-slate-400">Customize your feed</p>
              </div>
            </div>
          </div>
        </main>

        <Footer />
      </div>
    </div>
  );
};

export default DashboardPage;
