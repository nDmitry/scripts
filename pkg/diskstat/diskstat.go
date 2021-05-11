package diskstat

import "syscall"

type Diskstat struct {
	All   float64
	Avail float64
}

func Get() (*Diskstat, error) {
	stat := syscall.Statfs_t{}
	err := syscall.Statfs("/", &stat)

	if err != nil {
		return nil, err
	}

	return &Diskstat{
		All:   float64(stat.Blocks * uint64(stat.Bsize)),
		Avail: float64(stat.Bavail * uint64(stat.Bsize)),
	}, nil
}
