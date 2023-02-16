package version

import "fmt"

type Version struct {
	Major uint
	Minor uint
	Patch uint
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
