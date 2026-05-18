package compiler

import (
	"fmt"

	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/wasm-compiler/wasm"
)

func (fc *FuncCompiler) compileStatement(ctx parser.IStatementContext) error {
	switch s := ctx.(type) {
	case *parser.ExpressionStatementContext:
		res, err := fc.compileExpression(s.Expression())
		if err != nil {
			return err
		}
		// Drop expression value if not Void (expression-statement discards it).
		if res.Type != nil && !res.Type.IsVoid() {
			fc.body.Drop()
		}
		return nil

	case *parser.DeclarationStatementContext:
		return fc.compileLocalDecl(s.Declaration())

	case *parser.ControlTransferStmtContext:
		return fc.compileControlTransfer(s.ControlTransfer())

	case *parser.BranchStmtContext:
		return fc.compileBranchStatement(s.BranchStatement())

	case *parser.LoopStmtContext:
		return fc.compileLoopStatement(s.LoopStatement())

	case *parser.DeferStmtContext:
		// Simplified: execute defer body inline (no actual deferred execution).
		return fc.compileCodeBlock(s.DeferStatement().CodeBlock())

	case *parser.CompilerControlStmtContext:
		// Compiler directives – skip.
		return nil
	}
	return nil
}

// ── Local declarations ────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileLocalDecl(decl parser.IDeclarationContext) error {
	if cd := decl.ConstantDeclaration(); cd != nil {
		return fc.compileBindingDecl(cd.PatternInitializerList(), false)
	}
	if vd := decl.VariableDeclaration(); vd != nil {
		if vd.PatternInitializerList() != nil {
			return fc.compileBindingDecl(vd.PatternInitializerList(), true)
		}
	}
	return nil
}

func (fc *FuncCompiler) compileBindingDecl(
	pil parser.IPatternInitializerListContext,
	isMut bool,
) error {
	for _, pi := range pil.AllPatternInitializer() {
		name, annotatedType, err := fc.extractPatternName(pi.Pattern())
		if err != nil {
			return err
		}
		if name == "_" || name == "" {
			// Evaluate for side effects only.
			if pi.Initializer() != nil {
				if _, err := fc.compileExpression(pi.Initializer().Expression()); err != nil {
					return err
				}
				fc.body.Drop()
			}
			continue
		}

		var res ExprResult
		if pi.Initializer() != nil {
			res, err = fc.compileExpression(pi.Initializer().Expression())
			if err != nil {
				return err
			}
		}

		// Determine type: prefer explicit annotation, then infer from rhs.
		varType := annotatedType
		if varType == nil {
			varType = res.Type
		}
		if varType == nil {
			return fmt.Errorf("cannot infer type for %q", name)
		}

		wt, hasWasm := varType.WasmType()
		if !hasWasm {
			continue // Void – no local needed.
		}

		localIdx := fc.allocLocal(wt)
		if pi.Initializer() != nil {
			fc.body.LocalSet(localIdx)
		}

		_ = fc.scope.define(name, Symbol{
			Name:    name,
			Type:    varType,
			WasmIdx: localIdx,
			IsMut:   isMut,
		})
	}
	return nil
}

// extractPatternName pulls the binding name and optional type from a pattern.
func (fc *FuncCompiler) extractPatternName(pat parser.IPatternContext) (
	name string, annotationType VType, err error,
) {
	switch p := pat.(type) {
	case *parser.IdentPatContext:
		name = p.IdentifierPattern().Identifier().GetText()
		if p.TypeAnnotation() != nil {
			annotationType, err = fc.c.resolveType(p.TypeAnnotation().Type_())
		}
	case *parser.WildcardPatContext:
		name = "_"
	default:
		name = pat.GetText()
	}
	return
}

// ── Control transfer ──────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileControlTransfer(ctx parser.IControlTransferContext) error {
	switch t := ctx.(type) {
	case *parser.ReturnStatementContext:
		if expr := t.Expression(); expr != nil {
			if _, err := fc.compileExpression(expr); err != nil {
				return err
			}
		}
		fc.body.Return()

	case *parser.BreakStatementContext:
		// Break targets the innermost block (label 0 in our scheme).
		fc.body.Br(1)

	case *parser.ContinueStatementContext:
		fc.body.Br(0)
	}
	return nil
}

