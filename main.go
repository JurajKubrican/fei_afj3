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
type Solver map[string]State

type Action struct {
	goTo   string
	action byte
	reduce int
}

type State struct {
	id        string
	instances InstanceGroup
	action    map[string]Action
}

func (S Solver) solve(str string) bool {

	return true
}

func (state State) makeId() State {
	var arr []string
	for _, instance := range state.instances {
		arr = append(arr, instance.id)
	}
	sort.Strings(arr)
	state.id = strings.Join(arr, "-")
	return state
}

func (state State) isFinal() bool {
	if len(state.instances) == 1 && len(state.instances[0].rule.right) == state.instances[0].dot+1 {
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

func (R RuleGroup) follow(inChar string, done map[string]struct{}) []string {
	if _, ok := done[inChar]; ok {
		return make([]string, 0)
	}
	done[inChar] = struct{}{}

	follow := make(map[string]struct{})
	if inChar == "1" {
		return []string{"0"}
	}
	for _, rule := range R {
		for i, char := range rule.right {
			if inChar == string(char) && len(rule.right) > i+1 {
				for _, val := range R.first(string(rule.right[i+1])) {
					follow[val] = struct{}{}
				}
			} else if len(rule.right) == i+1 && char >= 'A' && char <= 'Z' {
				//for _, val := range R.follow(rule.left, done) {
				//	follow[val] = struct{}{}
				//}
			}
		}
	}

	var result []string
	for i := range follow {
		result = append(result, i)
	}

	return result
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

	return result.makeId()
}

func getStateZero(allRules []Rule) State {
	zeroInstance := Instance{
		rule: allRules[0],
		dot:  0,
	}.makeId()

	state := InstanceGroup{zeroInstance}.getFullState(allRules)
	fmt.Println(state)
	return state.makeId()
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
	stateZero := getStateZero(R)
	var states = Solver{stateZero.id: stateZero}
	stack := make([]State, 1)
	stack[0] = stateZero

	var err string
	var state State
	for ; len(stack) > 0; {
		state, stack = stack[0], stack[1:]
		for _, char := range append(append(NT, T...), "0") {
			subset := state.getSubSet(char)
			if len(subset) > 0 {
				newState := subset.getFullState(R)
				if _, ok := states[newState.id]; !ok {

					// GOTO / ACTION
					var action byte
					if char[0] >= 'A' && char[0] <= 'Z' {
						action = 'N'
					} else {
						action = 'S'
					}

					if state.isFinal() {
						for _, follow := range R.follow(state.instances[0].rule.left, make(map[string]struct{})) {
							if val, ok := state.action[follow]; ok {
								if val.action == 'R' {
									return nil, "RR ERROR"
								} else {
									return nil, "SR ERROR"
								}
							}
							state.action[follow] = Action{
								reduce: state.instances[0].rule.id,
								action: 'R',
							}
						}
					}

					state.action[char] = Action{
						goTo:   newState.id,
						action: action,
					}
					stack = append(stack, newState)
					states[newState.id] = newState
				}
			}
		}
	}

	return states, err
}

func main() {
	NT, T, R := readGrammar("./in3.txt")
	solver, err := makeSolver(NT, T, R)
	if len(err) > 0 {
		fmt.Println(err)
	}

	solver.solve("")

	//fmt.Println( R.follow("1", make(map[string]struct{})))

	fmt.Println("FOLLOW")
	for _, val := range NT {
		fmt.Println(val, R.follow(val, make(map[string]struct{})))
	}
	//
	//fmt.Println("FIRST")
	//for _, val := range NT {
	//	fmt.Println(val, R.first(val))
	//}

	fmt.Println(NT, T, R, "\n=========")
	for _, state := range solver {
		fmt.Println(state)
	}
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
