package castle_test

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/buarki/find-castles/castle"
)

func TestCreateId(t *testing.T) {
	testCases := []struct {
		c          castle.Model
		expectedID string
	}{
		{
			c: castle.Model{
				Country: castle.Ireland,
				Name:    "Athenry Castle",
			},
			expectedID: "ireland-athenry_castle",
		},
		{
			c: castle.Model{
				Country: castle.Portugal,
				Name:    "Mau Vizinho(Évora)",
			},
			expectedID: "portugal-mau_vizinhoevora",
		},
		{
			c: castle.Model{
				Country: castle.Portugal,
				Name:    "Montemor-o-Velho",
			},
			expectedID: "portugal-montemor-o-velho",
		},
		{
			c: castle.Model{
				Country: castle.Portugal,
				Name:    "São Martinho de Mouros ",
			},
			expectedID: "portugal-sao_martinho_de_mouros",
		},
		{
			c: castle.Model{
				Country: castle.Portugal,
				Name:    "Monforte (Fig.C.Rodrigo)",
			},
			expectedID: "portugal-monforte_fig.c.rodrigo",
		},
		{
			c: castle.Model{
				Country: castle.Portugal,
				Name:    "ç(oi)ChãoZêzere Freixo de Espada à Cinta",
			},
			expectedID: "portugal-coichaozezere_freixo_de_espada_a_cinta",
		},
	}

	for _, tt := range testCases {
		receivedID := tt.c.StringID()
		if receivedID != tt.expectedID {
			t.Errorf("expected to have ID [%s], got [%s]", tt.expectedID, receivedID)
		}
	}
}

func TestBInaryID(t *testing.T) {
	testCases := []struct {
		c            castle.Model
		expectedHash [32]byte
	}{
		{
			c: castle.Model{
				Country: castle.Portugal,
				Name:    "ç(oi)ChãoZêzere Freixo de Espada à Cinta",
			},
			expectedHash: sha256.Sum256([]byte("portugal-coichaozezere_freixo_de_espada_a_cinta")),
		},
	}

	for _, tt := range testCases {
		receivedID := tt.c.BinaryID()
		if !bytes.Equal(receivedID, tt.expectedHash[:]) {
			t.Errorf("expected BinaryID %x, got %x", tt.expectedHash[:], receivedID)
		}
	}
}
