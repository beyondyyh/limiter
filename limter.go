// @Copyright Weibo.Inc
// @Author Yehong.Yang<yehong@staff.weibo.com>
// @Date 2020/05/22 10:58:00
// @package limiter Common Rate Limit
// There are many ways, such as:
// - Leaky Bucket
// - Token BUcket

package limiter

// Limiter interface
type Limiter interface {
	Put(int64) error     // Put count token to bucket
	Get(int64) error     // Get count token from bucket
	Run(fn func()) error // Invoke fn if get one token success, otherwise do nothing and returns error
}

// Config config for Token Bucket
type Config struct {
	QPS      float64 // Query per second
	MaxCount int64   // Max token count
}
