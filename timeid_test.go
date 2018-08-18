package timeid

import (
	"math/rand"
	"testing"
	"time"
)

func TestMinimumIDValue(t *testing.T) {
	tid := TimeID{
		Random: 0,
		Time:   0,
	}
	id, err := tid.Int64()
	if err != nil {
		t.Errorf("unexpected error converting to Int64: %v", err)
	}
	if id != 0 {
		t.Errorf("expected zero value, got %v", id)
	}
}

func TestIDConversion(t *testing.T) {
	tests := []struct {
		time   uint32
		random uint32
	}{
		{
			time:   0,
			random: 0,
		},
		{
			time:   1,
			random: 2,
		},
		{
			time:   0,
			random: 65536,
		},
		{
			time:   65536,
			random: 65536,
		},
		{
			time:   4294967295,
			random: 4294967295,
		},
	}

	for index, test := range tests {
		input := TimeID{
			Time:   test.time,
			Random: test.random,
		}
		id, err := input.Int64()
		if err != nil {
			t.Errorf("%d: Unexpected error creating an int64 from the ID: %v", index, err)
		}
		output, err := FromInt64(id)
		if err != nil {
			t.Errorf("%d: Unexpected error creating an TimeID from id %v: %v", index, id, err)
		}
		if input.Random != output.Random {
			t.Errorf("%d: expected random %d from id %v, got %d", index, input.Random, id, output.Random)
		}
		if input.Time != output.Time {
			t.Errorf("%d: expected time %d from id %v, got %d", index, input.Time, id, output.Random)
		}
	}
}

func TestThatIDsCanBeCreatedFor100Years(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	now := func() time.Time { return time.Date(2118, time.January, 1, 0, 0, 0, 0, time.UTC) }
	_, err := NewTimeID(r, now)
	if err != nil {
		t.Fatalf("Unable to create an ID in year 2118 with error: %v", err)
	}
}

func TestTimeRetrieval(t *testing.T) {
	times := []time.Time{
		time.Date(2019, time.February, 20, 13, 12, 11, 0, time.UTC),
		time.Date(2036, time.June, 7, 6, 32, 19, 0, time.UTC),
	}

	random := rand.New(rand.NewSource(time.Now().Unix()))
	for index, testTime := range times {
		tid, err := NewTimeID(random, func() time.Time { return testTime })
		if err != nil {
			t.Errorf("%d: unexpected error creating time ID: %d", index, err)
		}
		if !tid.GetTime().Equal(testTime) {
			t.Errorf("%d: expected time %v, got %v", index, testTime, tid.GetTime())
		}
	}
}

func TestTooFarIntoTheFuture(t *testing.T) {
	random := rand.New(rand.NewSource(time.Now().Unix()))
	now := func() time.Time { return time.Date(2156, time.January, 1, 1, 0, 0, 0, time.UTC) }
	tid, err := NewTimeID(random, now)
	if err == nil {
		// If we don't get an error because you've changed the code to support dates that are further into
		// the future, great.
		t.Errorf("expected error creating time so far into the future, but didn't get one, got ID: %v", tid)
	}
}
