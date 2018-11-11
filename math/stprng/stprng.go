package stprng

type Counter struct {
	iv, is int64
	fv, fs float64
}

func NewCounter() *Counter {
	return &Counter{
		is: 1,
		fs: 1e-6,
	}
}

func (c *Counter) SeedInt(seed int64) {
	c.iv = seed
}

func (c *Counter) SeedFloat(seed float64) {
	c.fv = seed
}

func (c *Counter) Int() int {
	v := c.iv
	c.iv += c.is
	return int(v)
}

func (c *Counter) Float64() float64 {
	v := c.fv
	c.fv += c.fs
	return v
}
