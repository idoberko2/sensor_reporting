package main

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SensorsDaoSuite struct {
	suite.Suite
	dao SensorsDao
	cfg Config
}

type measureEntry struct {
	T time.Time
	MeasurePayload
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SensorsDaoSuite))
}

func (suite *SensorsDaoSuite) SetupSuite() {
	suite.Require().NoError(LoadDotEnv())
	cfg, err := ReadConfig(context.Background())
	suite.Require().NoError(err)

	dao := NewSensorsDao(cfg)
	suite.Require().NoError(dao.Init())

	suite.dao = dao
	suite.cfg = cfg
}

func (suite *SensorsDaoSuite) SetupTest() {
	suite.Require().NoError(CleanupDb(suite.cfg))
}

func (suite *SensorsDaoSuite) getAllMeasures() ([]measureEntry, error) {
	db, err := ConnectToDb(suite.cfg)
	suite.Require().NoError(err)
	defer db.Close()

	var res []measureEntry
	query := "SELECT t, sensor, value FROM sensors_data ORDER BY t;"
	if err := db.Select(&res, query); err != nil {
		return nil, err
	}

	return res, nil
}

func (suite *SensorsDaoSuite) TestWriteSensors() {
	now := time.Now()
	err := suite.dao.WriteMeasures(now, []MeasurePayload{
		{36.6, "bmp"},
		{1.12, "dust"},
	})
	suite.Require().NoError(err)

	measures, err := suite.getAllMeasures()
	suite.Require().NoError(err)

	suite.Assert().Len(measures, 2)

	var bmp, dust measureEntry
	for _, measure := range measures {
		if measure.Sensor == "bmp" {
			bmp = measure
		} else if measure.Sensor == "dust" {
			dust = measure
		} else {
			suite.Assert().Fail("unknown measure")
		}
	}

	suite.Assert().Equal(now.Unix(), bmp.T.Unix())
	suite.Assert().Equal(36.6, bmp.Value)

	suite.Assert().Equal(now.Unix(), dust.T.Unix())
	suite.Assert().Equal(1.12, dust.Value)
}
