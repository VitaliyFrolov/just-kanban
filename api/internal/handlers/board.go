package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"just-kanban/internal/access"
	"just-kanban/internal/config"
	"just-kanban/internal/contextkeys"
	"just-kanban/internal/services"
	"just-kanban/pkg/sqlddl"
	"just-kanban/pkg/validation"
)

// BoardHandler accepts http requests for working with methods of services.BoardService
type BoardHandler struct {
	*services.TaskService
	*services.BoardService
	*services.BoardMemberService
	*validation.Validate
}

// NewBoardHandler creates new instance of BoardHandler
func NewBoardHandler(
	ts *services.TaskService,
	bs *services.BoardService,
	bms *services.BoardMemberService,
	validate *validation.Validate,
) *BoardHandler {
	return &BoardHandler{ts, bs, bms, validate}
}

func (bh *BoardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	boardIdParam := r.PathValue(config.ParamBoardID)
	ctx := r.Context()
	userId, _ := contextkeys.GetUserId(ctx)
	if boardIdParam == "" {
		bh.handleMultipleBoards(ctx, w, r, userId)
	} else {
		bh.handleSingleBoard(ctx, w, r, userId, boardIdParam)
	}
}

func (bh *BoardHandler) handleMultipleBoards(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	userId sqlddl.ID,
) {
	switch r.Method {
	case http.MethodGet:
		boards, fetchErr := bh.FindAllBoards(ctx)
		if fetchErr != nil {
			http.Error(w, fetchErr.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(boards)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		var boardData services.CreateBoardData
		if err := json.NewDecoder(r.Body).Decode(&boardData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if validateErr := bh.Validate.Struct(&boardData); validateErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.FormatValidationErr(validateErr))
			return
		}
		board, creationErr := bh.CreateBoard(ctx, &boardData)
		if creationErr != nil {
			http.Error(w, creationErr.Error(), http.StatusInternalServerError)
			return
		}
		_, ownerSetErr := bh.CreateBoardMember(ctx, board.ID, &services.CreateBoardMemberData{
			UserId: userId,
			Role:   access.RoleOwner,
		})
		if ownerSetErr != nil {
			bh.DeleteBoard(ctx, board.ID)
			http.Error(w, ownerSetErr.Error(), http.StatusInternalServerError)
			return
		}
		updatedBoard, _ := bh.FindBoardByID(ctx, board.ID)
		w.WriteHeader(http.StatusCreated)
		encodeErr := json.NewEncoder(w).Encode(updatedBoard)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (bh *BoardHandler) handleSingleBoard(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	userId sqlddl.ID,
	boardIdParam string,
) {
	boardId := sqlddl.ID(boardIdParam)
	findBoard, boardSearchErr := bh.FindBoardByID(ctx, boardId)
	if boardSearchErr != nil {
		http.Error(w, boardNotExistErr.Error(), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		_, searchErr := bh.FindBoardMemberByUserID(ctx, boardId, userId)
		if searchErr != nil {
			http.Error(w, notAllowedRequester.Error(), http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(findBoard)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		if !bh.IsUserBoardOwner(ctx, userId, boardId) {
			http.Error(w, notAllowedRequester.Error(), http.StatusForbidden)
			return
		}
		deleteErr := bh.DeleteBoard(ctx, boardId)
		if deleteErr != nil {
			http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodPatch:
		if !bh.IsUserAllowedManageBoard(ctx, userId, boardId) {
			http.Error(w, notAllowedRequester.Error(), http.StatusForbidden)
			return
		}
		var updateData services.UpdateBoardData
		decodeErr := json.NewDecoder(r.Body).Decode(&updateData)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		if validationErr := bh.Validate.Struct(updateData); validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.FormatValidationErr(validationErr))
			return
		}
		updatedBoard, updateErr := bh.UpdateBoard(ctx, boardId, &updateData)
		if updateErr != nil {
			http.Error(w, updateErr.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(updatedBoard)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
