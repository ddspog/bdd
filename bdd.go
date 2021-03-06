package bdd

import (
	"fmt"
	"testing"

	"github.com/ddsgok/bdd/internal/golden"
	"github.com/ddsgok/bdd/spec"
)

// Given defines one Feature's specific context to be tested.
func Given(t *testing.T, given string, args ...interface{}) {
	gTestBodies, gTestCases := split(S(), args)
	whenFunc := gTestBodies.asWhenFunc()

	for _, gArgs := range gTestCases {
		// setup the testspec that we will be using
		testspec := spec.New(t, feature(), printf(given, gArgs))
		testspec.PrintFeature()
		testspec.PrintContext()

		if whenFunc != nil {
			whenFunc(func(when string, args ...interface{}) {
				wTestBodies, wTestCases := split(gArgs, args)
				itFunc := wTestBodies.asItFuncs()

				for _, wArgs := range wTestCases {
					testspec.When = printf(when, wArgs)
					testspec.PrintWhen()

					if itFunc != nil {
						itFunc(func(it string, args ...interface{}) {
							iTestBodies, iTestCases := split(wArgs, args)
							assertFunc := iTestBodies.asAssertFunc()

							for _, iArgs := range iTestCases {
								testspec.It = printf(it, iArgs)
								// It output is handled in the testspec.Run() below

								if assertFunc != nil {
									// Having at least 1 assert means we are implemented

									testspec.AssertFn = func(a Assert) {
										assertFunc(a, iArgs...)
									}

									testspec.NotImplemented = false
								} else {
									testspec.AssertFn = notImplemented()
									testspec.NotImplemented = true
								}

								// Run() handles contextual printing and some delegation
								// to the Assert's implementation for error handling
								testspec.Run()
							}
						}, wArgs...)
					}
				}
			}, gArgs...)
		}

		// reset to default
		spec.Config().ResetLasts()

		if spec.Config().Output != spec.OutputNone {
			fmt.Println()
		}
	}
}

// GivenWithGolden defines one Feature's specific context to be tested.
func GivenWithGolden(t *testing.T, given string, args ...interface{}) {
	goldenFunc := newTestFunc(args...).asGoldenFunc()
	feature := feature()
	gm := golden.NewManager(feature, given)

	if goldenFunc != nil {
		for i := 0; i < gm.NumGoldies(); i++ {
			testspec := spec.New(t, feature, gprintf(given, gm.Get(i)))
			testspec.PrintFeature()
			testspec.PrintContext()

			goldenFunc(func(when string, wTestBodies ...interface{}) {
				itFunc := newTestFunc(wTestBodies...).asItFuncs()
				testspec.When = gprintf(when, gm.Get(i))
				testspec.PrintWhen()

				if itFunc != nil {
					itFunc(func(it string, iTestBodies ...interface{}) {
						assertFunc := newTestFunc(iTestBodies...).asAssertFunc()
						testspec.It = gprintf(it, gm.Get(i))

						if assertFunc != nil {
							testspec.AssertFn = func(a Assert) {
								assertFunc(a)
							}
							testspec.NotImplemented = false
						} else {
							testspec.AssertFn = notImplemented()
							testspec.NotImplemented = true
						}

						testspec.Run()
					})
				}
			}, gm.Get(i))
		}
	}

	gm.Update()

	spec.Config().ResetLasts()
	if spec.Config().Output != spec.OutputNone {
		fmt.Println()
	}
}

// Setup is used to define before/after (setup/teardown) functions.
func Setup(before, after func()) (fn func(fn func(Assert)) func(Assert)) {
	fn = func(fn func(Assert)) func(Assert) {
		before()
		return func(assert Assert) {
			fn(assert)
			after()
		}
	}
	return
}
