import { FileText, Github, Mail, Heart, Lock } from "lucide-react";
import epitaLogo from "../assets/icons/epita.png";

export const Footer = () => {
  return (
    <footer style={{
      backgroundColor: 'var(--nav-bg)',
      borderTop: '1px solid var(--nav-border)',
    }} className="w-full mt-auto">
      <div  className="max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-6">
        <div  className="flex flex-col md:flex-row justify-between items-center gap-4">
          {/* Left side */}
          <div style={{ color: 'var(--nav-text-muted)' }} className="flex items-center gap-2">
            <FileText style={{ color: 'var(--nav-active-text)' }} className="w-5 h-5" />
            <span className="text-sm">
              © 2025 NewsBro. All rights reserved. version: {__APP_VERSION__}
            </span>
          </div>
          
          {/* Right side */}
          <div className="flex flex-wrap justify-center md:justify-end items-center gap-6">
            <a
              href="https://github.com/TheRayquaza/newsbro"
              target="_blank"
              rel="noopener noreferrer"
              style={{ color: 'var(--nav-text-muted)' }}
              onMouseEnter={(e) => e.currentTarget.style.color = 'var(--nav-active-text)'}
              onMouseLeave={(e) => e.currentTarget.style.color = 'var(--nav-text-muted)'}
              className="flex items-center gap-2 transition"
            >
              <Github className="w-5 h-5" />
              <span className="text-sm">GitHub</span>
            </a>
            
            <a
              href="https://www.epita.fr/en"
              target="_blank"
              rel="noopener noreferrer"
              style={{ color: 'var(--nav-text-muted)' }}
              onMouseEnter={(e) => e.currentTarget.style.color = 'var(--nav-active-text)'}
              onMouseLeave={(e) => e.currentTarget.style.color = 'var(--nav-text-muted)'}
              className="flex items-center gap-2 transition"
            >
              <img src={epitaLogo} alt="Epita Logo" className="w-7 h-5" />
              <span className="text-sm">Epita</span>
            </a>
            
            <a
              href="mailto:mateo.lelong@epita.fr"
              style={{ color: 'var(--nav-text-muted)' }}
              onMouseEnter={(e) => e.currentTarget.style.color = 'var(--nav-active-text)'}
              onMouseLeave={(e) => e.currentTarget.style.color = 'var(--nav-text-muted)'}
              className="flex items-center gap-2 transition"
            >
              <Mail className="w-5 h-5" />
              <span className="text-sm">Contact</span>
            </a>
            
            <a
              href="/privacy-policy"
              style={{ color: 'var(--nav-text-muted)' }}
              onMouseEnter={(e) => e.currentTarget.style.color = 'var(--nav-active-text)'}
              onMouseLeave={(e) => e.currentTarget.style.color = 'var(--nav-text-muted)'}
              className="flex items-center gap-2 transition"
            >
              <Lock className="w-4 h-4" />
              <span className="text-sm">Privacy</span>
            </a>
            
            <div style={{ color: 'var(--nav-text-muted)' }} className="flex items-center gap-2">
              <span className="text-sm">Made with</span>
              <Heart className="w-4 h-4 text-red-400 fill-current" />
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
};
