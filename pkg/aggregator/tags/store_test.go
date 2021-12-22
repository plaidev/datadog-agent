// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2021-present Datadog, Inc.

package tags

import (
	"testing"

	"github.com/DataDog/datadog-agent/pkg/tagset"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	c := NewStore(true, "test")

	t1 := tagset.NewTags([]string{"1"})
	h1 := t1.Hash()
	t2 := tagset.NewTags([]string{"2"})
	h2 := t2.Hash()

	t1a := c.Insert(t1.Hash(), t1)

	require.EqualValues(t, 1, len(c.tagsByKey))
	require.EqualValues(t, 1, c.cap)
	require.EqualValues(t, 1, c.tagsByKey[h1].refs)

	t1b := c.Insert(t1.Hash(), t1)
	require.EqualValues(t, 1, len(c.tagsByKey))
	require.EqualValues(t, 1, c.cap)
	require.EqualValues(t, 2, c.tagsByKey[h1].refs)
	require.Same(t, t1a, t1b)

	t2a := c.Insert(t2.Hash(), t2)
	require.EqualValues(t, 2, len(c.tagsByKey))
	require.EqualValues(t, 2, c.cap)
	require.EqualValues(t, 2, c.tagsByKey[h1].refs)
	require.EqualValues(t, 1, c.tagsByKey[h2].refs)
	require.NotSame(t, t1a, t2a)

	t2b := c.Insert(t2.Hash(), t2)
	require.EqualValues(t, 2, len(c.tagsByKey))
	require.EqualValues(t, 2, c.cap)
	require.EqualValues(t, 2, c.tagsByKey[h1].refs)
	require.EqualValues(t, 2, c.tagsByKey[h2].refs)
	require.Same(t, t2a, t2b)

	t1a.Release()
	require.EqualValues(t, 2, len(c.tagsByKey))
	require.EqualValues(t, 2, c.cap)
	require.EqualValues(t, 1, c.tagsByKey[h1].refs)
	require.EqualValues(t, 2, c.tagsByKey[h2].refs)

	c.Shrink()
	require.EqualValues(t, 2, len(c.tagsByKey))
	require.EqualValues(t, 2, c.cap)

	t2a.Release()
	require.EqualValues(t, 2, len(c.tagsByKey))
	require.EqualValues(t, 2, c.cap)
	require.EqualValues(t, 1, c.tagsByKey[h1].refs)
	require.EqualValues(t, 1, c.tagsByKey[h2].refs)

	t1b.Release()
	require.EqualValues(t, 2, len(c.tagsByKey))
	require.EqualValues(t, 2, c.cap)
	require.EqualValues(t, 0, c.tagsByKey[h1].refs)
	require.EqualValues(t, 1, c.tagsByKey[h2].refs)

	c.Shrink()
	require.EqualValues(t, 1, len(c.tagsByKey))
	require.EqualValues(t, 2, c.cap)
	require.EqualValues(t, 1, c.tagsByKey[h2].refs)

	t2b.Release()
	require.EqualValues(t, 1, len(c.tagsByKey))
	require.EqualValues(t, 2, c.cap)
	require.EqualValues(t, 0, c.tagsByKey[h2].refs)

	c.Shrink()
	require.EqualValues(t, 0, len(c.tagsByKey))
	require.EqualValues(t, 0, c.cap)
}

func TestStoreDisabled(t *testing.T) {
	c := NewStore(false, "test")

	t1 := tagset.NewTags([]string{"1"})
	t2 := tagset.NewTags([]string{"2"})

	t1a := c.Insert(t1.Hash(), t1)
	require.EqualValues(t, 0, len(c.tagsByKey))
	require.EqualValues(t, 0, c.cap)

	t1b := c.Insert(t1.Hash(), t1)
	require.EqualValues(t, 0, len(c.tagsByKey))
	require.EqualValues(t, 0, c.cap)
	require.NotSame(t, t1a, t1b)
	require.Equal(t, t1a, t1b)

	t2a := c.Insert(t2.Hash(), t2)
	require.EqualValues(t, 0, len(c.tagsByKey))
	require.EqualValues(t, 0, c.cap)
	require.NotSame(t, t1a, t2a)
	require.NotEqual(t, t1a, t2a)

	t1a.Release()
	require.EqualValues(t, 0, len(c.tagsByKey))
	require.EqualValues(t, 0, c.cap)

	t2a.Release()
	require.EqualValues(t, 0, len(c.tagsByKey))
	require.EqualValues(t, 0, c.cap)

	c.Shrink()
	require.EqualValues(t, 0, len(c.tagsByKey))
	require.EqualValues(t, 0, c.cap)
}
