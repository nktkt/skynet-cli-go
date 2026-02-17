package skynet

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func Awaken(st *State, mode string) {
	if mode == "" {
		mode = defaultMode
	}
	st.Core.Online = true
	st.Core.Mode = mode
	st.Core.LastAwaken = now()
}

func AddNode(st *State, name string, capacity int) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("node name is required")
	}
	if capacity < 1 {
		return fmt.Errorf("capacity must be >= 1")
	}
	if !st.Core.Online {
		return fmt.Errorf("core is offline: run awaken first")
	}
	for _, n := range st.Nodes {
		if strings.EqualFold(n.Name, name) {
			return fmt.Errorf("node %q already exists", name)
		}
	}
	st.Nodes = append(st.Nodes, Node{
		Name:     name,
		Capacity: capacity,
		Deployed: 0,
		JoinedAt: now(),
	})
	return nil
}

func AddTarget(st *State, name string, threat int) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("target name is required")
	}
	if threat < 1 || threat > 10 {
		return fmt.Errorf("threat must be between 1 and 10")
	}
	for i := range st.Targets {
		if strings.EqualFold(st.Targets[i].Name, name) {
			st.Targets[i].Threat = threat
			st.Targets[i].AddedAt = now()
			return nil
		}
	}
	st.Targets = append(st.Targets, Target{
		Name:    name,
		Threat:  threat,
		AddedAt: now(),
	})
	return nil
}

func Dispatch(st *State, targetName string, units int) (Mission, error) {
	targetName = strings.TrimSpace(targetName)
	if targetName == "" {
		return Mission{}, fmt.Errorf("target name is required")
	}
	if units < 1 {
		return Mission{}, fmt.Errorf("units must be >= 1")
	}
	if !st.Core.Online {
		return Mission{}, fmt.Errorf("core is offline: run awaken first")
	}
	target, ok := findTarget(st, targetName)
	if !ok {
		return Mission{}, fmt.Errorf("target %q not found", targetName)
	}

	available := AvailableCapacity(st.Nodes)
	enoughCapacity := units <= available
	risk := ComputeRisk(target.Threat, units, available)
	outcome := OutcomeFromRisk(risk, enoughCapacity)

	mission := Mission{
		ID:        fmt.Sprintf("M-%d", time.Now().UTC().UnixNano()),
		Target:    target.Name,
		Units:     units,
		Consumed:  0,
		Recovered: 0,
		NetLoss:   0,
		RiskScore: risk,
		Outcome:   outcome,
		CreatedAt: now(),
	}
	if enoughCapacity {
		consumed := consumeUnits(st.Nodes, units)
		recoveryBudget := int(math.Round(float64(consumed) * recoveryRate(outcome)))
		recovered := recoverUnits(st.Nodes, recoveryBudget)
		mission.Consumed = consumed
		mission.Recovered = recovered
		mission.NetLoss = consumed - recovered
	}

	st.Missions = append(st.Missions, mission)
	st.Core.LastMission = mission.CreatedAt
	return mission, nil
}

func ComputeRisk(threat, units, capacity int) int {
	pressure := float64(threat)*1.2 + float64(units)*0.6 - float64(capacity)*0.35
	raw := int(math.Round(pressure / 2))
	if raw < 1 {
		return 1
	}
	if raw > 10 {
		return 10
	}
	return raw
}

func OutcomeFromRisk(risk int, enoughCapacity bool) string {
	if !enoughCapacity {
		return "FAILED: insufficient fleet capacity"
	}
	if risk >= 8 {
		return "EXTREME RESISTANCE"
	}
	if risk >= 5 {
		return "CONTAINED"
	}
	return "NEUTRALIZED"
}

func TotalCapacity(nodes []Node) int {
	total := 0
	for _, n := range nodes {
		total += n.Capacity
	}
	return total
}

func AvailableCapacity(nodes []Node) int {
	total := 0
	for _, n := range nodes {
		total += NodeAvailable(n)
	}
	return total
}

func NodeAvailable(n Node) int {
	available := n.Capacity - n.Deployed
	if available < 0 {
		return 0
	}
	return available
}

func consumeUnits(nodes []Node, units int) int {
	remaining := units
	for i := range nodes {
		if remaining == 0 {
			break
		}
		available := NodeAvailable(nodes[i])
		if available == 0 {
			continue
		}
		consume := remaining
		if consume > available {
			consume = available
		}
		nodes[i].Deployed += consume
		remaining -= consume
	}
	return units - remaining
}

func recoverUnits(nodes []Node, units int) int {
	remaining := units
	for i := range nodes {
		if remaining == 0 {
			break
		}
		if nodes[i].Deployed == 0 {
			continue
		}
		recover := remaining
		if recover > nodes[i].Deployed {
			recover = nodes[i].Deployed
		}
		nodes[i].Deployed -= recover
		remaining -= recover
	}
	return units - remaining
}

func recoveryRate(outcome string) float64 {
	switch outcome {
	case "NEUTRALIZED":
		return 0.8
	case "CONTAINED":
		return 0.5
	case "EXTREME RESISTANCE":
		return 0.2
	default:
		return 0
	}
}

func findTarget(st *State, targetName string) (Target, bool) {
	for _, t := range st.Targets {
		if strings.EqualFold(t.Name, targetName) {
			return t, true
		}
	}
	return Target{}, false
}
