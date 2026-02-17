package child

import "embed-parent-directory/data"

func Hello() string {
	return data.Hello
}
