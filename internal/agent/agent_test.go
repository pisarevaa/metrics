package agent_test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"

	"github.com/pisarevaa/metrics/internal/agent"
)

type AgentTestSuite struct {
	suite.Suite
	client *resty.Client
	config agent.Config
}

func (suite *AgentTestSuite) SetupSuite() {
	suite.config = agent.GetConfig()
	suite.client = resty.New()
}

func (suite *AgentTestSuite) TestUpdateMetrics() {
	storage := agent.NewMemStorageRepo()
	service := agent.NewService(suite.client, storage, suite.config)
	errFirst := service.UpdateMetrics()
	suite.Require().NoError(errFirst)
	heapInuseFirst, heapInuseFirstErr := service.Storage.Get("gauge", "HeapInuse")
	suite.Require().NoError(heapInuseFirstErr)
	suite.Require().NotEmpty(heapInuseFirst)
	randomValueFirst, randomValueFirstErr := service.Storage.Get("gauge", "RandomValue")
	suite.Require().NoError(randomValueFirstErr)
	suite.Require().NotEmpty(randomValueFirst)
	pollCounterFirst, pollCounterFirstErr := service.Storage.Get("counter", "PollCount")
	suite.Require().NoError(pollCounterFirstErr)
	suite.Require().Equal("1", pollCounterFirst)
	errSecond := service.UpdateMetrics()
	suite.Require().NoError(errSecond)
	heapInuseSecond, heapInuseSecondErr := service.Storage.Get("gauge", "HeapInuse")
	suite.Require().NoError(heapInuseSecondErr)
	suite.Require().NotEmpty(heapInuseSecond)
	suite.Require().NotEqual(heapInuseSecond, heapInuseFirst)
	randomValueSecond, randomValueSecondErr := service.Storage.Get("gauge", "RandomValue")
	suite.Require().NoError(randomValueSecondErr)
	suite.Require().NotEmpty(randomValueSecond)
	suite.Require().NotEqual(randomValueSecond, randomValueFirst)
	pollCounterSecond, pollCounterSecondErr := service.Storage.Get("counter", "PollCount")
	suite.Require().NoError(pollCounterSecondErr)
	suite.Require().Equal("2", pollCounterSecond)
}

func (suite *AgentTestSuite) TestSendMetrics() {
	storage := agent.NewMemStorageRepo()
	service := agent.NewService(suite.client, storage, suite.config)
	service.SendMetrics()
	err := service.UpdateMetrics()
	suite.Require().NoError(err)
	service.SendMetrics()
	suite.Require().NotEmpty(service.Storage.GetAll())
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}
