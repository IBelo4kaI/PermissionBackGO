package grpcpb

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/token"
)

// Servers реализует оба gRPC-сервиса (PermissionService и UserService) —
// как два Python-класса PermissionGrpcService/UserGrpcService на одной БД.
type Servers struct {
	UnimplementedPermissionServiceServer
	UnimplementedUserServiceServer

	queries repo.Querier
}

func NewServers(queries repo.Querier) *Servers {
	return &Servers{queries: queries}
}

// ValidatePermission — аналог PermissionGrpcService.ValidatePermission:
// проверка сессии и разрешения "service:entity:action". Для доступа к
// чужому пользователю (передан user_id, отличный от session.user_id)
// запрашивается разрешение с wildcard-сегментом entity: entity + ".all".
func (s *Servers) ValidatePermission(ctx context.Context, req *PermissionRequest) (*PermissionResponse, error) {
	session, err := s.sessionByToken(ctx, req.GetSessionToken())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &PermissionResponse{IsAccess: false, Message: "Сессия не действительная", Code: 401}, nil
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	entity := req.GetEntity()
	if req.GetUserId() != "" && session.UserID != req.GetUserId() {
		entity += ".all"
	}

	exist, err := s.queries.ExistsUserPermission(ctx, repo.ExistsUserPermissionParams{
		UserID:  session.UserID,
		Service: req.GetService(),
		Entity:  entity,
		Action:  req.GetAction(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if !exist {
		return &PermissionResponse{
			IsAccess: false,
			Message:  "Нет разрешения",
			Code:     403,
			UserId:   session.UserID,
		}, nil
	}

	return &PermissionResponse{
		IsAccess: true,
		Message:  "Разрешение получено",
		Code:     200,
		UserId:   session.UserID,
	}, nil
}

// GetUsers — аналог UserGrpcService.GetUsers: сессия, права, список пользователей.
func (s *Servers) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error) {
	permReq := req.GetPermissionRequest()

	session, err := s.sessionByToken(ctx, permReq.GetSessionToken())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &GetUsersResponse{Message: strPtr("Ошибка авторизации"), Code: 401}, nil
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	exist, err := s.queries.ExistsUserPermission(ctx, repo.ExistsUserPermissionParams{
		UserID:  session.UserID,
		Service: permReq.GetService(),
		Entity:  permReq.GetEntity(),
		Action:  permReq.GetAction(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	if !exist {
		return &GetUsersResponse{Message: strPtr("Нет доступа"), Code: 403}, nil
	}

	rows, err := s.queries.ListAllUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	users := make([]*UserResponse, 0, len(rows))
	for _, u := range rows {
		users = append(users, &UserResponse{Id: u.ID, Name: u.Name, Surname: u.Surname})
	}

	return &GetUsersResponse{Users: users, Message: strPtr("Разрешение получено"), Code: 200}, nil
}

func (s *Servers) sessionByToken(ctx context.Context, rawToken string) (repo.GetSessionByTokenHashRow, error) {
	// токен хэшируется SHA-256, как в Python-версии (SessionRepository.get_by_token).
	return s.queries.GetSessionByTokenHash(ctx, token.Hash(rawToken))
}

func strPtr(s string) *string {
	return &s
}
