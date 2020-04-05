package api

import (
	"github.com/OhYee/blotter/api/pkg/comment"
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/output"
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
func Comments(context *register.HandleContext) (err error) {
	output.Debug("call friends")
	args := new(CommentsRequest)
	res := new(CommentsResponse)

	context.RequestParams(args)

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
func CommentAdd(context *register.HandleContext) (err error) {
	args := new(CommentAddRequest)
	context.RequestParams(args)

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
func AdminComments(context *register.HandleContext) (err error) {
	if !user.CheckToken(context.GetCookie("token")) {
		context.Forbidden()
		return
	}

	args := new(AdminCommentsRequest)
	res := new(AdminCommentsResponse)
	context.RequestParams(args)

	if res.Total, res.Comments, err = comment.GetAdmin(args.Offset, args.Number); err != nil {
		return
	}
	err = context.ReturnJSON(res)
	return
}
