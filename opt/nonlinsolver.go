// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"strings"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/la"
)

// NonLinSolver solves (unconstrained) nonlinear optimization problems
type NonLinSolver interface {
	Min(x la.Vector, params dbf.Params) (fmin float64) // computes minimum and updates x @ min
}

// nlsMaker defines a function that makes non-linear-solvers
type nlsMaker func(prob *Problem) NonLinSolver

// nlsMakersDB implements a database of non-linear-solver makers
var nlsMakersDB = make(map[string]nlsMaker)

// GetNonLinSolver finds a non-linear-solver in database or panic
//  kind -- e.g. conjgrad, powel, graddesc
func GetNonLinSolver(kind string, prob *Problem) NonLinSolver {
	strKind := strings.ToLower(kind)
	if maker, ok := nlsMakersDB[strKind]; ok {
		return maker(prob)
	}
	chk.Panic("cannot find NonLinSolver named %q in database\n", kind)
	return nil
}
