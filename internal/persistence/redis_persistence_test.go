package persistence

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisPersistence_GetCount(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	rp := &RedisPersistence{
		client: mockRedis,
		ctx:    context.Background(),
	}

	t.Run("should return count successfully", func(t *testing.T) {
		mock.ExpectGet("test-key").SetVal("5")
		count, err := rp.GetCount("test-key")
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
		mock.ExpectationsWereMet()
	})

	t.Run("should return 0 if key does not exist", func(t *testing.T) {
		mock.ExpectGet("test-key").RedisNil()
		count, err := rp.GetCount("test-key")
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
		mock.ExpectationsWereMet()
	})

	t.Run("should return error if Redis fails", func(t *testing.T) {
		mock.ExpectGet("test-key").SetErr(errors.New("redis error"))
		count, err := rp.GetCount("test-key")
		assert.Error(t, err)
		assert.Equal(t, 0, count)
		mock.ExpectationsWereMet()
	})
}

func TestRedisPersistence_Increment(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	rp := &RedisPersistence{
		client: mockRedis,
		ctx:    context.Background(),
	}

	t.Run("should increment count and set expiration", func(t *testing.T) {
		mock.ExpectIncr("test-key").SetVal(1)
		mock.ExpectExpire("test-key", time.Minute).SetVal(true)

		err := rp.Increment("test-key")
		assert.NoError(t, err)
		mock.ExpectationsWereMet()
	})

	t.Run("should return error if Incr fails", func(t *testing.T) {
		mock.ExpectIncr("test-key").SetErr(errors.New("redis error"))

		err := rp.Increment("test-key")
		assert.Error(t, err)
		mock.ExpectationsWereMet()
	})

	t.Run("should return error if Expire fails", func(t *testing.T) {
		mock.ExpectIncr("test-key").SetVal(1)
		mock.ExpectExpire("test-key", time.Minute).SetErr(errors.New("expire error"))

		err := rp.Increment("test-key")
		assert.Error(t, err)
		mock.ExpectationsWereMet()
	})
}

func TestRedisPersistence_Block(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	rp := &RedisPersistence{
		client: mockRedis,
		ctx:    context.Background(),
	}

	t.Run("should block key successfully", func(t *testing.T) {
		mock.ExpectSet("test-key:blocked", "1", time.Minute).SetVal("OK")

		err := rp.Block("test-key", time.Minute)
		assert.NoError(t, err)
		mock.ExpectationsWereMet()
	})

	t.Run("should return error if Set fails", func(t *testing.T) {
		mock.ExpectSet("test-key:blocked", "1", time.Minute).SetErr(errors.New("redis error"))

		err := rp.Block("test-key", time.Minute)
		assert.Error(t, err)
		mock.ExpectationsWereMet()
	})
}

func TestRedisPersistence_IsBlocked(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	rp := &RedisPersistence{
		client: mockRedis,
		ctx:    context.Background(),
	}

	t.Run("should return true if key is blocked", func(t *testing.T) {
		mock.ExpectGet("test-key:blocked").SetVal("1")

		isBlocked := rp.IsBlocked("test-key")
		assert.True(t, isBlocked)
		mock.ExpectationsWereMet()
	})

	t.Run("should return false if key is not blocked", func(t *testing.T) {
		mock.ExpectGet("test-key:blocked").RedisNil()

		isBlocked := rp.IsBlocked("test-key")
		assert.False(t, isBlocked)
		mock.ExpectationsWereMet()
	})
}
