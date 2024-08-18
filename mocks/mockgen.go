//go:generate mockgen --build_flags=--mod=mod -destination=mock_db_client.go -package=mocks zt-event-logger/pkg/db DB
//go:generate mockgen --build_flags=--mod=mod -destination=mock_processor.go -package=mocks zt-event-logger/pkg/events Processor

package mocks

import (
	_ "go.uber.org/mock/gomock"
)
