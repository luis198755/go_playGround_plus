import React, { useState, useRef, useEffect } from 'react';
import { Download, Copy, Trash2, Terminal } from 'lucide-react';
import { TerminalOutputProps } from '../../types';

export function TerminalOutput({ output, isLoading, onClear, isMarkdown = false }: TerminalOutputProps) {
  const terminalRef = useRef<HTMLDivElement>(null);
  const [userInput, setUserInput] = useState('');
  const [showPrompt, setShowPrompt] = useState(false);
  const [history, setHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);

  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [output, userInput]);

  const handleDownload = () => {
    try {
      const blob = new Blob([output], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'terminal-output.txt';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Failed to download output:', error);
    }
  };

  const handleCopy = async () => {
    try {
      const textToCopy = output.split('\n').join('\n');
      await navigator.clipboard.writeText(textToCopy);
    } catch (error) {
      console.error('Failed to copy output:', error);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && userInput) {
      setHistory(prev => [...prev, userInput]);
      setHistoryIndex(-1);
      setUserInput('');
      setShowPrompt(false);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      if (historyIndex < history.length - 1) {
        const newIndex = historyIndex + 1;
        setHistoryIndex(newIndex);
        setUserInput(history[history.length - 1 - newIndex]);
      }
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (historyIndex > 0) {
        const newIndex = historyIndex - 1;
        setHistoryIndex(newIndex);
        setUserInput(history[history.length - 1 - newIndex]);
      } else if (historyIndex === 0) {
        setHistoryIndex(-1);
        setUserInput('');
      }
    }
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2 text-sm text-gray-400 font-medium">
          <Terminal size={16} />
          Terminal Output
        </div>
      </div>
      <div className="flex-1 bg-[#0C0C0C] rounded-lg overflow-hidden border border-[#333] flex flex-col">
        <div className="bg-[#1A1A1A] px-4 py-2 flex items-center gap-2 border-b border-[#333]">
          <div className="flex gap-1.5">
            <div className="w-3 h-3 rounded-full bg-[#FF5F56]"></div>
            <div className="w-3 h-3 rounded-full bg-[#FFBD2E]"></div>
            <div className="w-3 h-3 rounded-full bg-[#27C93F]"></div>
          </div>
          <div className="flex-1 text-center text-[#666] text-xs">terminal@go-playground</div>
          <div className="flex gap-2">
            <button 
              onClick={onClear}
              className="text-[#666] hover:text-white transition-colors" 
              title="Clear Terminal"
            >
              <Trash2 size={14} />
            </button>
            <button 
              onClick={handleCopy}
              className="text-[#666] hover:text-white transition-colors" 
              title="Copy Output"
            >
              <Copy size={14} />
            </button>
            <button 
              onClick={handleDownload}
              className="text-[#666] hover:text-white transition-colors" 
              title="Download Output"
            >
              <Download size={14} />
            </button>
          </div>
        </div>

        <div 
          ref={terminalRef}
          className={`flex-1 font-mono text-sm overflow-y-auto p-4 h-full bg-[#0C0C0C] ${isMarkdown ? 'markdown-preview' : ''}`}
          style={{ 
            fontFamily: "'Fira Code', 'JetBrains Mono', monospace",
            backgroundImage: 'radial-gradient(#1A1A1A 1px, transparent 1px)',
            backgroundSize: '4px 4px',
          }}
        >
          {isMarkdown ? (
            <div 
              className="text-white prose prose-invert max-w-none"
              dangerouslySetInnerHTML={{ __html: output }}
            />
          ) : (
            <div className="text-[#00FF00] space-y-1">
              {output.split('\n')
                .filter(line => line.trim())
                .map((line, i) => (
                  <div 
                    key={i} 
                    className="whitespace-pre-wrap leading-5 opacity-90 [text-shadow:_0_0_10px_rgba(0,255,0,0.4)] flex items-start"
                    style={{
                      animation: 'fadeIn 0.1s ease-in-out',
                    }}
                  >
                    <span className="opacity-70 mr-2 text-[#00FF00]">›</span>
                    <span className="flex-1">{line}</span>
                  </div>
              ))}
              {isLoading ? (
                <div className="animate-pulse flex items-center gap-1">
                  <span className="inline-block w-2 h-4 bg-[#00FF00]"></span>
                </div>
              ) : (
                <div className="flex items-center">
                  {showPrompt ? (
                    <>
                      <span className="text-[#00FF00] mr-2">$</span>
                      <input
                        type="text"
                        value={userInput}
                        onChange={(e) => setUserInput(e.target.value)}
                        onKeyDown={handleKeyDown}
                        className="flex-1 bg-transparent border-none outline-none text-[#00FF00]"
                        autoFocus
                      />
                    </>
                  ) : (
                    <div className="flex items-center gap-2">
                      <span className="opacity-70">›</span>
                      <span className="terminal-cursor" />
                    </div>
                  )}
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}