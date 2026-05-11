// Copyright 2018-2025 Celer Network

package utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/rpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/anypb"
)

var protojsonMarshaler = protojson.MarshalOptions{
	EmitUnpopulated: true,
}

// Dec2HexStr decimal string to hex
func Dec2HexStr(dec string) string {
	i := new(big.Int)
	i.SetString(dec, 10)
	return i.Text(16)
}

// Hex2DecStr hex string to decimal
func Hex2DecStr(hex string) string {
	i := new(big.Int)
	i.SetString(hex, 16)
	return i.Text(10)
}

func BytesToBigInt(in []byte) *big.Int {
	ret := big.NewInt(0)
	ret.SetBytes(in)
	return ret
}

// convert decimal wei string to big.Int
func Wei2BigInt(wei string) *big.Int {
	i := big.NewInt(0)
	_, ok := i.SetString(wei, 10)
	if !ok {
		return nil
	}
	return i
}

// float in 10e18 wei to wei
func Float2Wei(f float64) *big.Int {
	if f < 0 {
		return nil
	}
	wei := decimal.NewFromFloat(f).Mul(decimal.NewFromFloat(10).Pow(decimal.NewFromFloat(18)))
	weiInt := new(big.Int)
	weiInt.SetString(wei.String(), 10)
	return weiInt
}

// left padding
func Pad(origin []byte, n int) []byte {
	m := len(origin)
	padded := make([]byte, n)
	pn := n - m
	for i := m - 1; i >= 0; i-- {
		padded[pn+i] = origin[i]
	}
	return padded
}

func TryLock(m *sync.Mutex) bool {
	const mutexLocked = 1 << iota
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(m)), 0, mutexLocked)
}

// use for celer client dialing celer server/proxy
// support os ca and celer ca
func GetClientTlsOption() grpc.DialOption {
	cpool, _ := x509.SystemCertPool()
	if cpool == nil {
		cpool = x509.NewCertPool()
	}
	cpool.AppendCertsFromPEM(CelerCA)
	if sdkCert != nil && sdkKey != nil {
		cert, _ := tls.X509KeyPair(sdkCert, sdkKey)
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      cpool,
		})
		return grpc.WithTransportCredentials(creds)
	}
	return grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cpool, ""))
}

func IsPermissiveClientTLS() bool {
	v := os.Getenv("AGENTPAY_INSECURE_TLS")
	return v == "1" || strings.EqualFold(v, "true")
}

// GetClientTlsOptionPermissive returns insecure transport credentials when the
// environment variable AGENTPAY_INSECURE_TLS is set ("1"/"true"). This is useful
// for local e2e tests where the server uses a self-signed localhost certificate
// that may not chain to CAs available to the client.
func GetClientTlsOptionPermissive() grpc.DialOption {
	if IsPermissiveClientTLS() {
		// Connect to a TLS server but skip certificate verification.
		// This is only for local e2e where server uses a self-signed localhost cert.
		return grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}))
	}
	return GetClientTlsOption()
}

