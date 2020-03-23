package system

// Level ...
type Level struct {
	User   string
	Editor string
	Admin  string
}

// GetLevel ...
func GetLevel() Level {
	level := Level{}
	level.User = "user"
	level.Editor = "editor"
	level.Admin = "admin"
	return level
}

// ValidLevel ..
func ValidLevel(currentLevel string, routeLevel string) bool {
	if routeLevel == GetLevel().User {
		return true
	} else if routeLevel == GetLevel().Editor {
		if currentLevel == GetLevel().Editor || currentLevel == GetLevel().Admin {
			return true
		}
	} else if routeLevel == GetLevel().Admin && currentLevel == GetLevel().Admin {
		return true
	}
	return false
}

// Status ...
type Status struct {
	Active string
	Lock   string
	Block  string

	// Post
	Private string
	Draft   string
	Posting string
}

//GetStatus returns mode
func GetStatus() Status {
	status := Status{}

	status.Active = "active"
	status.Lock = "lock"
	status.Block = "block"

	status.Private = "private"
	status.Draft = "draft"
	status.Posting = "posting"
	return status
}

// TypeHistory ...
type TypeHistory struct {
	AddPost    string
	AddComment string
	Lock       string
	Block      string

	// Post
	Public  string
	Private string
	Draft   string
	Posting string
}

//GetTypeHistory returns mode
func GetTypeHistory() TypeHistory {
	typeHistory := TypeHistory{}
	typeHistory.AddPost = "add-post"
	typeHistory.AddComment = "add-comment"
	return typeHistory
}
