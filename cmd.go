package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Cmd describes the command needed to be runned and the rules used to choose
// the selected contexts.
type Cmd struct {
	// Command to run.
	Command []string

	// Context rules.
	Matchs []string
}

// Execute found all the contexts matching the context rules and run the command
// agains them.
func Execute(cmd *Cmd) error {
	// Compose the context list
	contexts, err := RetrieveAllContexts()
	if err != nil {
		return err
	}
	if len(cmd.Matchs) > 0 {
		contexts = filterWithMatchs(contexts, cmd.Matchs)
	}

	// Run the commands
	outputs := make([]OutputDescription, 0, len(contexts))
	for _, context := range contexts {
		args := []string{"--context=" + context}
		args = append(args, cmd.Command...)

		kubeCmd := exec.Command("kubectl", args...)

		// Can be used for debugging purpose with not kubernetes available.
		//kubeCmd := exec.Command(cmd.Command[0], cmd.Command[1:]...)
		cmdOut, err := kubeCmd.StdoutPipe()
		if err != nil {
			return err
		}

		go func() {
			err := kubeCmd.Start()
			if err != nil {
				fmt.Println(err)
			}
		}()

		title := []string{"kubectl"}
		title = append(title, args...)
		outputs = append(outputs, OutputDescription{
			Title: strings.Join(title, " "),
			In:    cmdOut,
		})
	}

	return RunGUIOutput(outputs)
}

func filterWithMatchs(contexts []string, matchs []string) []string {
	// Use a map[string]struct{} in order to create a base Set data structure.
	// It allows to avoid ady duplicates.
	set := make(map[string]struct{}, len(contexts))

	for _, match := range matchs {
		for _, context := range contexts {
			matched, err := regexp.MatchString(match, context)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if matched {
				set[context] = struct{}{}
			}
		}
	}

	res := make([]string, 0, len(contexts))
	for context := range set {
		res = append(res, context)
	}

	return res
}
