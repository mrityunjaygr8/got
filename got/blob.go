package got

type blob struct {
	blobdata []byte
	repo     Repo
}

func (b blob) Serialize() ([]byte, error) {
	return b.blobdata, nil
}
func (b blob) Deserialize(data []byte) {
	b.blobdata = data
}

func (b blob) get_type() []byte {
	return []byte("blob")
}
func (b blob) get_repo() Repo {
	return b.repo
}
