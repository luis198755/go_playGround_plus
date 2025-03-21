import React, { useState, useEffect, useRef } from 'react';
import { Navbar } from './components/Navbar/Navbar';
import { Tabs } from './components/Editor/Tabs';
import Editor from '@monaco-editor/react';
import { Play, Share2, HelpCircle, Terminal as TerminalIcon, Download, Copy, FileX, Menu as MenuIcon, X as XIcon, FolderInput } from 'lucide-react';
import { TerminalOutput } from './components/Terminal/TerminalOutput';
import { EditorSettings } from './components/Editor/EditorSettings';
import { UserManual } from './components/UserManual/UserManual';
import { Footer } from './components/Footer/Footer';
import goExamples from './codes/examples_go.json';
import jsonExamples from './codes/examples_json.json';
import { EditorSettings as EditorSettingsType, Tab, EditorMode } from './types';
import { useEditor } from './hooks/useEditor';
import { useTerminal } from './hooks/useTerminal';
import { EditorHeader } from './components/Editor/EditorHeader';

function App() {
  const [isManualVisible, setIsManualVisible] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [editorSettings, setEditorSettings] = useState<EditorSettingsType>({
    fontSize: 14,
    theme: 'vs-dark',
    minimap: false,
    wordWrap: 'on',
    isSettingsOpen: false,
  });

  const {
    mode,
    tabs,
    activeTabId,
    selectedExamples,
    getCurrentExamples,
    getFileExtension,
    handleModeChange,
    setTabs,
    setActiveTabId,
    setSelectedExamples,
  } = useEditor();

  const {
    output,
    isLoading,
    setOutput,
    handleRunCode,
  } = useTerminal();

  useEffect(() => {
    setOutput('');
  }, [mode]);

  const handleModeChangeWithClear = (newMode: EditorMode) => {
    setOutput('');
    handleModeChange(newMode);
  };

  const activeTab = tabs.find(tab => tab.id === activeTabId);

  const handleCodeChange = (value: string | undefined) => {
    const newCode = value || '';
    setTabs(prev => prev.map(tab => 
      tab.id === activeTabId ? { ...tab, code: newCode } : tab
    ));
  };

  const handleRunAndUpdateCode = async () => {
    if (!activeTab) return;
    const formattedCode = await handleRunCode(activeTab.code, activeTab.mode);
    if (formattedCode !== activeTab.code) {
      handleCodeChange(formattedCode);
    }
  };

  const handleNewTab = () => {
    const newId = `tab-${Date.now()}`;
    const extension = getFileExtension(mode);
    const defaultCode = mode === 'json' 
      ? '{\n  \n}' 
      : mode === 'markdown'
      ? '# New Document\n\nStart writing here...'
      : 'package main\n\nfunc main() {\n\t\n}';

    const newTab = {
      id: newId,
      title: `untitled.${extension}`,
      code: defaultCode,
      isExample: false,
      mode: mode
    };
    
    setTabs(prev => [...prev, newTab]);
    setActiveTabId(newId);
    setSelectedExamples(prev => ({ ...prev, [mode]: 'user-code' }));
  };

  const handleCloseTab = (tabId: string) => {
    if (tabs.filter(tab => tab.mode === mode).length === 1) return;

    const tabIndex = tabs.findIndex(tab => tab.id === tabId);
    if (tabIndex === -1) return;

    if (activeTabId === tabId) {
      const currentModeTabs = tabs.filter(tab => tab.mode === mode);
      const currentTabIndex = currentModeTabs.findIndex(tab => tab.id === tabId);
      const nextTab = currentModeTabs[currentTabIndex + 1] || currentModeTabs[currentTabIndex - 1];
      setActiveTabId(nextTab.id);
    }

    setTabs(prev => prev.filter(tab => tab.id !== tabId));
  };

  const handleExampleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const newExample = event.target.value;
    
    if (newExample === 'user-code') {
      handleNewTab();
      return;
    }

    const examples = getCurrentExamples();
    const newId = `tab-${Date.now()}`;
    const extension = getFileExtension(mode);
    const exampleTab = {
      id: newId,
      title: `${newExample}.${extension}`,
      code: examples[newExample].code,
      isExample: true,
      exampleKey: newExample,
      mode: mode
    };

    setTabs(prev => [...prev, exampleTab]);
    setActiveTabId(newId);
    setSelectedExamples(prev => ({
      ...prev,
      [mode]: newExample
    }));
    setOutput('');
  };

  const handleCopyCode = () => {
    if (!activeTab) return;
    navigator.clipboard.writeText(activeTab.code);
  };

  const handleShareOnX = () => {
    const text = 'Check out my Code Playground! 🚀\n';
    const url = window.location.href;
    const shareUrl = `https://twitter.com/intent/tweet?text=${encodeURIComponent(text)}&url=${encodeURIComponent(url)}`;
    window.open(shareUrl, '_blank', 'width=550,height=420');
  };

  const handleDownloadCode = () => {
    if (!activeTab) return;
    const blob = new Blob([activeTab.code], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = activeTab.title;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const handleFileImport = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (!files || files.length === 0) return;

    const file = files[0];
    const fileExtension = file.name.split('.').pop()?.toLowerCase();
    const validExtension = getFileExtension(mode);

    if (fileExtension !== validExtension) {
      setOutput(`Error: Please import a ${mode.toUpperCase()} file in ${mode.toUpperCase()} mode`);
      return;
    }

    try {
      const content = await file.text();
      const newId = `tab-${Date.now()}`;
      const newTab = {
        id: newId,
        title: file.name,
        code: content,
        isExample: false,
        mode: mode
      };
      
      setTabs(prev => [...prev, newTab]);
      setActiveTabId(newId);
      setSelectedExamples(prev => ({ ...prev, [mode]: 'user-code' }));
      setOutput(`Successfully imported ${file.name}`);
    } catch (error) {
      setOutput(`Error reading file: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }

    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const examples = getCurrentExamples();
  const selectedExample = selectedExamples[mode];

  return (
    <div className="min-h-screen bg-[#1e1e1e] text-white flex flex-col">
      <Navbar mode={mode} onModeChange={handleModeChangeWithClear} />
      <div className="container mx-auto p-4 pt-6 flex-1">
        <div className="bg-[#2d2d2d] rounded-lg overflow-hidden flex flex-col h-[calc(100vh-12rem)]">
          <div className="flex flex-col md:flex-row flex-1 gap-4 overflow-hidden">
            {/* Editor Section */}
            <div className="w-full md:w-1/2 p-2 flex flex-col min-h-[500px] overflow-hidden">
              <EditorHeader
                mode={mode}
                settings={editorSettings}
                onSettingsChange={setEditorSettings}
                onImport={handleFileImport}
                onCopy={handleCopyCode}
                onDownload={handleDownloadCode}
                onClear={() => {
                  handleCodeChange('');
                  setOutput('');
                }}
                fileInputRef={fileInputRef}
              />
              <Tabs
                tabs={tabs}
                activeTabId={activeTabId}
                onTabChange={setActiveTabId}
                onTabClose={handleCloseTab}
                onNewTab={handleNewTab}
                mode={mode}
              />
              {activeTab && (
                <div className="flex-1 min-h-0">
                  <Editor
                    height="100%"
                    language={mode}
                    value={activeTab.code}
                    theme={editorSettings.theme}
                    onChange={handleCodeChange}
                    options={{
                      fontSize: editorSettings.fontSize,
                      minimap: { enabled: editorSettings.minimap },
                      wordWrap: editorSettings.wordWrap,
                      scrollBeyondLastLine: false,
                    }}
                  />
                </div>
              )}
            </div>

            {/* Terminal Output Section */}
            <div className="w-full md:w-1/2 p-2 flex flex-col min-h-[500px] overflow-hidden">
              <TerminalOutput
                output={output}
                isLoading={isLoading}
                onClear={() => setOutput('')}
                isMarkdown={mode === 'markdown'}
              />
            </div>
          </div>

          {/* Control Bar */}
          <div className="bg-[#2d2d2d] p-4 flex flex-col gap-6 border-t border-gray-700">
            <div className="w-full">
              <select 
                className="bg-[#1e1e1e] text-white px-4 py-2.5 rounded border border-gray-700 w-full text-sm"
                value={selectedExample}
                onChange={handleExampleChange}
              >
                <option value="user-code">User Code</option>
                {Object.entries(examples).map(([key, example]) => (
                  <option key={key} value={key}>
                    {example.name}
                  </option>
                ))}
              </select>
            </div>
            
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 w-full">
              <button
                onClick={handleRunAndUpdateCode}
                disabled={isLoading}
                className={`px-4 py-2.5 rounded flex items-center gap-2 justify-center text-sm font-medium transition-colors ${
                  isLoading 
                    ? 'bg-gray-600 cursor-not-allowed' 
                    : 'bg-[#007acc] hover:bg-[#0066aa]'
                }`}
              >
                <Play size={16} />
                {isLoading ? 'Running...' : mode === 'json' ? 'Format JSON' : 'Run'}
              </button>
              
              <button 
                onClick={handleShareOnX}
                className="px-4 py-2.5 rounded border border-gray-700 hover:bg-gray-700 flex items-center gap-2 justify-center text-sm font-medium transition-colors"
                title="Share on X (Twitter)"
              >
                <Share2 size={16} />
                Share on X
              </button>
              
              <button 
                onClick={() => setIsManualVisible(!isManualVisible)}
                className={`px-4 py-2.5 rounded border border-gray-700 hover:bg-gray-700 flex items-center gap-2 justify-center text-sm font-medium transition-colors ${isManualVisible ? 'bg-gray-700' : ''}`}
                title="Toggle User Manual"
              >
                <HelpCircle size={16} />
                {isManualVisible ? 'Hide Manual' : 'Show Manual'}
              </button>
            </div>
          </div>
        </div>

        <UserManual isVisible={isManualVisible} />
        <Footer />
      </div>
    </div>
  );
}

export default App;