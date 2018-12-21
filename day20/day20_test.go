package day20

import (
	"reflect"
	"strings"
	"testing"
)

func TestLiteral_Enumerate(t *testing.T) {
	tables := []struct{
		test Expr
		prefix string
		visits []string
		finals []string
	}{
		{
			Literal(""),
			"",
			[]string{},
			[]string{""},
		},
		{
			Literal("ABCD"),
			"",
			[]string{"A", "AB", "ABC", "ABCD"},
			[]string{"ABCD"},
		},

		{
			Literal("ABCD"),
			"XYZ",
			[]string{"XYZA", "XYZAB", "XYZABC", "XYZABCD"},
			[]string{"XYZABCD"},
		},
	}

	for _, table := range tables {
		visits := make([]string, 0)
		f := func(s string) {
			visits = append(visits, s)
		}
		finals := table.test.Enumerate(table.prefix, f)
		if !reflect.DeepEqual(visits, table.visits) {
			t.Errorf("%s visits: expected %v, got %v", table.test, table.visits, visits)
		}
		if !reflect.DeepEqual(finals, table.finals) {
			t.Errorf("%s finals: expected %v, got %v", table.test, table.finals, finals)
		}
	}
}

func TestExpr_String(t *testing.T) {
	tables := []struct{
		test Expr
		expected string
	}{
		{
			Literal("ABCD"),
			"ABCD",
		},
		{
			Choice{Literal("AB"), Literal("CD")},
			"(AB|CD)",
		},
		{
			Sequence{Choice{Literal("AB"), Literal("CD")}, Literal("EF")},
			"(AB|CD)EF",
		},
	}

	for _, table := range tables {
		got := table.test.String()
		if got != table.expected {
			t.Errorf("expected %v, got %v", table.expected, got)
		}
	}
}

func TestChoice_Enumerate(t *testing.T) {
	tables := []struct{
		test Expr
		prefix string
		visits []string
		finals []string
	}{
		{
			Choice{},
			"",
			[]string{},
			[]string{},
		},
		{
			Choice{Literal("AB"), Literal("CD")},
			"",
			[]string{"A", "AB", "C", "CD"},
			[]string{"AB", "CD"},
		},

		{
			Choice{Literal("AB"), Literal("CD")},
			"XYZ",
			[]string{"XYZA", "XYZAB", "XYZC", "XYZCD"},
			[]string{"XYZAB", "XYZCD"},
		},
		{
			Choice{Literal("AB"), Choice{Literal("CD"), Literal("EF")}},
			"",
			[]string{"A", "AB", "C", "CD", "E", "EF"},
			[]string{"AB", "CD", "EF"},
		},
		{
			Choice{Choice{Literal("CD"), Literal("EF")}, Literal("AB")},
			"",
			[]string{"C", "CD", "E", "EF", "A", "AB"},
			[]string{"CD", "EF", "AB"},
		},
	}

	for _, table := range tables {
		visits := make([]string, 0)
		f := func(s string) {
			visits = append(visits, s)
		}
		finals := table.test.Enumerate(table.prefix, f)
		if !reflect.DeepEqual(visits, table.visits) {
			t.Errorf("%s visits: expected %v, got %v", table.test, table.visits, visits)
		}
		if !reflect.DeepEqual(finals, table.finals) {
			t.Errorf("%s finals: expected %v, got %v", table.test, table.finals, finals)
		}
	}
}

func TestSequence_Enumerate(t *testing.T) {
	tables := []struct{
		test Expr
		prefix string
		visits []string
		finals []string
	}{
		{
			Sequence{},
			"",
			[]string{},
			[]string{""},
		},
		{
			Sequence{Literal("AB"), Literal("CD")},
			"",
			[]string{"A", "AB", "ABC", "ABCD"},
			[]string{"ABCD"},
		},
		{
			Sequence{Literal("A"), Choice{Literal("B"), Literal("C")}, Literal("D")},
			"",
			[]string{"A", "AB", "AC", "ABD", "ACD"},
			[]string{"ABD", "ACD"},
		},
		{
			Sequence{Literal("A"), Choice{Literal("B"), Choice{Literal("C"), Literal("D")}}, Literal("E")},
			"",
			[]string{"A", "AB", "AC", "AD", "ABE", "ACE", "ADE"},
			[]string{"ABE", "ACE", "ADE"},
		},
		{
			Sequence{Literal("A"), Choice{Literal("B"), Literal("")}, Literal("E")},
			"",
			[]string{"A", "AB", "ABE", "AE"},
			[]string{"ABE", "AE"},
		},
	}

	for _, table := range tables {
		visits := make([]string, 0)
		f := func(s string) {
			visits = append(visits, s)
		}
		finals := table.test.Enumerate(table.prefix, f)
		if !reflect.DeepEqual(visits, table.visits) {
			t.Errorf("%s visits: expected %v, got %v", table.test, table.visits, visits)
		}
		if !reflect.DeepEqual(finals, table.finals) {
			t.Errorf("%s finals: expected %v, got %v", table.test, table.finals, finals)
		}
	}
}

