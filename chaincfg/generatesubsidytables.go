// Copyright (c) 2019-2020 The Hdfchain developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//+build ignore

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/hdfchain/hdfd/chaincfg/v3"
)

type payout struct {
	offset int
	amount int64
}

type payouts []payout

func (p payouts) GoString() string {
	var s strings.Builder
	s.WriteString("[]blockOnePayout{\n")
	for i := range p {
		s.WriteString("\t")
		s.WriteString(fmt.Sprintf("{offset: %#v, amount: %#v}", p[i].offset, p[i].amount))
		s.WriteString(",\n")
	}
	s.WriteString("}")
	return s.String()
}

type printer func(format string, args ...interface{})

func main() {
	buf := new(bytes.Buffer)
	p := func(format string, args ...interface{}) {
		fmt.Fprintf(buf, format, args...)
	}

	p("// autogenerated by generatesubsidytable.go; do not edit\n\n")
	p("package chaincfg\n")
	for _, g := range []struct {
		name         string
		tokenPayouts []chaincfg.TokenPayout
	}{
		{"MainNetParams", chaincfg.MainNetSubsidyDefinition},
		{"TestNet3Params", chaincfg.TestNet3SubsidyDefinition},
		{"SimNetParams", chaincfg.SimNetSubsidyDefinition},
		{"RegNetParams", chaincfg.RegNetSubsidyDefinition},
	} {
		p("\n")
		defs(p, g.name, g.tokenPayouts)
	}

	out, err := os.Create("subsidytables.go")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(out, buf)
	if err != nil {
		log.Fatal(err)
	}
}

func defs(p printer, paramsName string, tokenPayouts []chaincfg.TokenPayout) {
	var scripts strings.Builder
	var payouts payouts
	var offset int
	for _, t := range tokenPayouts {
		if t.ScriptVersion != 0 {
			log.Fatalf("this tool only generates code where all scripts are version 0")
		}
		_, err := scripts.WriteString(hex.EncodeToString(t.Script))
		if err != nil {
			log.Panic(err)
		}
		offset += len(t.Script)
		payouts = append(payouts, payout{
			offset: offset,
			amount: t.Amount,
		})
	}

	p("const blockOnePayoutScripts_%s = %#v\n\n", paramsName, scripts.String())
	p("var blockOnePayouts_%s = %#v\n\n", paramsName, payouts)
	p("func tokenPayouts_%s() []TokenPayout {\n"+
		"\treturn tokenPayouts(blockOnePayoutScripts_%[1]s, blockOnePayouts_%[1]s)\n"+
		"}\n", paramsName)
}
