package gomspec

import (
	"fmt"
	"github.com/eduncan911/gomspec/colors"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
)

type formatter interface {
	PrintFeature()
	PrintContext()
	PrintWhen()
	PrintTitle()
	PrintTitleNotImplemented()
	PrintTitleWithError()
	PrintError(string)
}

type failingLine struct {
	prev     string
	content  string
	next     string
	filename string
	number   int
	lines    []string
}

func (spec *Specification) PrintFeature() {
	if MSpec.lastFeature == spec.Feature {
		return
	}
	fmt.Printf("\n%sFeature: %s%s\n", MSpec.AnsiOfFeature, spec.Feature, colors.Reset)
	MSpec.lastFeature = spec.Feature
}

func (spec *Specification) PrintContext() {
	if MSpec.lastGiven == spec.Given {
		return
	}
	fmt.Printf("\n%s  Given %s%s\n", MSpec.AnsiOfGiven, padLf(spec.Given, 2), colors.Reset)
	MSpec.lastGiven = spec.Given
}

func (spec *Specification) PrintWhen() {
	if MSpec.lastWhen == spec.When {
		return
	}
	fmt.Printf("\n%s    When %s%s\n", MSpec.AnsiOfWhen, spec.When, colors.Reset)
	MSpec.lastWhen = spec.When
}

func (spec *Specification) PrintSpec() {
	fmt.Printf("%s    » It %s %s\n", MSpec.AnsiOfThen, spec.Spec, colors.Reset)
	MSpec.lastSpec = spec.Spec
}

func (spec *Specification) PrintTitleWithError() {
	if MSpec.lastSpec == spec.Spec {
		return
	}
	fmt.Printf("%s    » It %s %s\n", MSpec.AnsiOfThenWithError, spec.Spec, colors.Reset)
	MSpec.lastSpec = spec.Spec
}

func (spec *Specification) PrintTitleNotImplemented() {
	fmt.Printf("%s    » It %s «-- NOT IMPLEMENTED%s\n", MSpec.AnsiOfThenNotImplemented, spec.Spec, colors.Reset)
	MSpec.lastSpec = spec.Spec
}

func (spec *Specification) PrintError(message string) {
	failingLine, err := getFailingLine()

	if err != nil {
		return
	}

	fmt.Printf("%s%s%s\n", MSpec.AnsiOfExpectedError, message, colors.Reset)
	fmt.Printf("%s        in %s:%d%s\n", MSpec.AnsiOfCode, path.Base(failingLine.filename), failingLine.number, colors.Reset)
	fmt.Printf("%s        ---------\n", MSpec.AnsiOfCode)
	spec.PrintFailingLine(&failingLine)
	spec.T.Fail()
}

func (spec *Specification) PrintFailingLine(failingLine *failingLine) {
	fmt.Printf("%s        %d. %s%s\n", MSpec.AnsiOfCode, failingLine.number-1, softTabs(failingLine.prev), colors.Reset)
	fmt.Printf("%s        %d. %s %s\n", MSpec.AnsiOfCodeError, failingLine.number, failingLine.content, colors.Reset)
	fmt.Printf("%s        %d. %s%s\n", MSpec.AnsiOfCode, failingLine.number+1, softTabs(failingLine.next), colors.Reset)
	fmt.Println()
}

func getFailingLine() (failingLine, error) {

	// this entire func is now a hack because of where it is being called,
	// which is now one caller higher.  previously it was being called in the
	// Expect struct which had the right caller info.  but now, it is being
	// called after the Assertion has been executed to print details to the
	// string.

	_, filename, ln, _ := runtime.Caller(5)

	// TODO: this is really hacky, need to find a way of not using magic numbers for runtime.Caller
	// If we are not in a test file, we must still be inside this package,
	// so we need to go up one more stack frame to get to the test file
	if !strings.HasSuffix(filename, "_test.go") {
		_, filename, ln, _ = runtime.Caller(6)
	}

	bf, err := ioutil.ReadFile(filename)

	if err != nil {
		return failingLine{}, fmt.Errorf("Failed to open %s", filename)
	}

	lines := strings.Split(string(bf), "\n")[ln-2 : ln+2]

	return failingLine{
		prev:     softTabs(lines[0]),
		content:  softTabs(lines[1]),
		next:     softTabs(lines[2]),
		filename: filename,
		number:   int(ln),
	}, nil

}

func softTabs(text string) string {
	return strings.Replace(text, "\t", "  ", -1)
}

func padLf(strToPad string, padding int) string {
	pad := func() string {
		s := "\n"
		for i := 0; i < padding; i++ {
			s = strings.Join([]string{s, " "}, "")
		}
		return s
	}
	return strings.Replace(
		strToPad,
		"\n",
		pad(),
		-1,
	)
}
