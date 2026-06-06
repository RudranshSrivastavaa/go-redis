package command

type CommandType string

const (
	SET     CommandType = "SET"
	GET     CommandType = "GET"
	DEL     CommandType = "DEL"
	EXISTS  CommandType = "EXISTS"

	LPUSH   CommandType = "LPUSH"
	RPUSH   CommandType = "RPUSH"

	LPOP    CommandType = "LPOP"
	RPOP    CommandType = "RPOP"

	LLEN    CommandType = "LLEN"

	HSET    CommandType = "HSET"
	HGET    CommandType = "HGET"
	HDEL    CommandType = "HDEL"
	HEXISTS CommandType = "HEXISTS"
	HLEN    CommandType = "HLEN"

	EXPIRE  CommandType = "EXPIRE"
	TTL     CommandType = "TTL"
	PERSIST CommandType = "PERSIST"

	SUBSCRIBE CommandType = "SUBSCRIBE"
	UNSUBSCRIBE CommandType = "UNSUBSCRIBE"
	PUBLISH CommandType = "PUBLISH"
)

type Command struct {
	Type CommandType
	Args []string
}