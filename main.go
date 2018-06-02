package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"strings"
	"sort"
)

type Instance struct {
	id   string
	rule Rule
	dot  int
}

type Rule struct {
	id    int
	left  string
	right string
}

type RuleGroup []Rule
type InstanceGroup []Instance
type Solver []State

type Action struct {
	goTo   int
	action byte
}

type State struct {
	i         int
	instances InstanceGroup
	action    map[string]Action
}

/* RULES */
func (R Rule) makeInstance(index int) Instance {
	return Instance{
		rule: R,
		dot:  index,
	}.makeId()
}

func (R RuleGroup) findRulesFor(char string) RuleGroup {
	var result RuleGroup
	for _, r := range R {
		if r.left == char {
			result = append(result, r)
		}
	}
	return result
}
func (R RuleGroup) first(inChar string) []string {
	next := make(map[string]struct{})
	stack := R.findRulesFor(inChar)
	if len(stack) == 0 {
		return []string{inChar}
	}

	var rule Rule
	for ; len(stack) > 0; {
		rule, stack = stack[0], stack[1:]
		char := string(rule.right[0])
		rules := R.findRulesFor(char)
		if len(rules) == 0 {
			next[char] = struct{}{}
		} else {
			stack = append(stack, rules...)
		}
	}
	var result []string
	for i := range next {
		result = append(result, i)
	}

	return result
}
func (R RuleGroup) follow(inChar string) []string {
	if inChar == "1" {
		return []string{"0"}
	}

	follow := make(map[string]struct{})
	done := map[string]struct{}{inChar: {}}
	stack := []string{inChar}

	var findChar string
	for ; len(stack) > 0; { // stack
		findChar, stack = stack[0], stack[1:]
		for _, rule := range R { // rules
			for i, tmpChar := range rule.right { // rule chars
				if findChar == "1" {
					follow["0"] = struct{}{}
				}
				if findChar == string(tmpChar) { // matches found S83 A65 B66 a97 b98
					if len(rule.right) > i+1 { // exists follow
						for _, val := range R.first(string(rule.right[i+1])) {
							follow[val] = struct{}{}
						}
					} else if len(rule.right) == i+1 {
						if _, ok := done[rule.left]; !ok {
							stack = append(stack, string(rule.left))
						}
					}
				}
			}
		}
	}

	var result []string
	for i := range follow {
		result = append(result, i)
	}

	return result
}
func (R RuleGroup) getStateZero() State {
	zeroInstance := Instance{
		rule: R[0],
		dot:  0,
	}.makeId()

	state := InstanceGroup{zeroInstance}.getFullState(R)
	state.i = 0
	fmt.Println(state)
	return state
}

/* INSTANCES */
func (instance Instance) makeId() Instance {
	instance.id = strconv.Itoa(instance.rule.id) + "." + strconv.Itoa(instance.dot)
	return instance
}
func (instance Instance) getNextChar() string {
	if instance.dot < len(instance.rule.right) {
		return string(instance.rule.right[instance.dot])
	}
	return ""
}

func (instances InstanceGroup) containsInstance(I Instance) bool {
	for _, r := range instances {
		if r == I {
			return true
		}
	}
	return false
}
func (instances InstanceGroup) getFullState(allRules RuleGroup) State {

	var result = State{
		instances: instances,
		action:    make(map[string]Action),
	}

	stack := make([]string, 0)
	for _, instance := range instances {
		stack = append(stack, instance.getNextChar())
	}

	var char string
	for ; len(stack) > 0; {
		char, stack = stack[0], stack[1:]
		for _, rule := range allRules.findRulesFor(char) {
			tmpInstance := rule.makeInstance(0)
			if !result.instances.containsInstance(tmpInstance) {
				result.instances = append(result.instances, tmpInstance)
				stack = append(stack, tmpInstance.getNextChar())
			}

		}
	}

	return result
}

/* STATES */
func (state State) getId() string {
	var arr []string
	for _, instance := range state.instances {
		arr = append(arr, instance.id)
	}
	sort.Strings(arr)
	return strings.Join(arr, "-")
}
func (state State) isFinal() bool {
	if len(state.instances) == 1 && len(state.instances[0].rule.right) == state.instances[0].dot {
		return true
	}
	return false
}
func (state State) isAcc() bool {
	if len(state.instances) == 1 && state.instances[0].rule.left == "1" && state.instances[0].rule.right == "S" && state.instances[0].dot == 1 {
		return true
	}
	return false
}
func (state State) getSubSet(char string) InstanceGroup {
	var result InstanceGroup
	for _, instance := range state.instances {
		if instance.getNextChar() == char {
			newInst := Instance{
				rule: instance.rule,
				dot:  instance.dot + 1,
			}.makeId()
			if newInst.dot <= len(newInst.rule.right) {
				result = append(result, newInst)
			}
		}
	}
	return result
}

