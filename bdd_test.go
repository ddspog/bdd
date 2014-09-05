package gomspec

import (
	"testing"
)

func Test_Bdd_Specifications(t *testing.T) {

	Given(t, "a unique scenerio", func(when When) {

		when("an event occurs", func(it It) {
			it("should evaluate 1s are equal", func(assert Assert) {

			})

			it("should also evaluate 3 and 4 are not equal", func(assert Assert) {
				assert.False(true, "Well darn!")
			})

			it("should perform another evaluation", func(assert Assert) {
				assert.Contains("shop", "gogo")
				assert.Equal(1, 2, "It fucked up")
			})

			it("should not have this implemented", NA())

			it("should also perform another evaluation", func(eassert Assert) {
				//expect("hellow").ToNotEqual("world")
			})
		})
	})

	Given(t, "a scenario that needs a setup and teardown", func(when When) {

		count := 0

		when("using the Setup() extension", func(it It) {

			before := func() {
				count++
			}

			after := func() {
				count--
			}

			setup := Setup(before, after)

			it("should increment count to 1", setup(func(assert Assert) {
				//expect(count).ToEqual(1)
			}))

			if count != 0 {
				t.Error("In BDD-specs, count should have been reset to zero by the teardown func")
			}
		})
	})

}
