package verifier

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

var ruleOperators = map[string]string{
	"и":   " && ",
	"или": " || ",
	"не":  "!",
}  

func Verify(doc *string, rule string) error {
	rule = evalRule(rule)

	verified, err := verifyLine(*doc, rule)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if !verified {
		fmt.Println("НЕ СООТВЕТСТВИЕ")
	} else {
		fmt.Println("СООТВЕТСТВИЕ")
	}

	return nil
}

func evalRule(rule string) string {
	words := strings.Split(rule, " ")

	for i := 0; i < len(words); i++ {
		if len(ruleOperators[words[i]]) > 0 {
			words[i] = ruleOperators[words[i]]
		} else {
			openBrackets := strings.Repeat("(", strings.Count(words[i], "("))
			closeBrackets := strings.Repeat(")", strings.Count(words[i], ")"))

			operand := openBrackets + fmt.Sprintf("contains(doc, \"%s\")", strings.Trim(words[i], "()")) + closeBrackets
			words[i] = operand
		}
	}
	return strings.Join(words, "")

}

func verifyLine(doc, rule string) (bool, error) {
	parameters := map[string]govaluate.ExpressionFunction{
		"contains": stringsContains,
	}

	expr, err := govaluate.NewEvaluableExpressionWithFunctions(rule, parameters)
	if err != nil {
		return false, err
	}

	result, err := expr.Evaluate(map[string]interface{}{
		"doc": doc,
	})
	if err != nil {
		return false, err
	}

	verified, ok := result.(bool)
	if !ok {
		return false, errors.New("it is not possible to convert the result to the bool type")
	}

	return verified, nil
}

func stringsContains(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("two arguments are expected")
	}

	line, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("a string is expected as the first argument")
	}

	substr, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("a string is expected as the second argument")
	}

	return strings.Contains(line, substr), nil
}