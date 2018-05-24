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

type State struct {
	id        string
	instances InstanceGroup
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

func (state State) getSubSet(char string) InstanceGroup {
	var result InstanceGroup
	for _, instance := range state.instances {
		if instance.getNextChar() == char {
			newInst := Instance{
				rule: instance.rule,
				dot:  instance.dot + 1,
			}.makeId()
			if newInst.dot < len(newInst.rule.right) {
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
	return instance.rule.right[instance.dot : instance.dot+1]
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

func readGrammar(file string) ([]string, []string, []Rule) {

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

func makeSolver(NT []string, T []string, R []Rule) []State {
	stateZero := getStateZero(R)
	var states = map[string]State{stateZero.id: stateZero}
	stack := make([]State, 1)
	stack[0] = stateZero

	var state State
	for ; len(stack) > 0; {
		state, stack = stack[0], stack[1:]
		for _, char := range append(NT, T...) {
			subset := state.getSubSet(char)
			if len(subset) > 0 {
				newState := subset.getFullState(R)
				if _, ok := states[newState.id]; !ok {
					states[newState.id] = newState
					stack = append(stack, newState)
				}
			}
		}

	}

	var keys []string
	for k := range states {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var result []State
	for _, k := range keys {
		result = append(result, states[k])
	}

	return result
}

func main() {
	NT, T, R := readGrammar("./in1.txt")
	solver := makeSolver(NT, T, R)
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
