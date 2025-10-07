import { useContext } from "react";
import { AuthContext } from "../contexts/Auth";

const DashboardPage = () => {
  const { user } = useContext(AuthContext);

  return (
    <div className="w-full px-4 sm:px-6 lg:px-8 py-8">
      <div className="max-w-7xl mx-auto">
        <div className="w-full bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-blue-500/20 p-8">
          <h2 className="text-3xl font-bold text-blue-400 mb-4 text-center">
            🎉 Welcome, {user?.username}!
          </h2>
          <div className="bg-blue-500/10 border border-blue-500/30 text-blue-300 px-6 py-4 rounded-xl mb-8 text-center">
            <p className="font-semibold">
              🚧 This site is currently under development. New features coming soon!
            </p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 w-full">
            <div className="bg-slate-800/40 border border-blue-500/20 rounded-xl p-6 hover:bg-slate-800/60 transition cursor-pointer w-full">
              <h3 className="text-xl font-semibold text-blue-400 mb-2">
                Latest Articles
              </h3>
              <p className="text-slate-400">Stay updated with trending news</p>
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
      </div>
    </div>
  );
};

export default DashboardPage;