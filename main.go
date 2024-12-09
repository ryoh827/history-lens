package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	historyFile := os.Getenv("HOME") + "/.zsh_history"

	file, err := os.Open(historyFile)
	if err != nil {
		fmt.Printf("Error opening history file: %v\n", err)
		return
	}
	defer file.Close()

	commandFreq := make(map[string]int)
	commandWithArgsFreq := make(map[string]int)

	// 正規表現で ": timestamp:ID;command" を解析
	re := regexp.MustCompile(`^: \d+:\d+;(.+)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // 空行をスキップ
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue // マッチしない場合はスキップ
		}

		fullCommand := matches[1]
		commandParts := strings.Fields(fullCommand)
		if len(commandParts) > 0 {
			commandFreq[commandParts[0]]++          // コマンド名の頻度
			commandWithArgsFreq[fullCommand]++     // 引数込みのコマンド頻度
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading history file: %v\n", err)
		return
	}

	// 結果を表示
	fmt.Println("Top Commands by Frequency:")
	printSortedResults(commandFreq)

	fmt.Println("\nTop Commands with Arguments by Frequency:")
	printSortedResults(commandWithArgsFreq)
}

func printSortedResults(freqMap map[string]int) {
	type commandCount struct {
		Command string
		Count   int
	}
	var sortedCommands []commandCount
	for cmd, count := range freqMap {
		sortedCommands = append(sortedCommands, commandCount{cmd, count})
	}
	sort.Slice(sortedCommands, func(i, j int) bool {
		return sortedCommands[i].Count > sortedCommands[j].Count
	})
	for i, cmd := range sortedCommands {
		if i >= 20 { // 上位10件を表示
			break
		}
		fmt.Printf("%d. %s (%d times)\n", i+1, cmd.Command, cmd.Count)
	}
}

