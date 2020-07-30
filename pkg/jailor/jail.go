/*
 * Copyright (c) 2015-2020 Christian Blichmann
 * Copyright (c) 2020 Deep Dhillon
 *
 * Create, manage and run chroot/jail environments
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package jailor

import (
	"log"
	"path/filepath"

	"github.com/dhillondeep/jailor/internal/action"
	"github.com/dhillondeep/jailor/pkg/copy"
	"github.com/dhillondeep/jailor/pkg/loader"
	"github.com/dhillondeep/jailor/pkg/spec"
)

type JailContext struct {
	UseCow bool
	Force  bool
	Clean  bool
}

func CreateJail(jailCtx JailContext, specStmts spec.Statements, jailDir string) error {
	reflinkOpt := copy.ReflinkNo
	if jailCtx.UseCow {
		reflinkOpt = copy.ReflinkAlways
	}

	for _, s := range spec.ExpandLexical(expandWithDependencies(specStmts)) {
		target := filepath.Join(jailDir, s.Target())

		switch stmt := s.(type) {
		case spec.Directory:
			if err := action.Directory(target, stmt); err != nil {
				return err
			}
		case spec.RegularFile:
			if err := action.RegularFile(target, stmt, &copy.Options{
				Force:             jailCtx.Force,
				Reflink:           reflinkOpt,
				RemoveDestination: jailCtx.Clean,
			}); err != nil {
				return err
			}
		case spec.Link:
			if err := action.Link(target, stmt); err != nil {
				return err
			}
		case spec.Device:
			if err := action.Device(target, stmt); err != nil {
				return err
			}
		case spec.Run:
			if err := action.Run(target, stmt, jailDir); err != nil {
				return err
			}
		}
	}
	return nil
}

func expandWithDependencies(stmts spec.Statements) spec.Statements {
	expanded := spec.ExpandLexical(stmts)
	for _, s := range expanded {
		switch stmt := s.(type) {
		case spec.RegularFile:
			deps, err := loader.ImportedLibraries(stmt.Source())
			if err != nil {
				log.Fatalf("%s\n", err)
			}
			attr := stmt.FileAttr()
			for _, d := range deps {
				f := spec.NewRegularFile(d, d)
				*f.FileAttr() = *attr
				expanded = append(expanded, f)
			}
		}
	}
	return expanded
}