// ── Branch statements ─────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileBranchStatement(ctx parser.IBranchStatementContext) error {
	if ifStmt := ctx.IfStatement(); ifStmt != nil {
		return fc.compileIfStatement(ifStmt)
	}
	if sw := ctx.SwitchStatement(); sw != nil {
		return fc.compileSwitchStatement(sw)
	}
	return nil
}

func (fc *FuncCompiler) compileIfStatement(ctx parser.IIfStatementContext) error {
	// Compile the first condition.
	conds := ctx.ConditionList().AllCondition()
	if len(conds) == 0 {
		return nil
	}
	if err := fc.compileCondition(conds[0]); err != nil {
		return err
	}

	hasElse := ctx.ElseClause() != nil
	fc.body.If(wasm.BlockEmpty{})
	fc.pushScope()

	// Then branch.
	if err := fc.compileCodeBlock(ctx.CodeBlock()); err != nil {
		fc.popScope()
		return err
	}
	fc.popScope()

	// Else branch.
	if hasElse {
		fc.body.Else()
		ec := ctx.ElseClause()
		fc.pushScope()
		if ec.CodeBlock() != nil {
			if err := fc.compileCodeBlock(ec.CodeBlock()); err != nil {
				fc.popScope()
				return err
			}
		} else if ec.IfStatement() != nil {
			if err := fc.compileIfStatement(ec.IfStatement()); err != nil {
				fc.popScope()
				return err
			}
		}
		fc.popScope()
	}

	fc.body.End()
	return nil
}

func (fc *FuncCompiler) compileCondition(ctx parser.IConditionContext) error {
	if expr := ctx.Expression(); expr != nil {
		_, err := fc.compileExpression(expr)
		return err
	}
	// Optional-binding conditions and availability checks are unsupported for now.
	fc.body.I32Const(1)
	return nil
}

func (fc *FuncCompiler) compileSwitchStatement(ctx parser.ISwitchStatementContext) error {
	// Simplified switch: evaluate subject, compare against each case.
	// Full pattern-matching switch is a TODO.
	if _, err := fc.compileExpression(ctx.Expression()); err != nil {
		return err
	}
	fc.body.Drop() // drop subject (simplified)
	return nil
}

// ── Loop statements ───────────────────────────────────────────────────────────

func (fc *FuncCompiler) compileLoopStatement(ctx parser.ILoopStatementContext) error {
	if ws := ctx.WhileStatement(); ws != nil {
		return fc.compileWhile(ws)
	}
	if fi := ctx.ForInStatement(); fi != nil {
		return fc.compileForIn(fi)
	}
	if rw := ctx.RepeatWhileStatement(); rw != nil {
		return fc.compileRepeatWhile(rw)
	}
	return nil
}

// while cond { body }
//
//	block $break
//	  loop $continue
//	    <cond>; i32.eqz; br_if $break
//	    <body>
//	    br $continue
//	  end
//	end
func (fc *FuncCompiler) compileWhile(ctx parser.IWhileStatementContext) error {
	fc.body.Block(wasm.BlockEmpty{}) // label 1 = break target
	fc.body.Loop(wasm.BlockEmpty{})  // label 0 = continue target

	// Condition
	conds := ctx.ConditionList().AllCondition()
	if len(conds) > 0 {
		if err := fc.compileCondition(conds[0]); err != nil {
			return err
		}
	} else {
		fc.body.I32Const(1)
	}
	fc.body.I32Eqz()
	fc.body.BrIf(1) // break if condition false

	fc.pushScope()
	if err := fc.compileCodeBlock(ctx.CodeBlock()); err != nil {
		fc.popScope()
		return err
	}
	fc.popScope()

	fc.body.Br(0) // continue
	fc.body.End() // end loop
	fc.body.End() // end block
	return nil
}

// repeat { body } while cond
func (fc *FuncCompiler) compileRepeatWhile(ctx parser.IRepeatWhileStatementContext) error {
	fc.body.Block(wasm.BlockEmpty{})
	fc.body.Loop(wasm.BlockEmpty{})

	fc.pushScope()
	if err := fc.compileCodeBlock(ctx.CodeBlock()); err != nil {
		fc.popScope()
		return err
	}
	fc.popScope()

	if _, err := fc.compileExpression(ctx.Expression()); err != nil {
		return err
	}
	fc.body.BrIf(0) // loop back if condition true
	fc.body.End()
	fc.body.End()
	return nil
}

