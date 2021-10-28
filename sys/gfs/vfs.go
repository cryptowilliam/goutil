package gfs

/*
https://github.com/minio/minfs
https://github.com/bazil/fuse
https://github.com/hanwen/go-fuse
https://github.com/kahing/goofys
https://github.com/fntlnz/gridfsmount
https://github.com/apparentlymart/go-fsutil
*/

// https://github.com/maidsafe-archive/drive
// https://github.com/Microsoft/Windows-Driver-Frameworks
/*
func TVFS() {
	// Create a vfs accessing the filesystem of the underlying OS
	var osfs vfs.Filesystem = vfs.OS()
	osfs.Mkdir("/tmp", 0777)

	// Make the filesystem read-only:
	osfs = vfs.ReadOnly(osfs) // Simply wrap filesystems to change its behaviour

	// os.O_CREATE will fail and return vfs.ErrReadOnly
	// os.O_RDWR is supported but Write(..) on the file is disabled
	f, _ := osfs.OpenFile("/tmp/example.txt", os.O_RDWR, 0)

	// Return vfs.ErrReadOnly
	_, err := f.Write([]byte("Write on readonly fs?"))
	if err != nil {
		fmt.Println(err)
	}

	// Create a fully writable filesystem in memory
	mfs := memfs.Create()
	mfs.Mkdir("/root", 0777)

	// Create a vfs supporting mounts
	// The root fs is accessing the filesystem of the underlying OS
	fs := mountfs.Create(osfs)

	// Mount a memfs inside /memfs
	// /memfs may not exist
	err = fs.Mount(mfs, "/memfs")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Mount OK")
	}

	// This will create /testdir inside the memfs
	fs.Mkdir("/memfs/testdir", 0777)

	// This would create /tmp/testdir inside your OS fs
	// But the rootfs `osfs` is read-only
	fs.Mkdir("/tmp/testdir", 0777)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR2)

	s := <-c
	fmt.Println("Got signal:", s)
}*/
