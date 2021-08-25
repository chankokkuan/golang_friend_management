package domain

type CreateUserRequest struct {
	Email string
	Name  string
}

type UpdateUserRequest struct {
	ID         string
	Email      string
	Name       string
	VersionRev string
}

type AddFriendRequest struct {
	ID               string
	FriendID         string
	VersionRev       string
	FriendVersionRev string
}
