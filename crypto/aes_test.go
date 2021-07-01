package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESGCMEncrypt(t *testing.T) {
	var (
		key      = "1234"
		plainTxt = "abcd"
	)

	encTxt, err := AESGCMEncrypt(plainTxt, key)
	if err != nil {
		t.Fatalf("Error: %+v", err)
	}
	t.Logf("Encrypted text: %s", encTxt)
}

func TestAESGCMDecrypt(t *testing.T) {
	var (
		key        = "1234"
		cipherText = "8dfb29477e928dff166c700cc97a933b0aeac54a"
	)

	plainText, err := AESGCMDecrypt(cipherText, key)
	if err != nil {
		t.Fatalf("Error: %+v", err)
	}
	t.Logf("Plain text: %s", plainText)

	assert.Equal(t, plainText, "abcd")
}

func TestAESDecrypt(t *testing.T) {
	var (
		key        = "1234567891011121"
		cipherText = "naqsXiO6195ImApi1K_d-KQiuhz9OqgZ2kUCF_XEgH7VHB-0cOf8Ko1N30PcI18BYTrqHdUS1C55_aFVhkJtCDOk4MoKk1neImcT9Ux4yKA"
	)

	plainTxt, err := AESDecrypt(cipherText, key)
	if err != nil {
		t.Fatalf("Error: %+v", err)
	}
	t.Logf("Plain text: %s", plainTxt)
}

func TestAESEncrypt(t *testing.T) {
	type args struct {
		plaintext string
		key       string
	}
	tests := []struct {
		name           string
		args           args
		wantCiphertext string
	}{
		{"happy_case",
			args{
				plaintext: "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
				key:       "1234567891011121",
			},
			"naqsXiO6195ImApi1K_d-KQiuhz9OqgZ2kUCF_XEgH7VHB-0cOf8Ko1N30PcI18BYTrqHdUS1C55_aFVhkJtCDOk4MoKk1neImcT9Ux4yKA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCiphertext, err := AESEncrypt(tt.args.plaintext, tt.args.key)
			if err != nil {
				t.Error(err)
			}
			//
			t.Log(gotCiphertext)
			if gotCiphertext != tt.wantCiphertext {
				t.Errorf("AESEncrypt() gotCiphertext = %v, want %v", gotCiphertext, tt.wantCiphertext)
			}
		})
	}
}
