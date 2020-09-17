package dcrutil_test

import (
	"fmt"
	"math"

	"github.com/hdfchain/hdfd/dcrutil"
)

func ExampleAmount() {

	a := dcrutil.Amount(0)
	fmt.Println("Zero Atom:", a)

	a = dcrutil.Amount(1e8)
	fmt.Println("100,000,000 Atoms:", a)

	a = dcrutil.Amount(1e5)
	fmt.Println("100,000 Atoms:", a)
	// Output:
	// Zero Atom: 0 HDF
	// 100,000,000 Atoms: 1 HDF
	// 100,000 Atoms: 0.001 HDF
}

func ExampleNewAmount() {
	amountOne, err := dcrutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := dcrutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := dcrutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := dcrutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 HDF
	// 0.01234567 HDF
	// 0 HDF
	// invalid coin amount
}

func ExampleAmount_unitConversions() {
	amount := dcrutil.Amount(44433322211100)

	fmt.Println("Atom to kCoin:", amount.Format(dcrutil.AmountKiloCoin))
	fmt.Println("Atom to Coin:", amount)
	fmt.Println("Atom to MilliCoin:", amount.Format(dcrutil.AmountMilliCoin))
	fmt.Println("Atom to MicroCoin:", amount.Format(dcrutil.AmountMicroCoin))
	fmt.Println("Atom to Atom:", amount.Format(dcrutil.AmountAtom))

	// Output:
	// Atom to kCoin: 444.333222111 kHDF
	// Atom to Coin: 444333.222111 HDF
	// Atom to MilliCoin: 444333222.111 mHDF
	// Atom to MicroCoin: 444333222111 μHDF
	// Atom to Atom: 44433322211100 Atom
}
