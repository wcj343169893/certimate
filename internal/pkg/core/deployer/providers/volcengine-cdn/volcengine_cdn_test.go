package volcenginecdn_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "certimate/internal/pkg/core/deployer/providers/volcengine-cdn"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fAccessKey     string
	fSecretKey     string
	fDomain        string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_VOLCENGINECDN_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fAccessKey, argsPrefix+"ACCESSKEY", "", "")
	flag.StringVar(&fSecretKey, argsPrefix+"SECRETKEY", "", "")
	flag.StringVar(&fDomain, argsPrefix+"DOMAIN", "", "")
}

/*
Shell command to run this test:

	go test -v volcengine_cdn_test.go -args \
	--CERTIMATE_DEPLOYER_VOLCENGINECDN_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_VOLCENGINECDN_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_VOLCENGINECDN_ACCESSKEY="your-access-key" \
	--CERTIMATE_DEPLOYER_VOLCENGINECDN_SECRETKEY="your-secret-key" \
	--CERTIMATE_DEPLOYER_VOLCENGINECDN_DOMAIN="example.com"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("ACCESSKEY: %v", fAccessKey),
			fmt.Sprintf("SECRETKEY: %v", fSecretKey),
			fmt.Sprintf("DOMAIN: %v", fDomain),
		}, "\n"))

		deployer, err := provider.New(&provider.VolcEngineCDNDeployerConfig{
			AccessKey: fAccessKey,
			SecretKey: fSecretKey,
			Domain:    fDomain,
		})
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		fInputCertData, _ := os.ReadFile(fInputCertPath)
		fInputKeyData, _ := os.ReadFile(fInputKeyPath)
		res, err := deployer.Deploy(context.Background(), string(fInputCertData), string(fInputKeyData))
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		t.Logf("ok: %v", res)
	})
}
