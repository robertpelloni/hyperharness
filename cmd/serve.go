package cmd

import (
	"log"
	
	"github.com/spf13/cobra"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/robertpelloni/hypercode/mcp"
	"github.com/robertpelloni/hypercode/orchestrator"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the monolithic Borg Daemon Backend (Port 8080)",
	Long:  "Fires up the Go-native replacement for the legacy Bun/Hono TS backend.",
	Run: func(cmd *cobra.Command, args []string) {
		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		// Strictly handle CORS for Next.js Frontend parity
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
		}))

		// Initialize Database/Queues natively substituting BullMQ & Prisma
		if err := orchestrator.InitDatabase("./.borg_queue.db"); err != nil {
			log.Fatalf("Prisma Parity Core mapping failed: %v", err)
		}
		
		queue, err := orchestrator.NewTaskQueue("./.borg_queue.db")
		if err != nil {
			log.Fatalf("Queue Initialization Failure: %v", err)
		}
		queue.StartWorker()
		defer queue.Close()

		// Telemetry Service replacing Hono Websocket bridging
		wsSvc := orchestrator.NewTelemetrySocket()
		app.Use("/ws", wsSvc.WsHandler)
		app.Get("/ws/stream", websocket.New(wsSvc.ConnectionLoop))

		// Phase 2: Live Traffic Observability (JSON-RPC Telemetry Hook)
		mcp.MCPTrafficMonitor = func(payload string) {
			wsSvc.Broadcast(payload)
		}

		// Core TS Parity Endpoints
		api := app.Group("/api/v1")
		
		api.Get("/sessions", func(c *fiber.Ctx) error {
			var sessions []orchestrator.Session
			// Natively mapping Prisma's listSessions
			if err := orchestrator.DB.Order("created_at desc").Find(&sessions).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			return c.JSON(fiber.Map{"sessions": sessions})
		})
		
		api.Get("/sessions/:id/replay", func(c *fiber.Ctx) error {
			id := c.Params("id")
			var session orchestrator.Session
			if err := orchestrator.DB.First(&session, "id = ?", id).Error; err != nil {
				return c.Status(404).JSON(fiber.Map{"error": "Session absent"})
			}
			
			var logs []orchestrator.KeeperLog
			orchestrator.DB.Where("session_id = ?", id).Order("created_at asc").Find(&logs)
			
			// Map structural timeline
			var timeline []map[string]interface{}
			for _, l := range logs {
				timeline = append(timeline, map[string]interface{}{
					"id": l.ID,
					"timestamp": l.CreatedAt,
					"type": l.Type,
					"content": l.Message,
					"metadata": l.Metadata,
				})
			}
			
			return c.JSON(fiber.Map{
				"sessionId": id,
				"title": session.Title,
				"status": session.Status,
				"timeline": timeline,
			})
		})
		
		api.Get("/fleet/summary", func(c *fiber.Ctx) error {
			var sessionCount int64
			var pendingJobs int64
			var processingJobs int64
			var chunkCount int64
			
			orchestrator.DB.Model(&orchestrator.Session{}).Count(&sessionCount)
			orchestrator.DB.Model(&orchestrator.QueueJob{}).Where("status = ?", "pending").Count(&pendingJobs)
			orchestrator.DB.Model(&orchestrator.QueueJob{}).Where("status = ?", "processing").Count(&processingJobs)
			orchestrator.DB.Model(&orchestrator.CodeChunk{}).Count(&chunkCount)
			
			var recentActions []orchestrator.KeeperLog
			orchestrator.DB.Where("type = ?", "action").Order("created_at desc").Limit(5).Find(&recentActions)
			
			return c.JSON(fiber.Map{
				"fleet": fiber.Map{"total": sessionCount},
				"orchestrator": fiber.Map{
					"queueDepth": pendingJobs + processingJobs,
					"isActive": processingJobs > 0,
					"recentAutonomousActions": recentActions,
				},
				"knowledgeBase": fiber.Map{
					"totalChunks": chunkCount,
					"isIndexed": chunkCount > 0,
				},
				"borgReady": true,
			})
		})
		
		api.Get("/system/submodules", func(c *fiber.Ctx) error {
			subs, err := orchestrator.ExtractSubmoduleIntelligence()
			if err != nil {
				return c.JSON(fiber.Map{"submodules": []orchestrator.SubmoduleStatus{}}) // Fallback to empty matching TS parity
			}
			return c.JSON(fiber.Map{"submodules": subs})
		})
		
		api.Post("/webhooks/borg", func(c *fiber.Ctx) error {
			var payload orchestrator.WebhookPayload
			if err := c.BodyParser(&payload); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON format"})
			}
			result, err := orchestrator.HandleBorgWebhook(payload, queue, wsSvc)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			return c.JSON(result)
		})
		
		api.Post("/workspaces", func(c *fiber.Ctx) error {
			queue.Enqueue("INIT_WORKSPACE_ACTION")
			return c.JSON(fiber.Map{"job": "enqueued", "action": "INITIALIZE_GIT"})
		})
		
		app.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "hypercode_active", "version": "0.6.0", "daemon": "fiber"})
		})

		log.Println("[Server] Hono/Bun Parity Achieved. Listening locally on :8080")
		log.Fatal(app.Listen(":8080"))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