func IsLoopbackTarget(target string) bool {
	host := target
	if parsedHost, _, err := net.SplitHostPort(target); err == nil {
		host = parsedHost
	}
	host = strings.Trim(host, "[]")
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func WrapLocalTLSDialError(target string, err error) error {
	if err == nil || IsPermissiveClientTLS() || !IsLoopbackTarget(target) {
		return err
	}
	return fmt.Errorf("%w; local TLS hint: %s uses the built-in self-signed localhost certificate, set AGENTPAY_INSECURE_TLS=1 on the dialing process or configure a trusted cert via -tlscert/-tlskey", err, target)
}

// GetClientTlsConfig returns tls.Config with system and celerCA, for https interaction
func GetClientTlsConfig() *tls.Config {
	cpool, _ := x509.SystemCertPool()
	if cpool == nil {
		cpool = x509.NewCertPool()
	}
	cpool.AppendCertsFromPEM(CelerCA)
	return &tls.Config{
		RootCAs: cpool,
	}
}

func ValidateAndFormatAddress(address string) (ctype.Addr, error) {
	if !ethcommon.IsHexAddress(address) {
		return ctype.ZeroAddr, errors.New("INVALID_ADDRESS")
	}
	return ctype.Hex2Addr(address), nil
}

// GetTokenAddr returns token address
func GetTokenAddr(tokenInfo *entity.TokenInfo) ctype.Addr {
	switch tktype := tokenInfo.TokenType; tktype {
	case entity.TokenType_NATIVE:
		return ctype.NativeTokenAddr
	case entity.TokenType_ERC20:
		return ctype.Bytes2Addr(tokenInfo.TokenAddress)
	}
	return ctype.InvalidTokenAddr
}

// GetTokenAddrStr returns string for tokenInfo
func GetTokenAddrStr(tokenInfo *entity.TokenInfo) string {
	return ctype.Addr2Hex(GetTokenAddr(tokenInfo))
}

func PrintToken(tokenInfo *entity.TokenInfo) string {
	if tokenInfo.GetTokenType() == entity.TokenType_NATIVE {
		return "NATIVE"
	}
	return GetTokenAddrStr(tokenInfo)
}

func PrintTokenAddr(tkaddr ctype.Addr) string {
	if tkaddr == ctype.NativeTokenAddr {
		return "NATIVE"
	}
	return ctype.Addr2Hex(tkaddr)
}

// GetTokenInfoFromAddress returns TokenInfo from tkaddr
// only support ERC20 for now
func GetTokenInfoFromAddress(tkaddr ctype.Addr) *entity.TokenInfo {
	tkInfo := new(entity.TokenInfo)
	if tkaddr == ctype.NativeTokenAddr {
		tkInfo.TokenType = entity.TokenType_NATIVE
	} else {
		tkInfo.TokenType = entity.TokenType_ERC20
		tkInfo.TokenAddress = tkaddr.Bytes()
	}
	return tkInfo
}

// Uint64ToBytes converts uint to bytes in big-endian order.
func Uint64ToBytes(i uint64) []byte {
	ret := make([]byte, 8) // 8 bytes for uint64
	binary.BigEndian.PutUint64(ret, i)
	return ret
}

// GetTsAndSig returns current time and signature of current time using sign param passed in.
func GetTsAndSig(sign func([]byte) []byte) (ts uint64, sig []byte) {
	ts = uint64(time.Now().Unix())
	sig = sign(Uint64ToBytes(ts))
	return ts, sig
}

// PbToJSONString marshals a protobuf msg to json string.
//
// This is primarily used for logging/debugging; it is intentionally tolerant
// of unknown google.protobuf.Any type URLs by falling back to an opaque
// base64-encoded representation.
func PbToJSONString(pb proto.Message) (string, error) {
	b, err := protojsonMarshaler.Marshal(pb)
	if err == nil {
		return string(b), nil
	}
	b, err2 := marshalJSONLenientAny(pb, true)
	if err2 == nil {
		return string(b), nil
	}
	return "", err
}

// PbToJSONHexBytes historically used a custom marshaler to render bytes in hex.
// It is now log-only and uses the same tolerant JSON path as PbToJSONString.
func PbToJSONHexBytes(pb proto.Message) (string, error) {
	return PbToJSONString(pb)
}

func marshalJSONLenientAny(pb proto.Message, emitUnpopulated bool) ([]byte, error) {
	v, err := protoMessageToInterface(pb.ProtoReflect(), emitUnpopulated)
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

func protoMessageToInterface(msg protoreflect.Message, emitUnpopulated bool) (any, error) {
	if !msg.IsValid() {
		return nil, nil
	}
	if msg.Descriptor().FullName() == "google.protobuf.Any" {
		if a, ok := msg.Interface().(*anypb.Any); ok {
			return anyToInterface(a, emitUnpopulated)
		}
	}

	out := make(map[string]any)
	fields := msg.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)

		// For real oneofs, only emit the selected field.
		if oneof := fd.ContainingOneof(); oneof != nil && !oneof.IsSynthetic() {
			if !msg.Has(fd) {
				continue
			}
		} else if !emitUnpopulated && !msg.Has(fd) {
			continue
		}

		val := msg.Get(fd)
		converted, err := protoValueToInterface(fd, val, emitUnpopulated)
		if err != nil {
			return nil, err
		}
		if converted == nil {
			continue
		}
		out[fd.JSONName()] = converted
	}
	return out, nil
}

