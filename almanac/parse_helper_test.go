package almanac

import "fmt"

func sscanStamp(s string, y, m, d, h, mi, sec *int) (int, error) {
	return fmt.Sscanf(s, "%d-%d-%d %d:%d:%d", y, m, d, h, mi, sec)
}
