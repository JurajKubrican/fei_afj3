package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"strings"
)

type rule struct {
	left  string
	right string
}

type state struct {
	id    string
	rules map[rule]struct{}
}

func getStateId(rules []rule) string(){

}

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

func readGrammar(file string) ([]string, []string, []rule) {

	lines, _ := readLines(file)
	nNT, _ := strconv.Atoi(lines[0])
	nT, _ := strconv.Atoi(lines[1])
	nR, _ := strconv.Atoi(lines[2])

	i := 3
	var NT []string
	var T []string
	var R = []rule{{left: "1", right: "S"}}

	for ; i < 3+nNT; i++ {
		NT = append(NT, lines[i])
	}

	for ; i < 3+nNT+nT; i++ {
		T = append(T, lines[i])
	}

	for ; i < 3+nNT+nT+nR; i++ {
		rLine := lines[i]
		ruleData := strings.Split(rLine, "->")
		R = append(R, rule{
			left:  ruleData[0],
			right: ruleData[1],
		})
	}

	return NT, T, R
}

func main() {
	NT, T, R := readGrammar("./in1.txt")
	makeClassifier()

	fmt.Println(NT, T, R)
}
