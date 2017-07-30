package newfile

import (
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newdesignpattern/observer"
)

type Transfer struct {
	newdesignpattern.Observable
	src newpattern.SourcePattern

}


func NewTransfer(pattern newpattern.SourcePattern) *Transfer {
	return &Transfer{
		src: pattern,
	}

}


func (t *Transfer) find() {

}


func (t *Transfer) copyTo(dst string) {

}

func (t *Transfer) moveTo(dst string) {

}

func (t *Transfer) remove(dst string) {

}