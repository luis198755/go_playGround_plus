import React from 'react';
import { HelpCircle } from 'lucide-react';
import { UserManualProps } from '../../types';

export function UserManual({ isVisible }: UserManualProps) {
  return (
    <div className={`bg-[#2d2d2d] rounded-lg p-6 mt-6 border border-gray-700 transition-all duration-300 ${isVisible ? 'opacity-100 max-h-[1000px]' : 'opacity-0 max-h-0 overflow-hidden p-0'}`}>
      <h2 className="text-xl font-semibold mb-4 text-white flex items-center gap-2">
        <HelpCircle size={20} />
        User Manual
      </h2>
      
      <div className="space-y-6">
        <section>
          <h3 className="text-lg font-medium mb-2 text-gray-300">Code Editor</h3>
          <ul className="list-disc list-inside space-y-2 text-gray-400">
            <li><span className="text-white">Predefined Examples:</span> Select Go code examples from the dropdown menu.</li>
            <li><span className="text-white">Editor Settings:</span> Customize theme, font size, minimap, and word wrap using the settings icon.</li>
            <li><span className="text-white">Layout:</span> Switch between horizontal/vertical view automatically.</li>
            <li><span className="text-white">Code Management:</span> Use buttons to copy, download, or create new code.</li>
          </ul>
        </section>

        <section>
          <h3 className="text-lg font-medium mb-2 text-gray-300">Terminal and Execution</h3>
          <ul className="list-disc list-inside space-y-2 text-gray-400">
            <li><span className="text-white">Run Code:</span> Press the 'Run' button to execute your Go code.</li>
            <li><span className="text-white">Terminal Output:</span> View your code's output and error messages.</li>
            <li><span className="text-white">Output Management:</span> Clear, copy, or download terminal output using header buttons.</li>
          </ul>
        </section>

        <section>
          <h3 className="text-lg font-medium mb-2 text-gray-300">Additional Features</h3>
          <ul className="list-disc list-inside space-y-2 text-gray-400">
            <li><span className="text-white">Share:</span> Share your code on X (Twitter) using the 'Share on X' button.</li>
            <li><span className="text-white">Themes:</span> Choose between light, dark, and high contrast themes for better visibility.</li>
            <li><span className="text-white">Minimap:</span> Toggle the code miniature view for better navigation.</li>
          </ul>
        </section>
      </div>
    </div>
  );
}