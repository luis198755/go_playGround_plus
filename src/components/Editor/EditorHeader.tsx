import { Copy, Download, FileX, FolderInput } from 'lucide-react';
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
  return (
    <div className="flex flex-col space-y-2 mb-2">
      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-400 font-medium">
          {mode === 'json' ? 'JSON Viewer' : mode === 'markdown' ? 'Markdown Editor' : 'Code Editor'}
        </div>
        <div className="flex items-center space-x-1">
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
        </div>
      </div>
    </div>
  );
}