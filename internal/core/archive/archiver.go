package archive

type Archiver interface {
  Archive ([]EntryFile) ([]byte, error)
  Unarchive ([]byte) ([]EntryFile, error)
  Extension () string
}
