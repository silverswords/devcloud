package rest

const (
	mongoURL = "mongo://devcloud/collection?id_field=Name"
	memURL   = "mem://collection/Name"
)

type Config struct {
}

type URL struct {
	DB         string
	Collection string
	IDField    string
}

func Mongo(u URL) string {
	return "mongo://" + u.DB + "/" + u.Collection + "?id_field=" + u.IDField
}
