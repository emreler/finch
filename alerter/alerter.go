package alerter

import "github.com/emreler/finch/storage"
import "github.com/emreler/finch/logger"
import "github.com/emreler/finch/queue"

type Alerter struct {
	lgr logger.InfoErrorLogger
	stg storage.Storage
	q   queue.Queue
}
