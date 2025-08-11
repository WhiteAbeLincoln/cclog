package filepicker

// Renderer defines a function that renders a given path to Markdown.
// includeAll corresponds to the CLI/TUI "--include-all" semantics
// (i.e., no filtering and show placeholders).
type Renderer func(path string, includeAll bool) (string, error)

var renderFn Renderer

// SetRenderer allows the host application to inject a renderer implementation.
// Passing nil resets to the default internal renderer.
func SetRenderer(r Renderer) { renderFn = r }

