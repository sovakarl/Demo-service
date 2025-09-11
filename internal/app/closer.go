package app

// import "io"

// func (c closeFunc) Close() error {
// 	return c.Close()
// }

// type appCloser struct {
// 	funClose []func()
// }

// func newCloser() *appCloser {
// 	return &appCloser{
// 		funClose: []io.Closer{},
// 	}
// }

// func (c *appCloser) add(object io.Closer) {
// 	c.funClose = append(c.funClose, object)
// }

// func (c *appCloser) Close() error {
// 	for i := len(c.funClose) - 1; i >= 0; i-- {
// 		objectCloseFunc := c.funClose[i]
// 		objectCloseFunc.Close()
// 	}
// 	return nil
// }
