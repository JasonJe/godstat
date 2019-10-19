package info

type SystemConfig interface {
    GetConfig(args ...interface{}) error
}
