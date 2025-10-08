import { useNavigate } from "react-router-dom";
import { Home, ArrowLeft, FileQuestion } from "lucide-react";

export default function NotFoundPage() {
  const navigate = useNavigate();

  return (
    <div className="w-full flex-1 flex items-center justify-center px-4 py-12">
      <div className="max-w-2xl w-full text-center">
        {/* Animated 404 */}
        <div className="relative mb-8">
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="w-48 h-48 bg-blue-500/10 rounded-full blur-3xl animate-pulse"></div>
          </div>
          <div className="relative">
            <h1 className="text-[120px] md:text-[160px] font-bold text-transparent bg-clip-text bg-gradient-to-br from-blue-400 via-cyan-400 to-blue-600 leading-none mb-4">
              404
            </h1>
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
              <FileQuestion className="w-20 h-20 text-blue-500/20 animate-bounce" />
            </div>
          </div>
        </div>

        {/* Message */}
        <div className="space-y-3 mb-8">
          <h2 className="text-2xl md:text-3xl font-bold text-slate-200">
            Page Not Found
          </h2>
          <p className="text-base md:text-lg text-slate-400 max-w-md mx-auto">
            Oops! The page you're looking for seems to have wandered off into the digital void.
          </p>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mb-12">
          <button
            onClick={() => navigate("/")}
            className="group flex items-center gap-3 px-6 py-3 bg-gradient-to-r from-blue-500 to-cyan-500 text-white font-semibold rounded-xl hover:from-blue-600 hover:to-cyan-600 transition-all shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40 hover:scale-105"
          >
            <Home className="w-5 h-5" />
            Go Home
          </button>
          
          <button
            onClick={() => navigate(-1)}
            className="flex items-center gap-3 px-6 py-3 bg-slate-800/50 text-slate-300 font-semibold rounded-xl border border-blue-500/20 hover:border-blue-500/40 hover:text-blue-400 hover:bg-slate-800/80 transition-all"
          >
            <ArrowLeft className="w-5 h-5" />
            Go Back
          </button>
        </div>

        {/* Decorative Elements */}
        <div className="grid grid-cols-3 gap-4 max-w-md mx-auto opacity-50">
          <div className="h-1.5 bg-gradient-to-r from-transparent via-blue-500/30 to-transparent rounded-full animate-pulse"></div>
          <div className="h-1.5 bg-gradient-to-r from-transparent via-cyan-500/30 to-transparent rounded-full animate-pulse"></div>
          <div className="h-1.5 bg-gradient-to-r from-transparent via-blue-500/30 to-transparent rounded-full animate-pulse"></div>
        </div>
      </div>
    </div>
  );
};
