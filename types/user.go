package types

type (
	UserSource = string
	UserStatus uint8
)

const (
	UserSourceInternal   UserSource = "internal"
	UserSourceThirdParty UserSource = "third_party"
)

const (
	UserStatusActive UserStatus = iota
	UserStatusDisabled
	UserStatusFrozen
)

func (s UserStatus) Int() int {
	return int(s)
}

// User struct
/*
DDL:
	create table if not exists c_user
	(
		id            varchar(50)            not null primary key,
		username      varchar(20)            not null,
		password      varchar(200)           not null,
		created_at    datetime               null,
		updated_at    datetime               null,
		administrator tinyint(1)  default 0  null,
		source        varchar(20) default '' null,
		status        tinyint     default 0  null,
		extra         json                   null,
		unique index c_user_uk_username (username)
	);
*/
type User struct {
	Id            string         `json:"id"`
	Username      string         `json:"username"`
	Password      string         `json:"password,omitempty"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at"`
	Administrator bool           `json:"administrator"`
	Source        UserSource     `json:"source"`
	Status        UserStatus     `json:"status"`
	Extra         map[string]any `json:"extra"`
}
