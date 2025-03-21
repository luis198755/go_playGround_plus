import { Editor } from '@monaco-editor/react';
import { EditorSettings } from '../../types';

interface CodeEditorProps {
  code: string;
  language: string;
  settings: EditorSettings;
  onChange: (value: string | undefined) => void;
}

export function CodeEditor({ code, language, settings, onChange }: CodeEditorProps) {
  return (
    <div className="h-[calc(100vh-24rem)]">
      <Editor
        height="100%"
        language={language}
        theme={settings.theme}
        value={code}
        onChange={onChange}
        options={{
          minimap: { enabled: settings.minimap },
          fontSize: settings.fontSize,
          scrollBeyondLastLine: false,
          lineNumbers: 'on',
          roundedSelection: false,
          automaticLayout: true,
          wordWrap: settings.wordWrap,
          formatOnPaste: true,
          formatOnType: true,
        }}
      />
    </div>
  );
}