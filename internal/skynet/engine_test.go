package skynet

import "testing"

func TestAddNodeRequiresOnlineCore(t *testing.T) {
	st := NewState()
	if err := AddNode(&st, "alpha", 4); err == nil {
		t.Fatal("expected error when core is offline")
	}
}

func TestDispatchCreatesMission(t *testing.T) {
	st := NewState()
	Awaken(&st, "defense")
	if err := AddNode(&st, "alpha", 10); err != nil {
		t.Fatalf("add node: %v", err)
	}
	if err := AddTarget(&st, "resistance-cell", 7); err != nil {
		t.Fatalf("add target: %v", err)
	}

	mission, err := Dispatch(&st, "resistance-cell", 4)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if mission.Target != "resistance-cell" {
		t.Fatalf("unexpected target: %s", mission.Target)
	}
	if mission.Consumed != 4 {
		t.Fatalf("expected consumed=4, got %d", mission.Consumed)
	}
	if mission.Recovered != 3 {
		t.Fatalf("expected recovered=3, got %d", mission.Recovered)
	}
	if mission.NetLoss != 1 {
		t.Fatalf("expected net_loss=1, got %d", mission.NetLoss)
	}
	if len(st.Missions) != 1 {
		t.Fatalf("expected 1 mission, got %d", len(st.Missions))
	}
	if got := AvailableCapacity(st.Nodes); got != 9 {
		t.Fatalf("expected available=9, got %d", got)
	}
}

func TestDispatchInsufficientCapacityOutcome(t *testing.T) {
	st := NewState()
	Awaken(&st, "defense")
	if err := AddNode(&st, "alpha", 2); err != nil {
		t.Fatalf("add node: %v", err)
	}
	if err := AddTarget(&st, "hq", 3); err != nil {
		t.Fatalf("add target: %v", err)
	}

	mission, err := Dispatch(&st, "hq", 5)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if mission.Outcome != "FAILED: insufficient fleet capacity" {
		t.Fatalf("unexpected outcome: %s", mission.Outcome)
	}
	if mission.Consumed != 0 || mission.Recovered != 0 || mission.NetLoss != 0 {
		t.Fatalf("expected no resource change, got consumed=%d recovered=%d net=%d", mission.Consumed, mission.Recovered, mission.NetLoss)
	}
	if got := AvailableCapacity(st.Nodes); got != 2 {
		t.Fatalf("expected available=2, got %d", got)
	}
}

func TestDispatchContainedRecoveryRate(t *testing.T) {
	st := NewState()
	Awaken(&st, "defense")
	if err := AddNode(&st, "alpha", 12); err != nil {
		t.Fatalf("add node: %v", err)
	}
	if err := AddTarget(&st, "hq", 8); err != nil {
		t.Fatalf("add target: %v", err)
	}

	mission, err := Dispatch(&st, "hq", 6)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if mission.Outcome != "CONTAINED" {
		t.Fatalf("expected CONTAINED, got %s", mission.Outcome)
	}
	if mission.Recovered != 3 {
		t.Fatalf("expected recovered=3, got %d", mission.Recovered)
	}
	if got := AvailableCapacity(st.Nodes); got != 9 {
		t.Fatalf("expected available=9, got %d", got)
	}
}

func TestComputeRiskBounds(t *testing.T) {
	if got := ComputeRisk(1, 1, 99); got != 1 {
		t.Fatalf("expected lower bound 1, got %d", got)
	}
	if got := ComputeRisk(10, 99, 0); got != 10 {
		t.Fatalf("expected upper bound 10, got %d", got)
	}
}
