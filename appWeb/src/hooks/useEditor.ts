import { useState } from 'react';
import { EditorMode, Tab } from '../types';
import goExamples from '../codes/examples_go.json';
import jsonExamples from '../codes/examples_json.json';
import markdownExamples from '../codes/examples_markdown.json';

export const useEditor = () => {
  const [mode, setMode] = useState<EditorMode>('go');
  const [selectedExamples, setSelectedExamples] = useState<Record<EditorMode, string>>({
    go: 'hello-world',
    json: Object.keys(jsonExamples)[0],
    markdown: Object.keys(markdownExamples)[0]
  });

  const [tabs, setTabs] = useState<Tab[]>([{
    id: 'hello-world',
    title: 'hello-world.go',
    code: goExamples['hello-world'].code,
    isExample: true,
    exampleKey: 'hello-world',
    mode: 'go'
  }]);

  const [activeTabId, setActiveTabId] = useState('hello-world');

  const getCurrentExamples = () => {
    switch (mode) {
      case 'json':
        return jsonExamples;
      case 'markdown':
        return markdownExamples;
      default:
        return goExamples;
    }
  };

  const getFileExtension = (tabMode: EditorMode) => {
    switch (tabMode) {
      case 'json':
        return 'json';
      case 'markdown':
        return 'md';
      default:
        return 'go';
    }
  };

  const getExamplesForMode = (mode: EditorMode) => {
    switch (mode) {
      case 'json':
        return jsonExamples;
      case 'markdown':
        return markdownExamples;
      default:
        return goExamples;
    }
  };

  const handleModeChange = (newMode: EditorMode, clearTerminal?: () => void) => {
    setMode(newMode);
    
    // Clear terminal output when changing modes
    if (clearTerminal) {
      clearTerminal();
    }
    
    const existingTab = tabs.find(tab => tab.mode === newMode);
    if (existingTab) {
      setActiveTabId(existingTab.id);
      return;
    }

    const examples = getExamplesForMode(newMode);
    const firstExampleKey = Object.keys(examples)[0];
    const exampleKey = selectedExamples[newMode] || firstExampleKey;
    
    if (!examples[exampleKey]) {
      setSelectedExamples(prev => ({
        ...prev,
        [newMode]: firstExampleKey
      }));
    }

    const extension = getFileExtension(newMode);
    const newTab = {
      id: `${firstExampleKey}-${Date.now()}`,
      title: `${firstExampleKey}.${extension}`,
      code: examples[firstExampleKey].code,
      isExample: true,
      exampleKey: firstExampleKey,
      mode: newMode
    };
    
    setTabs(prev => [...prev, newTab]);
    setActiveTabId(newTab.id);
  };

  return {
    mode,
    tabs,
    activeTabId,
    selectedExamples,
    getCurrentExamples,
    getFileExtension,
    handleModeChange,
    setMode,
    setTabs,
    setActiveTabId,
    setSelectedExamples,
  };
};