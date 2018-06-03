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

type StackAction struct {
	char  string
	state int
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
	done := make(map[int]struct{})
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
			if _, ok := done[rule.id]; !ok {
				stack = append(stack, rules...)
				done[rule.id] = struct{}{}
			}
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
					if len(rule.right) == i+1 { // final char -> follow of previous
						if _, ok := done[rule.left]; !ok {
							stack = append(stack, string(rule.left))
						}
					}
					for j := i + 1; len(rule.right) > j; j++ {
						hasEps := false
						for _, val := range R.first(string(rule.right[j])) { //add all firsts of
							if val == "0" {
								hasEps = true
							} else {
								follow[val] = struct{}{}
							}
						}
						if !hasEps {
							break
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
	var zeroInstance Instance
	if R[0].right == "0" {
		zeroInstance = Instance{
			rule: R[0],
			dot:  1,
		}.makeId()
	} else {
		zeroInstance = Instance{
			rule: R[0],
			dot:  0,
		}.makeId()
	}

	state := InstanceGroup{zeroInstance}.getFullState(R)
	state.i = 0
	return state
}
func (R RuleGroup) printResult() {
	if len(R) == 0 {
		fmt.Println("Nepatri do jazyka")
	} else {
		for i := len(R) - 1; i >= 0; i -= 1 {
			fmt.Println(R[i].left + "->" + R[i].right)
		}
	}

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
			var tmpInstance Instance
			if rule.right == "0" {
				tmpInstance = rule.makeInstance(1)
			} else {
				tmpInstance = rule.makeInstance(0)
			}
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
func (state State) getFinal() ([]int, bool) {
	var result []int
	for index, instance := range state.instances {
		if len(instance.rule.right) == instance.dot {
			result = append(result, index)
		}
	}
	return result, len(result) > 0
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
func (S Solver) solve(input string, R RuleGroup) RuleGroup {

	var result RuleGroup
	input += "0"
	stack := []StackAction{{char: input[0:1], state: 0}}
	curSate := 0

	for ; ; {
		top := input[0:1]
		if action, ok := S[curSate].action[top]; ok {
			switch action.action {
			case 'S':
				curSate = action.goTo
				stack = append(stack, StackAction{char: top, state: curSate})
				input = input[1:]
				break

			case 'R':
				rule := R[action.goTo]
				if rule.right != "0" {
					stack = stack[:len(stack)-len(rule.right)] // shorten stack
				}

				newAction := StackAction{
					char:  rule.left,
					state: S[stack[len(stack)-1].state].action[rule.left].goTo,
				}
				top = stack[len(stack)-1].char
				stack = append(stack, newAction)

				curSate = stack[len(stack)-1].state

				result = append(result, rule)
				break

			case 'A':
				if top == "0" {
					return result
				} else {
					return RuleGroup{}
				}
			}

		} else {
			break
		}
		//fmt.Println(stack)
	}

	if len(stack) == 1 {
		return result
	} else {
		return RuleGroup{}
	}
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

				state.action[char] = Action{
					goTo:   states[newState.getId()].i,
					action: action,
				}

			}
		}
	}

	for _, state = range states {

		if finalRules, hasFinal := state.getFinal(); hasFinal { // is final state is created
			for _, finalRule := range finalRules {
				for _, follow := range R.follow(state.instances[finalRule].rule.left) {
					if val, ok := state.action[follow]; ok {
						if val.action == 'R' {
							return nil, "Konflikt redukcia-redukcia"
						} else {
							return nil, "Konflikt presun-redukcia"
						}
					}
					state.action[follow] = Action{ // state from which it is created
						goTo: state.instances[finalRule].rule.id,
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

func readKeyboardLine(query string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(query)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text[:len(text)-1])
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

func main() {
	var file string
	if len(os.Args) < 2 {
		file = readKeyboardLine("Zadajte meno vstupneho suboru:")
		//file = "in1.txt"
	} else {
		file = os.Args[1]
	}
	NT, T, R := readGrammar(file)

	//fmt.Println("FIRST")
	//for _, val := range NT {
	//	fmt.Println(val, R.first(val))
	//}
	//fmt.Println("FOLLOW")
	//for _, val := range NT {
	//	fmt.Print(val)
	//	fmt.Println(R.follow(val))
	//}

	solver, err := makeSolver(NT, T, R)
	if len(err) > 0 {
		fmt.Println(err)
		return
	}

	//solver.print(NT, T)

	query := "Zadajte slovo, (iba enter pre ukoncenie): "
	for text := readKeyboardLine(query); len(text) > 1; text = readKeyboardLine(query) {
		result := solver.solve(text, R)
		result.printResult()
	}
	
}
