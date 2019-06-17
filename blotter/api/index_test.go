package api

import "testing"

func TestIndex(t *testing.T) {

	var tests = []struct {
		in       indexInputType  // input
		expected indexOutputType // expected result
	}{
		{indexInputType{Name: "OhYee"}, indexOutputType{Data: "Hi, OhYee"}},
		{indexInputType{Name: ""}, indexOutputType{Data: "Hi, "}},
		{indexInputType{Name: "Ab1"}, indexOutputType{Data: "Hi, Ab1"}},
	}

	for _, test := range tests {
		actual := Index(test.in)
		if actual != test.expected {
			t.Errorf("[×] in: %+v out: %+v expected: %+v\n", test.in, actual, test.expected)
		} else {
			t.Logf("[√] in: %+v out: %+v expected: %+v\n", test.in, actual, test.expected)
		}
	}
}
