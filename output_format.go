package main

import (
	"fmt"
	"strings"
)

const (
	OutputFormatMarkdownUnorderedList = "markdown-unordered_list"
	OutputFormatMarkdownTaskList      = "markdown-task_list"
	OutputFormatRegex                 = "regex"
)

func validateOutputFormatFlag(f string) error {
	formatCatalog := []string{
		OutputFormatMarkdownUnorderedList,
		OutputFormatMarkdownTaskList,
		OutputFormatRegex,
	}
	isValidFlag := false

	for _, expected := range formatCatalog {
		isValidFlag = f == expected
		if isValidFlag {
			break
		}
	}
	if !isValidFlag {
		return fmt.Errorf("-output-format flag is invalid format; (%s) required but got %s", strings.Join(formatCatalog, "|"), f)
	}

	return nil
}
