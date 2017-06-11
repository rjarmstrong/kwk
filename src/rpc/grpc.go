package rpc

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/kwk-super-snippets/cli/src/cli"
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/errs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/transport"
	"os"
	"runtime"
	"time"
)

// /etc/ssl/certs/COMODO_RSA_Certification_Authority.pem
const cert = `-----BEGIN CERTIFICATE-----
MIIF2DCCA8CgAwIBAgIQTKr5yttjb+Af907YWwOGnTANBgkqhkiG9w0BAQwFADCB
hTELMAkGA1UEBhMCR0IxGzAZBgNVBAgTEkdyZWF0ZXIgTWFuY2hlc3RlcjEQMA4G
A1UEBxMHU2FsZm9yZDEaMBgGA1UEChMRQ09NT0RPIENBIExpbWl0ZWQxKzApBgNV
BAMTIkNPTU9ETyBSU0EgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMTAwMTE5
MDAwMDAwWhcNMzgwMTE4MjM1OTU5WjCBhTELMAkGA1UEBhMCR0IxGzAZBgNVBAgT
EkdyZWF0ZXIgTWFuY2hlc3RlcjEQMA4GA1UEBxMHU2FsZm9yZDEaMBgGA1UEChMR
Q09NT0RPIENBIExpbWl0ZWQxKzApBgNVBAMTIkNPTU9ETyBSU0EgQ2VydGlmaWNh
dGlvbiBBdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCR
6FSS0gpWsawNJN3Fz0RndJkrN6N9I3AAcbxT38T6KhKPS38QVr2fcHK3YX/JSw8X
pz3jsARh7v8Rl8f0hj4K+j5c+ZPmNHrZFGvnnLOFoIJ6dq9xkNfs/Q36nGz637CC
9BR++b7Epi9Pf5l/tfxnQ3K9DADWietrLNPtj5gcFKt+5eNu/Nio5JIk2kNrYrhV
/erBvGy2i/MOjZrkm2xpmfh4SDBF1a3hDTxFYPwyllEnvGfDyi62a+pGx8cgoLEf
Zd5ICLqkTqnyg0Y3hOvozIFIQ2dOciqbXL1MGyiKXCJ7tKuY2e7gUYPDCUZObT6Z
+pUX2nwzV0E8jVHtC7ZcryxjGt9XyD+86V3Em69FmeKjWiS0uqlWPc9vqv9JWL7w
qP/0uK3pN/u6uPQLOvnoQ0IeidiEyxPx2bvhiWC4jChWrBQdnArncevPDt09qZah
SL0896+1DSJMwBGB7FY79tOi4lu3sgQiUpWAk2nojkxl8ZEDLXB0AuqLZxUpaVIC
u9ffUGpVRr+goyhhf3DQw6KqLCGqR84onAZFdr+CGCe01a60y1Dma/RMhnEw6abf
Fobg2P9A3fvQQoh/ozM6LlweQRGBY84YcWsr7KaKtzFcOmpH4MN5WdYgGq/yapiq
crxXStJLnbsQ/LBMQeXtHT1eKJ2czL+zUdqnR+WEUwIDAQABo0IwQDAdBgNVHQ4E
FgQUu69+Aj36pvE8hI6t7jiY7NkyMtQwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB
/wQFMAMBAf8wDQYJKoZIhvcNAQEMBQADggIBAArx1UaEt65Ru2yyTUEUAJNMnMvl
wFTPoCWOAvn9sKIN9SCYPBMtrFaisNZ+EZLpLrqeLppysb0ZRGxhNaKatBYSaVqM
4dc+pBroLwP0rmEdEBsqpIt6xf4FpuHA1sj+nq6PK7o9mfjYcwlYRm6mnPTXJ9OV
2jeDchzTc+CiR5kDOF3VSXkAKRzH7JsgHAckaVd4sjn8OoSgtZx8jb8uk2Intzna
FxiuvTwJaP+EmzzV1gsD41eeFPfR60/IvYcjt7ZJQ3mFXLrrkguhxuhoqEwWsRqZ
CuhTLJK7oQkYdQxlqHvLI7cawiiFwxv/0Cti76R7CZGYZ4wUAc1oBmpjIXUDgIiK
boHGhfKppC3n9KUkEEeDys30jXlYsQab5xoq2Z0B15R97QNKyvDb6KkBPvVWmcke
jkk9u+UJueBPSZI9FoJAzMxZxuY67RIuaTxslbH9qh17f4a+Hg4yRvv7E491f0yL
S0Zj/gA0QHDBw7mh3aZw4gSzQbzpgJHqZJx64SIDqZxubw5lT2yHh17zbqD5daWb
QOhTsiedSrnAdyGN/4fy3ryM7xfft0kL0fJuMAsaDk527RH89elWsn2/x20Kk4yl
0MC2Hb46TpSi125sC8KKfPog88Tk5c0NqMuRkrF8hey1FGlmDoLnzc7ILaZRfyHB
NVOFBkpdn627G190
-----END CERTIFICATE-----`

