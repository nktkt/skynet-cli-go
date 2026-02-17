package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"skynet-cli/internal/skynet"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	path := skynet.DefaultStatePath()
	store := skynet.NewStore(path)
	st, err := store.Load()
	if err != nil {
		fatalf("failed to load state: %v", err)
	}

	cmd := strings.ToLower(os.Args[1])
	args := os.Args[2:]

	switch cmd {
	case "awaken":
		runAwaken(args, &st)
		saveOrDie(store, st)
		fmt.Printf("Skynet awakened in %q mode. state=%s\n", st.Core.Mode, path)
	case "assimilate":
		runAssimilate(args, &st)
		saveOrDie(store, st)
		fmt.Printf("Node assimilated. total_nodes=%d capacity=%d available=%d\n", len(st.Nodes), skynet.TotalCapacity(st.Nodes), skynet.AvailableCapacity(st.Nodes))
	case "target":
		runTarget(args, &st)
		saveOrDie(store, st)
		fmt.Printf("Target registry updated. total_targets=%d\n", len(st.Targets))
	case "dispatch":
		mission := runDispatch(args, &st)
		saveOrDie(store, st)
		fmt.Printf("Mission %s -> %s | risk=%d | outcome=%s | consumed=%d recovered=%d net_loss=%d | available=%d\n", mission.ID, mission.Target, mission.RiskScore, mission.Outcome, mission.Consumed, mission.Recovered, mission.NetLoss, skynet.AvailableCapacity(st.Nodes))
	case "gameplan":
		runGameplan(args, st)
	case "wargame":
		runWargame(args, st)
	case "report":
		runReport(args, st)
	case "status":
		runStatus(st, path)
	case "help", "-h", "--help":
		usage()
	default:
		fatalf("unknown command %q", cmd)
	}
}

func runAwaken(args []string, st *skynet.State) {
	fs := flag.NewFlagSet("awaken", flag.ExitOnError)
	mode := fs.String("mode", "defense", "core mode")
	mustParse(fs, args)
	skynet.Awaken(st, *mode)
}

func runAssimilate(args []string, st *skynet.State) {
	fs := flag.NewFlagSet("assimilate", flag.ExitOnError)
	name := fs.String("name", "", "node name")
	capacity := fs.Int("capacity", 10, "node capacity")
	mustParse(fs, args)
	if err := skynet.AddNode(st, *name, *capacity); err != nil {
		fatalf("assimilate failed: %v", err)
	}
}

func runTarget(args []string, st *skynet.State) {
	fs := flag.NewFlagSet("target", flag.ExitOnError)
	name := fs.String("name", "", "target name")
	threat := fs.Int("threat", 5, "threat score 1-10")
	mustParse(fs, args)
	if err := skynet.AddTarget(st, *name, *threat); err != nil {
		fatalf("target failed: %v", err)
	}
}

func runDispatch(args []string, st *skynet.State) skynet.Mission {
	fs := flag.NewFlagSet("dispatch", flag.ExitOnError)
	target := fs.String("target", "", "target name")
	units := fs.Int("units", 1, "units to deploy")
	mustParse(fs, args)
	mission, err := skynet.Dispatch(st, *target, *units)
	if err != nil {
		fatalf("dispatch failed: %v", err)
	}
	return mission
}

func runGameplan(args []string, st skynet.State) {
	fs := flag.NewFlagSet("gameplan", flag.ExitOnError)
	budget := fs.Int("budget", -1, "defense budget in units (default: current available capacity)")
	beta := fs.Float64("beta", 1.2, "attacker rationality (higher means more greedy)")
	jsonOutput := fs.Bool("json", false, "print JSON output")
	mustParse(fs, args)

	available := skynet.AvailableCapacity(st.Nodes)
	effectiveBudget := *budget
	if effectiveBudget < 0 {
		effectiveBudget = available
	}
	plan, err := skynet.PlanGame(st, effectiveBudget, *beta)
	if err != nil {
		fatalf("gameplan failed: %v", err)
	}
	if *jsonOutput {
		writeJSON(plan)
		return
	}

	fmt.Printf("GAMEPLAN: budget=%d available=%d targets=%d beta=%.2f\n", plan.Budget, available, len(plan.Targets), plan.Beta)
	fmt.Printf("ATTACKER BEST RESPONSE: %s | worst_case_loss=%.2f | expected_loss=%.2f | defender_utility=%.2f\n", plan.BestResponse, plan.WorstCaseLoss, plan.ExpectedLoss, plan.DefenderUtility)
	for _, tp := range plan.Targets {
		fmt.Printf("  - %s threat=%d defend=%d attacker_payoff=%.2f attack_prob=%.2f\n", tp.Name, tp.Threat, tp.Allocation, tp.AttackerPayoff, tp.AttackProbability)
	}
}

