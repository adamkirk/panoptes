package util

func Map[T1 any, T2 any](convert func(T1) T2, source []T1) []T2 {
    result := make([]T2, len(source))
    for i, t := range source {
        result[i] = convert(t)
    }
    return result
}
