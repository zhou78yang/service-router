package policy

import "testing"

func TestGenPolicyKey(t *testing.T) {
	cases := []struct {
		args   []string
		expect string
	}{
		{args: []string{"service1", "cluster1", "tenant1"}, expect: "service1:cluster1:tenant1"},
		{args: []string{"service1", "cluster1", ""}, expect: "service1:cluster1"},
		{args: []string{"service1", "", "tenant1"}, expect: "service1"},
		{args: []string{"service1", "", ""}, expect: "service1"},
		{args: []string{"", "cluster1", "service1"}, expect: ""},
		{args: []string{"", "", ""}, expect: ""},
	}

	for i, c := range cases {
		t.Logf("[case-%d] args: %v, expect: \"%s\"", i, c.args, c.expect)
		result := genPolicyKey(c.args...)
		if result != c.expect {
			t.Errorf("[case-%d]unexcept result, expect \"%s\", got \"%s\"\n", i, c.expect, result)
		}
	}
}
