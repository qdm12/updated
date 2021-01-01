package funcs

import "os"

type OSOpenFile func(name string, flag int, perm os.FileMode) (*os.File, error)
