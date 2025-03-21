import { useState } from 'react';
import { marked } from 'marked';

export function useTerminal() {
  const [output, setOutput] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleRunCode = async (code: string, mode: string) => {
    if (mode === 'json') {
      try {
        const parsedJSON = JSON.parse(code);
        const formatted = JSON.stringify(parsedJSON, null, 2);
        setOutput('JSON formatted successfully');
        return formatted; // Return formatted JSON to update editor content
      } catch (error) {
        setOutput(error instanceof Error ? error.message : 'Invalid JSON');
        return code; // Return original code if parsing fails
      }
    }

    if (mode === 'markdown') {
      try {
        const html = marked(code, { breaks: true });
        setOutput(html);
        return code;
      } catch (error) {
        setOutput(`Error parsing Markdown: ${error instanceof Error ? error.message : 'Unknown error'}`);
        return code;
      }
    }

    setIsLoading(true);
    setOutput('');
    
    try {
      const response = await fetch('/api/execute', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ code }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body?.getReader();
      if (!reader) {
        throw new Error('ReadableStream not supported in this browser.');
      }

      const decoder = new TextDecoder();
      let done = false;
      while (!done) {
        const { value, done: doneReading } = await reader.read();
        done = doneReading;
        if (value) {
          const chunk = decoder.decode(value);
          setOutput(prev => prev + chunk);
        }
      }
    } catch (error) {
      setOutput(`Failed to execute code: ${error instanceof Error ? error.message : 'Unknown error'}`);
    } finally {
      setIsLoading(false);
    }

    return code;
  };

  return {
    output,
    isLoading,
    setOutput,
    handleRunCode,
  };
}