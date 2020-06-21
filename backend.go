package quadtree

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

type Backend interface {
	GetQuadtree(uuid uuid.UUID) (*Quadtree, error)
	SetQuadtree(qt *Quadtree) error
	DelQuadtree(uuid uuid.UUID) error
	GetNodes(parent uuid.UUID) ([]uuid.UUID, error)
	SetNodes(parent uuid.UUID, nodes []uuid.UUID) error
	DetachNodes(parent uuid.UUID) error
	GetObjects(parent uuid.UUID) ([]*Object, error)
	AddObject(parent uuid.UUID, object *Object) error
	DelObject(parent uuid.UUID, object *Object) error
	ClearObjects(parent uuid.UUID) error
}

func (qt *Quadtree) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(
		qt.UUID,
		qt.Bounds,
		qt.Level,
		qt.MaxLevel,
		qt.MaxObjects,
		qt.Total,
	)
}

func (qt *Quadtree) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.DecodeMulti(
		&qt.UUID,
		&qt.Bounds,
		&qt.Level,
		&qt.MaxLevel,
		&qt.MaxObjects,
		&qt.Total,
	)
}

type RedisBackend struct {
	redis *redis.Client
}

func NewRedisBackend(uri string) (*RedisBackend, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}
	db := redis.NewClient(opt)
	_, err = db.Ping(context.Background()).Result()
	return &RedisBackend{db}, err
}

func (r *RedisBackend) KeyOfNodes(uuid uuid.UUID) string {
	return "nodes_" + uuid.String()
}

func (r *RedisBackend) KeyOfObjects(uuid uuid.UUID) string {
	return "objects_" + uuid.String()
}

func (r *RedisBackend) GetQuadtree(uuid uuid.UUID) (*Quadtree, error) {
	value := r.redis.Get(context.Background(), uuid.String()).Val()

	var qt Quadtree
	err := msgpack.Unmarshal([]byte(value), &qt)
	return &qt, err
}

func (r *RedisBackend) SetQuadtree(qt *Quadtree) error {
	b, err := msgpack.Marshal(qt)
	if err != nil {
		return err
	}
	return r.redis.Set(context.Background(), qt.UUID.String(), b, 0).Err()
}

func (r *RedisBackend) DelQuadtree(uuid uuid.UUID) error {
	return r.redis.Del(context.Background(), uuid.String()).Err()
}

func (r *RedisBackend) GetNodes(id uuid.UUID) ([]uuid.UUID, error) {
	cmd := r.redis.LRange(context.Background(), r.KeyOfNodes(id), 0, -1)
	values, err := cmd.Val(), cmd.Err()
	if err != nil {
		return nil, err
	}

	res := make([]uuid.UUID, len(values))
	for i, v := range values {
		res[i] = uuid.MustParse(v)
	}
	return res, nil
}

func (r *RedisBackend) SetNodes(id uuid.UUID, nodes []uuid.UUID) error {
	req := make([]interface{}, len(nodes))
	for i := range nodes {
		req[i] = nodes[i]
	}
	return r.redis.LPush(context.Background(), r.KeyOfNodes(id), req...).Err()
}

func (r *RedisBackend) DetachNodes(uuid uuid.UUID) error {
	subs, err := r.GetNodes(uuid)
	if err != nil {
		return err
	}
	for _, subid := range subs {
		err := r.DetachNodes(subid)
		if err != nil {
			return err
		}
	}
	err = r.DelQuadtree(uuid)
	if err != nil {
		return err
	}

	return r.redis.Del(context.Background(), r.KeyOfNodes(uuid)).Err()
}

func (r *RedisBackend) GetObjects(uuid uuid.UUID) ([]*Object, error) {
	cmd := r.redis.SMembers(context.Background(), r.KeyOfObjects(uuid))
	values, err := cmd.Val(), cmd.Err()
	if err != nil {
		return nil, err
	}

	res := make([]*Object, len(values))
	for i, val := range values {
		var obj Object
		err := msgpack.Unmarshal([]byte(val), &obj)
		if err != nil {
			return nil, err
		}
		res[i] = &obj
	}
	return res, nil
}

func (r *RedisBackend) AddObject(uuid uuid.UUID, object *Object) error {
	b, err := msgpack.Marshal(object)
	if err != nil {
		return err
	}
	return r.redis.SAdd(context.Background(), r.KeyOfObjects(uuid), b).Err()
}

func (r *RedisBackend) DelObject(uuid uuid.UUID, object *Object) error {
	b, err := msgpack.Marshal(object)
	if err != nil {
		return err
	}
	return r.redis.SRem(context.Background(), r.KeyOfObjects(uuid), b).Err()
}

func (r *RedisBackend) ClearObjects(uuid uuid.UUID) error {
	return r.redis.Del(context.Background(), r.KeyOfObjects(uuid)).Err()
}
