import { FileText, Github, Mail, Heart, Lock } from "lucide-react";
import epitaLogo from "../assets/icons/epita.png";

export const Footer = () => {
  return (
    <footer className="w-full bg-slate-900 border-t border-blue-500/20 mt-auto">
      <div className="max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-6">
        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          {/* Left side */}
          <div className="flex items-center gap-2 text-slate-400">
            <FileText className="w-5 h-5 text-blue-400" />
            <span className="text-sm">
              © 2025 NewsBro. All rights reserved.
            </span>
          </div>

          {/* Right side */}
          <div className="flex flex-wrap justify-center md:justify-end items-center gap-6">
            <a
              href="https://github.com/TheRayquaza/newsbro"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-slate-400 hover:text-blue-400 transition"
            >
              <Github className="w-5 h-5" />
              <span className="text-sm">GitHub</span>
            </a>

            <a
              href="https://www.epita.fr/en"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-slate-400 hover:text-blue-400 transition"
            >
              <img src={epitaLogo} alt="Epita Logo" className="w-7 h-5" />
              <span className="text-sm">Epita</span>
            </a>

            <a
              href="mailto:mateo.lelong@epita.fr"
              className="flex items-center gap-2 text-slate-400 hover:text-blue-400 transition"
            >
              <Mail className="w-5 h-5" />
              <span className="text-sm">Contact</span>
            </a>

            <a
              href="/privacy-policy"
              className="flex items-center gap-2 text-slate-400 hover:text-blue-400 transition"
            >
              <Lock className="w-4 h-4" />
              <span className="text-sm">Privacy</span>
            </a>

            <div className="flex items-center gap-2 text-slate-400">
              <span className="text-sm">Made with</span>
              <Heart className="w-4 h-4 text-red-400 fill-current" />
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
};
