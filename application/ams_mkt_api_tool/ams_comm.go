package main

import (
	"fmt"
	"git.code.oa.com/going/going/cat/qzs"
	"git.code.oa.com/going/going/codec/qzh"
	"git.code.oa.com/tme/kg_golang_proj/jce_go/app_dcreport"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
	"strconv"
	"time"
)

func DoReport(ctx *qzs.Context, mData map[string]string) error {
	qcReq := &app_dcreport.DataReportReq{}
	qcRsp := &app_dcreport.DataReportRsp{}
	item := app_dcreport.DataReportItem{}
	item.MData = mData
	qcReq.VecReportItem = append(qcReq.VecReportItem, item)

	callDesc := qzh.CallDesc{
		CmdId:       app_dcreport.NUM_CMD_DCREPORT,
		SubCmdId:    app_dcreport.CMD_DCREPORT_NEW,
		AppProtocol: "qza",
	}
	var authInfo qzh.AuthInfo
	qc := qzh.NewQzClient(callDesc, authInfo, "dcreport", qcReq, qcRsp)
	err := qc.Do(ctx)
	if err != nil {
		err := err.(qzh.QzError)
		ctx.Error("dcreport ret:%d,%s", err.RspCode, err.Msg)
		return err
	}
	ctx.Debug("dcreport succ, req:%+v, rsp:%+v", qcReq, qcRsp)

	return nil
}

func HasNextPage(pageInfo *model.Conf) (bool, int64) {
	curIndex := (pageInfo.Page - 1) * pageInfo.PageSize
	if curIndex < pageInfo.TotalNumber {
		return true, pageInfo.Page + 1
	}
	return false, pageInfo.Page
}

func ZeroClockTime(tNow uint32) time.Time {
	tNow /= 86400
	tNow *= 86400
	return time.Unix(int64(tNow), 0).Add(time.Duration(-8) * time.Hour)
}

func GetCampaignInfo(taDs *ads.SDKClient, accountId, campaignId int64) (*model.CampaignsGetResponseData, error) {
	ctx := *taDs.Ctx
	rsp, _, err := taDs.Campaigns().Get(ctx, accountId, &api.CampaignsGetOpts{
		Filtering: optional.NewInterface([]model.FilteringStruct{
			{
				Field:    "campaign_id",
				Operator: "EQUALS",
				Values:   &[]string{strconv.FormatInt(campaignId, 10)},
			},
		}),
		Fields: optional.NewInterface([]string{"campaign_id", "campaign_name"}),
	})
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}

func GetAllCampaignInfo(tds *ads.SDKClient, accountId int64) (map[int64]string, error) {
	mapCampaign := make(map[int64]string)
	ctx := *tds.Ctx
	pageSize := int64(100)
	curPage := int64(1)
	for true {
		rsp := model.CampaignsGetResponseData{}
		var err error
		for i := 0; i < 3; i++ {
			tmpRsp, _, tmpErr := tds.Campaigns().Get(ctx, accountId, &api.CampaignsGetOpts{
				Page:     optional.NewInt64(curPage),
				PageSize: optional.NewInt64(pageSize),
				Fields:   optional.NewInterface([]string{"campaign_id", "campaign_name"}),
			})
			if tmpErr != nil {
				err = tmpErr
				fmt.Printf("account id:%d get campaign info err: %+v\n", accountId, tmpErr)
				continue
			} else {
				err = nil
				rsp = tmpRsp
				break
			}
		}
		if err != nil {
			fmt.Printf("account id:%d all 3 times retry error.\n", accountId)
			return nil, err
		}
		for _, campaignInfo := range *rsp.List {
			mapCampaign[campaignInfo.CampaignId] = campaignInfo.CampaignName
		}
		if hasNext, nextPage := HasNextPage(&model.Conf{Page: rsp.PageInfo.Page, PageSize: rsp.PageInfo.PageSize,
			TotalNumber: rsp.PageInfo.TotalNumber, TotalPage: rsp.PageInfo.TotalPage}); hasNext {
			curPage = nextPage
		} else {
			break
		}
	}

	return mapCampaign, nil
}
