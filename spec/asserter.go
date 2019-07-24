package spec

import (
	"fmt"
	"strings"

	"github.com/ddsgok/bdd/internal/assert"
	"github.com/ddsgok/bdd/internal/common"
)

// specificationTesting is the glue used to bind assertions to the
// spec package as the default asserting package.
//
// specificationTesting implements the common.Tester interface with
// Errorf(...) func.
//
// It is used to print the specification portion to the output when an
// error occurs.  Also, it sets a flag that is used by the bdd
// framework to know that an error has been printed and therefore do
// not print a normal specification title.
//
// Errors are only handled this way under one condition: that is that
// Errorf() be executed by the Assertion package.  Else, we do not get
// the flag to know.
//
// The current Testify assert package fires Errorf() on every Fail(),
// which all asserts will fire when an error occurs.  So, we just wrap
// that below.
type specificationTesting struct {
	spec *TestSpecification
}

// Errorf is called by the Testify's assertions to signal a fail condition.
// This specific method sets an internal spec flag so that the framework
// is aware that an error occurred.
func (m *specificationTesting) Errorf(format string, args ...interface{}) {
	// because we control the output of specification, we
	// need to store these details in a state for later use in
	// the bdd framework.  to do that, we use the
	// m.spec.AssertionFailed boolean.
	m.spec.AssertionFailed = true

	// parse out Testify's location info by removing the first
	// line and reformat their Error message to our liking
	// using string foo
	err := fmt.Sprintf(format, args...)
	err = strings.Replace(err, "\r", "", -1)
	err = strings.Replace(err, "        ", "\t\t\t", -1) // some errors are two-liners
	lines := strings.Split(err, "\n")
	out := ""

	for i := range lines {
		if !strings.Contains(lines[i], "Location:") && lines[i] != "" {
			if out == "" {
				out = lines[i]
			} else {
				out = strings.Join([]string{out, "\n", lines[i]}, "")
			}
		}
	}

	m.spec.PrintItWithError()
	// to properly set the caller used, we currently need to call
	// m.spec.PrinterError here to capture the proper line number.
	// and since PrintError() comes after PrintTitleWithError(),
	// we have the line above.
	//
	// TODO refactor to pass the caller information down along with
	// the custom error message parsing.  that way we can control the
	// printing internally and seal up these Print*() messages.
	m.spec.PrintError(out)
}

// newAsserter constructs a wrapper around Testify's asserts.
func newAsserter(s *TestSpecification) (a common.Assert) {
	a = assert.New(&specificationTesting{
		spec: s,
	})
	return
}