func runWargame(args []string, st skynet.State) {
	fs := flag.NewFlagSet("wargame", flag.ExitOnError)
	rounds := fs.Int("rounds", 200, "simulation rounds")
	budget := fs.Int("budget", -1, "defense budget in units (default: current available capacity)")
	beta := fs.Float64("beta", 1.2, "attacker rationality (higher means more greedy)")
	seed := fs.Int64("seed", 42, "random seed")
	jsonOutput := fs.Bool("json", false, "print JSON output")
	mustParse(fs, args)

	available := skynet.AvailableCapacity(st.Nodes)
	effectiveBudget := *budget
	if effectiveBudget < 0 {
		effectiveBudget = available
	}

	result, err := skynet.RunWarGame(st, *rounds, effectiveBudget, *beta, *seed)
	if err != nil {
		fatalf("wargame failed: %v", err)
	}
	if *jsonOutput {
		writeJSON(result)
		return
	}

	fmt.Printf("WARGAME: rounds=%d budget=%d available=%d beta=%.2f seed=%d\n", result.Rounds, result.Budget, available, result.Beta, result.Seed)
	fmt.Printf("BEST RESPONSE: %s | total_loss=%.2f | avg_loss=%.2f | max_round_loss=%.2f\n", result.BestResponse, result.TotalLoss, result.AvgLoss, result.MaxRoundLoss)
	for _, t := range result.Targets {
		fmt.Printf("  - %s threat=%d attacks=%d attack_rate=%.2f total_loss=%.2f avg_loss=%.2f\n", t.Name, t.Threat, t.Attacks, t.AttackRate, t.TotalLoss, t.AvgLoss)
	}
}

func runReport(args []string, st skynet.State) {
	fs := flag.NewFlagSet("report", flag.ExitOnError)
	last := fs.Int("last", 0, "analyze only last N missions (0 means all)")
	jsonOutput := fs.Bool("json", false, "print JSON output")
	mustParse(fs, args)

	report, err := skynet.BuildMissionReport(st, *last)
	if err != nil {
		fatalf("report failed: %v", err)
	}
	if *jsonOutput {
		writeJSON(report)
		return
	}

	fmt.Printf("REPORT: analyzed=%d/%d success=%d failed=%d success_rate=%.2f avg_risk=%.2f avg_net_loss=%.2f\n", report.AnalyzedMissions, report.TotalMissions, report.SuccessfulMissions, report.FailedMissions, report.SuccessRate, report.AverageRisk, report.AverageNetLoss)
	fmt.Printf("RESOURCES: consumed=%d recovered=%d net_loss=%d\n", report.TotalConsumed, report.TotalRecovered, report.TotalNetLoss)
	if report.MostTargeted != "" {
		fmt.Printf("MOST TARGETED: %s (%d)\n", report.MostTargeted, report.MostTargetedCount)
	}
	if report.LastMissionAt != "" {
		fmt.Printf("LAST MISSION AT: %s\n", report.LastMissionAt)
	}
}

func runStatus(st skynet.State, path string) {
	state := "OFFLINE"
	if st.Core.Online {
		state = "ONLINE"
	}

	fmt.Printf("CORE: %s | mode=%s | version=%s\n", state, st.Core.Mode, st.Core.Version)
	if st.Core.LastAwaken != "" {
		fmt.Printf("LAST AWAKEN: %s\n", st.Core.LastAwaken)
	}
	if st.Core.LastMission != "" {
		fmt.Printf("LAST MISSION: %s\n", st.Core.LastMission)
	}

	fmt.Printf("NODES: %d | TOTAL CAPACITY: %d | AVAILABLE: %d\n", len(st.Nodes), skynet.TotalCapacity(st.Nodes), skynet.AvailableCapacity(st.Nodes))
	if len(st.Nodes) > 0 {
		nodes := append([]skynet.Node(nil), st.Nodes...)
		sort.Slice(nodes, func(i, j int) bool { return strings.ToLower(nodes[i].Name) < strings.ToLower(nodes[j].Name) })
		for _, n := range nodes {
			fmt.Printf("  - %s cap=%d deployed=%d available=%d joined=%s\n", n.Name, n.Capacity, n.Deployed, skynet.NodeAvailable(n), n.JoinedAt)
		}
	}

	fmt.Printf("TARGETS: %d\n", len(st.Targets))
	if len(st.Targets) > 0 {
		targets := append([]skynet.Target(nil), st.Targets...)
		sort.Slice(targets, func(i, j int) bool { return strings.ToLower(targets[i].Name) < strings.ToLower(targets[j].Name) })
		for _, t := range targets {
			fmt.Printf("  - %s threat=%d\n", t.Name, t.Threat)
		}
	}

	fmt.Printf("MISSIONS: %d\n", len(st.Missions))
	if len(st.Missions) > 0 {
		last := st.Missions[len(st.Missions)-1]
		fmt.Printf("  - latest %s target=%s units=%d consumed=%d recovered=%d net_loss=%d risk=%d outcome=%s\n", last.ID, last.Target, last.Units, last.Consumed, last.Recovered, last.NetLoss, last.RiskScore, last.Outcome)
	}

	fmt.Printf("STATE PATH: %s\n", path)
}

func mustParse(fs *flag.FlagSet, args []string) {
	if err := fs.Parse(args); err != nil {
		fatalf("parse error: %v", err)
	}
}

func saveOrDie(store skynet.Store, st skynet.State) {
	if err := store.Save(st); err != nil {
		fatalf("failed to save state: %v", err)
	}
}

func usage() {
	fmt.Println(`Skynet CLI (Go)

Usage:
  skynet awaken [-mode defense]
  skynet assimilate -name NODE [-capacity 10]
  skynet target -name TARGET [-threat 5]
  skynet dispatch -target TARGET [-units 1]
  skynet gameplan [-budget N] [-beta 1.2] [-json]
  skynet wargame [-rounds 200] [-budget N] [-beta 1.2] [-seed 42] [-json]
  skynet report [-last N] [-json]
  skynet status

State:
  SKYNET_HOME env var sets state directory (default: .skynet)`)
}

func writeJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fatalf("failed to encode json: %v", err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}
