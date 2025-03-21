import { Menu, X } from 'lucide-react';
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
          <div className="hidden md:flex items-center space-x-6">
            <a 
              href="https://go.dev" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-gray-400 hover:text-white transition-colors text-sm"
            >
              Go Docs
            </a>
            <a 
              href="https://github.com/luis198755/go_playGround_Json" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-gray-400 hover:text-white transition-colors text-sm"
            >
              GitHub
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
                className="block px-4 py-2 text-gray-400 hover:text-white hover:bg-[#3d3d3d] rounded-md transition-colors text-sm"
              >
                Go Docs
              </a>
              <a 
                href="https://github.com/luis198755/go_playGround_Json" 
                target="_blank" 
                rel="noopener noreferrer"
                className="block px-4 py-2 text-gray-400 hover:text-white hover:bg-[#3d3d3d] rounded-md transition-colors text-sm"
              >
                GitHub
              </a>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};