package user

import (
	"context"
	"encoding/json"

	"github.com/dozto/dozto-aster-auth/pkg"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDoc struct {
	ID             primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	Email          string             `bson:"email" json:"email" validate:"email"`
	Phone          string             `bson:"phone" json:"phone"` // TODO: validate:"phone"
	HashedPassword string             `bson:"hashed_password" json:"-" validate:"required"`
	Roles          []string           `bson:"roles" json:"roles"`
	Balance        int64              `bson:"balance" json:"balance"`
	IsSuper        bool               `bson:"is_super" json:"-"`
	Meta           bson.M             `bson:"meta" json:"meta"`
	CreatedAt      primitive.DateTime `bson:"created_at" json:"createdAt" validate:"required"`
	UpdatedAt      primitive.DateTime `bson:"updated_at" json:"updatedAt"`
}

type UserModel struct {
	users *mongo.Collection
}

func NewUserModel(db *mongo.Database, collName string) *UserModel {
	return &UserModel{
		users: db.Collection(collName),
	}
}

func (m *UserModel) CreateUser(doc *UserDoc) (string, error) {
	err := pkg.Valtor.Validate(doc)
	if err != nil {
		return "", err
	}

	doc.UpdatedAt = primitive.NewDateTimeFromTime(doc.UpdatedAt.Time())

	res, err := m.users.InsertOne(context.Background(), doc)
	if err != nil {
		return "", err
	}
	id := res.InsertedID
	return id.(primitive.ObjectID).Hex(), nil
}

func (m *UserModel) GetUserByID(id string) (*UserDoc, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	var user UserDoc
	err = m.users.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetByEmail(email string) (*UserDoc, error) {
	filter := bson.M{"email": email}
	var user UserDoc
	err := m.users.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetByPhone(phone string) (*UserDoc, error) {
	filter := bson.M{"phone": phone}
	var user UserDoc
	err := m.users.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Query(filter interface{}) (*[]UserDoc, error) {
	cursor, err := m.users.Find(context.Background(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		panic(err)
	}

	log.Debug().Any("cursoe", cursor)

	var userDocs []UserDoc
	if err = cursor.All(context.TODO(), &userDocs); err != nil {
		panic(err)
	}

	for _, userDoc := range userDocs {
		cursor.Decode(&userDoc)
		output, err := json.MarshalIndent(userDoc, "", "    ")
		if err != nil {
			panic(err)
		}
		log.Debug().Msgf("userDoc: %s", output)
	}

	return &userDocs, nil
}

func (m *UserModel) ReplaceUser(updateDoc *UserDoc) (*mongo.UpdateResult, error) {

	result, err := m.users.ReplaceOne(context.Background(), updateDoc.ID, updateDoc)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *UserModel) DeleteByID(id string) (*mongo.DeleteResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	result, err := m.users.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
