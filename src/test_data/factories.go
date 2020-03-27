package test_data

import (
	"fmt"
	"github.com/bluele/factory-go/factory"
	"github.com/google/uuid"
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/models"
	"math/rand"
	"time"
)

var MeasurementFactory = factory.NewFactory(&models.Measurement{}).
	Attr("Uuid", func(args factory.Args) (interface{}, error) {
		return models.MeasurementUUID(fmt.Sprintf("%s", uuid.New())), nil
	}).
	Attr("TargetUuid", func(args factory.Args) (interface{}, error) {
		return models.TargetUUID(fmt.Sprintf("%s", uuid.New())), nil
	}).
	Attr("Type", func(args factory.Args) (interface{}, error) {
		return models.MeasurementTypeHeight, nil
	}).
	Attr("Value", func(args factory.Args) (interface{}, error) {
		return rand.Float32() * 100, nil
	}).
	Attr("Timestamp", func(args factory.Args) (interface{}, error) {
		return time.Now(), nil
	})

var MeasurementStoredFactory = factory.NewFactory(&models.Measurement{}).
	Attr("Uuid", func(args factory.Args) (interface{}, error) {
		return models.MeasurementUUID(fmt.Sprintf("%s", uuid.New())), nil
	}).
	Attr("TargetUuid", func(args factory.Args) (interface{}, error) {
		return models.TargetUUID(fmt.Sprintf("%s", uuid.New())), nil
	}).
	Attr("Type", func(args factory.Args) (interface{}, error) {
		return models.MeasurementTypeHeight, nil
	}).
	Attr("Value", func(args factory.Args) (interface{}, error) {
		return rand.Float32() * 100, nil
	}).
	Attr("Timestamp", func(args factory.Args) (interface{}, error) {
		return time.Now(), nil
	}).
	OnCreate(func(args factory.Args) error {
		return config.Config.DB.Create(args.Instance()).Error
	})
