package weather

import "testing"

func TestDescribeKnown(t *testing.T) {
	if got := describe(0); got != "Despejado" {
		t.Fatalf("describe(0)=%q, esperaba Despejado", got)
	}
	if got := describe(61); got != "Lluvia ligera" {
		t.Fatalf("describe(61)=%q, esperaba Lluvia ligera", got)
	}
}

func TestDescribeUnknown(t *testing.T) {
	if got := describe(1234); got != "Desconocido" {
		t.Fatalf("describe(1234)=%q, esperaba Desconocido", got)
	}
}
