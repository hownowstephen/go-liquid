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

//   def test_contains_works_on_arrays
//     @context = Liquid::Context.new
//     @context['array'] = [1, 2, 3, 4, 5]
//     array_expr = VariableLookup.new("array")

//     assert_evalutes_false array_expr, 'contains', 0
//     assert_evalutes_true array_expr, 'contains', 1
//     assert_evalutes_true array_expr, 'contains', 2
//     assert_evalutes_true array_expr, 'contains', 3
//     assert_evalutes_true array_expr, 'contains', 4
//     assert_evalutes_true array_expr, 'contains', 5
//     assert_evalutes_false array_expr, 'contains', 6
//     assert_evalutes_false array_expr, 'contains', "1"
//   end

//   def test_contains_returns_false_for_nil_operands
//     @context = Liquid::Context.new
//     assert_evalutes_false VariableLookup.new('not_assigned'), 'contains', '0'
//     assert_evalutes_false 0, 'contains', VariableLookup.new('not_assigned')
//   end

//   def test_contains_return_false_on_wrong_data_type
//     assert_evalutes_false 1, 'contains', 0
//   end

//   def test_contains_with_string_left_operand_coerces_right_operand_to_string
//     assert_evalutes_true ' 1 ', 'contains', 1
//     assert_evalutes_false ' 1 ', 'contains', 2
//   end

//   def test_or_condition
//     condition = Condition.new(1, '==', 2)

//     assert_equal false, condition.evaluate

//     condition.or Condition.new(2, '==', 1)

//     assert_equal false, condition.evaluate

//     condition.or Condition.new(1, '==', 1)

//     assert_equal true, condition.evaluate
//   end

//   def test_and_condition
//     condition = Condition.new(1, '==', 1)

//     assert_equal true, condition.evaluate

//     condition.and Condition.new(2, '==', 2)

//     assert_equal true, condition.evaluate

//     condition.and Condition.new(2, '==', 1)

//     assert_equal false, condition.evaluate
//   end

//   def test_should_allow_custom_proc_operator
//     Condition.operators['starts_with'] = proc { |cond, left, right| left =~ %r{^#{right}} }

//     assert_evalutes_true 'bob', 'starts_with', 'b'
//     assert_evalutes_false 'bob', 'starts_with', 'o'
//   ensure
//     Condition.operators.delete 'starts_with'
//   end

//   def test_left_or_right_may_contain_operators
//     @context = Liquid::Context.new
//     @context['one'] = @context['another'] = "gnomeslab-and-or-liquid"

//     assert_evalutes_true VariableLookup.new("one"), '==', VariableLookup.new("another")
//   end

//   private

//   def assert_evalutes_true(left, op, right)
//     assert Condition.new(left, op, right).evaluate(@context || Liquid::Context.new),
//       "Evaluated false: #{left} #{op} #{right}"
//   end

//   def assert_evalutes_false(left, op, right)
//     assert !Condition.new(left, op, right).evaluate(@context || Liquid::Context.new),
//       "Evaluated true: #{left} #{op} #{right}"
//   end

//   def assert_evaluates_argument_error(left, op, right)
//     assert_raises(Liquid::ArgumentError) do
//       Condition.new(left, op, right).evaluate(@context || Liquid::Context.new)
//     end
//   end
// end # ConditionTest
