package objects

type ObjectType int

const (
	BlobType ObjectType = iota
	TreeType
	UnknownType
)

func (ot ObjectType) String() string {
	switch ot {
	case BlobType:
		return "blob"
	case TreeType:
		return "tree"
	default:
		return "unknown"
	}
}

func GetObjectTypeFromString(typeStr string) ObjectType {
	switch typeStr {
	case "blob":
		return BlobType
	case "tree":
		return TreeType
	default:
		return UnknownType
	}
}
