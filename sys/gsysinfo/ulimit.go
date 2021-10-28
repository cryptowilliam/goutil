package gsysinfo

// WARN: 放弃在程序中做这个检测和修改吧, 整个模块都是无效的，乱七八糟
/*
func UlimitNGet() (cur, max uint64, err error) {
	var rLimit syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return 0, 0, err
	}
	return rLimit.Cur, rLimit.Max, nil
}

// note: sudo required
func UlimitNSet(cur, max uint64) error {
	var rLimit syscall.Rlimit
	rLimit.Max = max
	rLimit.Cur = cur
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err
	}
	return nil
}*/
