package fly

import "testing"

func TestFly(t *testing.T) {
	GET("/get.json", func(c *Context) error {
		c.Json(c.Request().URL)
		return nil
	})
	ListenAndServe(":7654")
}
