package hir

// Flatten turns a Func's structured Body into flat Blocks in vir's Join
// Convention shape: no phi nodes, values merge by same-name reassignment,
// every block ends in exactly one terminator.
//
// It moves instructions; it never rewrites them. Every *Instr in Blocks is
// the same pointer that was in Body, which is what makes the two shapes one
// representation rather than two.
func Flatten(f *Func) {
	fl := &flattener{fn: f}
	fl.open("") // the entry block is unlabeled and unbranchable-to

	// Hoisted allocas land first, before anything can branch past them.
	fl.cur.Lines = append(fl.cur.Lines, f.Allocas...)

	fl.seq(f.Body)
	if fl.cur != nil {
		// A body that fell out the bottom without a terminator. A void
		// function returns; anything else is unreachable, and the analyzer
		// already guaranteed every path returns.
		fl.term(Ret{})
	}
	f.Blocks = fl.blocks
	f.Body = nil
}

type flattener struct {
	fn     *Func
	blocks []*FlatBlock
	cur    *FlatBlock
	n      int

	// loops is the break/continue target stack. There are no loop labels, so
	// the innermost entry is always the right one.
	loops []loopTargets
}

type loopTargets struct{ head, exit string }

func (fl *flattener) label(hint string) string {
	fl.n++
	return hint + itoa(fl.n)
}

func (fl *flattener) open(label string) {
	b := &FlatBlock{Label: label}
	fl.blocks = append(fl.blocks, b)
	fl.cur = b
}

// term closes the current block. A second terminator on one block would be a
// verification error, so everything after it is dead and is dropped.
func (fl *flattener) term(t Term) {
	if fl.cur == nil {
		return
	}
	fl.cur.Term = t
	fl.cur = nil
}

func (fl *flattener) emit(in *Instr) {
	if fl.cur == nil {
		// Unreachable code after a terminator. The analyzer permits it —
		// there is no reachability rule in the language — so it is dropped
		// rather than diagnosed.
		return
	}
	fl.cur.Lines = append(fl.cur.Lines, in)
}

func (fl *flattener) seq(s *Seq) {
	if s == nil {
		return
	}
	for _, st := range s.List {
		fl.stmt(st)
	}
}

func (fl *flattener) stmt(s Stmt) {
	switch x := s.(type) {
	case *Instr:
		fl.emit(x)

	case *If:
		thenL := fl.label("then")
		elseL := fl.label("else")
		joinL := fl.label("join")
		if x.Else == nil {
			elseL = joinL
		}
		fl.term(BrIf{Cond: x.Cond, Then: thenL, Else: elseL})

		fl.open(thenL)
		fl.seq(x.Then)
		fl.term(Br{Label: joinL})

		if x.Else != nil {
			fl.open(elseL)
			fl.seq(x.Else)
			fl.term(Br{Label: joinL})
		}
		fl.open(joinL)

	case *Loop:
		headL := fl.label("loop")
		exitL := fl.label("done")
		fl.term(Br{Label: headL})

		fl.open(headL)
		fl.loops = append(fl.loops, loopTargets{head: headL, exit: exitL})
		fl.seq(x.Body)
		fl.loops = fl.loops[:len(fl.loops)-1]
		fl.term(Br{Label: headL}) // the back edge

		fl.open(exitL)

	case *SwitchStmt:
		defL := fl.label("default")
		joinL := fl.label("join")
		cases := make([]SwitchTermCase, 0, len(x.Cases))
		bodies := make([]*Seq, 0, len(x.Cases))
		labels := make([]string, 0, len(x.Cases))

		// One label per distinct body: two case values sharing a clause
		// share a block, which is what a comma-separated PatternList means.
		seen := map[*Seq]string{}
		for _, c := range x.Cases {
			lbl, ok := seen[c.Body]
			if !ok {
				lbl = fl.label("case")
				seen[c.Body] = lbl
				bodies = append(bodies, c.Body)
				labels = append(labels, lbl)
			}
			cases = append(cases, SwitchTermCase{Value: c.Value, Label: lbl})
		}
		fl.term(SwitchTerm{Value: x.Value, Default: defL, Cases: cases})

		for i, body := range bodies {
			fl.open(labels[i])
			fl.seq(body)
			fl.term(Br{Label: joinL})
		}
		fl.open(defL)
		fl.seq(x.Default)
		fl.term(Br{Label: joinL})
		fl.open(joinL)

	case *Break:
		if len(fl.loops) == 0 {
			return
		}
		fl.term(Br{Label: fl.loops[len(fl.loops)-1].exit})

	case *Continue:
		if len(fl.loops) == 0 {
			return
		}
		fl.term(Br{Label: fl.loops[len(fl.loops)-1].head})

	case *ReturnStmt:
		fl.term(Ret{Value: x.Value})

	case *TrapStmt:
		fl.term(TrapTerm{})

	case *UnreachableStmt:
		fl.term(UnreachTerm{})
	}
}