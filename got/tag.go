package got

import (
	"log"
)

type tag struct {
	repo Repo
	klvm kvlm
}

func (t tag) Serialize() ([]byte, error) {
	return kvln_serialize(t.klvm)
}
func (t *tag) Deserialize(data []byte) {
	var in kvlm
	k, err := kvlm_parse(data, 0, in)
	if err != nil {
		log.Fatal(err)
	}

	t.klvm = k
}

func (t tag) get_type() []byte {
	return []byte("tag")
}
func (t tag) get_repo() Repo {
	return t.repo
}
