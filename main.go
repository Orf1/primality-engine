package main

import (
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"
)

var (
	zero        = big.NewInt(0)
	one         = big.NewInt(1)
	two         = big.NewInt(2)
	mu          sync.Mutex
	numWorkers  = 32
	smallPrimes = [...]uint{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101}
	currentP    = uint(1)
	start       = time.Now()
)

const (
	iterationsPerJob = 16
)

func main() {
	fmt.Println("Primality Engine 1.0")
	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go spawnWorkers()
	}

	wg.Wait()
}

func spawnWorkers() {
	runtime.LockOSThread()
	for {
		mu.Lock()
		initialP := currentP
		currentP += iterationsPerJob
		mu.Unlock()
		fmt.Printf(".")
		for i := uint(0); i < iterationsPerJob; i++ {
			myP := initialP + i
			if MightBePrime(myP) && LucasLehmer(myP) {
				fmt.Printf("\n2^%d-1 is prime, elapsed: %v\n", myP, time.Since(start))
			}
		}
	}
}

func LucasLehmer(p uint) (isPrime bool) {
	var dummy1, dummy2 big.Int
	s := big.NewInt(4)
	m := big.NewInt(0)
	m = m.Sub(m.Lsh(one, p), one)

	for i := 0; i < int(p)-2; i++ {
		s = s.Sub(s.Mul(s, s), two)

		for s.Cmp(m) == 1 {
			s.Add(dummy1.And(s, m), dummy2.Rsh(s, p))
		}

		if s.Cmp(m) == 0 {
			s = zero
		}
	}

	return s.Cmp(zero) == 0
}

func MightBePrimeAlt(prime uint) bool {
	if big.NewInt(int64(prime)).ProbablyPrime(0) {
		return true
	}
	return false
}

func MightBePrime(prime uint) bool {
	for j := uint(0); j < uint(len(smallPrimes)); j++ {
		if prime%smallPrimes[j] == 0 {
			return false
		}
	}
	return true
}
