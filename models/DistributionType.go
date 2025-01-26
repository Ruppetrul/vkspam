package models

type DistributionType int

const (
	AnyPublic = iota + 1
)

func (d DistributionType) String() string {
	return [...]string{"Any public"}[d-1]
}

func (d DistributionType) EnumIndex() int {
	return int(d)
}
