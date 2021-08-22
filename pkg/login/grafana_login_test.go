package login

import (
	"testing"

	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginUsingGrafinsightDB(t *testing.T) {
	grafinsightLoginScenario(t, "When login with non-existing user", func(sc *grafinsightLoginScenarioContext) {
		sc.withNonExistingUser()
		err := loginUsingGrafinsightDB(sc.loginUserQuery)
		require.EqualError(t, err, models.ErrUserNotFound.Error())

		assert.False(t, sc.validatePasswordCalled)
		assert.Nil(t, sc.loginUserQuery.User)
	})

	grafinsightLoginScenario(t, "When login with invalid credentials", func(sc *grafinsightLoginScenarioContext) {
		sc.withInvalidPassword()
		err := loginUsingGrafinsightDB(sc.loginUserQuery)

		require.EqualError(t, err, ErrInvalidCredentials.Error())

		assert.True(t, sc.validatePasswordCalled)
		assert.Nil(t, sc.loginUserQuery.User)
	})

	grafinsightLoginScenario(t, "When login with valid credentials", func(sc *grafinsightLoginScenarioContext) {
		sc.withValidCredentials()
		err := loginUsingGrafinsightDB(sc.loginUserQuery)
		require.NoError(t, err)

		assert.True(t, sc.validatePasswordCalled)

		require.NotNil(t, sc.loginUserQuery.User)
		assert.Equal(t, sc.loginUserQuery.Username, sc.loginUserQuery.User.Login)
		assert.Equal(t, sc.loginUserQuery.Password, sc.loginUserQuery.User.Password)
	})

	grafinsightLoginScenario(t, "When login with disabled user", func(sc *grafinsightLoginScenarioContext) {
		sc.withDisabledUser()
		err := loginUsingGrafinsightDB(sc.loginUserQuery)
		require.EqualError(t, err, ErrUserDisabled.Error())

		assert.False(t, sc.validatePasswordCalled)
		assert.Nil(t, sc.loginUserQuery.User)
	})
}

type grafinsightLoginScenarioContext struct {
	loginUserQuery         *models.LoginUserQuery
	validatePasswordCalled bool
}

type grafinsightLoginScenarioFunc func(c *grafinsightLoginScenarioContext)

func grafinsightLoginScenario(t *testing.T, desc string, fn grafinsightLoginScenarioFunc) {
	t.Helper()

	t.Run(desc, func(t *testing.T) {
		origValidatePassword := validatePassword

		sc := &grafinsightLoginScenarioContext{
			loginUserQuery: &models.LoginUserQuery{
				Username:  "user",
				Password:  "pwd",
				IpAddress: "192.168.1.1:56433",
			},
			validatePasswordCalled: false,
		}

		t.Cleanup(func() {
			validatePassword = origValidatePassword
		})

		fn(sc)
	})
}

func mockPasswordValidation(valid bool, sc *grafinsightLoginScenarioContext) {
	validatePassword = func(providedPassword string, userPassword string, userSalt string) error {
		sc.validatePasswordCalled = true

		if !valid {
			return ErrInvalidCredentials
		}

		return nil
	}
}

func (sc *grafinsightLoginScenarioContext) getUserByLoginQueryReturns(user *models.User) {
	bus.AddHandler("test", func(query *models.GetUserByLoginQuery) error {
		if user == nil {
			return models.ErrUserNotFound
		}

		query.Result = user
		return nil
	})
}

func (sc *grafinsightLoginScenarioContext) withValidCredentials() {
	sc.getUserByLoginQueryReturns(&models.User{
		Id:       1,
		Login:    sc.loginUserQuery.Username,
		Password: sc.loginUserQuery.Password,
		Salt:     "salt",
	})
	mockPasswordValidation(true, sc)
}

func (sc *grafinsightLoginScenarioContext) withNonExistingUser() {
	sc.getUserByLoginQueryReturns(nil)
}

func (sc *grafinsightLoginScenarioContext) withInvalidPassword() {
	sc.getUserByLoginQueryReturns(&models.User{
		Password: sc.loginUserQuery.Password,
		Salt:     "salt",
	})
	mockPasswordValidation(false, sc)
}

func (sc *grafinsightLoginScenarioContext) withDisabledUser() {
	sc.getUserByLoginQueryReturns(&models.User{
		IsDisabled: true,
	})
}