func TestReadLiteral(t *testing.T) {
	tables := []struct{
		input string
		output Expr
		remainingBytes int
	}{
		{
			"",
			Literal(""),
			0,
		},
		{
			"ABCDE",
			Literal("ABCDE"),
			0,
		},
		{
			"ABC(D|E)",
			Literal("ABC"),
			5,
		},
		{
			"ABC|DE",
			Literal("ABC"),
			3,
		},
	}

	for _, table := range tables {
		reader := strings.NewReader(table.input)
		output := ReadLiteral(reader)
		remainingBytes := reader.Len()
		if !reflect.DeepEqual(output, table.output) {
			t.Errorf("%s output: expected %v, got %v", table.input, table.output, output)
		}
		if !reflect.DeepEqual(output, table.output) {
			t.Errorf("%s remainingBytes: expected %v, got %v", table.input, table.remainingBytes, remainingBytes)
		}
	}
}

func TestReadExpression(t *testing.T) {
	tables := []struct{
		input string
		output Expr
	}{
		//{
		//	"",
		//	Literal(""),
		//},
		{
			"ABCDE",
			Literal("ABCDE"),
		},
		{
			"^ABCDE$",
			Literal("ABCDE"),
		},
		{
			"(AB|CD)",
			Choice{Literal("AB"), Literal("CD")},
		},
		{
			"ABC(D|E)",
			Sequence{Literal("ABC"), Choice{Literal("D"), Literal("E")}},
		},
		{
			"A(||)B",
			Sequence{Literal("A"), Choice{Literal(""), Literal(""), Literal("")}, Literal("B")},
		},
		{
			"((A|B)C|D)E",
			Sequence{
				Choice{
					Sequence{
						Choice{
							Literal("A"),
							Literal("B"),
						},
						Literal("C"),
					},
					Literal("D"),
				},
				Literal("E"),
			},
		},
	}

	for _, table := range tables {
		reader := strings.NewReader(table.input)
		output := ReadExpression(reader)
		if !reflect.DeepEqual(output, table.output) {
			t.Errorf("%s output: expected %+v, got %+v", table.input, table.output, output)
		}
	}
}

func TestEnumerateExpression(t *testing.T) {
	tables := []struct{
		input string
		visits []string
		finals []string
	}{
		{
			"A(B|C)D",
			[]string{"A", "AB", "AC", "ABD", "ACD"},
			[]string{"ABD", "ACD"},
		},
		{
			"A(B|(C|D))E",
			[]string{"A", "AB", "AC", "AD", "ABE", "ACE", "ADE"},
			[]string{"ABE", "ACE", "ADE"},
		},
		{
			"A(B|)E",
			[]string{"A", "AB", "ABE", "AE"},
			[]string{"ABE", "AE"},
		},
	}

	for _, table := range tables {
		visits := make([]string, 0)
		visit := func(s string) {
			visits = append(visits, s)
		}
		finals := make([]string, 0)
		final := func(s string) {
			finals = append(finals, s)
		}
		expr := ReadExpression(strings.NewReader(table.input))
		EnumerateExpression(expr, visit, final)
		if !reflect.DeepEqual(visits, table.visits) {
			t.Errorf("%s visits: expected %v, got %v", table.input, table.visits, visits)
		}
		if !reflect.DeepEqual(finals, table.finals) {
			t.Errorf("%s finals: expected %v, got %v", table.input, table.finals, finals)
		}
	}
}

func TestSimplifyPath(t *testing.T) {
	tables := []struct{
		input string
		output string
	}{
		{"ENNWSWWNEWSSSSEENEESWENNNN", "ENNWSWWSSSEENEENNN"},
	}

	for _, table := range tables {
		output := SimplifyPath(table.input)
		if output != table.output {
			t.Errorf("%s: expected %s, got %s", table.input, table.output, output)
		}
	}
}

func TestFurthestRoom(t *testing.T) {
	tables := []struct{
		input string
		result int
	}{
		{"^WNE$", 3},
		{"^ENWWW(NEEE|SSE(EE|N))$", 10},
		{"^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$", 18},
	}

	for _, table := range tables {
		result := FurthestRoom(table.input)
		if result != table.result {
			t.Errorf("expected %d, got %d", table.result, result)
		}
	}
}
