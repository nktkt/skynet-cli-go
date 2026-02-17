package skynet

import (
	"fmt"
	"sort"
	"strings"
)

type MissionReport struct {
	TotalMissions      int     `json:"total_missions"`
	AnalyzedMissions   int     `json:"analyzed_missions"`
	SuccessfulMissions int     `json:"successful_missions"`
	FailedMissions     int     `json:"failed_missions"`
	SuccessRate        float64 `json:"success_rate"`
	AverageRisk        float64 `json:"average_risk"`
	AverageNetLoss     float64 `json:"average_net_loss"`
	TotalConsumed      int     `json:"total_consumed"`
	TotalRecovered     int     `json:"total_recovered"`
	TotalNetLoss       int     `json:"total_net_loss"`
	MostTargeted       string  `json:"most_targeted"`
	MostTargetedCount  int     `json:"most_targeted_count"`
	LastMissionAt      string  `json:"last_mission_at"`
}

func BuildMissionReport(st State, last int) (MissionReport, error) {
	if last < 0 {
		return MissionReport{}, fmt.Errorf("last must be >= 0")
	}

	total := len(st.Missions)
	if total == 0 {
		return MissionReport{TotalMissions: 0}, nil
	}

	start := 0
	if last > 0 && last < total {
		start = total - last
	}
	missions := st.Missions[start:]

	report := MissionReport{
		TotalMissions:    total,
		AnalyzedMissions: len(missions),
		LastMissionAt:    st.Missions[total-1].CreatedAt,
	}

	targetCount := map[string]int{}
	for _, m := range missions {
		if isMissionSuccess(m) {
			report.SuccessfulMissions++
		} else {
			report.FailedMissions++
		}
		report.TotalConsumed += m.Consumed
		report.TotalRecovered += m.Recovered
		report.TotalNetLoss += m.NetLoss
		report.AverageRisk += float64(m.RiskScore)
		targetCount[m.Target]++
	}

	if report.AnalyzedMissions > 0 {
		count := float64(report.AnalyzedMissions)
		report.SuccessRate = float64(report.SuccessfulMissions) / count
		report.AverageRisk = report.AverageRisk / count
		report.AverageNetLoss = float64(report.TotalNetLoss) / count
	}

	report.MostTargeted, report.MostTargetedCount = findMostTargeted(targetCount)
	return report, nil
}

func isMissionSuccess(m Mission) bool {
	return !strings.HasPrefix(strings.ToUpper(m.Outcome), "FAILED")
}

func findMostTargeted(targetCount map[string]int) (string, int) {
	if len(targetCount) == 0 {
		return "", 0
	}
	type kv struct {
		name  string
		count int
	}
	items := make([]kv, 0, len(targetCount))
	for name, count := range targetCount {
		items = append(items, kv{name: name, count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].count != items[j].count {
			return items[i].count > items[j].count
		}
		return items[i].name < items[j].name
	})
	return items[0].name, items[0].count
}
