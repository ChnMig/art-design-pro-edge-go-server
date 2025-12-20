package user

import (
	"unicode/utf8"

	"api-server/config"
	"api-server/db/pgdb/system"
)

type LoginInput struct {
	TenantCode string
	Account    string
	Password   string
}

func VerifyLogin(input LoginInput) (system.SystemUser, system.SystemTenant, error) {
	user, tenant, err := system.VerifyUser(input.TenantCode, input.Account, input.Password)
	if err != nil {
		return system.SystemUser{}, system.SystemTenant{}, err
	}
	if user.ID == 0 {
		return system.SystemUser{}, tenant, ErrInvalidCredentials
	}
	if user.Status != system.StatusEnabled {
		return system.SystemUser{}, tenant, ErrUserDisabled
	}
	return user, tenant, nil
}

func CreateLoginLog(item *system.SystemUserLoginLog) error {
	return system.CreateLoginLog(item)
}

type LoginLogInput struct {
	TenantCode  string
	UserName    string
	IP          string
	LoginStatus string
}

func CreateLoginLogFromInput(input LoginLogInput) error {
	log := system.SystemUserLoginLog{
		TenantCode:  input.TenantCode,
		UserName:    input.UserName,
		Password:    "",
		IP:          input.IP,
		LoginStatus: input.LoginStatus,
	}
	return CreateLoginLog(&log)
}

type FindLoginLogQuery struct {
	IP       string
	Username string
}

func FindLoginLogList(query FindLoginLogQuery, page, pageSize int) ([]system.SystemUserLoginLog, int64, error) {
	filter := system.SystemUserLoginLog{
		IP:       query.IP,
		UserName: query.Username,
	}
	return system.FindLoginLogList(&filter, page, pageSize)
}

type TenantSuggestion struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func SuggestTenantForLogin(code string, limit int) ([]TenantSuggestion, error) {
	if utf8.RuneCountInString(code) < config.TenantMinQueryLength {
		return nil, ErrTenantQueryTooShort
	}

	tenants, err := system.SuggestTenantByCode(code, limit)
	if err != nil {
		return nil, err
	}

	result := make([]TenantSuggestion, 0, len(tenants))
	for _, t := range tenants {
		result = append(result, TenantSuggestion{
			ID:   t.ID,
			Code: t.Code,
			Name: t.Name,
		})
	}
	return result, nil
}
