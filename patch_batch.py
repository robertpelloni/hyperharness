import sys

with open("tools/goose_opencode_kimi_parity.go", "r") as f:
    content = f.read()

batch_old = """		Execute: func(args map[string]interface{}) (string, error) {
			toolCalls, ok := args["tool_calls"].([]interface{})
			if !ok || len(toolCalls) == 0 {
				return "", fmt.Errorf("tool_calls must be a non-empty array")
			}
			if len(toolCalls) > 25 {
				toolCalls = toolCalls[:25]
			}

			var results []string
			for i, callRaw := range toolCalls {
				call, ok := callRaw.(map[string]interface{})
				if !ok {
					results = append(results, fmt.Sprintf("Call %d: invalid format", i))
					continue
				}
				toolName, _ := call["tool"].(string)
				params, _ := call["parameters"].(map[string]interface{})
				if toolName == "" {
					results = append(results, fmt.Sprintf("Call %d: missing tool name", i))
					continue
				}
				// Look up and execute the tool
				if t, ok := args["_registry"].(*Registry); ok {
					if found, ok := t.Find(toolName); ok {
						result, err := found.Execute(params)
						if err != nil {
							results = append(results, fmt.Sprintf("Call %d (%s): error: %v", i, toolName, err))
						} else {
							results = append(results, fmt.Sprintf("Call %d (%s): %s", i, toolName, result))
						}
					} else {
						results = append(results, fmt.Sprintf("Call %d: tool %q not found", i, toolName))
					}
				} else {
					results = append(results, fmt.Sprintf("Call %d (%s): [batch execution requires registry context]", i, toolName))
				}
			}

			return strings.Join(results, "\\n\\n---\\n\\n"), nil
		}"""

batch_new = """		Execute: func(args map[string]interface{}) (string, error) {
			toolCalls, ok := args["tool_calls"].([]interface{})
			if !ok || len(toolCalls) == 0 {
				return "", fmt.Errorf("tool_calls must be a non-empty array")
			}
			if len(toolCalls) > 25 {
				toolCalls = toolCalls[:25]
			}

			regRaw, hasReg := args["_registry"]
			if !hasReg {
				return "", fmt.Errorf("batch execution requires _registry context injection")
			}
			reg, ok := regRaw.(*Registry)
			if !ok {
				return "", fmt.Errorf("invalid registry context")
			}

			type batchResult struct {
				index  int
				result string
			}

			resultCh := make(chan batchResult, len(toolCalls))
			var wg sync.WaitGroup

			for i, callRaw := range toolCalls {
				wg.Add(1)
				go func(idx int, raw interface{}) {
					defer wg.Done()

					call, ok := raw.(map[string]interface{})
					if !ok {
						resultCh <- batchResult{idx, fmt.Sprintf("Call %d: invalid format", idx)}
						return
					}
					toolName, _ := call["tool"].(string)
					params, _ := call["parameters"].(map[string]interface{})
					if toolName == "" {
						resultCh <- batchResult{idx, fmt.Sprintf("Call %d: missing tool name", idx)}
						return
					}

					// For safety, re-inject the registry into params if sub-tools need it
					if params == nil {
						params = make(map[string]interface{})
					}
					params["_registry"] = reg

					if found, ok := reg.Find(toolName); ok {
						res, err := found.Execute(params)
						if err != nil {
							resultCh <- batchResult{idx, fmt.Sprintf("Call %d (%s): error: %v", idx, toolName, err)}
						} else {
							resultCh <- batchResult{idx, fmt.Sprintf("Call %d (%s): %s", idx, toolName, res)}
						}
					} else {
						resultCh <- batchResult{idx, fmt.Sprintf("Call %d: tool %q not found", idx, toolName)}
					}
				}(i, callRaw)
			}

			wg.Wait()
			close(resultCh)

			// Collect results and sort by original index
			results := make([]string, len(toolCalls))
			for res := range resultCh {
				results[res.index] = res.result
			}

			return strings.Join(results, "\\n\\n---\\n\\n"), nil
		}"""

content = content.replace(batch_old.replace("\\n", "\n"), batch_new.replace("\\n", "\n"))

with open("tools/goose_opencode_kimi_parity.go", "w") as f:
    f.write(content)
print("Updated batch tool in goose_opencode_kimi_parity.go")
