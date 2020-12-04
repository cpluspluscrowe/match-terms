package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"bufio"
	"os"
	"sort"
)

func getFileContents() string {
	contents, _ := ioutil.ReadFile("contexts.txt")
	var s string = string(contents)
	return s
}

func getMatches(contents, regex string) []string {
	re := regexp.MustCompile(regex)
	return re.FindAllString(contents, -1)
}

func replaceQuote(old string) string {
	updated := strings.Replace(old, "\"", "", 2)
	return updated
}

func splitIntoWords(toSplit string) []string {
	return strings.Split(toSplit, "_")
}

func getUniqueMatches(matches []string) map[string]bool {
	set := make(map[string]bool)
	for _, s := range matches {
		noQuote := replaceQuote(s)
		set[noQuote] = true
	}
	return set
}

func readTargets() string {
	contents, _ := ioutil.ReadFile("targets.txt")
	var s string = string(contents)
	return s
}

func prepareTargets(targetString string) [][]string {
	targets := [][]string{}
	var lines []string = strings.Split(targetString, "\n")
	for _, line := range lines {
		pieces := strings.Split(line, " ")
		targets = append(targets, pieces)
	}
	return targets
}

func prepareInputs(uniqueMatches map[string]bool) [][]string {
	inputs := [][]string{}
	for match, _ := range uniqueMatches {
		inputs = append(inputs, splitIntoWords(match))
	}
	return inputs
}

func contains(searchWithin []string, lookingFor string) bool {
	for _, searchValue := range searchWithin {
		if searchValue == lookingFor {
			return true
		}
	}
	return false
}

func deleteEmptyCounts(counts map[string]int) map[string]int {
	for key, amount := range counts {
		if amount <= 0 {
			delete(counts, key)
		}
	}
	return counts
}

func countMatches(target_array []string, inputs [][]string) map[string]int {
	counts := make(map[string]int)
	for _, input_array := range inputs {
		input_string := strings.Join(input_array, "_")
		counts[input_string] = 0
		for _, target := range target_array {
			for _, input := range input_array {
				if target == input {
					counts[input_string] = counts[input_string] + 1
				}
			}
		}
	}
	return deleteEmptyCounts(counts)
}

func reverseMap(counts map[string]int) map[int]string {
	byCount := map[int]string{}
	for key, value := range counts {
		byCount[value] = key
	}
	return byCount
}

func getKeys(data map[int]string) []int {
	keys := []int{}
	for key, _ := range data {
		keys = append(keys, key)
	}
	return keys
}

func matchTheInputsAndTargets(targets, inputs [][]string) []string {
	answers := []string{}
	for _, target_array := range targets {
		target_string := strings.Join(target_array, "_")
		counts := reverseMap(countMatches(target_array, inputs))
		keys := getKeys(counts)
		sort.Ints(keys)
		if len(counts) > 0 {			
			fmt.Println(target_string, ":")
			for _, key := range keys {
				input_key := counts[key]
				fmt.Println("\t", input_key, key)
			}
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			answers = append(answers, target_string + "|" + text)
		}
	}
	return answers
}

func main() {
	regex := "\".+?\""
	targets := prepareTargets(readTargets())
	text := getFileContents()
	matches := getMatches(text, regex)
	uniqueMatches := getUniqueMatches(matches)
	inputs := prepareInputs(uniqueMatches)
	answers := matchTheInputsAndTargets(targets, inputs)
	for _, answer := range answers {
		fmt.Println(answer)
	}
}