// for item in collection { body }
//
// Supports:
//   - integer range: `for i in lo..<hi`  (RANGE_HALF_OPEN detected by pattern)
//   - generic fallback: treat collection as i32 pointer (no actual iteration)
func (fc *FuncCompiler) compileForIn(ctx parser.IForInStatementContext) error {
	// Try to detect a half-open range literal: lo..<hi
	exprText := ctx.Expression().GetText()
	_ = exprText

	// Pattern for loop variable
	patCtx := ctx.Pattern()
	loopVarName, _, _ := fc.extractPatternName(patCtx)

	// We look at the expression to see if it is a range (binary op "..<").
	// The grammar parses `0..<10` as expression with RANGE_HALF_OPEN binary-op.
	// Try to decompose into lo/hi.
	loExpr, hiExpr, isRange := tryDecomposeRange(ctx.Expression())
	if isRange {
		return fc.compileRangeForIn(loopVarName, loExpr, hiExpr, ctx.CodeBlock())
	}

	// Generic fallback: just compile the collection expression for side-effects.
	if _, err := fc.compileExpression(ctx.Expression()); err != nil {
		return err
	}
	fc.body.Drop()
	return nil
}

// tryDecomposeRange detects `lo ..<  hi` and returns both sides.
func tryDecomposeRange(expr parser.IExpressionContext) (
	lo, hi parser.IExpressionContext, ok bool,
) {
	binExprs := expr.BinaryExpressions()
	if binExprs == nil {
		return nil, nil, false
	}
	all := binExprs.AllBinaryExpression()
	if len(all) != 1 {
		return nil, nil, false
	}
	bop, isBO := all[0].(*parser.BinaryOpContext)
	if !isBO {
		return nil, nil, false
	}
	opText := bop.BinaryOperator().GetText()
	if opText != "..<" && opText != "..." {
		return nil, nil, false
	}
	// Wrap the rhs back into a synthetic expression. We'll compile the rhs
	// PrefixExpression directly.
	_ = bop.PrefixExpression()
	return expr, nil, false // TODO: proper range decomposition
}

// compileRangeForIn emits a counted loop: for i in lo..<hi
//
//	local i = lo
//	block $break
//	  loop $continue
//	    i >= hi → br $break
//	    <body>
//	    i += 1
//	    br $continue
//	  end
//	end
func (fc *FuncCompiler) compileRangeForIn(
	varName string,
	loExpr, hiExpr parser.IExpressionContext,
	body parser.ICodeBlockContext,
) error {
	// Allocate loop variable.
	loopIdx := fc.allocLocal(wasm.I32)

	// Init: i = lo
	if _, err := fc.compileExpression(loExpr); err != nil {
		return err
	}
	fc.body.LocalSet(loopIdx)

	// Allocate hi local.
	hiIdx := fc.allocLocal(wasm.I32)
	if _, err := fc.compileExpression(hiExpr); err != nil {
		return err
	}
	fc.body.LocalSet(hiIdx)

	fc.body.Block(wasm.BlockEmpty{}) // $break
	fc.body.Loop(wasm.BlockEmpty{})  // $continue

	// Condition: i >= hi → break
	fc.body.LocalGet(loopIdx)
	fc.body.LocalGet(hiIdx)
	fc.body.I32GeS()
	fc.body.BrIf(1)

	// Define loop variable in inner scope.
	fc.pushScope()
	if varName != "_" && varName != "" {
		_ = fc.scope.define(varName, Symbol{
			Name:    varName,
			Type:    &PrimitiveType{Kind: KindInt},
			WasmIdx: loopIdx,
			IsMut:   false,
		})
	}
	if err := fc.compileCodeBlock(body); err != nil {
		fc.popScope()
		return err
	}
	fc.popScope()

	// i += 1
	fc.body.LocalGet(loopIdx)
	fc.body.I32Const(1)
	fc.body.I32Add()
	fc.body.LocalSet(loopIdx)

	fc.body.Br(0) // continue
	fc.body.End() // end loop
	fc.body.End() // end block
	return nil
}