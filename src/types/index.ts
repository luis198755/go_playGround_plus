export interface EditorSettings {
  fontSize: number;
  theme: string;
  minimap: boolean;
  wordWrap: 'on' | 'off';
  isSettingsOpen: boolean;
}

export interface Tab {
  id: string;
  title: string;
  code: string;
  isExample: boolean;
  exampleKey?: string;
  mode: EditorMode;
}

export interface CodeExample {
  name: string;
  code: string;
  mode: EditorMode;
}

export type EditorMode = 'go' | 'json' | 'markdown';

export interface TerminalOutputProps {
  output: string;
  isLoading: boolean;
  onClear: () => void;
  isMarkdown?: boolean;
}

export interface UserManualProps {
  isVisible: boolean;
}