package fun

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/lovego/gospec/problems"
)

func checkFuncType(thing string, typ *ast.FuncType, fileSet *token.FileSet) {
	rule, ruleName := getSizeRule(typ.Func, fileSet)

	if typ.Params != nil && uint(typ.Params.NumFields()) > rule.MaxParams {
		problems.Add(
			fileSet.Position(typ.Params.Pos()), fmt.Sprintf(
				`%s params size: %d, limit: %d`, thing, typ.Params.NumFields(), rule.MaxParams,
			), ruleName+`.maxParams`,
		)
	}

	if typ.Results != nil && uint(typ.Results.NumFields()) > rule.MaxResults {
		problems.Add(
			fileSet.Position(typ.Results.Pos()), fmt.Sprintf(
				`%s results size: %d, limit: %d`, thing, typ.Results.NumFields(), rule.MaxResults,
			), ruleName+`.maxResults`,
		)
	}
}

func checkFuncBody(thing string, body *ast.BlockStmt, fileSet *token.FileSet) {
	rule, ruleName := getSizeRule(body.Pos(), fileSet)

	if size := stmtsCount(body); size > rule.MaxStatements {
		problems.Add(
			fileSet.Position(body.Pos()), fmt.Sprintf(
				`%s body size: %d statements, limit: %d`, thing, size, rule.MaxStatements,
			), ruleName+`.maxStatements`,
		)
	}
}

func getSizeRule(pos token.Pos, fileSet *token.FileSet) (sizeRule, string) {
	if strings.HasSuffix(fileSet.Position(pos).Filename, "_test.go") {
		return RuleInTest.Size, "funcInTest.size"
	} else {
		return Rule.Size, "func.size"
	}
}

func stmtsCount(node ast.Node) uint {
	w := &stmtsWalker{}
	ast.Walk(w, node)
	return w.count
}

type stmtsWalker struct {
	count uint
}

func (w *stmtsWalker) Visit(node ast.Node) ast.Visitor {
	if _, ok := node.(ast.Stmt); ok {
		if _, ok := node.(*ast.BlockStmt); !ok {
			w.count++
		}
	}
	return w
}
