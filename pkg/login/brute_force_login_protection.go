package login

import (
	"time"

	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/models"
)

var (
	maxInvalidLoginAttempts int64 = 5
	loginAttemptsWindow           = time.Minute * 5
)

var validateLoginAttempts = func(query *models.LoginUserQuery) error {
	if query.Cfg.DisableBruteForceLoginProtection {
		return nil
	}

	loginAttemptCountQuery := models.GetUserLoginAttemptCountQuery{
		Username: query.Username,
		Since:    time.Now().Add(-loginAttemptsWindow),
	}

	if err := bus.Dispatch(&loginAttemptCountQuery); err != nil {
		return err
	}

	if loginAttemptCountQuery.Result >= maxInvalidLoginAttempts {
		return ErrTooManyLoginAttempts
	}

	return nil
}

var saveInvalidLoginAttempt = func(query *models.LoginUserQuery) error {
	if query.Cfg.DisableBruteForceLoginProtection {
		return nil
	}

	loginAttemptCommand := models.CreateLoginAttemptCommand{
		Username:  query.Username,
		IpAddress: query.IpAddress,
	}

	return bus.Dispatch(&loginAttemptCommand)
}
