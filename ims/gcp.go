//go:build gcp

package ims

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alexbezu/gobol/pl"
)

func DLI(IMS_FUNC string, pcb pl.Objer, IO_AREA pl.Objer, SSAs ...string) {
	if conn == nil {
		connect2db()
	}
}

func TDLI(parcnt int, IMS_FUNC string, pcb pl.Objer, IO_AREA pl.Objer, SSAs ...string) {
	DLI(IMS_FUNC, pcb, IO_AREA, SSAs...)
}