func anyToInterface(a *anypb.Any, emitUnpopulated bool) (any, error) {
	if a == nil {
		return nil, nil
	}
	if a.TypeUrl == "" {
		return map[string]any{"@type": "", "value": base64.StdEncoding.EncodeToString(a.Value)}, nil
	}

	inner, err := anypb.UnmarshalNew(a, proto.UnmarshalOptions{AllowPartial: true, Resolver: protoregistry.GlobalTypes})
	if err != nil {
		return map[string]any{"@type": a.TypeUrl, "value": base64.StdEncoding.EncodeToString(a.Value)}, nil
	}

	innerV, err := protoMessageToInterface(inner.ProtoReflect(), emitUnpopulated)
	if err != nil {
		return nil, err
	}
	if m, ok := innerV.(map[string]any); ok {
		m["@type"] = a.TypeUrl
		return m, nil
	}
	return map[string]any{"@type": a.TypeUrl, "value": innerV}, nil
}

func protoValueToInterface(fd protoreflect.FieldDescriptor, v protoreflect.Value, emitUnpopulated bool) (any, error) {
	if fd.IsList() {
		list := v.List()
		arr := make([]any, 0, list.Len())
		for i := 0; i < list.Len(); i++ {
			elem, err := protoScalarToInterface(fd, list.Get(i), emitUnpopulated)
			if err != nil {
				return nil, err
			}
			arr = append(arr, elem)
		}
		return arr, nil
	}
	if fd.IsMap() {
		m := v.Map()
		out := make(map[string]any)
		m.Range(func(k protoreflect.MapKey, mv protoreflect.Value) bool {
			keyStr := mapKeyToString(k)
			val, err := protoScalarToInterface(fd.MapValue(), mv, emitUnpopulated)
			if err != nil {
				// Range doesn't allow returning an error; capture via closure.
				out = nil
				return false
			}
			out[keyStr] = val
			return true
		})
		if out == nil {
			return nil, errors.New("failed to marshal protobuf map field")
		}
		return out, nil
	}
	return protoScalarToInterface(fd, v, emitUnpopulated)
}

func mapKeyToString(k protoreflect.MapKey) string {
	switch k.Interface().(type) {
	case string:
		return k.String()
	case bool:
		if k.Bool() {
			return "true"
		}
		return "false"
	case int32, int64:
		return strconv.FormatInt(k.Int(), 10)
	case uint32, uint64:
		return strconv.FormatUint(k.Uint(), 10)
	default:
		return k.String()
	}
}

func protoScalarToInterface(fd protoreflect.FieldDescriptor, v protoreflect.Value, emitUnpopulated bool) (any, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return v.Bool(), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return int32(v.Int()), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		// JSON mapping uses strings for 64-bit integers.
		return strconv.FormatInt(v.Int(), 10), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return uint32(v.Uint()), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return strconv.FormatUint(v.Uint(), 10), nil
	case protoreflect.FloatKind:
		return float32(v.Float()), nil
	case protoreflect.DoubleKind:
		return v.Float(), nil
	case protoreflect.StringKind:
		return v.String(), nil
	case protoreflect.BytesKind:
		return base64.StdEncoding.EncodeToString(v.Bytes()), nil
	case protoreflect.EnumKind:
		ed := fd.Enum().Values().ByNumber(v.Enum())
		if ed != nil {
			return string(ed.Name()), nil
		}
		return int32(v.Enum()), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return protoMessageToInterface(v.Message(), emitUnpopulated)
	default:
		return nil, nil
	}
}

func GetAddressFromKeystore(ksBytes []byte) (string, error) {
	type ksStruct struct {
		Address string
	}
	var ks ksStruct
	if err := json.Unmarshal(ksBytes, &ks); err != nil {
		return "", err
	}
	return ks.Address, nil
}

