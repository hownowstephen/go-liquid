package liquid

import (
	"fmt"
	"reflect"
	"strings"
)

type ErrInvalidOperator string

func (e ErrInvalidOperator) Error() string {
	return fmt.Sprintf("Liquid::InvalidOperator: %v", string(e))
}

type ErrBadArgument struct {
	args []Expression
}

func (e ErrBadArgument) Error() string {
	return fmt.Sprintf("Liquid::ArgumentError: %v", e.args)
}

// Condition is a callable func that wraps an operator
type Condition struct {
	a        Expression
	operator operator
	b        Expression
}

// Evaluate the supplied condition
func (c Condition) Evaluate(ctx Context) (bool, error) {
	return c.operator(c.a.Evaluate(ctx), c.b.Evaluate(ctx))
}

type operator func(a, b Expression) (bool, error)

func equal(a, b Expression) (bool, error) {
	return reflect.DeepEqual(a, b), nil
}

func notEqual(a, b Expression) (bool, error) {
	r, err := equal(a, b)
	return !r, err
}

func lt(a, b Expression) (bool, error) {
	if t := reflect.TypeOf(a); t == reflect.TypeOf(b) {
		switch a.(type) {
		case integerExpr:
			return a.(integerExpr) < b.(integerExpr), nil
		case floatExpr:
			return a.(floatExpr) < b.(floatExpr), nil
		case stringExpr:
			return a.(stringExpr) < b.(stringExpr), nil
		}
	}
	return false, ErrBadArgument{}
}

func gt(a, b Expression) (bool, error) {
	if t := reflect.TypeOf(a); t == reflect.TypeOf(b) {
		switch a.(type) {
		case integerExpr:
			return a.(integerExpr) > b.(integerExpr), nil
		case floatExpr:
			return a.(floatExpr) > b.(floatExpr), nil
		case stringExpr:
			return a.(stringExpr) > b.(stringExpr), nil
		}
	}
	return false, ErrBadArgument{}
}

func contains(a, b Expression) (bool, error) {

	switch a.(type) {
	case arrayExpr:
		for _, value := range a.(arrayExpr) {
			if reflect.DeepEqual(interfaceToExpression(value), b) {
				return true, nil
			}
		}
		return false, nil
	case stringExpr:
		return strings.Contains(string(a.(stringExpr)), fmt.Sprintf("%v", b)), nil
	}

	return false, fmt.Errorf("unimplemented")
}

var operators = map[string]operator{
	"==": equal,
	"!=": notEqual,
	"<>": notEqual,
	"<":  lt,
	">":  gt,
	">=": func(a, b Expression) (bool, error) {
		if g, err := gt(a, b); err != nil || g {
			return g, err
		}
		return equal(a, b)
	},
	"<=": func(a, b Expression) (bool, error) {
		if l, err := lt(a, b); err != nil || l {
			return l, err
		}
		return equal(a, b)
	},
	"contains": contains,
}

func NewCondition(op1 Expression, operator string, op2 Expression) (*Condition, error) {

	if found, ok := operators[operator]; ok {
		return &Condition{op1, found, op2}, nil
	}

	//       # If the operator is empty this means that the decision statement is just
	//       # a single variable. We can just poll this variable from the context and
	//       # return this as the result.
	//       return context.evaluate(left) if op.nil?

	//       left = context.evaluate(left)
	//       right = context.evaluate(right)

	//       operation = self.class.operators[op] || raise(Liquid::ArgumentError.new("Unknown operator #{op}"))

	//       if operation.respond_to?(:call)
	//         operation.call(self, left, right)
	//       elsif left.respond_to?(operation) && right.respond_to?(operation)
	//         begin
	//           left.send(operation, right)
	//         rescue ::ArgumentError => e
	//           raise Liquid::ArgumentError.new(e.message)
	//         end
	//       end

	//  raise(Liquid::ArgumentError.new("Unknown operator #{op}"))
	return nil, ErrInvalidOperator(operator)
}

// module Liquid
//   # Container for liquid nodes which conveniently wraps decision making logic
//   #
//   # Example:
//   #
//   #   c = Condition.new(1, '==', 1)
//   #   c.evaluate #=> true
//   #
//   class Condition #:nodoc:
//     @@operators = {
//       '=='.freeze => ->(cond, left, right) {  cond.send(:equal_variables, left, right) },
//       '!='.freeze => ->(cond, left, right) { !cond.send(:equal_variables, left, right) },
//       '<>'.freeze => ->(cond, left, right) { !cond.send(:equal_variables, left, right) },
//       '<'.freeze  => :<,
//       '>'.freeze  => :>,
//       '>='.freeze => :>=,
//       '<='.freeze => :<=,
//       'contains'.freeze => lambda do |cond, left, right|
//         if left && right && left.respond_to?(:include?)
//           right = right.to_s if left.is_a?(String)
//           left.include?(right)
//         else
//           false
//         end
//       end
//     }

//     def self.operators
//       @@operators
//     end

//     attr_reader :attachment
//     attr_accessor :left, :operator, :right

//     def initialize(left = nil, operator = nil, right = nil)
//       @left = left
//       @operator = operator
//       @right = right
//       @child_relation  = nil
//       @child_condition = nil
//     end

//     def evaluate(context = Context.new)
//       result = interpret_condition(left, right, operator, context)

//       case @child_relation
//       when :or
//         result || @child_condition.evaluate(context)
//       when :and
//         result && @child_condition.evaluate(context)
//       else
//         result
//       end
//     end

//     def or(condition)
//       @child_relation = :or
//       @child_condition = condition
//     end

//     def and(condition)
//       @child_relation = :and
//       @child_condition = condition
//     end

//     def attach(attachment)
//       @attachment = attachment
//     end

//     def else?
//       false
//     end

//     def inspect
//       "#<Condition #{[@left, @operator, @right].compact.join(' '.freeze)}>"
//     end

//     private

//     def equal_variables(left, right)
//       if left.is_a?(Liquid::Expression::MethodLiteral)
//         if right.respond_to?(left.method_name)
//           return right.send(left.method_name)
//         else
//           return nil
//         end
//       end

//       if right.is_a?(Liquid::Expression::MethodLiteral)
//         if left.respond_to?(right.method_name)
//           return left.send(right.method_name)
//         else
//           return nil
//         end
//       end

//       left == right
//     end

//     def interpret_condition(left, right, op, context)
//       # If the operator is empty this means that the decision statement is just
//       # a single variable. We can just poll this variable from the context and
//       # return this as the result.
//       return context.evaluate(left) if op.nil?

//       left = context.evaluate(left)
//       right = context.evaluate(right)

//       operation = self.class.operators[op] || raise(Liquid::ArgumentError.new("Unknown operator #{op}"))

//       if operation.respond_to?(:call)
//         operation.call(self, left, right)
//       elsif left.respond_to?(operation) && right.respond_to?(operation)
//         begin
//           left.send(operation, right)
//         rescue ::ArgumentError => e
//           raise Liquid::ArgumentError.new(e.message)
//         end
//       end
//     end
//   end

//   class ElseCondition < Condition
//     def else?
//       true
//     end

//     def evaluate(_context)
//       true
//     end
//   end
// end
