package rdb

import (
	"fmt"
	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/pkg/errors"

	"github.com/eiicon-company/go-utils/util"
)

type (
	// RDB is ....
	RDB interface {
		Get(key string) (string, error)
		GetInt(key string) (int64, error)
		GetFloat(key string) (float64, error)
		Set(key string, val string) (string, error)
		LRange(key string, first, last int) ([]string, error)
		RPush(key string, val ...string) (int64, error)
		IncrBy(key string, val int64) (int64, error)
		IncrByFloat(key string, val float64) (float64, error)
		HIncrBy(key, field string, val int64) (int64, error)
		HIncrByFloat(key, field string, val float64) (float64, error)
		HGet(key, field string) (string, error)
		HGetAll(key string) (map[string]string, error)
		Del(key string) (int64, error)
	}

	rdb struct {
		Env  util.Environment `inject:""`
		Pool *pool.Pool       `inject:""`
	}
)

func (rc *rdb) getCmd(cmd string, args ...interface{}) (*redis.Resp, error) {
	conn, err := rc.Pool.Get()
	if err != nil {
		return nil, err
	}
	defer rc.Pool.Put(conn)

	return conn.Cmd(cmd, args...), nil
}

func (rc *rdb) setCmd(cmd string, args ...interface{}) (*redis.Resp, error) {
	conn, err := rc.Pool.Get()
	if err != nil {
		return nil, err
	}
	defer rc.Pool.Put(conn)

	// return conn.Cmd(cmd, args[:len(args)-1]...), nil
	return conn.Cmd(cmd, args...), nil
}

func (rc *rdb) Get(key string) (string, error) {
	resp, err := rc.getCmd("GET", key)
	if err != nil {
		return "", err
	}

	return resp.Str()
	// switch {
	// default:
	// return resp.Str()
	// case resp.IsType(redis.Int):
	// return resp.Int()
	// }
}

func (rc *rdb) GetInt(key string) (int64, error) {
	numeric, err := rc.Get(key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(numeric, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("%s cloud not cast", key)
		return 0, errors.Wrap(err, msg)
	}

	return n, nil
}

func (rc *rdb) GetFloat(key string) (float64, error) {
	numeric, err := rc.Get(key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseFloat(numeric, 64)
	if err != nil {
		msg := fmt.Sprintf("%s cloud not cast", key)
		return 0, errors.Wrap(err, msg)
	}

	return n, nil
}

func (rc *rdb) Set(key string, val string) (string, error) {
	resp, err := rc.setCmd("SET", key, val)
	if err != nil {
		return "", err
	}

	return resp.Str()
}

func (rc *rdb) LRange(key string, first, last int) ([]string, error) {
	resp, err := rc.getCmd("LRANGE", key, first, last)
	if err != nil {
		return nil, err
	}

	ary, err := resp.Array()
	if err != nil {
		return nil, err
	}
	r := make([]string, len(ary))
	for i, a := range ary {
		r[i], _ = a.Str()
	}

	return r, nil
}

func (rc *rdb) RPush(key string, val ...string) (int64, error) {
	values := append([]string{key}, val...)

	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}

	resp, err := rc.setCmd("RPUSH", args...)
	if err != nil {
		return 0, err
	}

	return resp.Int64()
}

func (rc *rdb) IncrBy(key string, val int64) (int64, error) {
	resp, err := rc.setCmd("INCRBY", key, val)
	if err != nil {
		return 0, err
	}

	return resp.Int64()
}

func (rc *rdb) IncrByFloat(key string, val float64) (float64, error) {
	resp, err := rc.setCmd("INCRBYFLOAT", key, val)
	if err != nil {
		return 0.0, err
	}

	return resp.Float64()
}

func (rc *rdb) HIncrBy(key, field string, val int64) (int64, error) {
	resp, err := rc.setCmd("HINCRBY", key, field, val)
	if err != nil {
		return 0, err
	}

	return resp.Int64()
}

func (rc *rdb) HIncrByFloat(key, field string, val float64) (float64, error) {
	resp, err := rc.setCmd("HINCRBYFLOAT", key, field, val)
	if err != nil {
		return 0.0, err
	}

	return resp.Float64()
}

func (rc *rdb) HGet(key, field string) (string, error) {
	resp, err := rc.getCmd("HGET", key, field)
	if err != nil {
		return "", err
	}

	return resp.Str()
}

func (rc *rdb) HGetAll(key string) (map[string]string, error) {
	resp, err := rc.getCmd("HGETALL", key)
	if err != nil {
		return nil, err
	}

	return resp.Map()
}

func (rc *rdb) Del(key string) (int64, error) {
	resp, err := rc.setCmd("DEL", key)
	if err != nil {
		return 0, err
	}

	return resp.Int64()
}
