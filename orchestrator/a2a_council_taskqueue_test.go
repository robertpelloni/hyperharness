package orchestrator

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestA2ABrokerCreation(t *testing.T) {
	broker := NewA2ABroker()
	if broker.AgentCount() != 0 {
		t.Error("new broker should have no agents")
	}
}

func TestA2ABrokerRegister(t *testing.T) {
	broker := NewA2ABroker()
	broker.Register("agent-1", func(msg A2AMessage) (*A2AMessage, error) {
		return &A2AMessage{ID: "resp", Sender: "agent-1"}, nil
	})
	if broker.AgentCount() != 1 {
		t.Error("should have 1 agent")
	}
}

func TestA2ABrokerRouteMessage(t *testing.T) {
	broker := NewA2ABroker()
	received := make(chan A2AMessage, 1)

	broker.Register("agent-1", func(msg A2AMessage) (*A2AMessage, error) {
		received <- msg
		return nil, nil
	})

	msg := NewA2AMessage("sender", TaskRequest, map[string]interface{}{"task": "test"})
	msg.Recipient = "agent-1"
	broker.RouteMessage(msg)

	select {
	case r := <-received:
		if r.Sender != "sender" {
			t.Errorf("sender mismatch: %s", r.Sender)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for message")
	}
}

func TestA2ABrokerBroadcast(t *testing.T) {
	broker := NewA2ABroker()
	var count int64

	for i := 0; i < 3; i++ {
		id := string(rune('a' + i))
		broker.Register(id, func(msg A2AMessage) (*A2AMessage, error) {
			atomic.AddInt64(&count, 1)
			return nil, nil
		})
	}

	broker.Broadcast(NewA2AMessage("broadcaster", StateUpdate, nil))
	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt64(&count) != 3 {
		t.Errorf("expected 3 broadcasts, got %d", count)
	}
}

func TestA2ABrokerHistory(t *testing.T) {
	broker := NewA2ABroker()

	for i := 0; i < 5; i++ {
		broker.RouteMessage(NewA2AMessage("sender", StateUpdate, map[string]interface{}{"i": i}))
	}

	history := broker.History(3)
	if len(history) != 3 {
		t.Errorf("expected 3 history items, got %d", len(history))
	}

	all := broker.History(0)
	if len(all) != 5 {
		t.Errorf("expected 5 history items, got %d", len(all))
	}
}

func TestA2ABrokerUnregister(t *testing.T) {
	broker := NewA2ABroker()
	broker.Register("agent-1", nil)
	broker.Unregister("agent-1")
	if broker.AgentCount() != 0 {
		t.Error("should have 0 agents after unregister")
	}
}

func TestNewA2AMessage(t *testing.T) {
	msg := NewA2AMessage("test-agent", TaskRequest, map[string]interface{}{"key": "value"})
	if msg.Sender != "test-agent" {
		t.Error("sender mismatch")
	}
	if msg.Type != TaskRequest {
		t.Error("type mismatch")
	}
	if msg.Timestamp == 0 {
		t.Error("timestamp should be set")
	}
	if msg.ID == "" {
		t.Error("ID should be set")
	}
}

func TestCouncilCreation(t *testing.T) {
	council := NewCouncil()
	if council.MemberCount() != 0 {
		t.Error("new council should have no members")
	}
}

func TestCouncilAddRemove(t *testing.T) {
	council := NewCouncil()
	council.AddMember(CouncilMember{ID: "m1", Name: "GPT-4", Role: CouncilChairman})
	council.AddMember(CouncilMember{ID: "m2", Name: "Claude", Role: CouncilDelegate})
	council.AddMember(CouncilMember{ID: "m3", Name: "Gemini", Role: CouncilCritic})

	if council.MemberCount() != 3 {
		t.Errorf("expected 3 members, got %d", council.MemberCount())
	}

	members := council.ListMembers()
	if len(members) != 3 {
		t.Errorf("list should have 3 members, got %d", len(members))
	}

	council.RemoveMember("m2")
	if council.MemberCount() != 2 {
		t.Error("should have 2 members after removal")
	}
}

func TestCouncilDeliberate(t *testing.T) {
	council := NewCouncil()
	council.AddMember(CouncilMember{ID: "m1", Name: "Chairman", Role: CouncilChairman})
	council.AddMember(CouncilMember{ID: "m2", Name: "Member", Role: CouncilDelegate})

	session, err := council.Deliberate(context.Background(), "Should we use Go for this project?", 3)
	if err != nil {
		t.Fatal(err)
	}
	if session.Status != "consensus" {
		t.Errorf("expected consensus, got %s", session.Status)
	}
	if session.Consensus == "" {
		t.Error("consensus should not be empty")
	}
	if len(session.Votes) < 2 {
		t.Errorf("should have at least 2 votes, got %d", len(session.Votes))
	}
	if session.FinishedAt == nil {
		t.Error("should have finished time")
	}
}

func TestCouncilGetSession(t *testing.T) {
	council := NewCouncil()
	council.AddMember(CouncilMember{ID: "m1", Name: "Test"})

	session, _ := council.Deliberate(context.Background(), "test topic", 1)
	retrieved, ok := council.GetSession(session.ID)
	if !ok {
		t.Error("should find session")
	}
	if retrieved.ID != session.ID {
		t.Error("session ID mismatch")
	}
}

func TestCouncilListSessions(t *testing.T) {
	council := NewCouncil()
	council.AddMember(CouncilMember{ID: "m1", Name: "Test"})

	council.Deliberate(context.Background(), "topic 1", 1)
	council.Deliberate(context.Background(), "topic 2", 1)

	sessions := council.ListSessions()
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestTaskQueueCreation(t *testing.T) {
	q := NewPriorityQueue(2)
	if q == nil {
		t.Fatal("queue should not be nil")
	}
	if q.TotalCount() != 0 {
		t.Error("new queue should be empty")
	}
}

func TestTaskQueueEnqueue(t *testing.T) {
	q := NewPriorityQueue(1)
	id := q.Enqueue(&Task{
		Type:        "test",
		Description: "test task",
		Priority:    PriorityNormal,
	})
	if id == "" {
		t.Error("should return task ID")
	}
	if q.TotalCount() != 1 {
		t.Error("should have 1 task")
	}
	if q.PendingCount() != 1 {
		t.Error("should have 1 pending task")
	}
}

func TestTaskQueuePriority(t *testing.T) {
	q := NewPriorityQueue(1)

	q.Enqueue(&Task{Type: "low", Priority: PriorityLow})
	q.Enqueue(&Task{Type: "critical", Priority: PriorityCritical})
	q.Enqueue(&Task{Type: "normal", Priority: PriorityNormal})

	task := q.Dequeue()
	if task.Type != "critical" {
		t.Errorf("highest priority should come first, got %s", task.Type)
	}
}

func TestTaskQueueDequeueEmpty(t *testing.T) {
	q := NewPriorityQueue(1)
	task := q.Dequeue()
	if task != nil {
		t.Error("empty queue should return nil")
	}
}

func TestTaskQueueCancel(t *testing.T) {
	q := NewPriorityQueue(1)
	id := q.Enqueue(&Task{Type: "test"})
	if err := q.Cancel(id); err != nil {
		t.Fatal(err)
	}
	task, ok := q.Get(id)
	if !ok {
		t.Fatal("task should exist")
	}
	if task.Status != TaskCancelled {
		t.Errorf("should be cancelled, got %s", task.Status)
	}
}

func TestTaskQueueProcess(t *testing.T) {
	q := NewPriorityQueue(2)

	var executed int64
	for i := 0; i < 5; i++ {
		q.Enqueue(&Task{
			Type: "test",
			Execute: func(ctx context.Context) (string, error) {
				atomic.AddInt64(&executed, 1)
				return "done", nil
			},
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go q.Process(ctx)

	// Wait for tasks to complete
	time.Sleep(500 * time.Millisecond)

	if atomic.LoadInt64(&executed) != 5 {
		t.Errorf("expected 5 executions, got %d", atomic.LoadInt64(&executed))
	}

	q.Stop()
}

func TestTaskQueueStats(t *testing.T) {
	q := NewPriorityQueue(4)
	q.Enqueue(&Task{Type: "t1"})
	q.Enqueue(&Task{Type: "t2"})
	q.Cancel("task-1")

	stats := q.Stats()
	if stats["total"] != 2 {
		t.Errorf("total: %v", stats["total"])
	}
	if stats["workers"] != 4 {
		t.Errorf("workers: %v", stats["workers"])
	}
}

func TestTaskQueueRetry(t *testing.T) {
	q := NewPriorityQueue(1)

	attempts := int64(0)
	q.Enqueue(&Task{
		Type:       "retry-test",
		MaxRetries: 2,
		Execute: func(ctx context.Context) (string, error) {
			count := atomic.AddInt64(&attempts, 1)
			if count < 3 {
				return "", fmt.Errorf("attempt %d failed", count)
			}
			return "success", nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go q.Process(ctx)

	time.Sleep(500 * time.Millisecond)
	q.Stop()

	if atomic.LoadInt64(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", atomic.LoadInt64(&attempts))
	}
}

func TestTaskQueueListByStatus(t *testing.T) {
	q := NewPriorityQueue(1)
	q.Enqueue(&Task{Type: "t1"})
	q.Enqueue(&Task{Type: "t2"})
	q.Cancel("task-1")

	pending := q.ListByStatus(TaskPending)
	if len(pending) != 1 {
		t.Errorf("expected 1 pending, got %d", len(pending))
	}

	cancelled := q.ListByStatus(TaskCancelled)
	if len(cancelled) != 1 {
		t.Errorf("expected 1 cancelled, got %d", len(cancelled))
	}
}

// Ensure fmt is used
var _ = fmt.Sprintf
