package got

import (
	"errors"
	"fmt"
	"log"

	mapset "github.com/deckarep/golang-set/v2"
)

type commit struct {
	repo Repo
	klvm kvlm
}

func (c commit) Serialize() ([]byte, error) {
	return kvln_serialize(c.klvm)
}
func (c *commit) Deserialize(data []byte) {
	var in kvlm
	k, err := kvlm_parse(data, 0, in)
	if err != nil {
		log.Fatal(err)
	}

	c.klvm = k
}

func (c commit) get_type() []byte {
	return []byte("commit")
}
func (c commit) get_repo() Repo {
	return c.repo
}
