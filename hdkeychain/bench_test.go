// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hdkeychain

import (
	"testing"
)

// bip0032MasterPriv1 is the master private extended key from the first set of
// test vectors in BIP0032.
const bip0032MasterPriv1 = "dprv3hCznBesA6jBtmoyVFPfyMSZ1qYZ3WdjdebquvkEfmRfx" +
	"C9VFEFi2YDaJqHnx7uGe75eGSa3Mn3oHK11hBW7KZUrPxwbCPBmuCi1nwm182s"

// BenchmarkDeriveHardened benchmarks how long it takes to derive a hardened
// child from a master private extended key.
func BenchmarkDeriveHardened(b *testing.B) {
	b.StopTimer()
	masterKey, err := NewKeyFromString(bip0032MasterPriv1, mockMainNetParams())
	if err != nil {
		b.Errorf("Failed to decode master seed: %v", err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		masterKey.Child(HardenedKeyStart)
	}
}

// BenchmarkDeriveNormal benchmarks how long it takes to derive a normal
// (non-hardened) child from a master private extended key.
func BenchmarkDeriveNormal(b *testing.B) {
	b.StopTimer()
	masterKey, err := NewKeyFromString(bip0032MasterPriv1, mockMainNetParams())
	if err != nil {
		b.Errorf("Failed to decode master seed: %v", err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		masterKey.Child(0)
	}
}

// BenchmarkPrivToPub benchmarks how long it takes to convert a private extended
// key to a public extended key.
func BenchmarkPrivToPub(b *testing.B) {
	b.StopTimer()
	masterKey, err := NewKeyFromString(bip0032MasterPriv1, mockMainNetParams())
	if err != nil {
		b.Errorf("Failed to decode master seed: %v", err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		masterKey.Neuter()
	}
}

// BenchmarkDeserialize benchmarks how long it takes to deserialize a private
// extended key.
func BenchmarkDeserialize(b *testing.B) {
	mainNetParams := mockMainNetParams()
	for i := 0; i < b.N; i++ {
		NewKeyFromString(bip0032MasterPriv1, mainNetParams)
	}
}

// BenchmarkSerialize benchmarks how long it takes to serialize a private
// extended key.
func BenchmarkSerialize(b *testing.B) {
	b.StopTimer()
	masterKey, err := NewKeyFromString(bip0032MasterPriv1, mockMainNetParams())
	if err != nil {
		b.Errorf("Failed to decode master seed: %v", err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = masterKey.String()
	}
}
