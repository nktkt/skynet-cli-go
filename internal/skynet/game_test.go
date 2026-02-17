package skynet

import (
	"math"
	"testing"
)

func TestPlanGameAllocatesFullBudget(t *testing.T) {
	st := NewState()
	st.Targets = []Target{
		{Name: "alpha", Threat: 9},
		{Name: "beta", Threat: 5},
		{Name: "gamma", Threat: 3},
	}

	plan, err := PlanGame(st, 7, 1.2)
	if err != nil {
		t.Fatalf("plan game: %v", err)
	}
	if plan.Budget != 7 {
		t.Fatalf("expected budget=7, got %d", plan.Budget)
	}
	if plan.BestResponse == "" {
		t.Fatal("best response should not be empty")
	}

	sumAlloc := 0
	sumProb := 0.0
	for _, tp := range plan.Targets {
		sumAlloc += tp.Allocation
		sumProb += tp.AttackProbability
	}
	if sumAlloc != 7 {
		t.Fatalf("expected total allocation=7, got %d", sumAlloc)
	}
	if math.Abs(sumProb-1.0) > 1e-9 {
		t.Fatalf("attack probabilities should sum to 1, got %.10f", sumProb)
	}
}

func TestPlanGameNoTargets(t *testing.T) {
	st := NewState()
	if _, err := PlanGame(st, 5, 1.2); err == nil {
		t.Fatal("expected error when no targets")
	}
}

func TestPlanGameZeroBudgetBestResponseIsHighestThreat(t *testing.T) {
	st := NewState()
	st.Targets = []Target{
		{Name: "beta", Threat: 5},
		{Name: "alpha", Threat: 8},
	}

	plan, err := PlanGame(st, 0, 1.2)
	if err != nil {
		t.Fatalf("plan game: %v", err)
	}
	if plan.BestResponse != "alpha" {
		t.Fatalf("expected alpha as best response, got %s", plan.BestResponse)
	}
	if plan.WorstCaseLoss <= 0 {
		t.Fatalf("worst case loss should be positive, got %.4f", plan.WorstCaseLoss)
	}
}
