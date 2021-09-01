package kvstore

import (
	"context"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/sqlstore"
)

func ProvideService(sqlStore *sqlstore.SQLStore) KVStore {
	return &kvStoreSQL{
		sqlStore: sqlStore,
		log:      log.New("infra.kvstore.sql"),
	}
}

// KVStore is an interface for k/v store.
type KVStore interface {
	Get(ctx context.Context, orgId int64, namespace string, key string) (string, bool, error)
	Set(ctx context.Context, orgId int64, namespace string, key string, value string) error
	Del(ctx context.Context, orgId int64, namespace string, key string) error
}

// WithNamespace returns a kvstore wrapper with fixed orgId and namespace.
func WithNamespace(kv KVStore, orgId int64, namespace string) *NamespacedKVStore {
	return &NamespacedKVStore{
		kvStore:   kv,
		orgId:     orgId,
		namespace: namespace,
	}
}

// NamespacedKVStore is a KVStore wrapper with fixed orgId and namespace.
type NamespacedKVStore struct {
	kvStore   KVStore
	orgId     int64
	namespace string
}

func (kv *NamespacedKVStore) Get(ctx context.Context, key string) (string, bool, error) {
	return kv.kvStore.Get(ctx, kv.orgId, kv.namespace, key)
}

func (kv *NamespacedKVStore) Set(ctx context.Context, key string, value string) error {
	return kv.kvStore.Set(ctx, kv.orgId, kv.namespace, key, value)
}

func (kv *NamespacedKVStore) Del(ctx context.Context, key string) error {
	return kv.kvStore.Del(ctx, kv.orgId, kv.namespace, key)
}
