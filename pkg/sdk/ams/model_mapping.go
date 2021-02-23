package ams

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
)

var (
	AccountRoleTypeMapping = map[model.AccountRoleType]sdk.AccountRoleType{
		model.AccountRoleType_ADVERTISER:       sdk.AccountRoleTypeAdvertiser,
		model.AccountRoleType_AGENCY:           sdk.AccountRoleTypeAgency,
		model.AccountRoleType_T1:               sdk.AccountRoleTypeT1,
		model.AccountRoleType_BUSINESS_MANAGER: sdk.AccountRoleTypeBusinessManager,
	}

	AccountTypeMapping = map[model.AccountType]sdk.AccountType{
		model.AccountType_UNKNOWN:           sdk.AccountTypeUnknown,
		model.AccountType_ADVERTISER:        sdk.AccountTypeAdvertiser,
		model.AccountType_AGENCY:            sdk.AccountTypeAgency,
		model.AccountType_DSP:               sdk.AccountTypeDSP,
		model.AccountType_DEVELOPER:         sdk.AccountTypeDeveloper,
		model.AccountType_MEMBER:            sdk.AccountTypeMember,
		model.AccountType_EXTERNAL_SUPPLIER: sdk.AccountTypeExternalSupplier,
		model.AccountType_TDC:               sdk.AccountTypeTDC,
		model.AccountType_TONE:              sdk.AccountTypeTone,
		model.AccountType_BM:                sdk.AccountTypeBM,
	}

	RoleTypeMapping = map[model.RoleType]sdk.RoleType{
		model.RoleType_UNKNOWN:                  sdk.RoleTypeUnknown,
		model.RoleType_ADMIN:                    sdk.RoleTypeAdmin,
		model.RoleType_OBSERVER:                 sdk.RoleTypeObserver,
		model.RoleType_OPERATOR:                 sdk.RoleTypeOperator,
		model.RoleType_TREASURER:                sdk.RoleTypeTreasurer,
		model.RoleType_ASSISTANT:                sdk.RoleTypeAssistant,
		model.RoleType_SELF_OPERATOR:            sdk.RoleTypeSelfOperator,
		model.RoleType_ROOT:                     sdk.RoleTypeRoot,
		model.RoleType_AGENCY_BOSS:              sdk.RoleTypeAgencyBoss,
		model.RoleType_AGENCY_ADMIN:             sdk.RoleTypeAgencyAdmin,
		model.RoleType_AGENCY_OBSERVER:          sdk.RoleTypeAgencyObserver,
		model.RoleType_AGENCY_TREASURER:         sdk.RoleTypeAgencyTreasurer,
		model.RoleType_AGENCY_OPERATOR:          sdk.RoleTypeAgencyOperator,
		model.RoleType_AGENCY_PROVIDER_OPERATOR: sdk.RoleTypeAgencyProviderOperator,
		model.RoleType_AGENCY_PROVIDER_OBSERVER: sdk.RoleTypeAgencyProviderObserver,
		model.RoleType_AGENCY_YYB:               sdk.RoleTypeAgencyYYB,
		model.RoleType_AGENCY_AGENT_OPERATOR:    sdk.RoleTypeAgencyAgentOperator,
		model.RoleType_AGENCY_SELF_OPERATOR:     sdk.RoleTypeAgencySelfOperator,
		model.RoleType_AGENCY_MDM_BOSS:          sdk.RoleTypeAgencyMDMBoss,
		model.RoleType_AGENCY_MDM_ADMIN:         sdk.RoleTypeAgencyMDMAdmin,
		model.RoleType_AGENCY_MDM_TREASURER:     sdk.RoleTypeAgencyMDMTreasurer,
		model.RoleType_AGENCY_MDM_OBSERVER:      sdk.RoleTypeAgencyMDMObserver,
		model.RoleType_AGENCY_MDM_OPERATOR:      sdk.RoleTypeAgencyMDMOperator,
	}

	RoleTypeReverseMapping = make(map[sdk.RoleType]model.RoleType)
)

func init() {
	reverseRoleTypeMapping(RoleTypeMapping, RoleTypeReverseMapping)
}

func reverseRoleTypeMapping(src map[model.RoleType]sdk.RoleType, dst map[sdk.RoleType]model.RoleType) {
	for k, v := range src {
		dst[v]= k
	}
}

