package cicada

import (
	"fmt"
	"log"
	_ "os"

	"github.com/urfave/cli"
	"github.com/garyburd/redigo/redis"
)

func ClearRedis6310Action(ctx *cli.Context) error {
	clusters := []string{
		"10.105.203.106:16301",
		"10.105.201.120:16306",
		"10.105.201.120:16303",
		"10.105.201.232:16307",
		"10.105.201.232:16314",
		"10.105.201.107:16317",
		"10.105.201.107:16310",
		"10.105.201.107:16321",
		"10.105.247.94:16330",
		"10.105.219.21:16331",
		"10.105.200.80:16331",
	}
	rkeys := []string{
		"picktop:*",
		"stat:*",
		"pickchance:*",
		"pubtimings:*",
		"statadv:*",
		"channelSum:*",
		"puv:*",
		"endTimes:*",
	}

	for _, v := range clusters {
		rc, err := redis.Dial("tcp", v)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		for _, vv := range rkeys {
			keys, err := redis.Strings(rc.Do("keys", vv))
			if err != nil {
				return cli.NewExitError(err.Error(), 2)
			}

			if len(keys) == 0 { continue }

			// batch delete
			for i := 0; i < (len(keys)/1000 + 1); i++ {
				start := i * 1000
				end := start + 1000
				if end >= len(keys) {
					end = len(keys)
				}
				partKeys := keys[start:end]

				args := make([]interface{}, len(partKeys))
				for ki, kv := range partKeys {
					args[ki] = kv
				}
				if _, err := rc.Do("del", args...); err != nil {
					return cli.NewExitError(err.Error(), 3)
				} else {
					log.Printf("%s: del %v", v, partKeys)
				}
			}
			
		}

		rc.Close()
	}

	return nil
}
