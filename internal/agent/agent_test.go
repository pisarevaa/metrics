package agent_test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/pisarevaa/metrics/internal/agent"
)

type AgentTestSuite struct {
	suite.Suite
	client *resty.Client
	config agent.Config
	logger *zap.SugaredLogger
}

func (suite *AgentTestSuite) SetupSuite() {
	suite.config = agent.GetConfig()
	suite.logger = agent.GetLogger()
	suite.client = resty.New()
}

func (suite *AgentTestSuite) TestUpdateMetrics() {
	storage := agent.NewMemStorageRepo()
	service := agent.NewService(suite.client, storage, suite.config, suite.logger)
	errFirst := service.UpdateMetrics()
	suite.Require().NoError(errFirst)
	allocFirst, allocFirstErr := service.Storage.Get("gauge", "Alloc")
	suite.Require().NoError(allocFirstErr)
	suite.Require().NotEmpty(allocFirst)
	randomValueFirst, randomValueFirstErr := service.Storage.Get("gauge", "RandomValue")
	suite.Require().NoError(randomValueFirstErr)
	suite.Require().NotEmpty(randomValueFirst)
	pollCounterFirst, pollCounterFirstErr := service.Storage.Get("counter", "PollCount")
	suite.Require().NoError(pollCounterFirstErr)
	suite.Require().Equal("1", pollCounterFirst)
	errSecond := service.UpdateMetrics()
	suite.Require().NoError(errSecond)
	allocSecond, allocSecondErr := service.Storage.Get("gauge", "Alloc")
	suite.Require().NoError(allocSecondErr)
	suite.Require().NotEmpty(allocSecond)
	suite.Require().NotEqual(allocSecond, allocFirst)
	// randomValueSecond, randomValueSecondErr := service.Storage.Get("gauge", "RandomValue")
	// suite.Require().NoError(randomValueSecondErr)
	// suite.Require().NotEmpty(randomValueSecond)
	// suite.Require().NotEqual(randomValueSecond, randomValueFirst)
	pollCounterSecond, pollCounterSecondErr := service.Storage.Get("counter", "PollCount")
	suite.Require().NoError(pollCounterSecondErr)
	suite.Require().Equal("2", pollCounterSecond)
}

func (suite *AgentTestSuite) TestSendMetrics() {
	storage := agent.NewMemStorageRepo()
	service := agent.NewService(suite.client, storage, suite.config, suite.logger)
	service.SendMetrics()
	err := service.UpdateMetrics()
	suite.Require().NoError(err)
	service.SendMetrics()
	suite.Require().NotEmpty(service.Storage.GetAll())
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}
