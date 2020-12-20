package phelp

import (
	"io"
	"os"
	"os/exec"

	"github.com/nickwells/param.mod/v5/param"
	"golang.org/x/term"
)

const pagerCmd = "less"

type pager struct {
	pagerIn io.WriteCloser
	cmd     *exec.Cmd
}

// pagerStart returns a pager which should have done() called on it after any
// output is complete.
func pagerStart(ps *param.PSet) *pager {
	pagerPath, err := exec.LookPath(pagerCmd)
	if err != nil {
		return nil
	}
	outIsTty := isWriterATerminal(ps.StdWriter())
	errIsTty := isWriterATerminal(ps.ErrWriter())

	if !outIsTty && !errIsTty {
		return nil
	}
	cmd := exec.Command(pagerPath)

	pagerIn, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}

	cmd.Stdout = ps.StdWriter()
	cmd.Stderr = ps.ErrWriter()
	err = cmd.Start()
	if err != nil {
		return nil
	}

	if outIsTty {
		_ = param.SetStdWriter(pagerIn)(ps)
	}
	if errIsTty {
		_ = param.SetErrWriter(pagerIn)(ps)
	}

	return &pager{
		pagerIn: pagerIn,
		cmd:     cmd,
	}
}

// isWriterATerminal returns true if the io.Writer is a Terminal
func isWriterATerminal(w io.Writer) bool {
	if outFile, ok := w.(*os.File); ok {
		if term.IsTerminal(int(outFile.Fd())) {
			return true
		}
	}
	return false
}

// done will wait for the pager to complete. Note that it is safe to call
// with a nil pointer
func (p *pager) done() {
	if p == nil {
		return
	}

	p.pagerIn.Close()
	_ = p.cmd.Wait()
}
