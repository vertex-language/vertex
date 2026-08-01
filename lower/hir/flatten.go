package hir

// Flatten converts a Func's structured Body into flat Blocks in vir's Join
// Convention shape: every block ends in exactly one terminator, values
// merge across blocks by same-name reassignment, and there are no phi
// nodes. This is what makes lower/vir legitimately mechanical afterward.
//
// Note what Flatten does *not* decide: vir's switch terminator takes a
// uniform operand/label list regardless of density, so jump-table-versus-
// compare-chain is cpu/lower/<arch>'s decision, not this package's.
func Flatten(f *Func) {
	if f.Body == nil {
		return
	}
	fl := &flattener{fn: f}
	fl.open("") // the entry block is implicit, unlabeled, unbranchable-to
	fl.seq(f.Body)
	fl.terminate(TermReturn{})
	f.Blocks = fl.blocks
	f.Body = nil
}

type flattener struct {
	fn     *Func
	blocks []*Block
	cur    *Block
	n      int
	loops  []loopCtx
}

type loopCtx struct{ head, exit string }

func (f *flattener) label(base string) string {
	f.n++
	return base + "_" + itoa(f.n)
}

func (f *flattener) open(label string) {
	b := &Block{Label: label}
	f.blocks = append(f.blocks, b)
	f.cur = b
}

func (f *flattener) terminate(t Terminator) {
	if f.cur == nil {
		return
	}
	f.cur.Term = t
	f.cur = nil
}

// ensure opens a fresh unreachable block when code follows a terminator.
// The block is dead by construction, but vir requires every instruction to
// live in a terminated block, so it gets one.
func (f *flattener) ensure() {
	if f.cur == nil {
		f.open(f.label("dead"))
	}
}

func (f *flattener) branch(label string) {
	if f.cur == nil {
		return
	}
	f.terminate(TermBranch{Label: label})
}

func (f *flattener) seq(s *Seq) {
	if s == nil {
		return
	}
	for _, st := range s.List {
		f.stmt(st)
	}
}

func (f *flattener) stmt(s Stmt) {
	switch s := s.(type) {
	case *Seq:
		f.seq(s)

	case *Instrs:
		f.ensure()
		f.cur.Instr = append(f.cur.Instr, s.List...)

	case *If:
		f.ensure()
		then, els, join := f.label("then"), f.label("else"), f.label("join")
		hasElse := s.Else != nil
		target := join
		if hasElse {
			target = els
		}
		f.terminate(TermBranchIf{Cond: s.Cond, Then: then, Else: target})

		f.open(then)
		f.seq(s.Then)
		f.branch(join)

		if hasElse {
			f.open(els)
			f.seq(s.Else)
			f.branch(join)
		}
		f.open(join)

	case *Loop:
		f.ensure()
		head, exit := f.label("loop"), f.label("endloop")
		f.branch(head)
		f.open(head)
		f.loops = append(f.loops, loopCtx{head: head, exit: exit})
		f.seq(s.Body)
		f.loops = f.loops[:len(f.loops)-1]
		f.branch(head) // fall-through re-enters the head, re-testing the condition
		f.open(exit)

	case *Switch:
		f.ensure()
		join := f.label("endswitch")
		def := join
		var cases []TermCase
		type pending struct {
			label string
			body  *Seq
		}
		var bodies []pending
		for _, c := range s.Cases {
			lbl := f.label("case")
			for _, v := range c.Values {
				cases = append(cases, TermCase{Value: v, Label: lbl})
			}
			bodies = append(bodies, pending{lbl, c.Body})
		}
		if s.Default != nil {
			def = f.label("default")
			bodies = append(bodies, pending{def, s.Default})
		}
		f.terminate(TermSwitch{Value: s.Tag, Default: def, Cases: cases})
		for _, p := range bodies {
			f.open(p.label)
			f.seq(p.body)
			f.branch(join)
		}
		f.open(join)

	case *Break:
		if n := len(f.loops); n > 0 {
			f.ensure()
			f.terminate(TermBranch{Label: f.loops[n-1].exit})
		}

	case *Continue:
		if n := len(f.loops); n > 0 {
			f.ensure()
			f.terminate(TermBranch{Label: f.loops[n-1].head})
		}

	case *Return:
		f.ensure()
		f.terminate(TermReturn{Value: s.Value})

	case *Trap:
		f.ensure()
		f.terminate(TermTrap{})

	case *Unreachable:
		f.ensure()
		f.terminate(TermUnreachable{})
	}
}