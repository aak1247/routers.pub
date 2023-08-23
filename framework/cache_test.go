package framework

import (
	"testing"
)

func TestCached(t *testing.T) {
	ctx := NewRouterCtx()
	k := "s1"
	v := "v1"
	ctx.SetToLocalCache(k, v)
	v1 := Cached(ctx, k, func() interface{} {
		t.Errorf("should not call getter")
		return v
	})
	if v1 != v {
		t.Errorf("v1 should be %s, but got %s", v, v1)
	}
}

func TestCachedMiss(t *testing.T) {
	ctx := NewRouterCtx()
	k := "s1"
	v := "v1"
	v1 := Cached(ctx, k, func() interface{} {
		t.Log("cache miss")
		return v
	})
	if v1 != v {
		t.Errorf("v1 should be %s, but got %s", v, v1)
	}
}

func TestCachedMiss2(t *testing.T) {
	ctx := NewRouterCtx()
	k := "s1"
	v1 := "v1"
	v2 := "v2"
	res := Cached(ctx, k, func() interface{} {
		t.Log("cache miss")
		return []string{v1, v2}
	})
	if v1 != res.([]string)[0] {
		t.Errorf("v1 should be %s, but got %s", res.([]string)[0], v1)
	}
}

func TestCachedMissMulti(t *testing.T) {
	ctx := NewRouterCtx()
	k := "s1"
	v1 := "v1"
	v2 := "v2"
	res := Cached(ctx, k, func() (string, string) {
		t.Log("cache miss")
		return v1, v2
	})
	if v1 != res.([]interface{})[0] {
		t.Errorf("v1 should be %s, but got %s", res.([]string)[0], v1)
	}
}
