package sdk

// AccountRoleType : 授权账号身份类型
type AccountRoleType string

// List of AccountRoleType
const (
	AccountRoleTypeAdvertiser      AccountRoleType = "ACCOUNT_ROLE_TYPE_ADVERTISER"
	AccountRoleTypeAgency          AccountRoleType = "ACCOUNT_ROLE_TYPE_AGENCY"
	AccountRoleTypeT1              AccountRoleType = "ACCOUNT_ROLE_TYPE_T1"
	AccountRoleTypeBusinessManager AccountRoleType = "ACCOUNT_ROLE_TYPE_BUSINESS_MANAGER"
)

// AMSSystemType AMS的投放平台类型
type AMSSystemType string
const (
	AmsEqq AMSSystemType = "ams_eqq"
	AmsMp  AMSSystemType ="ams_mp"
)


// AccountType : 账号类型
type AccountType string

// List of AccountType
const (
	AccountTypeUnknown          AccountType = "ACCOUNT_TYPE_UNKNOWN"
	AccountTypeAdvertiser       AccountType = "ACCOUNT_TYPE_ADVERTISER"
	AccountTypeAgency           AccountType = "ACCOUNT_TYPE_AGENCY"
	AccountTypeDSP              AccountType = "ACCOUNT_TYPE_DSP"
	AccountTypeDeveloper        AccountType = "ACCOUNT_TYPE_DEVELOPER"
	AccountTypeMember           AccountType = "ACCOUNT_TYPE_MEMBER"
	AccountTypeExternalSupplier AccountType = "ACCOUNT_TYPE_EXTERNAL_SUPPLIER"
	AccountTypeTDC              AccountType = "ACCOUNT_TYPE_TDC"
	AccountTypeTone             AccountType = "ACCOUNT_TYPE_TONE"
	AccountTypeBM               AccountType = "ACCOUNT_TYPE_BM"
)

// RoleType : 角色
type RoleType string

// List of RoleType
const (
	RoleTypeUnknown                RoleType = "ROLE_TYPE_UNKNOWN"
	RoleTypeAdmin                  RoleType = "ROLE_TYPE_ADMIN"
	RoleTypeObserver               RoleType = "ROLE_TYPE_OBSERVER"
	RoleTypeOperator               RoleType = "ROLE_TYPE_OPERATOR"
	RoleTypeTreasurer              RoleType = "ROLE_TYPE_TREASURER"
	RoleTypeAssistant              RoleType = "ROLE_TYPE_ASSISTANT"
	RoleTypeSelfOperator           RoleType = "ROLE_TYPE_SELF_OPERATOR"
	RoleTypeRoot                   RoleType = "ROLE_TYPE_ROOT"
	RoleTypeAgencyBoss             RoleType = "ROLE_TYPE_AGENCY_BOSS"
	RoleTypeAgencyAdmin            RoleType = "ROLE_TYPE_AGENCY_ADMIN"
	RoleTypeAgencyObserver         RoleType = "ROLE_TYPE_AGENCY_OBSERVER"
	RoleTypeAgencyTreasurer        RoleType = "ROLE_TYPE_AGENCY_TREASURER"
	RoleTypeAgencyOperator         RoleType = "ROLE_TYPE_AGENCY_OPERATOR"
	RoleTypeAgencyProviderOperator RoleType = "ROLE_TYPE_AGENCY_PROVIDER_OPERATOR"
	RoleTypeAgencyProviderObserver RoleType = "ROLE_TYPE_AGENCY_PROVIDER_OBSERVER"
	RoleTypeAgencyYYB              RoleType = "ROLE_TYPE_AGENCY_YYB"
	RoleTypeAgencyAgentOperator    RoleType = "ROLE_TYPE_AGENCY_AGENT_OPERATOR"
	RoleTypeAgencySelfOperator     RoleType = "ROLE_TYPE_AGENCY_SELF_OPERATOR"
	RoleTypeAgencyMDMBoss          RoleType = "ROLE_TYPE_AGENCY_MDM_BOSS"
	RoleTypeAgencyMDMAdmin         RoleType = "ROLE_TYPE_AGENCY_MDM_ADMIN"
	RoleTypeAgencyMDMTreasurer     RoleType = "ROLE_TYPE_AGENCY_MDM_TREASURER"
	RoleTypeAgencyMDMObserver      RoleType = "ROLE_TYPE_AGENCY_MDM_OBSERVER"
	RoleTypeAgencyMDMOperator      RoleType = "ROLE_TYPE_AGENCY_MDM_OPERATOR"
)
