import { Menu, X, BookOpen } from 'lucide-react';

const GithubIcon = ({ size = 20, className = '' }) => (
  <svg
    viewBox="0 0 24 24"
    width={size}
    height={size}
    className={className}
    fill="currentColor"
  >
    <path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.462-1.11-1.462-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.308.678.917.678 1.847 0 1.333-.012 2.409-.012 2.736 0 .267.18.578.688.48C19.138 20.165 22 16.418 22 12c0-5.523-4.477-10-10-10z" />
  </svg>
);
import { useState } from 'react';
import gopherImage from '../../assets/images/gopher.svg';
import { EditorMode } from '../../types';

interface NavbarProps {
  mode: EditorMode;
  onModeChange: (mode: EditorMode) => void;
}

export const Navbar: React.FC<NavbarProps> = ({ mode, onModeChange }) => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const handleModeChange = (newMode: EditorMode) => {
    onModeChange(newMode);
    setIsMenuOpen(false);
  };

  return (
    <nav className="bg-[#2d2d2d] border-b border-gray-700">
      <div className="container mx-auto px-4">
        <div className="h-16 flex items-center">
          {/* Logo Section */}
          <div className="flex items-center space-x-3">
            <img src={gopherImage} alt="Go Gopher" className="h-8 w-auto" />
            <span className="text-xl font-semibold text-white">Go Playground</span>
          </div>

          {/* Mobile Menu Button */}
          <div className="md:hidden ml-auto">
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="text-gray-400 hover:text-white p-2 rounded-md transition-colors"
            >
              {isMenuOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>

          {/* Desktop Navigation - Centered */}
          <div className="hidden md:flex flex-1 justify-center">
            <div className="flex rounded-lg overflow-hidden bg-[#1e1e1e] p-1">
              <button
                onClick={() => handleModeChange('go')}
                className={`px-6 py-2 rounded-md text-sm font-medium transition-colors ${
                  mode === 'go'
                    ? 'bg-[#007acc] text-white shadow-md'
                    : 'text-gray-400 hover:text-white hover:bg-[#3d3d3d]'
                }`}
              >
                Go Playground
              </button>
              <button
                onClick={() => handleModeChange('json')}
                className={`px-6 py-2 rounded-md text-sm font-medium transition-colors ${
                  mode === 'json'
                    ? 'bg-[#007acc] text-white shadow-md'
                    : 'text-gray-400 hover:text-white hover:bg-[#3d3d3d]'
                }`}
              >
                JSON Viewer
              </button>
              <button
                onClick={() => handleModeChange('markdown')}
                className={`px-6 py-2 rounded-md text-sm font-medium transition-colors ${
                  mode === 'markdown'
                    ? 'bg-[#007acc] text-white shadow-md'
                    : 'text-gray-400 hover:text-white hover:bg-[#3d3d3d]'
                }`}
              >
                Markdown Editor
              </button>
            </div>
          </div>

          {/* Desktop Links - Right Side */}
          <div className="hidden md:flex items-center space-x-1">
            <a 
              href="https://go.dev" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-gray-400 hover:text-white transition-colors p-2 rounded-md"
              title="Go Documentation"
            >
              <BookOpen size={20} strokeWidth={1.5} />
            </a>
            <a 
              href="https://github.com/luis198755/go_playGround_Json" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-gray-400 hover:text-white transition-colors p-2 rounded-md"
              title="View on GitHub"
            >
              <GithubIcon size={20} />
            </a>
          </div>
        </div>

        {/* Mobile Navigation */}
        <div className={`md:hidden ${isMenuOpen ? 'block' : 'hidden'}`}>
          <div className="px-2 pt-2 pb-3 space-y-2 border-t border-gray-700">
            <button
              onClick={() => handleModeChange('go')}
              className={`w-full text-left px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                mode === 'go'
                  ? 'bg-[#007acc] text-white'
                  : 'text-gray-400 hover:text-white hover:bg-[#3d3d3d]'
              }`}
            >
              Go Playground
            </button>
            <button
              onClick={() => handleModeChange('json')}
              className={`w-full text-left px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                mode === 'json'
                  ? 'bg-[#007acc] text-white'
                  : 'text-gray-400 hover:text-white hover:bg-[#3d3d3d]'
              }`}
            >
              JSON Viewer
            </button>
            <button
              onClick={() => handleModeChange('markdown')}
              className={`w-full text-left px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                mode === 'markdown'
                  ? 'bg-[#007acc] text-white'
                  : 'text-gray-400 hover:text-white hover:bg-[#3d3d3d]'
              }`}
            >
              Markdown Editor
            </button>
            <div className="border-t border-gray-700 pt-2 space-y-2">
              <a 
                href="https://go.dev" 
                target="_blank" 
                rel="noopener noreferrer"
                className="flex items-center gap-2 px-4 py-2 text-gray-400 hover:text-white hover:bg-[#3d3d3d] rounded-md transition-colors text-sm"
                title="Go Documentation"
              >
                <BookOpen size={20} strokeWidth={1.5} />
                <span>Documentation</span>
              </a>
              <a 
                href="https://github.com/luis198755/go_playGround_Json" 
                target="_blank" 
                rel="noopener noreferrer"
                className="flex items-center gap-2 px-4 py-2 text-gray-400 hover:text-white hover:bg-[#3d3d3d] rounded-md transition-colors text-sm"
                title="View on GitHub"
              >
                <GithubIcon size={20} />
                <span>GitHub</span>
              </a>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};