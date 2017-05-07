package gokwk

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/types"
	"bitbucket.com/sharingmachine/types/errs"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func GetConn(serverAddress string, trustAllCerts bool) (*grpc.ClientConn, error) {
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
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))
	opts = append(opts, grpc.WithTimeout(time.Second*10))
	opts = append(opts, grpc.WithBlock())
	grpclog.SetLogger(models.DebugLogger)

	conn, err := grpc.Dial(serverAddress, opts...)
	return conn, err
}

type Headers struct {
	version string
}

func (i Headers) Context() context.Context {
	if models.Principal == nil {
		return context.Background()
	} else {
		hostname, _ := os.Hostname()
		ctx := metadata.NewContext(
			context.Background(),
			metadata.Pairs(
				types.TokenHeaderName,
				models.Principal.Token,
				"host", hostname,
				"os", runtime.GOOS,
				"agnt", "<not implemented>", //agent //ps -p $$ | tail -1 | awk '{print $NF}'
				"v", i.version,
			),
		)
		return ctx
	}
}

var noAuthMethods = map[string]bool{
	"/usersRpc.UsersRpc/SignIn": true,
	"/usersRpc.UsersRpc/SignUp": true,
}

func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	models.Debug("GRPC: %s %+v", method, req)
	if models.Principal.Token == "" && !noAuthMethods[method] {
		models.Debug("AUTH: No token in request.")
		return errs.NotAuthenticated
	}
	return translateGrpcErr(invoker(ctx, method, req, reply, cc, opts...))
}

// ParseGrpcErr should be used at RPC service call level. i.e. the errors
// returned by the GRPC stubs need to be converted to local errors.
func translateGrpcErr(e error) error {
	if e == nil {
		return nil
	}
	se, _ := status.FromError(e)
	models.Debug("API ERROR:%v", se.Message())
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