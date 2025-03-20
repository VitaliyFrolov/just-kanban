package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"just-kanban/internal/config"
	"just-kanban/internal/contextkeys"
	"just-kanban/internal/services"
	"just-kanban/pkg/sqlddl"
	"just-kanban/pkg/validation"
)

// BoardMemberHandler handles http requests for working with methods of services.BoardMemberService
type BoardMemberHandler struct {
	*services.BoardMemberService
	*validation.Validate
}

var (
	badBoardIdErr       = errors.New("bad board identifier")
	boardNotExistErr    = errors.New("requested board does not exist")
	notAllowedRequester = errors.New("not allowed action for requester user")
)

// NewBoardMemberHandler creates new instance of BoardMemberHandler
func NewBoardMemberHandler(bms *services.BoardMemberService, validator *validation.Validate) *BoardMemberHandler {
	return &BoardMemberHandler{bms, validator}
}

func (bmh *BoardMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	boardIdParam := r.PathValue(config.ParamBoardID)
	if boardIdParam == "" {
		http.Error(w, badBoardIdErr.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	boardId := sqlddl.ID(boardIdParam)
	_, searchErr := bmh.FindBoardByID(ctx, boardId)
	if searchErr != nil {
		http.Error(w, boardNotExistErr.Error(), http.StatusBadRequest)
		return
	}
	memberIdParam := r.PathValue(config.ParamBoardMemberID)
	userId := ctx.Value(contextkeys.KeyUserId).(sqlddl.ID)
	if memberIdParam == "" {
		bmh.handleMultipleMembers(ctx, w, r, userId, boardId)
	} else {
		bmh.handleSingleMember(ctx, w, r, userId, boardId, memberIdParam)
	}
}

func (bmh *BoardMemberHandler) handleMultipleMembers(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	userId sqlddl.ID,
	boardId sqlddl.ID,
) {
	switch r.Method {
	case http.MethodGet:
		_, searchErr := bmh.FindBoardMemberByUserID(ctx, boardId, userId)
		if searchErr != nil {
			http.Error(w, notAllowedRequester.Error(), http.StatusBadRequest)
			return
		}
		members, err := bmh.ListBoardMembers(ctx, boardId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(members)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		if !bmh.BoardMemberService.IsUserAllowedManageBoard(ctx, userId, boardId) {
			http.Error(w, notAllowedRequester.Error(), http.StatusBadRequest)
			return
		}
		var creationData services.CreateBoardMemberData
		decodeErr := json.NewDecoder(r.Body).Decode(&creationData)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		if err := bmh.Validate.Struct(creationData); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.FormatValidationErr(err))
			return
		}
		createdMember, creationErr := bmh.CreateBoardMember(ctx, boardId, &services.CreateBoardMemberData{
			UserId: creationData.UserId,
			Role:   creationData.Role,
		})
		if creationErr != nil {
			http.Error(w, creationErr.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		encodeErr := json.NewEncoder(w).Encode(createdMember)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (bmh *BoardMemberHandler) handleSingleMember(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	userId,
	boardId sqlddl.ID,
	memberIdParam string,
) {
	memberId := sqlddl.ID(memberIdParam)
	switch r.Method {
	case http.MethodGet:
		_, searchErr := bmh.FindBoardMemberByUserID(ctx, boardId, userId)
		if searchErr != nil {
			http.Error(w, notAllowedRequester.Error(), http.StatusForbidden)
			return
		}
		member, err := bmh.FindBoardMemberByID(ctx, memberId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(member)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		if !bmh.IsUserAllowedManageBoard(ctx, userId, boardId) {
			http.Error(w, notAllowedRequester.Error(), http.StatusForbidden)
			return
		}
		var updateData services.UpdateBoardMemberData
		decodeErr := json.NewDecoder(r.Body).Decode(&updateData)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		if validationErr := bmh.Validate.Struct(updateData); validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.FormatValidationErr(validationErr))
			return
		}
		updatedMember, roleChangeErr := bmh.ChangeBoardMemberRole(ctx, memberId, updateData.Role)
		if roleChangeErr != nil {
			http.Error(w, roleChangeErr.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(updatedMember)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		if !bmh.IsUserAllowedDeleteMember(ctx, userId, memberId) {
			http.Error(w, notAllowedRequester.Error(), http.StatusForbidden)
			return
		}
		removeErr := bmh.RemoveBoardMember(ctx, memberId)
		if removeErr != nil {
			http.Error(w, removeErr.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
