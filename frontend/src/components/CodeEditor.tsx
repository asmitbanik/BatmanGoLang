import React from 'react'
import MonacoEditor from 'react-monaco-editor'

type Props = {
  code: string
  setCode: (s: string) => void
  language?: string
}

export default function CodeEditor({ code, setCode, language = 'javascript' }: Props) {
  const options = {
    selectOnLineNumbers: true,
    automaticLayout: true
  }

  return (
    <div className="editor">
      <MonacoEditor
        width="100%"
        height="480"
        language={language}
        theme="vs-dark"
        value={code}
        options={options}
        onChange={(newVal) => setCode(newVal)}
        editorDidMount={(editor) => {
          editor.focus()
        }}
      />
    </div>
  )
}