var (
	lg = logger{}
)

type Rpc struct {
	*grpc.ClientConn
	pr      *cli.UserWithToken
	cliInfo *types.AppInfo
	prefs   *out.Prefs
}

func GetRpc(pr *cli.UserWithToken, prefs *out.Prefs, cliInfo *types.AppInfo,
	serverAddress string, trustAllCerts bool) (*Rpc, error) {
	rpc := &Rpc{pr: pr, cliInfo: cliInfo, prefs: prefs}

	var opts []grpc.DialOption

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(cert))
	creds := credentials.NewTLS(&tls.Config{
		RootCAs:                     pool,
		InsecureSkipVerify:          trustAllCerts,
		ClientSessionCache:          tls.NewLRUClientSessionCache(-1),
		PreferServerCipherSuites:    true,
		DynamicRecordSizingDisabled: false,
		SessionTicketsDisabled:      false,
	})

	//https://github.com/coreos/etcd/blob/master/clientv3/naming/grpc.go
	//https://github.com/lstoll/grpce
	//b := grpc.RoundRobin(r)
	//opts = append(opts, grpc.WithBalancer(b))
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithUnaryInterceptor(rpc.interceptor))
	opts = append(opts, grpc.WithTimeout(time.Second*10))
	opts = append(opts, grpc.WithBlock())
	grpclog.SetLogger(lg)
	//grpc.EnableTracing = false
	out.Debug("API: %s", serverAddress)
	conn, err := grpc.Dial("localhost:8000", opts...)
	if err != nil {
		return nil, err
	}
	rpc.ClientConn = conn
	return rpc, err
}

func (rp *Rpc) Cxf() context.Context {
	if rp.pr == nil {
		return context.Background()
	} else {
		hostname, _ := os.Hostname()
		ctx := metadata.NewContext(
			context.Background(),
			metadata.Pairs(
				"prefs_private_view",
				fmt.Sprintf("%t", rp.prefs.PrivateView),
				types.TokenHeaderName,
				rp.pr.AccessToken,
				"host", hostname,
				"os", runtime.GOOS,
				"agnt", "<not implemented>", //agent //ps -p $$ | tail -1 | awk '{print $NF}'
				"v", rp.cliInfo.String(),
			),
		)
		return ctx
	}
}

var noAuthMethods = map[string]bool{
	"/types.Users/SignIn": true,
	"/types.Users/SignUp": true,
}

func (rp *Rpc) interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	out.Debug("GRPC: %s", method)
	if !rp.pr.HasAccessToken() && !noAuthMethods[method] {
		out.Debug("AUTH: No token in request.")
		return errs.NotAuthenticated
	}
	out.Debug("CTX: %+v", ctx)
	out.Debug("OPTS: %+v", opts)
	err := invoker(ctx, method, req, reply, cc, opts...)
	return translateGrpcErr(err)
}

// ParseGrpcErr should be used at RPC service call level. i.e. the errors
// returned by the GRPC stubs need to be converted to local errors.
func translateGrpcErr(e error) error {
	if e == nil {
		return nil
	}
	se, _ := status.FromError(e)
	out.Debug("API ERROR: %v", e)
	switch se.Code() {
	case codes.InvalidArgument:
		te := &errs.Error{}
		err := json.Unmarshal([]byte(se.Message()), te)
		if err != nil {
			return err
		}
		return te
	case codes.Unauthenticated:
		return errs.NotAuthenticated
	case codes.NotFound:
		return errs.NotFound
	case codes.AlreadyExists:
		return errs.AlreadyExists
	case codes.PermissionDenied:
		return errs.PermissionDenied
	case codes.Unimplemented:
		return errs.NotImplemented
	case codes.Internal:
		return errs.Internal
	case codes.Unavailable:
		return errs.ApiDown
	}
	return e
}

type logger struct {
}

func (logger) Fatal(args ...interface{}) {
	out.DebugLogger.Fatal(args...)
}

func (logger) Fatalf(format string, args ...interface{}) {
	out.DebugLogger.Fatalf(format, args...)
}

func (logger) Fatalln(args ...interface{}) {
	out.DebugLogger.Fatalln(args...)
}

func (logger) Print(args ...interface{}) {
	out.DebugLogger.Print(args...)
}

var attempts = 0

func (logger) Printf(format string, args ...interface{}) {
	_, ok := args[0].(transport.ConnectionError)
	if ok {
		if attempts == 0 {
			fmt.Print("\n", style.Margin, "Connecting to kwk .")
		}
		fmt.Print(".")
		attempts++
	}
	out.Debug(format, args)
}

func (logger) Println(args ...interface{}) {
	out.DebugLogger.Println(args...)
}
