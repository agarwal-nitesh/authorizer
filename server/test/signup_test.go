package test

import (
	"testing"

	"github.com/authorizerdev/authorizer/server/constants"
	"github.com/authorizerdev/authorizer/server/db"
	"github.com/authorizerdev/authorizer/server/enum"
	"github.com/authorizerdev/authorizer/server/graph/model"
	"github.com/authorizerdev/authorizer/server/resolvers"
	"github.com/stretchr/testify/assert"
)

func commonSignupTest(s TestSetup, t *testing.T) {
	email := "signup." + s.TestInfo.Email
	res, err := resolvers.Signup(s.Ctx, model.SignUpInput{
		Email:           email,
		Password:        s.TestInfo.Password,
		ConfirmPassword: s.TestInfo.Password + "s",
	})
	assert.NotNil(t, err, "invalid password errors")

	res, err = resolvers.Signup(s.Ctx, model.SignUpInput{
		Email:           email,
		Password:        s.TestInfo.Password,
		ConfirmPassword: s.TestInfo.Password,
	})

	user := *res.User
	assert.Equal(t, email, user.Email)
	assert.Nil(t, res.AccessToken, "access token should be nil")

	res, err = resolvers.Signup(s.Ctx, model.SignUpInput{
		Email:           email,
		Password:        s.TestInfo.Password,
		ConfirmPassword: s.TestInfo.Password,
	})

	assert.NotNil(t, err, "should throw duplicate email error")

	verificationRequest, err := db.Mgr.GetVerificationByEmail(email, enum.BasicAuthSignup.String())
	assert.Nil(t, err)
	assert.Equal(t, email, verificationRequest.Email)
	cleanData(email)
}

func TestSignUp(t *testing.T) {
	s := testSetup()
	defer s.Server.Close()

	if s.TestInfo.ShouldExecuteForSQL {
		t.Run("signup for sql dbs should pass", func(t *testing.T) {
			constants.DATABASE_URL = s.TestInfo.SQL
			constants.DATABASE_TYPE = enum.Sqlite.String()
			db.InitDB()
			commonSignupTest(s, t)
		})
	}

	if s.TestInfo.ShouldExecuteForArango {
		t.Run("signup for arangodb should pass", func(t *testing.T) {
			constants.DATABASE_URL = s.TestInfo.ArangoDB
			constants.DATABASE_TYPE = enum.Arangodb.String()
			db.InitDB()
			commonSignupTest(s, t)
		})
	}

	if s.TestInfo.ShouldExecuteForMongo {
		t.Run("signup for mongodb should pass", func(t *testing.T) {
			constants.DATABASE_URL = s.TestInfo.MongoDB
			constants.DATABASE_TYPE = enum.Mongodb.String()
			db.InitDB()
			commonSignupTest(s, t)
		})
	}
}
