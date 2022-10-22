package secp256k1

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

type PrivateKey struct {
	Secret       string
	SecretBigInt *big.Int
	Point        *S256Point
}

func MakePrivateKeyFromBigInt(secret *big.Int) (*PrivateKey, error) {
	pk := &PrivateKey{}
	var err error
	pk.Point, err = RMultiply(*GetGeneratorPoint(), *secret)

	if err != nil {
		return nil, err
	}

	// store the secret locally
	pk.Secret = string(secret.Bytes())
	pk.SecretBigInt = secret

	return pk, nil
}

func MakePrivateKey(secret string) (*PrivateKey, error) {
	pk := &PrivateKey{}
	s := new(big.Int).SetBytes([]byte(secret))
	pk.Secret = secret
	var err error
	g := GetGeneratorPoint()
	pk.Point, err = RMultiply(*g, *s)

	if err != nil {
		return nil, err
	}

	return pk, nil
}

func (p PrivateKey) Hex() string {
	return fmt.Sprintf("0x%.64x", []byte(p.Secret))
}

func (p PrivateKey) GetSecretBytes() []byte {
	b := make([]byte, 32)
	for i := 0; i < len(p.Secret); i++ {
		index := 32 - len(p.Secret) + i
		b[index] = p.Secret[i]
	}
	return b
}

// Get a unique determansitic k value for the specified sig (z)
func (p *PrivateKey) GetDeterministsicK(z *big.Int) *big.Int {

	// k=b'\x00'*32
	// v=b'\x01'*32
	k := make([]byte, 32)
	v := make([]byte, 32)
	for i := 0; i < 32; i++ {
		v[i] = 0x01
	}

	if z.Cmp(GetNonce()) > 0 {
		z = z.Sub(z, GetNonce())
	}

	// z_bytes = z.to_bytes(32, 'big')
	// secret_bytes = self.secret.to_bytes(32, 'big')
	zb := z.Bytes()
	zBytes := make([]byte, 32)
	for i := 0; i < len(zb); i++ {
		index := 32 - len(zb) + i
		zBytes[index] = zb[i]
	}
	secretBytes := p.GetSecretBytes()

	// k = hmac.new(k, v + b'\x00' + secret_bytes + z_bytes, s256).digest()
	b := v
	b = append(b, 0x00)
	b = append(b, secretBytes...)
	b = append(b, zBytes...)
	mac := hmac.New(sha256.New, k)
	mac.Write(b)
	k = mac.Sum(nil)

	// v = hmac.new(k, v, s256).digest()
	mac = hmac.New(sha256.New, k)
	mac.Write(v)
	v = mac.Sum(nil)
	fmt.Printf("%x\n", v)

	b = v
	b = append(b, 0x01)
	b = append(b, secretBytes...)
	b = append(b, zBytes...)
	mac = hmac.New(sha256.New, k)
	mac.Write(b)
	k = mac.Sum(nil)

	mac = hmac.New(sha256.New, k)
	mac.Write(v)
	v = mac.Sum(nil)
	fmt.Printf("%x\n", v)

	for {
		mac = hmac.New(sha256.New, k)
		mac.Write(v)
		v = mac.Sum(nil)

		candidate := new(big.Int).SetBytes(v)
		if candidate.Cmp(big.NewInt(1)) >= 0 && candidate.Cmp(GetNonce()) < 0 {
			return candidate
		}

		b = v
		b = append(b, 0x00)
		mac = hmac.New(sha256.New, k)
		mac.Write(b)
		k = mac.Sum(nil)

		mac = hmac.New(sha256.New, k)
		mac.Write(v)
		v = mac.Sum(nil)
	}
}

func (pk PrivateKey) Sign(z *big.Int) (*Signature, error) {

	// k - 32 bytes = 256 bit K
	// b := make([]byte, 32)
	// _, err := rand.Read(b)
	// if err != nil {
	// 	return nil, err
	// }
	// k := new(big.Int).SetBytes(b)

	k := pk.GetDeterministsicK(z)
	fmt.Printf("%x\n", k)
	// // r
	// kG, err := RMultiply(*GetGeneratorPoint(), *k)
	// if err != nil {
	// 	return nil, err
	// }
	// r := kG.X.Num

	// // kInt
	// n2 := new(big.Int).Sub(GetNonce(), big.NewInt(2))
	// kInv := new(big.Int).Exp(k, n2, GetNonce())

	// // s
	// secret := new(big.Int).SetBytes([]byte(p.Secret))
	// rs := new(big.Int).Mul(r, secret)
	// rs = rs.Add(rs, z)
	// rs = rs.Mul(kInv, rs)
	// s := rs.Mod(rs, GetNonce())

	// // if s > N/2:
	// if s.Cmp(new(big.Int).Div(GetNonce(), big.NewInt(2))) > 0 {
	// 	s = new(big.Int).Sub(GetNonce(), s)
	// }

	// return &Signature{
	// 		R: r,
	// 		S: s,
	// 	},
	// 	nil

	e := new(big.Int).SetBytes([]byte(pk.Secret))
	G := GetGeneratorPoint()
	N := GetNonce()

	rPoint, err := RMultiply(*G, *k)
	if err != nil {
		return nil, err
	}
	r := rPoint.Point.X.Num

	n2 := new(big.Int).Sub(N, big.NewInt(2))

	tmp := big.NewInt(0)
	kInv := tmp.Exp(k, n2, N)
	s := new(big.Int).Mul(r, e)
	s = new(big.Int).Add(s, z)
	s = s.Mul(s, kInv)
	s = s.Mod(s, N)

	fmt.Printf("\n%x\n", r)
	fmt.Printf("%x\n", s)

	return &Signature{
			R: r,
			S: s,
		},
		nil
}

func (p PrivateKey) Wif(compressed, testnet bool) []byte {
	var prefix []byte
	var suffix []byte
	if testnet {
		prefix = append(prefix, 0xef)
	} else {
		prefix = append(prefix, 0x80)
	}
	if compressed {
		suffix = append(suffix, 0x01)
	}

	//prefix + secret + suffix
	// secretBytes := make([]byte, 32)
	var secretBytes []byte
	secretBytes = append(secretBytes, p.GetSecretBytes()...)
	s := append(prefix, secretBytes...)
	s = append(s, suffix...)

	return utils.EncodeBase58Checksum(s)
}
