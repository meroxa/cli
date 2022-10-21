package hotfix

import (
	"fmt"
	"regexp"
	"strings"
)

func Apply(cfg map[string]interface{}, resourceType string, stream string, targetCollection string) error {
	switch resourceType {
	case "redshift", "postgres", "mysql", "sqlserver": // JDBC sink
		cfg["table.name.format"] = strings.ToLower(targetCollection)
	case "mongodb":
		cfg["targetCollectionlection"] = strings.ToLower(targetCollection)
	case "s3":
		cfg["aws_s3_prefix"] = strings.ToLower(targetCollection) + "/"
	case "snowflakedb":
		r := regexp.MustCompile("^[a-zA-Z]{1}[a-zA-Z0-9_]*$")
		matched := r.MatchString(targetCollection)
		if !matched {
			return fmt.Errorf("%q is an invalid Snowflake name - must start with "+
				"a letter and contain only letters, numbers, and underscores", targetCollection)
		}
		cfg["snowflake.topic2table.map"] =
			fmt.Sprintf("%s:%s", stream, targetCollection)
	}
	return nil
}
