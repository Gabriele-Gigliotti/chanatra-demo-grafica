package document

import (
	"drawino/rmm"
)

func init() {

}

func NewDocument() {
	rmm.SetRawMode()
	rmm.ResetTerm()
}
