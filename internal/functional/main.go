// These functions come from:
//		https://bitfieldconsulting.com/golang/functional

package functional

type keepFunc[E any] func(E) bool

func Filter[E any](s []E, f keepFunc[E]) []E {
	result := []E{}
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
