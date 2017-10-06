package data

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"github.com/spf13/viper"
)

// DynamoDB wraps the dynamo DB library
type DynamoDB struct {
	db *dynamo.DB
}

// Connect begins the connection with the database
func (d *DynamoDB) Connect() error {
	var awsConfig *aws.Config
	if viper.GetString("DynamoEndpoint") == "" {
		// AWS environment.
		awsConfig = &aws.Config{}
	} else {
		// Local development environment.
		myTrue := true
		awsConfig = &aws.Config{
			Endpoint:    aws.String(viper.GetString("DynamoEndpoint")),
			Credentials: credentials.NewStaticCredentials("foo", "foo", "foo"),
			DisableSSL:  &myTrue,
		}
	}
	awsConfig.WithRegion(viper.GetString("DynamoRegion"))
	if gin.Mode() == gin.DebugMode {
		awsConfig.WithLogLevel(aws.LogDebugWithHTTPBody)
	}

	d.db = dynamo.New(session.Must(session.NewSession()), awsConfig)

	return nil
}

func (d *DynamoDB) FindByHash(hash string) (*Avatar, error) {
	var avatar Avatar

	if err := d.getTable().Get("Hash", hash).One(&avatar); err != nil {
		return nil, fmt.Errorf("Cannot find avatar with hash %s", hash)
	}

	return &avatar, nil
}

func (d *DynamoDB) Save(a *Avatar) error {

	if err := d.getTable().Put(a).Run(); err != nil {
		return err
	}

	return nil
}

func (d *DynamoDB) Migrate() error {
	// Nothing to do
	return nil
}

func (d *DynamoDB) getTable() dynamo.Table {
	return d.db.Table(viper.GetString("TableName"))
}
