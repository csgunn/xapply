package dicer

import (
	"testing"
)

type DicerTest struct {
	template string
	inputs   []string
	output   string
	err      string
}

func TestDicer(T *testing.T) {
	tests := []DicerTest{
		// Simple expansions
		DicerTest{
			template: "hello %1",
			inputs:   []string{"world"},
			output:   "hello world",
		},
		DicerTest{
			template: "%1 world",
			inputs:   []string{"hello"},
			output:   "hello world",
		},
		DicerTest{
			template: "hello %[1]",
			inputs:   []string{"world"},
			output:   "hello world",
		},
		DicerTest{
			template: "%1 %2",
			inputs:   []string{"hello", "world"},
			output:   "hello world",
		},

		// Auto-append %1
		DicerTest{
			template: "",
			inputs:   []string{"world"},
			output:   "world",
		},
		DicerTest{
			template: "hello",
			inputs:   []string{"world"},
			output:   "hello world",
		},
		DicerTest{
			template: "%%hello",
			inputs:   []string{"world"},
			output:   "%hello world",
		},

		// Multi-character index
		DicerTest{
			template: "%10",
			inputs:   []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			output:   "10",
		},
		DicerTest{
			template: "%[10]",
			inputs:   []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			output:   "10",
		},
		DicerTest{
			template: "%10test",
			inputs:   []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			output:   "10test",
		},

		// Must be one input
		DicerTest{
			template: "",
			inputs:   []string{},
			err:      "at least one input must be specified",
		},

		// Out of bound index
		DicerTest{
			template: "hello %2",
			inputs:   []string{"a"},
			err:      "index 2: out of bounds (inputs size 1)",
		},

		// Missing ]
		DicerTest{
			template: "xhello %[1 world",
			inputs:   []string{"hello"},
			err:      "char 8: dicer expression missing closing ]",
		},

		DicerTest{
			template: "hello %[",
			inputs:   []string{"a"},
			err:      "char 7: dicer expression missing closing ]",
		},

		// %% escape
		DicerTest{
			template: "test%%%1",
			inputs:   []string{"hello"},
			output:   "test%hello",
		},

		// EOL %
		DicerTest{
			template: "test %1 %",
			inputs:   []string{"a"},
			output:   "test a %",
		},

		// % followed by a non-number
		DicerTest{
			template: "test %1 %q",
			inputs:   []string{"a"},
			output:   "test a %q",
		},

		// dice selection
		DicerTest{
			template: "%[1/2]",
			inputs:   []string{"1/2/3"},
			output:   "2",
		},

		// multi-level dice selections
		DicerTest{
			template: "%[1/2,1]",
			inputs:   []string{"1,a/2,b/3,c"},
			output:   "2",
		},

		// dice removal
		DicerTest{
			template: "%[1.-2]",
			inputs:   []string{"aa.bb.cc.dd"},
			output:   "aa.cc.dd",
		},
		DicerTest{
			template: "%[1.-2.-$]",
			inputs:   []string{"aa.bb.cc.dd"},
			output:   "aa.cc",
		},
		DicerTest{
			template: "%[1.-2.-$.-1]",
			inputs:   []string{"aa.bb.cc.dd"},
			output:   "cc",
		},
		DicerTest{
			template: "%[1.-2]z",
			inputs:   []string{"a..c"},
			output:   "a.cz",
		},
		DicerTest{
			template: "%[1.-4]",
			inputs:   []string{"test.example.com."},
			output:   "test.example.com",
		},
		DicerTest{
			template: "%[1.-$]",
			inputs:   []string{"test.example.com."},
			output:   "test.example.com",
		},

		// TODO: support '@' joining
		/*
			DicerTest{
				template: "%[1/@.2]",
				inputs:   []string{"1.2.3.4/A.B.C.D/I.II.III.IV"},
				output:   "2/B/II",
			},
			DicerTest{
				template: "%[1/@.-2]",
				inputs:   []string{"1.2.3.4/A.B.C.D/I.II.III.IV"},
				output:   "1.3.4/A.C.D/I.III.IV",
			},
			DicerTest{
				template: "%[1/@.-$].four",
				inputs:   []string{"1.2.3.4/A.B.C/I.II"},
				output:   "1.2.3/A.B/I",
			},
			DicerTest{
				template: "%[ @.2]",
				inputs:   []string{"1.2.3.4   A.B.C.D  I.II.III.IV"},
				output:   "2 B II",
			},
			DicerTest{
				template: "%[ @]",
				inputs:   []string{"1.2.3.4   A.B.C.D  I.II.III.IV"},
				output:   "1.2.3.4 A.B.C.D I.II.III.IV",
			},
		*/

		// dice removal of an index that doesn't exist
		DicerTest{
			template: "%[1.-5]",
			inputs:   []string{"aa.bb.cc.dd"},
			output:   "aa.bb.cc.dd",
		},

		// dice selection and removal
		DicerTest{
			template: "%[1/2.-2]",
			inputs:   []string{"aa.bb.cc/dd.ee.ff"},
			output:   "dd.ff",
		},

		// $ dice position
		DicerTest{
			template: "%[1.$]",
			inputs:   []string{"aa.bb.cc.dd"},
			output:   "dd",
		},

		// $ dice position removal
		DicerTest{
			template: "%[1.-$]",
			inputs:   []string{"aa.bb.cc.dd"},
			output:   "aa.bb.cc",
		},

		// dice position out of bounds
		DicerTest{
			template: "%[1/4]",
			inputs:   []string{"1/2/3"},
			output:   "",
		},
	}

	for _, test := range tests {
		output, err := Expand(test.template, test.inputs)

		if err != nil && test.err == "" {
			T.Errorf("template %v with %v: got unexpected error %q", test.template, test.inputs, err.Error())
		} else if test.err != "" {
			if err == nil {
				T.Errorf("template %v with %v: expected error %q, did not get any error", test.template, test.inputs, test.err)
			} else if err.Error() != test.err {
				T.Errorf("template %v with %v: expected error %q, got %q", test.template, test.inputs, test.err, err.Error())
			}
		} else if output != test.output {
			T.Errorf("template %v with %v: expected %v, got %v", test.template, test.inputs, test.output, output)
		}
	}
}
