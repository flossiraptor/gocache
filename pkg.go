package cache

import (
	"context"
	"time"
)

type MatchPredicate[V comparable] func(entry *V) bool;

type Entry[V comparable] interface {
	Age() time.Duration;
	Equal(e2 Entry[V]) bool;
	Expires(expires time.Time);
	Get() *V;
	IsValid() bool;
	Matches(predicate MatchPredicate[V]) bool;
	Ttl(ttl time.Duration);
}

type Cache[K comparable, V comparable] interface {
	Clear();
	Expires(key K, expires time.Time);
	Filter(filter MatchPredicate[V]) []*V;
	Gc(ctx context.Context);
	Get(key K) *V
	GetEntry(key K) Entry[V];
	Has(key K) bool;
	Keys() []K;
	Set(key K, val V);
	Size() int;
	Ttl(key K, ttl time.Duration);
	Unset(key K);
	Values() []*V;
}
