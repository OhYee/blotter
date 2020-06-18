package api

import (
	"regexp"

	"github.com/OhYee/blotter/api/pkg/comment"
	"github.com/OhYee/blotter/register"
)

// CommentsRequest request of comments api
type CommentsRequest struct {
	URL string `json:"url"`
}

// CommentsResponse response of comments api
type CommentsResponse struct {
	Total    int64           `json:"total"`
	Comments []*comment.Type `json:"comments"`
}

// Comments get comments of url, return comments and total comment number
func Comments(context register.HandleContext) (err error) {
	args := new(CommentsRequest)
	res := new(CommentsResponse)

	context.RequestArgs(args)

	var comments []comment.TypeDB
	if res.Total, comments, err = comment.Get(args.URL); err != nil {
		return
	}

	res.Comments = comment.MakeRelation(comments)

	err = context.ReturnJSON(res)

	return
}

// CommentAddRequest request of CommentAdd api
type CommentAddRequest struct {
	URL   string `json:"url"`
	Reply string `json:"reply"`
	Email string `json:"email"`
	Recv  bool   `json:"recv"`
	Raw   string `json:"raw"`
}

// CommentAdd add comment api
func CommentAdd(context register.HandleContext) (err error) {
	args := new(CommentAddRequest)
	context.RequestArgs(args)

	if m, e := regexp.MatchString(
		"^([A-Za-z0-9_\\-\\.\u4e00-\u9fa5])+\\@([A-Za-z0-9_\\-\\.])+\\.([A-Za-z]{2,8})$",
		args.Email,
	); e != nil || args.URL == "" || m == false || args.Raw == "" {
		context.Forbidden()
		return
	}

	if err = comment.Add(args.URL, args.Reply, args.Email, args.Recv, args.Raw); err != nil {
		return
	}

	context.ReturnJSON(SimpleResponse{Success: true, Title: "评论发布成功"})

	return
}

// AdminCommentsRequest request for AdminComments api
type AdminCommentsRequest struct {
	Number int64 `json:"number"`
	Offset int64 `json:"Offset"`
}

// AdminCommentsResponse response for AdminComments api
type AdminCommentsResponse struct {
	Total    int64           `json:"total"`
	Comments []comment.Admin `json:"comments"`
}

// AdminComments api for admin comments page
func AdminComments(context register.HandleContext) (err error) {
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

	args := new(AdminCommentsRequest)
	res := new(AdminCommentsResponse)
	context.RequestArgs(args)

	if res.Total, res.Comments, err = comment.GetAdmin(args.Offset, args.Number); err != nil {
		return
	}
	err = context.ReturnJSON(res)
	return
}

// AdminCommentSetRequest request for AdminCommentSet api
type AdminCommentSetRequest struct {
	ID   string `json:"id" bson:"id"`
	Recv bool   `json:"recv" bson:"recv"`
	Show bool   `json:"show" bson:"show"`
	Ad   bool   `json:"ad" bson:"ad"`
}

// AdminCommentSetResponse response for AdminCommentSet api
type AdminCommentSetResponse SimpleResponse

// AdminCommentSet api for updating admin comments page
func AdminCommentSet(context register.HandleContext) (err error) {
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

	args := new(AdminCommentSetRequest)
	res := new(AdminCommentSetResponse)
	context.RequestArgs(args)

	if err = comment.Set(args.ID, args.Ad, args.Show, args.Recv); err != nil {
		return
	}

	res.Success = true
	res.Title = "修改成功"

	err = context.ReturnJSON(res)
	return
}
