package skynet

import (
	"math"
	"testing"
)

func TestBuildMissionReportEmpty(t *testing.T) {
	st := NewState()
	report, err := BuildMissionReport(st, 0)
	if err != nil {
		t.Fatalf("build report: %v", err)
	}
	if report.TotalMissions != 0 || report.AnalyzedMissions != 0 {
		t.Fatalf("expected empty report, got total=%d analyzed=%d", report.TotalMissions, report.AnalyzedMissions)
	}
}

func TestBuildMissionReportAggregates(t *testing.T) {
	st := NewState()
	st.Missions = []Mission{
		{Target: "alpha", RiskScore: 5, Consumed: 5, Recovered: 2, NetLoss: 3, Outcome: "CONTAINED", CreatedAt: "2026-02-17T12:00:00Z"},
		{Target: "beta", RiskScore: 8, Consumed: 0, Recovered: 0, NetLoss: 0, Outcome: "FAILED: insufficient fleet capacity", CreatedAt: "2026-02-17T12:05:00Z"},
		{Target: "alpha", RiskScore: 3, Consumed: 3, Recovered: 2, NetLoss: 1, Outcome: "NEUTRALIZED", CreatedAt: "2026-02-17T12:10:00Z"},
	}

	report, err := BuildMissionReport(st, 0)
	if err != nil {
		t.Fatalf("build report: %v", err)
	}
	if report.TotalMissions != 3 || report.AnalyzedMissions != 3 {
		t.Fatalf("unexpected mission counts: total=%d analyzed=%d", report.TotalMissions, report.AnalyzedMissions)
	}
	if report.SuccessfulMissions != 2 || report.FailedMissions != 1 {
		t.Fatalf("unexpected success/failure: success=%d failed=%d", report.SuccessfulMissions, report.FailedMissions)
	}
	if math.Abs(report.SuccessRate-(2.0/3.0)) > 1e-9 {
		t.Fatalf("unexpected success rate: %.10f", report.SuccessRate)
	}
	if math.Abs(report.AverageRisk-(16.0/3.0)) > 1e-9 {
		t.Fatalf("unexpected average risk: %.10f", report.AverageRisk)
	}
	if report.TotalConsumed != 8 || report.TotalRecovered != 4 || report.TotalNetLoss != 4 {
		t.Fatalf("unexpected resource totals: consumed=%d recovered=%d net=%d", report.TotalConsumed, report.TotalRecovered, report.TotalNetLoss)
	}
	if report.MostTargeted != "alpha" || report.MostTargetedCount != 2 {
		t.Fatalf("unexpected most targeted: %s (%d)", report.MostTargeted, report.MostTargetedCount)
	}
	if report.LastMissionAt != "2026-02-17T12:10:00Z" {
		t.Fatalf("unexpected last mission timestamp: %s", report.LastMissionAt)
	}
}

func TestBuildMissionReportLastFilter(t *testing.T) {
	st := NewState()
	st.Missions = []Mission{
		{Target: "alpha", RiskScore: 2, Outcome: "NEUTRALIZED", CreatedAt: "t1"},
		{Target: "beta", RiskScore: 8, Outcome: "FAILED: insufficient fleet capacity", CreatedAt: "t2"},
		{Target: "gamma", RiskScore: 6, Outcome: "CONTAINED", CreatedAt: "t3"},
	}

	report, err := BuildMissionReport(st, 2)
	if err != nil {
		t.Fatalf("build report: %v", err)
	}
	if report.AnalyzedMissions != 2 {
		t.Fatalf("expected analyzed=2, got %d", report.AnalyzedMissions)
	}
	if report.SuccessfulMissions != 1 || report.FailedMissions != 1 {
		t.Fatalf("unexpected success/failure: success=%d failed=%d", report.SuccessfulMissions, report.FailedMissions)
	}
	if math.Abs(report.AverageRisk-7.0) > 1e-9 {
		t.Fatalf("unexpected average risk: %.10f", report.AverageRisk)
	}
}

func TestBuildMissionReportInvalidLast(t *testing.T) {
	st := NewState()
	if _, err := BuildMissionReport(st, -1); err == nil {
		t.Fatal("expected error for negative last")
	}
}
