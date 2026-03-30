package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

var daemonRunning bool

// StartKeeperDaemon natively replicates `daemon.ts` running an asynchronous evaluation sweeping TS legacy boundaries.
func StartKeeperDaemon(queue *TaskQueue, ws *TelemetrySocket) {
	log.Println("[Daemon] Starting Native Watchdog Monitoring Sequence...")
	if daemonRunning {
		return
	}
	daemonRunning = true

	go func() {
		for {
			// Pull bounds explicitly every minute dynamically defining timeouts Native to Go
			runLoop(queue, ws)

			// Extract timeout boundaries defined by operators dynamically via DB SQLite rows
			var settings KeeperSettings
			if err := DB.First(&settings, "id = ?", "default").Error; err == nil && settings.CheckIntervalSeconds > 0 {
				time.Sleep(time.Duration(settings.CheckIntervalSeconds) * time.Second)
			} else {
				time.Sleep(60 * time.Second)
			}
		}
	}()
}

func runLoop(queue *TaskQueue, ws *TelemetrySocket) {
	var settings KeeperSettings
	if err := DB.First(&settings, "id = ?", "default").Error; err != nil {
		log.Printf("[Daemon] Settings inaccessible skipping background sequence...")
		return
	}

	if !settings.IsEnabled {
		return
	}

	apiKey := settings.JulesApiKey
	if apiKey == "" || apiKey == "placeholder" {
		return
	}

	log.Println("[Daemon] Background Sweeper Scanning Active Fleet Bounds...")

	// 1. Session Monitoring
	var activeSessions []Session
	DB.Find(&activeSessions)

	queuedCount := 0
	for _, s := range activeSessions {
		// Mock strict JULES API mapping via queue inserts natively
		payloadBytes, _ := json.Marshal(map[string]interface{}{
			"session": s.ID,
		})

		DB.Create(&QueueJob{
			ID:      fmt.Sprintf("chk-%s-%s", s.ID, uuid.New().String()[:8]),
			Type:    "check_session",
			Payload: string(payloadBytes),
			Status:  "pending",
		})
		queuedCount++
	}

	if queuedCount > 0 {
		log.Printf("[Daemon] Push mapped %d native SQLite bounds dynamically targeting QueueWorkers...", queuedCount)
	}

	// 2. Periodic RAG Indexing Hook mimicking queue.ts handling singleton `index_codebase` jobs strictly
	var indexingJobs int64
	DB.Model(&QueueJob{}).Where("type = ? AND status = ?", "index_codebase", "pending").Count(&indexingJobs)

	if indexingJobs == 0 {
		DB.Create(&QueueJob{
			ID:      fmt.Sprintf("index-%s", uuid.New().String()[:8]),
			Type:    "index_codebase",
			Payload: "{}",
			Status:  "pending",
		})
		log.Printf("[Daemon] Spawned codebase vector chunk execution natively mapping ORM boundaries.")
	}

	// Ping Dashboard strictly replicating native emits safely via Telemetry
	ws.Broadcast(`{"event": "daemon_swept", "payload": {"status": "ok"}}`)
}
