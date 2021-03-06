package spec

import (
	"strings"

	"github.com/ddsgok/bdd/colors"
	"github.com/ddsgok/bdd/internal/common"
)

const (
	// OutputNone sets the system to print nothing.
	OutputNone outputType = 1 << iota
	// OutputStdout sets the system to print as usual on Stdout.
	OutputStdout
	// OutputStderr
	// OutputHTML
)

var (
	// config represents the current configuration for system.
	config *Configuration
)

// outputType defines the desired output for the specificationTesting information.
type outputType int

// Configuration defines the configuration used by the package.
type Configuration struct {
	Output outputType

	AnsiOfFeature            string
	AnsiOfGiven              string
	AnsiOfWhen               string
	AnsiOfThen               string
	AnsiOfThenNotImplemented string
	AnsiOfThenWithError      string
	AnsiOfCode               string
	AnsiOfCodeError          string
	AnsiOfExpectedError      string

	assertFn func(*TestSpecification) common.Assert

	LastFeature string
	LastGiven   string
	LastWhen    string
	LastIt      string
}

// ResetLasts sets all last variables to empty. This makes config ready
// to print information about another context.
func (c *Configuration) ResetLasts() {
	c.ResetWhen()
	c.LastGiven = ""
}

// ResetWhen sets when and it last variables to empty. This makes
// config ready to print information about another condition.
func (c *Configuration) ResetWhen() {
	c.ResetIt()
	c.LastWhen = ""
}

// ResetIt sets it last variables to empty. This makes config ready to
// print information about another verification.
func (c *Configuration) ResetIt() {
	c.LastIt = ""
}

// init the configuration and assertions.
func init() {
	ResetConfig()

	// set to verbose output by default
	SetVerbose()

	// register the default Assertions package
	SetAssertionsFn(func(s *TestSpecification) (a common.Assert) {
		a = newAsserter(s)
		return
	})
}

// Config returns current configuration for system.
func Config() *Configuration {
	return config
}

// SetAssertionsFn will assign the assertions used for all tests.
// The specified struct must implement the spec.Assert interface.
//
//    spec.SetAssertionsFn(func(s *TestSpecification) Assert {
//	    return &MyCustomAssertions{}
//    })
func SetAssertionsFn(fn func(s *TestSpecification) common.Assert) {
	config.assertFn = fn
}

// SetConfig takes a Config instance and will be used for all tests
// until ResetConfig() is called.
//
//    spec.SetConfig(Config{
//      AnsiOfFeature: "",	// remove color coding for Feature
//    })
//
//noinspection GoUnusedExportedFunction
func SetConfig(c Configuration) {
	config = &c
}

// ResetConfig will reset all options back to their default configuration.
// Useful for custom colors in the middle of a specification.
func ResetConfig() {
	// setup a default configuration
	config = &Configuration{
		AnsiOfFeature:            strings.Join([]string{colors.White}, ""),
		AnsiOfGiven:              strings.Join([]string{colors.Grey}, ""),
		AnsiOfWhen:               strings.Join([]string{colors.LightGreen}, ""),
		AnsiOfThen:               strings.Join([]string{colors.Green}, ""),
		AnsiOfThenNotImplemented: strings.Join([]string{colors.LightYellow}, ""),
		AnsiOfThenWithError:      strings.Join([]string{colors.RegBg, colors.White, colors.Bold}, ""),
		AnsiOfCode:               strings.Join([]string{colors.Grey}, ""),
		AnsiOfCodeError:          strings.Join([]string{colors.White, colors.Bold}, ""),
		AnsiOfExpectedError:      strings.Join([]string{colors.Red}, ""),
	}
}
