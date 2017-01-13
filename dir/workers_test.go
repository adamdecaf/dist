package main

import (
	"github.com/adamdecaf/dist/dist"
	"net"
	"testing"
)

func expect(t *testing.T, wanted, size int) {
	if len(findWorkers(wanted)) != size {
		t.Fatalf("got %d workers, expected %d", wanted, size)
	}
}

func TestFindWorkers__Empty(t *testing.T) {
	expect(t, 0, 0)
	expect(t, 1, 0)
	expect(t, 100, 0)
}

func TestFindWorkers__Populated(t *testing.T) {
	a1 := dist.Address{IP: net.ParseIP("1.2.3.4"), Port: 1234}
	a2 := dist.Address{IP: net.ParseIP("2.3.4.5"), Port: 2345}
	a3 := dist.Address{IP: net.ParseIP("3.4.5.6"), Port: 3456}

	Register(a1)
	expect(t, 0, 0)
	expect(t, 1, 1)

	Register(a2)
	Register(a3)
	expect(t, 0, 0)
	expect(t, 1, 1)
	expect(t, 3, 3)
	expect(t, 4, 3)
	expect(t, 100, 3)
}
