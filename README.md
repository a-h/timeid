# timeid

Creates numeric IDs using a mixture of time and random values which make it unlikely that they have previously been used.

In case of legacy systems which assume an `int64` database ID, we can use the current number of seconds since the minimum date / time as a `uint32` for the first half of the underlying `int64` bytes and a random `uint32` number between 0-4294967295 for the remainder.

The numbers are positive (greater than zero) and monotonic (keep getting bigger).

This provides around 130 years worth of IDs.

For any given second the first insert is unique, then afterwards there's a 1/4294967295, 1/4294967294, 1/4294967293 etc. chance that the ID is not unique.

This was designed for relatively low volume insertion.

## Usage

```
random := rand.New(rand.NewSource(time.Now().Unix()))
id, err := NewTimeID(random, time.Now)
```