package repository

import (
	"context"
	"errors"
	"friend-management/internal/core/domain"
	"friend-management/internal/log"
	"friend-management/internal/repository/util"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection string = "user"

type UserRepo struct {
	db     *mongo.Database
	client *mongo.Client
	logger log.Logger
}

type userUpdate struct {
	Email      string    `bson:"email"`
	Name       string    `bson:"name"`
	UpdatedAt  time.Time `bson:"updated_at"`
	VersionRev string    `bson:"version_rev"`
	VersionSeq int       `bson:"version_seq"`
}

func NewUserRepo(db *mongo.Database, client *mongo.Client, logger log.Logger) *UserRepo {
	return &UserRepo{db, client, logger}
}

// ConnectionCheck pings the DB to check if the connection is working
func (r *UserRepo) ConnectionCheck(c context.Context) bool {
	err := r.client.Ping(c, nil)
	if err != nil {
		r.logger.WithCtx(c).Error("function", "ConnectionCheck", "msg", "Ping to MongoDB failed", "error", err)
		return false
	}
	return true
}

func (r *UserRepo) CreateUser(c context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	var (
		collection = r.db.Collection(userCollection)
		setAtUTC   = time.Now().UTC()
	)

	user := domain.User{
		Email:      req.Email,
		Name:       req.Name,
		Friends:    make([]domain.Friend, 0),
		VersionSeq: 0,
		VersionRev: util.GenerateVersionRev(0),
		CreatedAt:  setAtUTC,
		UpdatedAt:  setAtUTC,
	}

	for {
		user.ID = uuid.New().String()
		_, err := collection.InsertOne(c, user)

		if !mongo.IsDuplicateKeyError(err) {
			return &user, err
		}
	}
}

func (r *UserRepo) GetUsers(c context.Context, query domain.UserQuery) (*domain.MetaUsers, error) {
	var (
		results    domain.MetaUsers
		collection = r.db.Collection(userCollection)
		skip       = int64(0)
		sort       = util.GenerateSortFilter(query)
	)

	if query.Page > 1 {
		skip = (query.Page - 1) * query.PerPage
	}

	cur, err := collection.Find(
		c,
		bson.M{},
		options.Find().SetSkip(skip).SetLimit(query.PerPage).SetSort(sort),
	)

	if err != nil {
		r.logger.WithCtx(c).Error("function", "GetUsers", "error", err)
		return nil, domain.ErrUnknown
	}

	defer cur.Close(c)
	users := make([]domain.User, 0)
	if err = cur.All(c, &users); err != nil {
		r.logger.WithCtx(c).Error("function", "GetUsers", "error", err)
		return nil, domain.ErrUnknown
	}

	results.Meta.PerPage = query.PerPage
	results.Meta.Page = query.Page
	results.Users = users
	return &results, nil
}

func (r *UserRepo) GetUser(c context.Context, id string) (*domain.User, error) {
	var (
		collection = r.db.Collection(userCollection)
		user       *domain.User
	)

	err := collection.FindOne(c, bson.M{
		"_id": id,
	}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		r.logger.WithCtx(c).Error("function", "GetUser", "error", err)
		return nil, domain.ErrUnknown
	}

	return user, nil
}

func (r *UserRepo) UpdateUser(c context.Context, req domain.UpdateUserRequest) (*domain.User, error) {
	var (
		collection    = r.db.Collection(userCollection)
		setAtUTC      = time.Now().UTC()
		user          domain.User
		intVer, _     = strconv.Atoi(strings.Split(req.VersionRev, "-")[0])
		newVersionSeq = intVer + 1
		newVersionRev = util.GenerateVersionRev(newVersionSeq)
	)

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err := collection.FindOneAndUpdate(
		c,
		bson.M{"_id": req.ID, "version_rev": req.VersionRev},
		bson.M{
			"$set": userUpdate{
				Name:       req.Name,
				Email:      req.Email,
				UpdatedAt:  setAtUTC,
				VersionRev: newVersionRev,
				VersionSeq: newVersionSeq,
			},
		},
		&opt,
	).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = collection.FindOne(c, bson.M{"_id": req.ID}).Decode(&domain.User{})

			if err == nil {
				return nil, domain.ErrVersionConflict
			}
			return nil, domain.ErrNotFound
		}
		r.logger.WithCtx(c).Error("function", "UpdateUser", "error", err)
		return nil, domain.ErrUnknown
	}

	return &user, nil
}

func (r *UserRepo) AddFriend(c context.Context, req domain.AddFriendRequest) error {
	var (
		collection = r.db.Collection(userCollection)
		setAtUTC   = time.Now().UTC()
	)

	findUser := func(c context.Context, id string) (*domain.User, error) {
		user := domain.User{}

		err := collection.FindOne(c, bson.M{"_id": id}).Decode(&user)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, domain.ErrNotFound
			}
			r.logger.WithCtx(c).Error("function", "UpdateUser", "error", err)
			return nil, domain.ErrUnknown
		}

		return &user, nil
	}

	addFriendFunc := func(c context.Context, id string, friend domain.Friend, versionRev string) error {
		var (
			intVer, _     = strconv.Atoi(strings.Split(versionRev, "-")[0])
			newVersionSeq = intVer + 1
			newVersionRev = util.GenerateVersionRev(newVersionSeq)
		)

		err := collection.FindOneAndUpdate(
			c,
			bson.M{"_id": id, "version_rev": versionRev},
			bson.M{
				"$set": userUpdate{
					UpdatedAt:  setAtUTC,
					VersionRev: newVersionRev,
					VersionSeq: newVersionSeq,
				},
				"$push": bson.M{
					"friends": friend,
				},
			},
		).Decode(&domain.User{})
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				err = collection.FindOne(c, bson.M{"_id": id}).Decode(&domain.User{})

				if err == nil {
					return domain.ErrVersionConflict
				}
				return domain.ErrNotFound
			}
			r.logger.WithCtx(c).Error("function", "AddFriend", "error", err)
			return domain.ErrUnknown
		}

		return nil
	}

	currentUser, err := findUser(c, req.ID)
	if err != nil {
		return err
	}

	for _, friend := range currentUser.Friends {
		if friend.ID == req.FriendID {
			return domain.ErrFriendConflict
		}
	}

	friendUser, err := findUser(c, req.FriendID)
	if err != nil {
		return err
	}

	for _, friend := range friendUser.Friends {
		if friend.ID == req.ID {
			return domain.ErrFriendConflict
		}
	}

	err = addFriendFunc(
		c,
		req.ID,
		domain.Friend{
			ID:    friendUser.ID,
			Name:  friendUser.Name,
			Email: friendUser.Email,
		},
		req.VersionRev,
	)

	if err != nil {
		return err
	}

	err = addFriendFunc(
		c,
		req.FriendID,
		domain.Friend{
			ID:    currentUser.ID,
			Name:  currentUser.Name,
			Email: currentUser.Email,
		},
		req.FriendVersionRev,
	)

	return err
}
