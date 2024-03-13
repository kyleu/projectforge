package queue

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/queries"
)

type Queue struct {
	db     *database.Service
	logger util.Logger

	name    string
	limit   int
	timeout time.Duration
	table   string
	started time.Time

	sent     map[string]int
	received map[string]int
}

func New(ctx context.Context, name string, limit int, timeout time.Duration, table string, db *database.Service, logger util.Logger) (*Queue, error) {
	if limit == 0 {
		limit = 100
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if table == "" {
		table = "queue"
	}
	count, err := initDatabaseIfNeeded(ctx, name, db, table, logger)
	if err != nil {
		return nil, err
	}
	logger.Infof("queue [%s] created; [%d] pending messages", name, count)
	return &Queue{
		db: db, logger: logger, name: name, limit: limit, timeout: timeout, table: table,
		started: time.Now(), sent: map[string]int{}, received: map[string]int{},
	}, nil
}

func (q *Queue) Send(ctx context.Context, tx *sqlx.Tx, m *Message, logger util.Logger) error {
	_, err := q.db.Exec(ctx, queries.QueueWrite(q.table), tx, -1, logger, m.ID, m.Topic, util.ToJSON(m.Param))
	q.sent[m.Topic]++
	return err
}

func (q *Queue) Receive(ctx context.Context, tx *sqlx.Tx, topic string, logger util.Logger) (*Message, error) {
	ownTx := tx == nil
	if ownTx {
		var err error
		tx, err = q.db.StartTransaction(logger)
		if err != nil {
			return nil, err
		}
	}
	now := time.Now()
	nowStr := util.TimeToRFC3339(&now)
	timeout := now.Add(q.timeout)
	timeoutStr := util.TimeToRFC3339(&timeout)
	ret := &struct {
		ID      uuid.UUID `db:"id"`
		Topic   string    `db:"topic"`
		Param   string    `db:"param"`
		Retries int       `db:"retries"`
	}{}
	err := q.db.Get(ctx, ret, queries.QueueRead(q.table), tx, logger, timeoutStr, topic, nowStr, q.limit)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	param, err := util.FromJSONAny([]byte(ret.Param))
	if err != nil {
		return nil, err
	}

	err = q.db.DeleteOne(ctx, queries.QueueDelete(q.table), tx, logger, ret.ID)
	if err != nil {
		return nil, err
	}

	if ownTx {
		if err = tx.Commit(); err != nil {
			return nil, err
		}
	}
	q.received[ret.Topic]++
	return &Message{ID: ret.ID, Topic: ret.Topic, Param: param, Retries: ret.Retries}, nil
}

func (q *Queue) ReceiveBlocking(ctx context.Context, tx *sqlx.Tx, topic string, maxAge time.Duration, logger util.Logger) (*Message, error) {
	ticker := time.NewTicker(maxAge)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			m, err := q.Receive(ctx, tx, topic, logger)
			if err != nil {
				return nil, err
			}
			if m != nil {
				return m, nil
			}
		}
	}
}

func (q *Queue) Remove(ctx context.Context, tx *sql.Tx, id uuid.UUID, logger util.Logger) error {
	_, err := tx.ExecContext(ctx, queries.QueueDelete(q.table), id)
	return err
}

func initDatabaseIfNeeded(ctx context.Context, name string, db *database.Service, table string, logger util.Logger) (int, error) {
	count, err := db.SingleInt(ctx, queries.QueueCount(table), nil, logger, "startup")
	if err != nil {
		logger.Infof("attempting to create table [%s] for queue [%s]", table, name)
		_, err = db.Exec(ctx, queries.QueueCreateTable(table), nil, -1, logger)
		if err != nil {
			return 0, nil
		}
	}
	count, err = db.SingleInt(ctx, queries.QueueCount(table), nil, logger, "startup")
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
