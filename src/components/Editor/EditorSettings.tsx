import { Settings } from 'lucide-react';
import { EditorSettings as EditorSettingsType } from '../../types';

interface EditorSettingsProps {
  settings: EditorSettingsType;
  onSettingsChange: (settings: EditorSettingsType) => void;
}

export function EditorSettings({ settings, onSettingsChange }: EditorSettingsProps) {
  return (
    <div className="relative">
      <button 
        onClick={() => onSettingsChange({ ...settings, isSettingsOpen: !settings.isSettingsOpen })}
        className="text-[#007acc] hover:text-white hover:bg-[#0066aa] transition-colors rounded p-2"
        title="Editor Settings"
        aria-expanded={settings.isSettingsOpen}
        aria-haspopup="true"
      >
        <Settings size={16} />
      </button>
      <div 
        className={`absolute right-0 mt-2 w-48 bg-[#1A1A1A] rounded-lg shadow-lg border border-[#333] z-10 ${settings.isSettingsOpen ? 'block' : 'hidden'}`}
        role="menu"
        aria-orientation="vertical"
        aria-labelledby="editor-settings-button"
      >
        <div className="p-2 space-y-3">
          <div>
            <label className="block text-xs text-gray-400 mb-1">Theme</label>
            <select
              value={settings.theme}
              onChange={(e) => onSettingsChange({ 
                ...settings, 
                theme: e.target.value as EditorSettingsType['theme'],
                isSettingsOpen: true
              })}
              className="w-full bg-[#1e1e1e] text-white text-xs rounded border border-gray-700 p-1"
            >
              <option value="vs-dark">Dark</option>
              <option value="light">Light</option>
              <option value="hc-black">High Contrast Dark</option>
              <option value="hc-light">High Contrast Light</option>
            </select>
          </div>
          <div>
            <label className="block text-xs text-gray-400 mb-1">Font Size</label>
            <input
              type="number"
              value={settings.fontSize}
              onChange={(e) => onSettingsChange({ 
                ...settings, 
                fontSize: Math.max(8, Math.min(32, parseInt(e.target.value) || 14)),
                isSettingsOpen: true
              })}
              className="w-full bg-[#1e1e1e] text-white text-xs rounded border border-gray-700 p-1"
              min="8"
              max="32"
            />
          </div>
          <div className="flex items-center justify-between">
            <label className="text-xs text-gray-400">Minimap</label>
            <input
              type="checkbox"
              checked={settings.minimap}
              onChange={(e) => onSettingsChange({ 
                ...settings, 
                minimap: e.target.checked,
                isSettingsOpen: true
              })}
              className="rounded border-gray-700"
            />
          </div>
          <div>
            <label className="block text-xs text-gray-400 mb-1">Word Wrap</label>
            <select
              value={settings.wordWrap}
              onChange={(e) => onSettingsChange({ 
                ...settings, 
                wordWrap: e.target.value as EditorSettingsType['wordWrap'],
                isSettingsOpen: true
              })}
              className="w-full bg-[#1e1e1e] text-white text-xs rounded border border-gray-700 p-1"
            >
              <option value="on">On</option>
              <option value="off">Off</option>
            </select>
          </div>
        </div>
      </div>
    </div>
  );
}