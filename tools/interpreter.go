package tools

import (
	"sync"
	"github.com/robertpelloni/hyperharness/repl"
)

var (
	replSessions = make(map[string]*repl.Session)
	replMu       sync.RWMutex
)

func (r *Registry) registerInterpreterTools() {
	r.Tools = append(r.Tools, Tool{
		Name:        "code_interpreter",
		Description: "Executes code statefully in a persistent session. Arguments: language (string: 'python' or 'node'), code (string)",
		Execute: func(args map[string]interface{}) (string, error) {
			lang, _ := args["language"].(string)
			code, _ := args["code"].(string)

			replMu.RLock()
			session, ok := replSessions[lang]
			replMu.RUnlock()

			if !ok {
				var err error
				session, err = repl.NewSession(lang)
				if err != nil {
					return "", err
				}
				replMu.Lock()
				replSessions[lang] = session
				replMu.Unlock()
			}

			return session.Execute(code)
		},
	})
}
