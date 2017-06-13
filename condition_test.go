package liquid

import (
	"reflect"
	"testing"
)

func checkConditionContext(t *testing.T, op1 Expression, operator string, op2 Expression, want bool, ctx Context) error {
	cond, err := NewCondition(op1, operator, op2)
	if err != nil {
		return err
	}

	got, err := cond.Evaluate(ctx)
	if err != nil {
		return err
	}

	if got != want {
		t.Errorf(`Condition "%v %v %v" evaluated wrongly, want: %v, got:%v`, op1, operator, op2, want, got)
	}

	return nil
}

func checkCondition(t *testing.T, op1 Expression, operator string, op2 Expression, want bool) error {
	return checkConditionContext(t, op1, operator, op2, want, Context{})
}

func checkIntCondition(t *testing.T, op1 int, operator string, op2 int, want bool) {
	err := checkCondition(t, integerExpr(op1), operator, integerExpr(op2), want)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func checkFloatCondition(t *testing.T, op1 float64, operator string, op2 float64, want bool) {
	err := checkCondition(t, floatExpr(op1), operator, floatExpr(op2), want)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func checkStringCondition(t *testing.T, op1, operator, op2 string, want bool) {
	err := checkCondition(t, stringExpr(op1), operator, stringExpr(op2), want)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestBasicCondition(t *testing.T) {
	checkIntCondition(t, 1, "==", 2, false)
	checkIntCondition(t, 1, "==", 1, true)
}

func TestDefaultOperatorsEvaluateTrue(t *testing.T) {
	checkIntCondition(t, 1, "==", 1, true)
	checkIntCondition(t, 1, "!=", 2, true)
	checkIntCondition(t, 1, "<>", 2, true)
	checkIntCondition(t, 1, "<", 2, true)
	checkIntCondition(t, 2, ">", 1, true)
	checkIntCondition(t, 1, ">=", 1, true)
	checkIntCondition(t, 2, ">=", 1, true)
	checkIntCondition(t, 1, "<=", 2, true)
	checkIntCondition(t, 1, "<=", 1, true)
	// negative numbers
	checkIntCondition(t, 1, ">", -1, true)
	checkIntCondition(t, -1, "<", 1, true)
	checkFloatCondition(t, 1.0, ">", -1.0, true)
	checkFloatCondition(t, -1.0, "<", 1.0, true)
}

func TestDefaultOperatorsEvaluateFalse(t *testing.T) {
	checkIntCondition(t, 1, "==", 2, false)
	checkIntCondition(t, 1, "!=", 1, false)
	checkIntCondition(t, 1, "<>", 1, false)
	checkIntCondition(t, 1, "<", 0, false)
	checkIntCondition(t, 2, ">", 4, false)
	checkIntCondition(t, 1, ">=", 3, false)
	checkIntCondition(t, 2, ">=", 4, false)
	checkIntCondition(t, 1, "<=", 0, false)
}

func TestContainsWorksOnStrings(t *testing.T) {
	checkStringCondition(t, "bob", "contains", "o", true)
	checkStringCondition(t, "bob", "contains", "b", true)
	checkStringCondition(t, "bob", "contains", "bo", true)
	checkStringCondition(t, "bob", "contains", "ob", true)
	checkStringCondition(t, "bob", "contains", "bob", true)

	checkStringCondition(t, "bob", "contains", "bob2", false)
	checkStringCondition(t, "bob", "contains", "a", false)
	checkStringCondition(t, "bob", "contains", "---", false)
}

func TestInvalidComparisonOperator(t *testing.T) {
	err := checkCondition(t, integerExpr(1), "~~", integerExpr(0), false)
	if !reflect.DeepEqual(err, ErrInvalidOperator("~~")) {
		t.Errorf("Bad error for operator, want: %v got: %v", ErrInvalidOperator("~~"), err)
	}
}

func TestComparisonOfIntAndString(t *testing.T) {
	err := checkCondition(t, stringExpr("1"), ">", integerExpr(0), false)
	if !reflect.DeepEqual(err, ErrBadArgument{}) {
		t.Errorf("wrong error type, want: ErrBadArgment, got: %v", err)
	}

	err = checkCondition(t, stringExpr("1"), "<", integerExpr(0), false)
	if !reflect.DeepEqual(err, ErrBadArgument{}) {
		t.Errorf("wrong error type, want: ErrBadArgment, got: %v", err)
	}

	err = checkCondition(t, stringExpr("1"), ">=", integerExpr(0), false)
	if !reflect.DeepEqual(err, ErrBadArgument{}) {
		t.Errorf("wrong error type, want: ErrBadArgment, got: %v", err)
	}

	err = checkCondition(t, stringExpr("1"), "<=", integerExpr(0), false)
	if !reflect.DeepEqual(err, ErrBadArgument{}) {
		t.Errorf("wrong error type, want: ErrBadArgment, got: %v", err)
	}
}

func TestContainsWorksOnArrays(t *testing.T) {
	ctx := Context{
		vars: map[string]interface{}{
			"array": []interface{}{1, 2, 3, 4, 5},
		},
	}

	arrayExpr := ParseVariableLookup("array")

	tests := []struct {
		expr Expression
		want bool
	}{
		{integerExpr(0), false},
		{integerExpr(1), true},
		{integerExpr(2), true},
		{integerExpr(3), true},
		{integerExpr(4), true},
		{integerExpr(5), true},
		{integerExpr(6), false},
		{stringExpr("1"), false},
	}

	for i, test := range tests {
		if err := checkConditionContext(t, arrayExpr, "contains", test.expr, test.want, ctx); err != nil {
			t.Errorf("check %v failed with %v", i, err)
		}
	}
}

func TestContainsReturnsFalseForNilOperands(t *testing.T) {
	ctx := Context{}

	if err := checkConditionContext(t, ParseVariableLookup("not_assigned"), "contains", stringExpr("0"), false, ctx); err != nil {
		t.Errorf("expected false, got: %v", err)
	}

	if err := checkConditionContext(t, integerExpr(0), "contains", ParseVariableLookup("not_assigned"), false, ctx); err != nil {
		t.Errorf("expected false, got: %v", err)
	}
}

func TestContainsReturnsFalseOnWrongDataType(t *testing.T) {
	if err := checkCondition(t, integerExpr(1), "contains", integerExpr(0), false); err != nil {
		t.Errorf("expected false, got: %v", err)
	}
}

func TestContainsWithStringLeftOperandCoercesRightOperandToString(t *testing.T) {
	if err := checkCondition(t, stringExpr(` 1 `), "contains", integerExpr(1), true); err != nil {
		t.Errorf("expected true, got: %v", err)
	}
	if err := checkCondition(t, stringExpr(` 1 `), "contains", integerExpr(2), false); err != nil {
		t.Errorf("expected false, got: %v", err)
	}
}

func TestOrCondition(t *testing.T) {
	cond, err := NewCondition(integerExpr(1), "==", integerExpr(2))
	if err != nil {
		t.Error(err)
		return
	}

	got, err := cond.Evaluate(Context{})
	if err != nil {
		t.Error(err)
		return
	}

	if got != false {
		t.Errorf("condition evaluated wrong, want: false, got: %v", got)
	}

	if err := cond.Or(integerExpr(2), "==", integerExpr(1)); err != nil {
		t.Errorf("error adding OR condition: %v", err)
	}

	got, err = cond.Evaluate(Context{})
	if err != nil {
		t.Error(err)
		return
	}

	if got != false {
		t.Errorf("condition evaluated wrong, want: false, got: %v", got)
	}

	if err := cond.Or(integerExpr(1), "==", integerExpr(1)); err != nil {
		t.Errorf("error adding OR condition: %v", err)
	}

	got, err = cond.Evaluate(Context{})
	if err != nil {
		t.Error(err)
		return
	}

	if got != true {
		t.Errorf("condition evaluated wrong, want: true, got: %v", got)
	}

}

func TestAndCondition(t *testing.T) {
	cond, err := NewCondition(integerExpr(1), "==", integerExpr(1))
	if err != nil {
		t.Error(err)
		return
	}

	got, err := cond.Evaluate(Context{})
	if err != nil {
		t.Error(err)
		return
	}

	if got != true {
		t.Errorf("condition evaluated wrong, want: true, got: %v", got)
	}

	if err := cond.And(integerExpr(2), "==", integerExpr(2)); err != nil {
		t.Errorf("error adding OR condition: %v", err)
	}

	got, err = cond.Evaluate(Context{})
	if err != nil {
		t.Error(err)
		return
	}

	if got != true {
		t.Errorf("condition evaluated wrong, want: true, got: %v", got)
	}

	if err := cond.And(integerExpr(1), "==", integerExpr(2)); err != nil {
		t.Errorf("error adding OR condition: %v", err)
	}

	got, err = cond.Evaluate(Context{})
	if err != nil {
		t.Error(err)
		return
	}

	if got != false {
		t.Errorf("condition evaluated wrong, want: false, got: %v", got)
	}

}

func TestShouldAllowCustomProcOperator(t *testing.T) {
	t.Skip("unimplemented")
	//   def test_should_allow_custom_proc_operator
	//     Condition.operators['starts_with'] = proc { |cond, left, right| left =~ %r{^#{right}} }

	//     assert_evalutes_true 'bob', 'starts_with', 'b'
	//     assert_evalutes_false 'bob', 'starts_with', 'o'
	//   ensure
	//     Condition.operators.delete 'starts_with'
	//   end
}

func TestLeftOrRightMayContainOperators(t *testing.T) {
	t.Skip("unimplemented")
	//   def test_left_or_right_may_contain_operators
	//     @context = Liquid::Context.new
	//     @context['one'] = @context['another'] = "gnomeslab-and-or-liquid"

	//     assert_evalutes_true VariableLookup.new("one"), '==', VariableLookup.new("another")
	//   end
}
