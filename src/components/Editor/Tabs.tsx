import { X } from 'lucide-react';
import { Tab, EditorMode } from '../../types';

interface TabsProps {
  tabs: Tab[];
  activeTabId: string;
  onTabChange: (tabId: string) => void;
  onTabClose: (tabId: string) => void;
  onNewTab: () => void;
  mode: EditorMode;
}

export const Tabs = ({ tabs, activeTabId, onTabChange, onTabClose, onNewTab, mode }: TabsProps) => {
  // Filter tabs by current mode
  const filteredTabs = tabs.filter(tab => tab.mode === mode);

  return (
    <div className="flex items-center border-b border-gray-700 bg-[#252526] overflow-x-auto">
      {filteredTabs.map((tab) => (
        <div
          key={tab.id}
          className={`group flex items-center min-w-[120px] max-w-[200px] h-9 px-3 
            ${activeTabId === tab.id ? 'bg-[#1e1e1e] text-white' : 'bg-[#2d2d2d] text-gray-400'} 
            cursor-pointer border-r border-gray-700 hover:text-white transition-colors`}
          onClick={() => onTabChange(tab.id)}
        >
          <span className="truncate flex-1 text-sm">{tab.title}</span>
          <button
            className={`flex items-center justify-center ml-2 p-1 rounded opacity-0 group-hover:opacity-100 
              ${activeTabId === tab.id ? 'hover:bg-[#3d3d3d]' : 'hover:bg-[#4d4d4d]'} 
              text-gray-400 hover:text-white transition-all`}
            onClick={(e) => {
              e.stopPropagation();
              onTabClose(tab.id);
            }}
            title="Close tab"
          >
            <X size={14} />
          </button>
        </div>
      ))}
      <button
        onClick={onNewTab}
        className="flex items-center justify-center w-9 h-9 bg-[#2d2d2d] text-gray-400 
          hover:text-white transition-colors border-r border-gray-700"
        title="New tab"
      >
        +
      </button>
    </div>
  );
};