package document

import (
	"drawino/lib/rmm"
)

func init() {

}

func NewDocument() {
	rmm.SetRawMode()
	rmm.ResetTerm()
}
