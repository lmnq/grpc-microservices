package db

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MySQLContainer struct {
	*mysql.MySQLContainer
	DataSourceURL string
}

func CreateMySqlContainer(ctx context.Context) (*MySQLContainer, error) {
	container, err := mysql.Run(ctx,
		"mysql",
		mysql.WithDatabase("order"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server - GPL").
				WithOccurrence(1).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := container.ConnectionString(ctx,
		"charset=utf8mb4",
		"parseTime=true",
		"loc=Local",
		"tls=skip-verify",
	)
	if err != nil {
		return nil, err
	}

	return &MySQLContainer{
		MySQLContainer: container,
		DataSourceURL:  connStr,
	}, nil
}
