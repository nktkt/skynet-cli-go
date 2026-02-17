package skynet

import (
	"fmt"
	"math"
	"sort"
)

const (
	defaultAttackBeta = 1.2
	defenseElasticity = 0.18
)

type GameTargetPlan struct {
	Name              string  `json:"name"`
	Threat            int     `json:"threat"`
	Allocation        int     `json:"allocation"`
	AttackerPayoff    float64 `json:"attacker_payoff"`
	AttackProbability float64 `json:"attack_probability"`
}

type GamePlan struct {
	Budget          int              `json:"budget"`
	Beta            float64          `json:"beta"`
	BestResponse    string           `json:"best_response"`
	WorstCaseLoss   float64          `json:"worst_case_loss"`
	ExpectedLoss    float64          `json:"expected_loss"`
	DefenderUtility float64          `json:"defender_utility"`
	Targets         []GameTargetPlan `json:"targets"`
}

func PlanGame(st State, budget int, beta float64) (GamePlan, error) {
	if budget < 0 {
		return GamePlan{}, fmt.Errorf("budget must be >= 0")
	}
	if len(st.Targets) == 0 {
		return GamePlan{}, fmt.Errorf("at least one target is required")
	}
	if beta <= 0 {
		beta = defaultAttackBeta
	}

	targets := make([]GameTargetPlan, 0, len(st.Targets))
	for _, t := range st.Targets {
		targets = append(targets, GameTargetPlan{Name: t.Name, Threat: t.Threat})
	}
	sort.Slice(targets, func(i, j int) bool {
		left := targets[i]
		right := targets[j]
		if left.Threat != right.Threat {
			return left.Threat > right.Threat
		}
		return left.Name < right.Name
	})

	for i := 0; i < budget; i++ {
		idx := argmaxAttackerPayoff(targets)
		targets[idx].Allocation++
	}

	for i := range targets {
		targets[i].AttackerPayoff = attackerPayoff(targets[i].Threat, targets[i].Allocation)
	}
	probs := attackProbabilities(targets, beta)
	for i := range targets {
		targets[i].AttackProbability = probs[i]
	}

	bestIdx := argmaxAttackerPayoff(targets)
	worst := targets[bestIdx].AttackerPayoff
	expected := 0.0
	for i := range targets {
		expected += targets[i].AttackProbability * targets[i].AttackerPayoff
	}

	return GamePlan{
		Budget:          budget,
		Beta:            beta,
		BestResponse:    targets[bestIdx].Name,
		WorstCaseLoss:   worst,
		ExpectedLoss:    expected,
		DefenderUtility: -expected,
		Targets:         targets,
	}, nil
}

func attackerPayoff(threat, allocation int) float64 {
	return float64(threat) * math.Exp(-defenseElasticity*float64(allocation))
}

func argmaxAttackerPayoff(targets []GameTargetPlan) int {
	best := 0
	bestScore := attackerPayoff(targets[0].Threat, targets[0].Allocation)
	for i := 1; i < len(targets); i++ {
		score := attackerPayoff(targets[i].Threat, targets[i].Allocation)
		if score > bestScore {
			best = i
			bestScore = score
			continue
		}
		if score == bestScore {
			if targets[i].Threat > targets[best].Threat {
				best = i
				bestScore = score
				continue
			}
			if targets[i].Threat == targets[best].Threat && targets[i].Name < targets[best].Name {
				best = i
				bestScore = score
			}
		}
	}
	return best
}

func attackProbabilities(targets []GameTargetPlan, beta float64) []float64 {
	exps := make([]float64, len(targets))
	sum := 0.0
	for i := range targets {
		v := beta * targets[i].AttackerPayoff
		e := math.Exp(v)
		exps[i] = e
		sum += e
	}

	probs := make([]float64, len(targets))
	if sum == 0 {
		uniform := 1.0 / float64(len(targets))
		for i := range probs {
			probs[i] = uniform
		}
		return probs
	}
	for i := range probs {
		probs[i] = exps[i] / sum
	}
	return probs
}
