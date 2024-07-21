package repo

import "errors"

var ErrNoMoney = errors.New("not enough money")

var ErrUniqConstrain = errors.New("object already exists")

var ErrOrderAlreadyRegisteredByUser = errors.New("order already registered by user")
var ErrOrderAlreadyRegisteredByOther = errors.New("order already registered by other user")
