import sys

with open("tools/goose_opencode_kimi_parity.go", "r") as f:
    content = f.read()

# Add subagents import if missing
if '"github.com/robertpelloni/hyperharness/internal/subagents"' not in content:
    content = content.replace(
        '"github.com/robertpelloni/hyperharness/internal/extensions"',
        '"github.com/robertpelloni/hyperharness/internal/extensions"\n\t"github.com/robertpelloni/hyperharness/internal/subagents"'
    )

delegate_old = """		Execute: func(args map[string]interface{}) (string, error) {
			task, _ := args["task"].(string)
			description, _ := args["description"].(string)
			subagentType, _ := args["subagent_type"].(string)
			if subagentType == "" {
				subagentType = "code"
			}
			mode, _ := args["mode"].(string)
			if mode == "" {
				mode = "sync"
			}

			taskID := fmt.Sprintf("task_%d", time.Now().UnixNano())

			if mode == "async" {
				// Store the task for async execution
				GlobalJobManager.mu.Lock()
				GlobalJobManager.jobs[taskID] = &Job{
					ID:          taskID,
					Command:     fmt.Sprintf("[delegate] %s", description),
					Description: task,
					StartTime:   time.Now(),
					Output:      &strings.Builder{},
				}
				GlobalJobManager.mu.Unlock()

				return fmt.Sprintf("Delegated task %s to %s subagent (async)\\nTask ID: %s\\nUse job_output to check progress.", description, subagentType, taskID), nil
			}

			return fmt.Sprintf("Delegated task '%s' to %s subagent.\\nTask ID: %s\\nStatus: completed synchronously.\\n\\n<task_result>\\nTask executed: %s\\n</task_result>", description, subagentType, taskID, task), nil
		}"""

delegate_new = """		Execute: func(args map[string]interface{}) (string, error) {
			task, _ := args["task"].(string)
			description, _ := args["description"].(string)
			subagentType, _ := args["subagent_type"].(string)
			if subagentType == "" {
				subagentType = "code"
			}
			mode, _ := args["mode"].(string)
			if mode == "" {
				mode = "sync"
			}

			sType := subagents.TypeResearch
			switch strings.ToLower(subagentType) {
			case "plan":
				sType = subagents.TypePlan
			case "research":
				sType = subagents.TypeResearch
			case "test":
				sType = subagents.TypeTest
			case "code":
				sType = subagents.TypeCode
			case "debug":
				sType = subagents.TypeDebug
			case "review":
				sType = subagents.TypeReview
			}

			subtask := subagents.NewManager().CreateTask(sType, description, task, "")
			taskID := subtask.ID

			if mode == "async" {
				// We won't block, so we wrap it inside a background Job
				// Note: A real implementation would run the subagent in a goroutine
				GlobalJobManager.mu.Lock()
				job := &Job{
					ID:          taskID,
					Command:     fmt.Sprintf("[delegate] %s", description),
					Description: task,
					StartTime:   time.Now(),
					Output:      &strings.Builder{},
					Status:      "running",
				}
				GlobalJobManager.jobs[taskID] = job
				GlobalJobManager.mu.Unlock()

				go func() {
					result, _ := subagents.NewManager().ExecuteTask(context.Background(), subtask)
					job.mu.Lock()
					job.Output.WriteString(result)
					job.Status = "completed"
					job.EndTime = time.Now()
					job.mu.Unlock()
				}()

				return fmt.Sprintf("Delegated task %s to %s subagent (async)\\nTask ID: %s\\nUse job_output to check progress.", description, subagentType, taskID), nil
			}

			result, err := subagents.NewManager().ExecuteTask(context.Background(), subtask)
			if err != nil {
				return "", fmt.Errorf("failed to delegate task: %w", err)
			}

			return fmt.Sprintf("Delegated task '%s' to %s subagent.\\nTask ID: %s\\nStatus: completed synchronously.\\n\\n<task_result>\\n%s\\n</task_result>", description, subagentType, taskID, result), nil
		}"""

content = content.replace(delegate_old.replace("\\n", "\n"), delegate_new.replace("\\n", "\n"))

with open("tools/goose_opencode_kimi_parity.go", "w") as f:
    f.write(content)
print("Updated delegate")
