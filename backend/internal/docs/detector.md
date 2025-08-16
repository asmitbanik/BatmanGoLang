# detector.go

This file implements the core detection logic for code analysis, including:
- Security analysis (static checks for deprecated APIs, insecure patterns, hardcoded secrets, and known CVEs)
- Purpose guessing (using ML/keyword/context rules to infer what a code snippet does)
- Complexity and style analysis (cyclomatic complexity, style metrics)

## Key Functions

### `AnalyzeSecurity(code, language, framework string) []map[string]interface{}`
Scans the code for known insecure patterns, deprecated APIs, and hardcoded secrets. It uses language-specific and generic rules, e.g.:
- For Python: detects use of `pickle.load`, `os.system`, old PyYAML, and insecure requests versions.
- For Go: flags deprecated `ioutil.ReadAll`, use of `os/exec`.
- For JS/PHP: flags `eval()`, `child_process.exec`, and other dangerous patterns.
- Generic: flags hardcoded secrets like `password=`, `api_key=`, etc.
Returns a list of issues with severity and CVE where possible.

### `GuessPurpose(code, language, framework string) string`
Uses keyword and context rules to guess the purpose of a snippet, e.g.:
- Detects authentication logic, API handlers, database access, CLI tools, data processing, web scraping, tests, config loaders, logging, entry points, and more.
- Uses both code content and detected framework for more accurate guesses.

### `AnalyzeComplexityAndStyle(code, language string) (int, float64)`
Estimates cyclomatic complexity and a style score:
- Counts long lines, comment lines, and function definitions.
- For Go, counts control flow keywords (`if`, `for`, `case`, etc.) to estimate complexity.
- Style score is based on comment density, line length, and other heuristics.

## How it fits in the backend
These functions are called by the API handler (`/analyze` endpoint) to:
- Provide security warnings for uploaded code
- Guess the code's purpose for the user
- Show code quality metrics
This makes the backend more than just a language detector—it provides real, actionable insights for users and demonstrates advanced static analysis techniques.
