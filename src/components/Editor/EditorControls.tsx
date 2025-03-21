import { Play, Share2, HelpCircle } from 'lucide-react';
import { EditorMode } from '../../types';

interface EditorControlsProps {
  mode: EditorMode;
  isLoading: boolean;
  onRun: () => void;
  onShare: () => void;
  onToggleManual: () => void;
  isManualVisible: boolean;
}

export function EditorControls({
  mode,
  isLoading,
  onRun,
  onShare,
  onToggleManual,
  isManualVisible,
}: EditorControlsProps) {
  const getRunButtonText = () => {
    if (isLoading) return 'Running...';
    switch (mode) {
      case 'json': return 'Format JSON';
      case 'markdown': return 'Preview';
      default: return 'Run';
    }
  };

  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 w-full">
      <button
        onClick={onRun}
        disabled={isLoading}
        className={`px-4 py-2 rounded flex items-center gap-2 justify-center text-sm font-medium transition-colors ${
          isLoading 
            ? 'bg-gray-600 cursor-not-allowed' 
            : 'bg-[#007acc] hover:bg-[#0066aa]'
        }`}
      >
        <Play size={16} />
        {getRunButtonText()}
      </button>
      
      <button 
        onClick={onShare}
        className="px-4 py-2 rounded border border-gray-700 hover:bg-gray-700 flex items-center gap-2 justify-center text-sm font-medium transition-colors"
        title="Share on X (Twitter)"
      >
        <Share2 size={16} />
        Share on X
      </button>
      
      <button 
        onClick={onToggleManual}
        className={`px-4 py-2 rounded border border-gray-700 hover:bg-gray-700 flex items-center gap-2 justify-center text-sm font-medium transition-colors ${isManualVisible ? 'bg-gray-700' : ''}`}
        title="Toggle User Manual"
      >
        <HelpCircle size={16} />
        {isManualVisible ? 'Hide Manual' : 'Show Manual'}
      </button>
    </div>
  );
}