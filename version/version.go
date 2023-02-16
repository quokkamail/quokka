package version

import "strconv"

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (v *Version) ToString() string {
	return strconv.FormatUint(v.Major, 10) + "." + strconv.FormatUint(v.Minor, 10) + "." + strconv.FormatUint(v.Patch, 10)
}
