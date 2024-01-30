package localcache

type Hook interface {
	IncrementHit()
	IncrementMiss()
	IncrementDelHit()
	IncrementDelMiss()
}
