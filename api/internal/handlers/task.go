package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"just-kanban/internal/config"
	"just-kanban/internal/services"
	"just-kanban/pkg/sqlddl"
	"just-kanban/pkg/validation"
)

// TaskHandler handles http requests for working with methods of services.TaskService
type TaskHandler struct {
	*services.TaskService
	*validation.Validate
}

// NewTaskHandler creates new instance of TaskHandler
func NewTaskHandler(
	ts *services.TaskService,
	validate *validation.Validate,
) *TaskHandler {
	return &TaskHandler{ts, validate}
}

func (th *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	boardIdParam := r.PathValue(config.ParamBoardID)
	taskOrderParam := r.PathValue(config.ParamTaskOrder)
	ctx := r.Context()
	boardId := sqlddl.ID(boardIdParam)
	if taskOrderParam == "" {
		th.handleMultipleTasks(ctx, w, r, boardId)
	} else {
		th.handleSingleTask(ctx, w, r, boardId, taskOrderParam)
	}
}

func (th *TaskHandler) handleMultipleTasks(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	boardId sqlddl.ID,
) {
	switch r.Method {
	case http.MethodPost:
		var createData services.CreateTaskData
		if decodeErr := json.NewDecoder(r.Body).Decode(&createData); decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		createData.BoardID = boardId
		if validationErr := th.Validate.Struct(createData); validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.FormatValidationErr(validationErr))
			return
		}
		createdTask, creationErr := th.TaskService.CreateTask(ctx, &createData)
		if creationErr != nil {
			http.Error(w, creationErr.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		encodeErr := json.NewEncoder(w).Encode(createdTask)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusBadRequest)
			return
		}
	case http.MethodGet:
		tasks, searchErr := th.TaskService.FindAllByBoardId(ctx, boardId)
		if searchErr != nil {
			http.Error(w, searchErr.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(tasks)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (th *TaskHandler) handleSingleTask(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	boardId sqlddl.ID,
	taskOrderParam string,
) {
	order, parseErr := strconv.Atoi(taskOrderParam)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusBadRequest)
		return
	}
	task, searchErr := th.TaskService.FindByOrder(ctx, boardId, uint(order))
	if searchErr != nil {
		http.Error(w, searchErr.Error(), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(task)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		var updateData services.UpdateTaskData
		if decodeErr := json.NewDecoder(r.Body).Decode(&updateData); decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		if validationErr := th.Validate.Struct(updateData); validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.FormatValidationErr(validationErr))
			return
		}
		updatedTask, updateErr := th.TaskService.UpdateTask(ctx, task.ID, &updateData)
		if updateErr != nil {
			http.Error(w, updateErr.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(updatedTask)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		deleteErr := th.TaskService.Delete(ctx, task.ID)
		if deleteErr != nil {
			http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
