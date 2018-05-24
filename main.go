package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"strings"
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

func (state State) getStateId() State {
	return state
}

func (instance Instance) getInstaceId() Instance {
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
	}.getInstaceId()
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

func (R InstanceGroup) containsInstance(I Instance) bool {
	for _, r := range R {
		if r == I {
			return true
		}
	}
	return false
}

func getFullState(instances []Instance, allRules RuleGroup) State {

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

	return result.getStateId()
}

func getStateZero(allRules []Rule) State {
	zeroInstance := Instance{
		rule: allRules[0],
		dot:  0,
	}.getInstaceId()

	state := getFullState([]Instance{zeroInstance}, allRules)
	fmt.Println(state)
	return state
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



func main() {
	NT, T, R := readGrammar("./in1.txt")
	stateZero := getStateZero(R)
	fmt.Println(NT, T, R, stateZero)
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
