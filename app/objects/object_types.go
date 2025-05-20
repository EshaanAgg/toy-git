package objects

type ObjectType int

const (
	BlobType ObjectType = iota
	TreeType
	CommitType
	UnknownType
)

func (ot ObjectType) String() string {
	switch ot {
	case BlobType:
		return "blob"
	case TreeType:
		return "tree"
	case CommitType:
		return "commit"
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
	case "commit":
		return CommitType
	default:
		return UnknownType
	}
}
