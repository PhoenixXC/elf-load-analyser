package module

type vmAreaStruct struct {
    vmStart uint64
    vmEnd   uint64
    vmFlags uint64
    vmPgoff uint64
}

type mmStruct struct {
    stackVm uint64
    totalVm uint64
}