/* SOLVER */
func (S Solver) Len() int {
	return len(S)
}
func (S Solver) Swap(i, j int) {
	S[i], S[j] = S[j], S[i]
}
func (S Solver) Less(i, j int) bool {
	return S[i].i < S[j].i
}
func (S Solver) solve(str string) bool {

	return true
}
func (S Solver) print(NT []string, T []string) bool {
	sort.Strings(NT)
	fmt.Print("   ")
	for _, char := range append(append(T, "0"), NT...) {
		fmt.Print(" " + char + "  ")
	}
	fmt.Println()

	for _, state := range S {
		fmt.Printf("%2d", state.i)
		for _, char := range append(append(T, "0"), NT...) {
			if val, ok := state.action[char]; ok {
				fmt.Printf(" %2d%s", val.goTo, string(val.action))
			} else {
				fmt.Print("  . ")
			}
		}
		fmt.Println()
	}

	return true
}

/* UTILS */

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
func readKeyboardLine() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	return text
}

func readGrammar(file string) ([]string, []string, RuleGroup) {

	lines, _ := readLines(file)
	nNT, _ := strconv.Atoi(lines[0])
	nT, _ := strconv.Atoi(lines[1])
	nR, _ := strconv.Atoi(lines[2])

	i := 3
	var NT []string
	var T []string
	var R = []Rule{{id: 0, left: "1", right: "S"}}

	for ; i < 3+nNT; i++ {
		NT = append(NT, lines[i])
	}

	for ; i < 3+nNT+nT; i++ {
		T = append(T, lines[i])
	}

	for id := 1; i < 3+nNT+nT+nR; i++ {
		rLine := lines[i]
		ruleData := strings.Split(rLine, "->")
		R = append(R, Rule{
			id:    id,
			left:  ruleData[0],
			right: ruleData[1],
		})
		id++
	}

	return NT, T, R
}

func makeSolver(NT []string, T []string, R RuleGroup) (Solver, string) {
	stateZero := R.getStateZero()
	var states = map[string]State{stateZero.getId(): stateZero}
	stack := make([]State, 1)
	stack[0] = stateZero

	i := 1
	var err string
	var state State
	for ; len(stack) > 0; {
		state, stack = stack[0], stack[1:]
		for _, char := range append(append(NT, T...), "0") {
			subset := state.getSubSet(char)
			if len(subset) > 0 {
				newState := subset.getFullState(R)
				if _, ok := states[newState.getId()]; !ok { // is new
					newState.i = i
					i += 1

					stack = append(stack, newState)
					states[newState.getId()] = newState
				}
				// GOTO / ACTION
				var action byte
				if char[0] >= 'A' && char[0] <= 'Z' {
					action = 'N'
				} else {
					action = 'S'
				}

				newState = states[newState.getId()]
				state.action[char] = Action{
					goTo:   newState.i,
					action: action,
				}

			}
		}
	}

	for _, state := range states {
		if state.isFinal() { // is final state is created
			for _, follow := range R.follow(state.instances[0].rule.left) {
				if val, ok := state.action[follow]; ok {
					if val.action == 'R' {
						return nil, "RR ERROR"
					} else {
						return nil, "SR ERROR"
					}
				}
				state.action[follow] = Action{ // state from which it is created
					goTo: state.instances[0].rule.id,
					action: 'R',
				}
			}
			if state.isAcc() {
				state.action["0"] = Action{
					goTo:   -1,
					action: 'A',
				}
			}
		}
	}

	keys := make([]string, len(states))
	i = 0
	for key := range states {
		keys[i] = key
		i += 1
	}
	result := make(Solver, 0)
	for _, key := range keys {
		result = append(result, states[key])
	}
	sort.Sort(result)

	return result, err
}

func main() {
	NT, T, R := readGrammar("./in3.txt")
	solver, err := makeSolver(NT, T, R)
	if len(err) > 0 {
		fmt.Println(err)
	}

	solver.print(NT, T)
	solver.solve("")

	//fmt.Println( R.follow("1", make(map[string]struct{})))

	fmt.Println("FOLLOW")
	for _, val := range NT {
		fmt.Print(val)
		fmt.Println(R.follow(val))
	}

	fmt.Println("FIRST")
	for _, val := range NT {
		fmt.Println(val, R.first(val))
	}
	//
	//fmt.Println(NT, T, R, "\n=========")
	//for _, state := range solver {
	//	fmt.Println(state)
	//}

	//readKeyboardLine()
}
