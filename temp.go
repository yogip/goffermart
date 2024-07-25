package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Создаем регулярное выражение
	// "UPDATE balance SET current = current - $1, withdrawn = withdrawn + $1 WHERE user_id=$2 RETURNING current"
	// "UPDATE balance SET current = current - \$1, withdrawn = withdrawn + \$1 WHERE user_id=\$2 RETURNING current"
	re := regexp.MustCompile(
		`UPDATE balance SET current = current - \$1, withdrawn = withdrawn \+ \$1 WHERE user_id=\$2 RETURNING current`,
	)

	// Строка, которую мы хотим проверить
	str := "UPDATE balance SET current = current - $1, withdrawn = withdrawn + $1 WHERE user_id=$2 RETURNING current"

	// Проверяем, соответствует ли строка регулярному выражению
	match := re.MatchString(str)

	// Выводим результатs
	fmt.Println(match)
}
