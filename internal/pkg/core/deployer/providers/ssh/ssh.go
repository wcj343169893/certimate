package ssh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	xerrors "github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"certimate/internal/pkg/core/deployer"
	"certimate/internal/pkg/utils/x509"
)

type SshDeployerConfig struct {
	// SSH 主机。
	// 零值时默认为 "localhost"。
	SshHost string `json:"sshHost,omitempty"`
	// SSH 端口。
	// 零值时默认为 22。
	SshPort int32 `json:"sshPort,omitempty"`
	// SSH 登录用户名。
	SshUsername string `json:"sshUsername,omitempty"`
	// SSH 登录密码。
	SshPassword string `json:"sshPassword,omitempty"`
	// SSH 登录私钥。
	SshKey string `json:"sshKey,omitempty"`
	// SSH 登录私钥口令。
	SshKeyPassphrase string `json:"sshKeyPassphrase,omitempty"`
	// 前置命令。
	PreCommand string `json:"preCommand,omitempty"`
	// 后置命令。
	PostCommand string `json:"postCommand,omitempty"`
	// 输出证书格式。
	OutputFormat OutputFormatType `json:"outputFormat,omitempty"`
	// 输出证书文件路径。
	OutputCertPath string `json:"outputCertPath,omitempty"`
	// 输出私钥文件路径。
	OutputKeyPath string `json:"outputKeyPath,omitempty"`
	// PFX 导出密码。
	// 证书格式为 PFX 时必填。
	PfxPassword string `json:"pfxPassword,omitempty"`
	// JKS 别名。
	// 证书格式为 JKS 时必填。
	JksAlias string `json:"jksAlias,omitempty"`
	// JKS 密钥密码。
	// 证书格式为 JKS 时必填。
	JksKeypass string `json:"jksKeypass,omitempty"`
	// JKS 存储密码。
	// 证书格式为 JKS 时必填。
	JksStorepass string `json:"jksStorepass,omitempty"`
}

type SshDeployer struct {
	config *SshDeployerConfig
	logger deployer.Logger
}

var _ deployer.Deployer = (*SshDeployer)(nil)

func New(config *SshDeployerConfig) (*SshDeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *SshDeployerConfig, logger deployer.Logger) (*SshDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	return &SshDeployer{
		logger: logger,
		config: config,
	}, nil
}

func (d *SshDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 连接
	client, err := createSshClient(
		d.config.SshHost,
		d.config.SshPort,
		d.config.SshUsername,
		d.config.SshPassword,
		d.config.SshKey,
		d.config.SshKeyPassphrase,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssh client")
	}
	defer client.Close()

	d.logger.Logt("SSH connected")

	// 执行前置命令
	if d.config.PreCommand != "" {
		stdout, stderr, err := execSshCommand(client, d.config.PreCommand)
		if err != nil {
			return nil, xerrors.Wrapf(err, "failed to run pre-command: stdout: %s, stderr: %s", stdout, stderr)
		}

		d.logger.Logt("SSH pre-command executed", stdout)
	}

	// 上传证书和私钥文件
	switch d.config.OutputFormat {
	case OUTPUT_FORMAT_PEM:
		if err := writeSftpFileString(client, d.config.OutputCertPath, certPem); err != nil {
			return nil, xerrors.Wrap(err, "failed to upload certificate file")
		}

		d.logger.Logt("certificate file uploaded")

		if err := writeSftpFileString(client, d.config.OutputKeyPath, privkeyPem); err != nil {
			return nil, xerrors.Wrap(err, "failed to upload private key file")
		}

		d.logger.Logt("private key file uploaded")

	case OUTPUT_FORMAT_PFX:
		pfxData, err := x509.TransformCertificateFromPEMToPFX(certPem, privkeyPem, d.config.PfxPassword)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to transform certificate to PFX")
		}

		d.logger.Logt("certificate transformed to PFX")

		if err := writeSftpFile(client, d.config.OutputCertPath, pfxData); err != nil {
			return nil, xerrors.Wrap(err, "failed to upload certificate file")
		}

		d.logger.Logt("certificate file uploaded")

	case OUTPUT_FORMAT_JKS:
		jksData, err := x509.TransformCertificateFromPEMToJKS(certPem, privkeyPem, d.config.JksAlias, d.config.JksKeypass, d.config.JksStorepass)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to transform certificate to JKS")
		}

		d.logger.Logt("certificate transformed to JKS")

		if err := writeSftpFile(client, d.config.OutputCertPath, jksData); err != nil {
			return nil, xerrors.Wrap(err, "failed to upload certificate file")
		}

		d.logger.Logt("certificate file uploaded")

	default:
		return nil, fmt.Errorf("unsupported output format: %s", d.config.OutputFormat)
	}

	// 执行后置命令
	if d.config.PostCommand != "" {
		stdout, stderr, err := execSshCommand(client, d.config.PostCommand)
		if err != nil {
			return nil, xerrors.Wrapf(err, "failed to run command, stdout: %s, stderr: %s", stdout, stderr)
		}

		d.logger.Logt("SSH post-command executed", stdout)
	}

	return &deployer.DeployResult{}, nil
}

func createSshClient(host string, port int32, username string, password string, key string, keyPassphrase string) (*ssh.Client, error) {
	if host == "" {
		host = "localhost"
	}

	if port == 0 {
		port = 22
	}

	var authMethod ssh.AuthMethod
	if key != "" {
		var signer ssh.Signer
		var err error

		if keyPassphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(key), []byte(keyPassphrase))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(key))
		}

		if err != nil {
			return nil, err
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		authMethod = ssh.Password(password)
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
}

func execSshCommand(sshCli *ssh.Client, command string) (string, string, error) {
	session, err := sshCli.NewSession()
	if err != nil {
		return "", "", err
	}

	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf
	err = session.Run(command)
	if err != nil {
		return "", "", err
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}

func writeSftpFileString(sshCli *ssh.Client, path string, content string) error {
	return writeSftpFile(sshCli, path, []byte(content))
}

func writeSftpFile(sshCli *ssh.Client, path string, data []byte) error {
	sftpCli, err := sftp.NewClient(sshCli)
	if err != nil {
		return xerrors.Wrap(err, "failed to create sftp client")
	}
	defer sftpCli.Close()

	if err := sftpCli.MkdirAll(filepath.Dir(path)); err != nil {
		return xerrors.Wrap(err, "failed to create remote directory")
	}

	file, err := sftpCli.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return xerrors.Wrap(err, "failed to open remote file")
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return xerrors.Wrap(err, "failed to write to remote file")
	}

	return nil
}
