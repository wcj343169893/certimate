package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	"certimate/internal/domain"
	"certimate/internal/pkg/utils/maps"

	auto "certimate/internal/automation/unicloud"

	xerrors "github.com/pkg/errors"
)

type UnicloudDeployer struct {
	option *DeployerOption
	infos  []string
}

func NewUnicloudDeployer(option *DeployerOption) (Deployer, error) {
	return &UnicloudDeployer{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *UnicloudDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *UnicloudDeployer) GetInfos() []string {
	return d.infos
}

type unicloudData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Id       string `json:"id"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}

func (d *UnicloudDeployer) Deploy(ctx context.Context) error {
	access := &domain.UnicloudAccess{}
	d.infos = append(d.infos, toStr("Unicloud Access", d.option.Access))
	if err := json.Unmarshal([]byte(d.option.Access), access); err != nil {
		return xerrors.Wrap(err, "failed to get access")
	}

	// 打印日志
	d.infos = append(d.infos, toStr("Unicloud Access", access.Username))

	err := auto.UpdateCert(
		maps.GetValueAsString(d.option.DeployConfig.Config, "spaceId"),
		access.Username,
		access.Password,
		d.option.Certificate.Certificate,
		d.option.Certificate.PrivateKey,
	)
	if err != nil {
		return fmt.Errorf("unicloud deploy failed: %w", err)
	}

	d.infos = append(d.infos, toStr("Unicloud Response", string("success")))

	return nil
}
