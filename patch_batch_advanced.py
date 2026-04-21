import sys

with open("tools/advanced_parity.go", "r") as f:
    content = f.read()

batch_old = """		Execute: func(args map[string]interface{}) (string, error) {
			invocations, _ := args["invocations"].([]interface{})
			if len(invocations) == 0 {
				return "", fmt.Errorf("invocations array is empty or missing")
			}
			return fmt.Sprintf("[batch] %d invocations queued (requires runtime integration)", len(invocations)), nil
		}"""

batch_new = """		Execute: func(args map[string]interface{}) (string, error) {
			invocations, _ := args["invocations"].([]interface{})
			if len(invocations) == 0 {
				return "", fmt.Errorf("invocations array is empty or missing")
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

			resultCh := make(chan batchResult, len(invocations))
			var wg sync.WaitGroup

			for i, callRaw := range invocations {
				wg.Add(1)
				go func(idx int, raw interface{}) {
					defer wg.Done()

					call, ok := raw.(map[string]interface{})
					if !ok {
						resultCh <- batchResult{idx, fmt.Sprintf("Call %d: invalid format", idx)}
						return
					}
					toolName, _ := call["tool_name"].(string)
					params, _ := call["parameters"].(map[string]interface{})
					if toolName == "" {
						resultCh <- batchResult{idx, fmt.Sprintf("Call %d: missing tool name", idx)}
						return
					}

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

			results := make([]string, len(invocations))
			for res := range resultCh {
				results[res.index] = res.result
			}

			return strings.Join(results, "\\n\\n---\\n\\n"), nil
		}"""

content = content.replace(batch_old.replace("\\n", "\n"), batch_new.replace("\\n", "\n"))

# Ensure sync is imported
if '"sync"' not in content:
    content = content.replace('"strings"', '"strings"\n\t"sync"')

with open("tools/advanced_parity.go", "w") as f:
    f.write(content)
print("Updated batch tool in advanced_parity.go")