func UnmarshalDelegationDescription(proof *rpc.DelegationProof) (*rpc.DelegationDescription, error) {
	if proof == nil {
		return nil, errors.New("nil delegation proof")
	}

	var desc rpc.DelegationDescription
	err := proto.Unmarshal(proof.GetDelegationDescriptionBytes(), &desc)
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

// CelerCA root CA file, generated via certstrap
// new CA file w/ 10yr expiration
var CelerCA = []byte(`-----BEGIN CERTIFICATE-----
MIIE5DCCAsygAwIBAgIBATANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDEwdDZWxl
ckNBMB4XDTE5MDkxNzIwMTEyNFoXDTI5MDkxNzIwMTEyNFowEjEQMA4GA1UEAxMH
Q2VsZXJDQTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAJMk0gwPrlv/
/8aVaKAoURvPMHEgpDyj8c08A55kjNlStgjtpBpBEfNiLg3nIBeyrrmt/9QXLexD
RxewLgm/9uUB/U1Q4ZBQdAdsSHTyU+wap1HKJ7GX1jhZMjY95vIfTbrSIbo2clym
zKlbwvvBlNYgweHz0YiyWUCsqi2wH++ybUNzmgs0qI+lE/Fg4k8sReVcix5rUvNF
na9tpvGdV9u+iZNlwkeb3Hp9Ank5MR0830LzG2uf95p+d0fXmfl92wxdAFWnEhWi
uPK4Zfqt2orTIpY1uhiDl4d4kf1p0Niowf9FNOHMURYbTQqFMGFLOZI7+dOPW8Wy
AfkcZcgBfEQ2rGkd3+kb8A2pOTBaFqG9HkspKe9d/dXKKZMW63nLU1MIERWmj2S+
uEBAObnNCpmWPDDFnUpNaACt76tqRP+jYaOaEdPp8svOB6mD47lQemb2SEuWd8oa
afbv/tS6rdRvJEPJ2PgHSzIrYG3cTrLrNDhjFa3CoPwilWtP414+AfvwZTRfjNWy
kcVjK8kurHmrhNzhYNC/rtaQquU2NwS0UYf7+sRLytgEK6+6HahHXJD7/6aRyzrm
1bwYANmgrWf7LfvEf4ezGb1m+qLU3B2/Lzh3HpEBw/ySTFXNHKtv7ZbsN+ccjez6
au3fgns4jO8nCF7rAJLChopKHMkGeLkdAgMBAAGjRTBDMA4GA1UdDwEB/wQEAwIB
BjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBSoQm3esw5att/hencMPoPd
qhyvXTANBgkqhkiG9w0BAQsFAAOCAgEAIJATh5yJd7XzfnfM8f2MNKRbUWzdDhE8
AHPjtKoOCsKOyua41xIpPvM+emdg8oTOZdNRlaoDiO1/8DB7PU1k1iXFaZ/MrgeM
Cz8pP9MvXLSXmg039hYREWV7pFvdbqhvfnOU+pj/uMwif1pl6+CRDxxSdwqUNeJr
gmqbDFBvdRa5DQJm7rbIYpSMc5P/GHZcVgOb+g3y6iODaPL/VR7Uo1xVvxzjgxpI
09QcYiDNK5vPondgaoh7W3c+KuaEKO18G8TEN0NGFOadk5ZjJ9uq+8aGfy51qny8
SMOI5/wW+7HODeQmSqtaxVlhZdmWa/iIzya/NGe+5JhRKgBKN9BIysEiVc4i4ver
utwMnSqqDSCZKUD6FeJn+CUimDf9nb9xbsZ8a+5pw2D6/iaZ+mJd1Pv0vHX5NxMJ
36Rj0MMB9I1xY9C2/ugiP1a/JG+Ve4n1r4GX1S2MfYH/k8wYcs4cVLQ21nphvTW5
osJePOuWfBuWD77selYHU/PhlzNVq2bSWHDQlQJoQr12dGk0NiYAf0FtTWRQoMkq
nwCu157ZSeK2bWffJUcLnTFV63ftZmsqEjYHHVQrbthc+LBpTT4ZYOBerEfQfoKi
c6Z63v4v3R9WB5VYZWH7nh+lMJBPhhL1043iN4Be3Z27GJ0jKIPQAL2gfNiukuz/
D/qaayaXjbo=
-----END CERTIFICATE-----`)

var sdkCert []byte
var sdkKey []byte
