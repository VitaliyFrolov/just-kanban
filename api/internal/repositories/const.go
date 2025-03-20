package repositories

import "just-kanban/pkg/sqlddl"

const (
	ColumnName        = "name"
	ColumnDescription = "description"
	ColumnUserID      = "user_id"
	ColumnBoardID     = "board_id"
	ColumnRole        = "role"
	ColumnToken       = "token"
	ColumnEmail       = "email"
	ColumnPassword    = "password"
	ColumnAvatar      = "avatar"
	ColumnUsername    = "username"
	ColumnFirstName   = "first_name"
	ColumnsLastName   = "last_name"
	ColumnStatus      = "status"
	ColumnOrder       = `"order"`
	ColumnAssigneeID  = "assignee_id"
	ColumnCreatorID   = "creator_id"
)

const (
	TableUsers         = "users"
	TableBoards        = "boards"
	TableBoardMembers  = "board_members"
	TableRefreshTokens = "refresh_tokens"
	TableTasks         = "tasks"
)

// Tables defines structure of generating migration script files
var Tables = []sqlddl.SchemaTable{
	{
		Name: TableUsers,
		Columns: []sqlddl.SchemaColumn{
			{
				Name:        ColumnEmail,
				Type:        sqlddl.TypeText,
				Constraints: []string{sqlddl.ConstraintNotNull, sqlddl.ConstraintUnique},
			},
			{
				Name:        ColumnUsername,
				Type:        sqlddl.TypeVarchar(30),
				Constraints: []string{sqlddl.ConstraintNotNull, sqlddl.ConstraintUnique},
			},
			{
				Name:        ColumnPassword,
				Type:        sqlddl.TypeVarchar(70),
				Constraints: []string{sqlddl.ConstraintNotNull},
			},
			{
				Name:        ColumnFirstName,
				Type:        sqlddl.TypeVarchar(50),
				Constraints: []string{sqlddl.ConstraintNotNull},
			},
			{
				Name:        ColumnsLastName,
				Type:        sqlddl.TypeVarchar(50),
				Constraints: []string{sqlddl.ConstraintNotNull},
			},
			{
				Name: ColumnAvatar,
				Type: sqlddl.TypeText,
			},
		},
	},
	{
		Name: TableBoards,
		Columns: []sqlddl.SchemaColumn{
			{
				Name:        ColumnName,
				Type:        sqlddl.TypeVarchar(150),
				Constraints: []string{sqlddl.ConstraintNotNull},
			},
			{
				Name: ColumnDescription,
				Type: sqlddl.TypeText,
			},
		},
	},
	{
		Name: TableBoardMembers,
		Columns: []sqlddl.SchemaColumn{
			{
				Name:        ColumnRole,
				Type:        sqlddl.TypeVarchar(100),
				Constraints: []string{sqlddl.ConstraintNotNull},
			},
		},
		ForeignKeys: []sqlddl.SchemaForeignKey{
			{
				ColumnName:      ColumnBoardID,
				ReferenceTable:  TableBoards,
				ReferenceColumn: sqlddl.ColumnID,
				OnDelete:        sqlddl.ConstraintOnDeleteCascade,
			},
			{
				ColumnName:      ColumnUserID,
				ReferenceTable:  TableUsers,
				ReferenceColumn: sqlddl.ColumnID,
				OnDelete:        sqlddl.ConstraintOnDeleteCascade,
			},
		},
	},
	{
		Name: TableRefreshTokens,
		Columns: []sqlddl.SchemaColumn{
			{
				Name:        ColumnToken,
				Type:        sqlddl.TypeText,
				Constraints: []string{sqlddl.ConstraintNotNull},
			},
		},
		ForeignKeys: []sqlddl.SchemaForeignKey{
			{
				ColumnName:      ColumnUserID,
				ReferenceTable:  TableUsers,
				ReferenceColumn: sqlddl.ColumnID,
				OnDelete:        sqlddl.ConstraintOnDeleteCascade,
			},
		},
	},
	{
		Name: TableTasks,
		Columns: []sqlddl.SchemaColumn{
			{
				Name:        ColumnName,
				Type:        sqlddl.TypeVarchar(100),
				Constraints: []string{sqlddl.ConstraintNotNull},
			},
			{
				Name:        ColumnStatus,
				Type:        sqlddl.TypeInt,
				Constraints: []string{sqlddl.ConstraintNotNull},
			},
			{
				Name:        ColumnOrder,
				Type:        sqlddl.TypeInt,
				Constraints: []string{sqlddl.ConstraintNotNull},
			},
			{
				Name: ColumnDescription,
				Type: sqlddl.TypeText,
			},
		},
		ForeignKeys: []sqlddl.SchemaForeignKey{
			{
				ColumnName:      ColumnCreatorID,
				ReferenceTable:  TableUsers,
				ReferenceColumn: sqlddl.ColumnID,
				OnDelete:        sqlddl.ConstraintOnDeleteCascade,
			},
			{
				ColumnName:      ColumnAssigneeID,
				ReferenceTable:  TableUsers,
				ReferenceColumn: sqlddl.ColumnID,
				OnDelete:        sqlddl.ConstraintOnDeleteCascade,
			},
			{
				ColumnName:      ColumnBoardID,
				ReferenceTable:  TableBoards,
				ReferenceColumn: sqlddl.ColumnID,
				OnDelete:        sqlddl.ConstraintOnDeleteCascade,
			},
		},
	},
}
