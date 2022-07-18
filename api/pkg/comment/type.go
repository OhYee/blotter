package comment

import (
	"github.com/OhYee/blotter/api/pkg/avatar"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/goutils/time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TypeDB comment type in database
type TypeDB struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Avatar   string             `json:"avatar" bson:"avatar"`
	Time     int64              `json:"time" bson:"time"`
	Raw      string             `json:"raw" bson:"raw"`
	Content  string             `json:"content" bson:"content"`
	Reply    primitive.ObjectID `json:"reply" bson:"reply"`
	URL      string             `json:"url" bson:"url"`
	Ad       bool               `json:"ad" bson:"ad"`
	Show     bool               `json:"show" bson:"show"`
	Recv     bool               `json:"recv" bson:"recv"`
	IP       string             `json:"ip" bson:"ip"`
	Position string             `json:"position" bson:"position"`
}

// ToComment transfer CommentDB to comment
func (cm *TypeDB) ToComment() *Type {
	return &Type{
		ID:       cm.ID.Hex(),
		Email:    cm.Email,
		Avatar:   cm.Avatar,
		Time:     time.ToString(cm.Time),
		Content:  cm.Content,
		Reply:    cm.Reply,
		Children: []*Type{},
		Ad:       cm.Ad,
		Show:     cm.Show,
		Recv:     cm.Recv,
		Position: cm.Position,
	}
}

// Type type
type Type struct {
	ID       string             `json:"id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Avatar   string             `json:"avatar" bson:"avatar"`
	Time     string             `json:"time" bson:"time"`
	Content  string             `json:"content" bson:"content"`
	Reply    primitive.ObjectID `json:"reply" bson:"reply"`
	Children []*Type            `json:"children" bson:"children"`
	Ad       bool               `json:"ad" bson:"ad"`
	Show     bool               `json:"show" bson:"show"`
	Recv     bool               `json:"recv" bson:"recv"`
	Position string             `json:"position" bson:"position"`
}

func (c *Type) UpdateAvatar() (err error) {
	oldAvatar := c.Avatar
	c.Avatar = avatar.Get(c.Email)
	if c.Avatar != "" && c.Avatar != oldAvatar && c.Avatar != avatar.DefaultAvatar {
		objectID, err := primitive.ObjectIDFromHex(c.ID)
		if err != nil {
			return err
		}
		result, err := mongo.Update(
			DatabaseName, CollectionName,
			bson.M{"_id": objectID},
			bson.M{
				"$set": bson.M{
					"avatar": c.Avatar,
				},
			}, nil,
		)
		output.DebugOutput.Printf("Update avatar %s => %s\n%+v %+v\n%+v", oldAvatar, c.Avatar, result, err, c)

	}
	return
}

// AdminDB comment database Type
type AdminDB struct {
	TypeDB       `bson:",inline"`
	ReplyComment TypeDB `json:"reply_comment" bson:"reply_comment"`
	Title        string `json:"title" bson:"title"`
}

// ToAdmin transfer AdminDB to Admin
func (comment AdminDB) ToAdmin() *Admin {
	return &Admin{
		Type:         *comment.TypeDB.ToComment(),
		ReplyComment: *comment.ReplyComment.ToComment(),
		Title:        comment.Title,
		URL:          comment.TypeDB.URL,
		IP:           comment.IP,
	}
}

// Admin comment Type
type Admin struct {
	Type         `bson:",inline"`
	ReplyComment Type   `json:"reply_comment" bson:"reply_comment"`
	URL          string `json:"url" bson:"url"`
	Title        string `json:"title" bson:"title"`
	IP           string `json:"ip" bson:"ip"`
}

// Info base info of comment
type Info struct {
	Recv  bool   `json:"recv" bson:"recv"`
	Email string `json:"email" bson:"email"`
}
