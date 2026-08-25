package cache

import (
	"context"
	"time"
)

func NewCacheEntry[V comparable](value *V) Entry[V] {
	return &cacheEntry[V] {
		value: value,
		created: time.Now(),
	};
}
type cacheEntry[V comparable] struct {
	value *V;
	created time.Time;
	expires time.Time;
}
func (e *cacheEntry[V]) Age() time.Duration {
	return time.Since(e.created);
}
func (e *cacheEntry[V]) Get() *V {
	return e.value;
}
func (e1 *cacheEntry[V]) Equal(e2 Entry[V]) bool {
	return e1.value == e2.Get();
}
func (e *cacheEntry[V]) Expires(expires time.Time)  {
	e.expires = expires;
}
func (e *cacheEntry[V]) IsValid() bool {
	return time.Now().After(e.expires);
}
func (e *cacheEntry[V]) Matches(predicate MatchPredicate[V]) bool {
	return predicate(e.value);
}
func (e *cacheEntry[V]) Ttl(ttl time.Duration)  {
	e.expires = e.created.Add(ttl);
}

func NewCache[K comparable, V comparable](ttl time.Duration, ctx context.Context) Cache[K, V] {
	if (ctx == nil) {
		ctx = context.Background();
	}
	cache := &cache[K, V] {
		entries: make(map[K]Entry[V]),
		ttl: ttl,
	};
	defer cache.Gc(ctx);
	return cache;
}
type cache[K comparable, V comparable] struct {
	ttl time.Duration;
	entries map[K]Entry[V];
}
func (c *cache[K, V]) Clear() {
	clear(c.entries);
}
func (c *cache[K, V]) Expires(key K, expires time.Time)  {
	e, ok := c.entries[key];
	if (ok) {
		e.Expires(expires)
	}
}
func (c *cache[K, V]) Filter(filter MatchPredicate[V]) []*V {
	matches := make([]*V, 0, len(c.entries))
	for _, entry := range c.entries {
		if (filter(entry.Get())) {
			matches = append(matches, entry.Get())
		}
	}
	return matches;
}
func (c *cache[K, V]) Gc(ctx context.Context) {
	for {
		select {
			case <-ctx.Done():
				return;

			default:
				c.gc();
		}
	}
}
func (c *cache[K, V]) Get(key K) *V {
	c.checkEviction(key);
	entry, _ := c.entries[key];
	return entry.Get();
}
func (c *cache[K, V]) GetEntry(key K) Entry[V] {
	c.checkEviction(key);
	entry, _ := c.entries[key];
	return entry;
}
func (c *cache[K, V]) Has(key K) bool {
	c.checkEviction(key);
	_, ok := c.entries[key];
	return ok;
}
func (c *cache[K, V]) Keys() []K {
	c.gc();
	keys := make([]K, len(c.entries))
	i := 0
	for k := range c.entries {
		keys[i] = k
		i++
	}
	return keys;
}
func (c *cache[K, V]) Set(key K, val V) {
	c.entries[key] = NewCacheEntry(&val);
	c.entries[key].Ttl(c.ttl);
}
func (c *cache[K, V]) Size() int {
	c.gc();
	return len(c.entries);
}
func (c *cache[K, V]) Ttl(key K, ttl time.Duration)  {
	e, ok := c.entries[key];
	if (ok) {
		e.Ttl(time.Duration(ttl));
	}
}
func (c *cache[K, V]) Unset(key K) {
	delete(c.entries, key);
}
func (c *cache[K, V]) Values() []*V {
	c.gc();
	values := make([]*V, len(c.entries))
	i := 0
	for _, entry := range c.entries {
		values[i] = entry.Get();
		i++
	}
	return values;
}

func (c *cache[K, V]) checkEviction(key K) {
	entry, ok := c.entries[key];
	if (ok && !entry.IsValid()) {
		delete(c.entries, key);
	}
}
func (c *cache[K, V]) gc() {
	for key := range(c.entries) {
		c.checkEviction(key);
	}
}
