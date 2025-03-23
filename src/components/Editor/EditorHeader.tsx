import { Copy, Download, FileX, FolderInput, Menu, X } from 'lucide-react';
import { useState, useRef, useEffect } from 'react';
import { EditorSettings } from './EditorSettings';
import { EditorMode, EditorSettings as EditorSettingsType } from '../../types';

interface EditorHeaderProps {
  mode: EditorMode;
  settings: EditorSettingsType;
  onSettingsChange: (settings: EditorSettingsType) => void;
  onImport: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onCopy: () => void;
  onDownload: () => void;
  onClear: () => void;
  fileInputRef: React.RefObject<HTMLInputElement>;
}

export function EditorHeader({
  mode,
  settings,
  onSettingsChange,
  onImport,
  onCopy,
  onDownload,
  onClear,
  fileInputRef,
}: EditorHeaderProps) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className="flex flex-col space-y-2 mb-2 relative">
      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-400 font-medium">
          {mode === 'json' ? 'JSON Viewer' : mode === 'markdown' ? 'Markdown Editor' : 'Code Editor'}
        </div>
        <div className="flex items-center space-x-1">
          {/* Mobile Menu Button */}
          <button 
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            className="sm:hidden p-1.5 rounded hover:bg-[#3d3d3d] text-gray-400 hover:text-white transition-colors"
          >
            {isMenuOpen ? <X size={16} /> : <Menu size={16} />}
          </button>

          {/* Mobile Menu */}
          {isMenuOpen && (
            <div ref={menuRef} className="absolute right-0 top-full mt-1 bg-[#2d2d2d] rounded-lg shadow-lg border border-gray-700 py-2 px-1 sm:hidden z-50">
              <button 
                onClick={() => {
                  fileInputRef.current?.click();
                  setIsMenuOpen(false);
                }}
                className="w-full p-2 rounded hover:bg-[#3d3d3d] text-gray-400 hover:text-white transition-colors flex items-center gap-2"
                title="Import File"
              >
                <FolderInput size={16} />
                <span>Import File</span>
              </button>
              <button 
                onClick={() => {
                  onCopy();
                  setIsMenuOpen(false);
                }}
                className="w-full p-2 rounded hover:bg-[#3d3d3d] text-gray-400 hover:text-white transition-colors flex items-center gap-2"
                title="Copy Code"
              >
                <Copy size={16} />
                <span>Copy Code</span>
              </button>
              <button 
                onClick={() => {
                  onDownload();
                  setIsMenuOpen(false);
                }}
                className="w-full p-2 rounded hover:bg-[#3d3d3d] text-gray-400 hover:text-white transition-colors flex items-center gap-2"
                title="Download Code"
              >
                <Download size={16} />
                <span>Download Code</span>
              </button>
              <button 
                onClick={() => {
                  onClear();
                  setIsMenuOpen(false);
                }}
                className="w-full p-2 rounded hover:bg-[#3d3d3d] text-gray-400 hover:text-white transition-colors flex items-center gap-2"
                title="Clear Code and Terminal"
              >
                <FileX size={16} />
                <span>Clear Code</span>
              </button>
            </div>
          )}

          {/* Desktop Menu */}
          <div className="hidden sm:flex items-center space-x-1 border-r border-gray-700 pr-2 mr-1">
            <button 
              onClick={() => fileInputRef.current?.click()}
              className="p-1.5 rounded hover:bg-[#3d3d3d] text-gray-400 hover:text-white transition-colors"
              title="Import File"
            >
              <FolderInput size={16} />
            </button>
            <input
              ref={fileInputRef}
              type="file"
              accept={mode === 'go' ? '.go' : mode === 'json' ? '.json' : '.md'}
              onChange={onImport}
              className="hidden"
            />
            <button 
              onClick={onCopy}
              className="p-1.5 rounded hover:bg-[#3d3d3d] text-gray-400 hover:text-white transition-colors"
              title="Copy Code"
            >
              <Copy size={16} />
            </button>
            <button 
              onClick={onDownload}
              className="p-1.5 rounded hover:bg-[#3d3d3d] text-gray-400 hover:text-white transition-colors"
              title="Download Code"
            >
              <Download size={16} />
            </button>
          </div>
          <div className="hidden sm:flex items-center space-x-1">
            <button 
              onClick={onClear}
              className="p-1.5 rounded hover:bg-[#3d3d3d] text-gray-400 hover:text-white transition-colors"
              title="Clear Code and Terminal"
            >
              <FileX size={16} />
            </button>
            <EditorSettings 
              settings={settings}
              onSettingsChange={onSettingsChange}
            />
          </div>

          {/* Mobile Settings */}
          <div className="sm:hidden">
            <EditorSettings 
              settings={settings}
              onSettingsChange={onSettingsChange}
            />
          </div>
        </div>
      </div>
    </div>
  );
}