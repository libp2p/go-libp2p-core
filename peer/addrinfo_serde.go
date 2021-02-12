package peer

import (
	"encoding/json"
	"errors"

	ma "github.com/multiformats/go-multiaddr"
)

func (pi AddrInfo) MarshalJSON() ([]byte, error) {
	out := make(map[string]interface{})
	out["ID"] = pi.ID.Pretty()
	var addrs []string
	for _, a := range pi.Addrs {
		addrs = append(addrs, a.String())
	}
	out["Addrs"] = addrs
	return json.Marshal(out)
}

func (pi *AddrInfo) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	id, ok := data["ID"]
	if !ok {
		return errors.New("no peer ID")
	}
	if idString, ok := id.(string); ok {
		pid, err := IDB58Decode(idString)
		if err != nil {
			return err
		}
		pi.ID = pid
	}
	if addrsEntry, ok := data["Addrs"]; ok {
		if addrs, ok := addrsEntry.([]interface{}); ok {
			for _, a := range addrs {
				pi.Addrs = append(pi.Addrs, ma.StringCast(a.(string)))
			}
		}
	}
	return nil
}
