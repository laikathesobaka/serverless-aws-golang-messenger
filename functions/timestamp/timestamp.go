package timestamp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

const DATE_FMT = "02-01-2002"

func TimeToDB(t time.Time) *string {
	fmt.Println("add pkg")

	tn := t.Unix()
	return aws.String(strconv.FormatInt(tn, 10))
}

func DBtoTime(s *string) time.Time {
	n, _ := strconv.ParseInt(*s, 10, 64)
	return time.Unix(n, 0)
}
