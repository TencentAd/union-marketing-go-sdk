package api

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/api/io"
)

var (
	instance *SDK
)

type SDK struct {

}

func InitSDK(config *config.Config) error {
	return nil
}

func Call() error {
	return nil
}

func GetReport(input *io.GetReportInput) (*io.GetReportOutput, error) {
	return nil, nil
}
