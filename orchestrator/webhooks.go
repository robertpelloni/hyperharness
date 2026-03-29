package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type WebhookPayload struct {
	Type   string          `json:"type"`
	Source string          `json:"source"`
	Data   json.RawMessage `json:"data"`
}

// HandleBorgWebhook processes real-time synchronization routes mirroring the JULES TS implementation.
func HandleBorgWebhook(payload WebhookPayload, queue *TaskQueue, ws *TelemetrySocket) (map[string]interface{}, error) {
	source := payload.Source
	if source == "" {
		source = "unknown"
	}

	// 1. Log natively into GORM Prisma mapped model
	actionLog := KeeperLog{
		ID:        uuid.New().String(),
		SessionId: "global",
		Type:      "info",
		Message:   fmt.Sprintf("Received Borg signal: %s from %s", payload.Type, source),
		Metadata:  string(payload.Data),
	}
	if err := DB.Create(&actionLog).Error; err != nil {
		log.Printf("[Webhooks] ORM Insert skipped: %v", err)
	}

	// 2. Queue specific tasks based on the webhook type
	switch payload.Type {
	case "repo_updated", "reindex_all":
		log.Printf("[Webhooks] Borg requested automated RAG re-index...")
		queue.Enqueue("index_codebase")
	case "issue_detected":
		log.Printf("[Webhooks] Borg reported an active issue ticket. Firing check_issues...")
		queue.Enqueue("check_issues")
	case "clear_logs":
		// Direct administrative truncate request mapping TS Prisma
		DB.Where("1 = 1").Delete(&KeeperLog{})
	default:
		log.Printf("[Webhooks] Unmapped signal subtype ignored: %s", payload.Type)
	}

	// 3. Emit live telemetry notification explicitly matching DaemonEventType logic in 'shared'
	emitPayload := map[string]interface{}{
		"type":      payload.Type,
		"source":    source,
		"data":      payload.Data,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	rawJson, _ := json.Marshal(emitPayload)
	ws.Broadcast(fmt.Sprintf(`{"event": "borg_signal_received", "payload": %s}`, string(rawJson)))

	return map[string]interface{}{"success": true, "processed": true}, nil
}
