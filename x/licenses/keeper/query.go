package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) LicenseType(goCtx context.Context, req *types.QueryLicenseTypeRequest) (*types.QueryLicenseTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	lt, found, err := k.GetLicenseType(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "license type %s not found", req.Id)
	}
	return &types.QueryLicenseTypeResponse{LicenseType: lt}, nil
}

func (k Keeper) LicenseTypes(goCtx context.Context, req *types.QueryLicenseTypesRequest) (*types.QueryLicenseTypesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	lts, err := k.GetAllLicenseTypes(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	limit, offset := paginationParams(req.Pagination)
	paginated, nextKey := paginateLicenseTypes(lts, offset, limit)

	return &types.QueryLicenseTypesResponse{
		LicenseTypes: paginated,
		Pagination:   &query.PageResponse{NextKey: nextKey, Total: uint64(len(lts))},
	}, nil
}

func (k Keeper) License(goCtx context.Context, req *types.QueryLicenseRequest) (*types.QueryLicenseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	l, found, err := k.GetLicense(ctx, req.TypeId, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "license (type=%s, id=%d) not found", req.TypeId, req.Id)
	}
	return &types.QueryLicenseResponse{License: l}, nil
}

func (k Keeper) LicensesByType(goCtx context.Context, req *types.QueryLicensesByTypeRequest) (*types.QueryLicensesByTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	licenses, err := k.GetLicensesByType(ctx, req.TypeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	limit, offset := paginationParams(req.Pagination)
	paginated, nextKey := paginateLicenses(licenses, offset, limit)

	return &types.QueryLicensesByTypeResponse{
		Licenses:   paginated,
		Pagination: &query.PageResponse{NextKey: nextKey, Total: uint64(len(licenses))},
	}, nil
}

func (k Keeper) LicensesByHolder(goCtx context.Context, req *types.QueryLicensesByHolderRequest) (*types.QueryLicensesByHolderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	licenses, err := k.GetLicensesByHolder(ctx, req.Holder)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	limit, offset := paginationParams(req.Pagination)
	paginated, nextKey := paginateLicenses(licenses, offset, limit)

	return &types.QueryLicensesByHolderResponse{
		Licenses:   paginated,
		Pagination: &query.PageResponse{NextKey: nextKey, Total: uint64(len(licenses))},
	}, nil
}

func (k Keeper) LicensesByHolderAndType(goCtx context.Context, req *types.QueryLicensesByHolderAndTypeRequest) (*types.QueryLicensesByHolderAndTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	licenses, err := k.GetLicensesByHolderAndType(ctx, req.Holder, req.TypeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	limit, offset := paginationParams(req.Pagination)
	paginated, nextKey := paginateLicenses(licenses, offset, limit)

	return &types.QueryLicensesByHolderAndTypeResponse{
		Licenses:   paginated,
		Pagination: &query.PageResponse{NextKey: nextKey, Total: uint64(len(licenses))},
	}, nil
}

func (k Keeper) AdminKey(goCtx context.Context, req *types.QueryAdminKeyRequest) (*types.QueryAdminKeyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	ak, found, err := k.GetAdminKey(ctx, req.Address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "admin key for address %s not found", req.Address)
	}
	return &types.QueryAdminKeyResponse{AdminKey: ak}, nil
}

func (k Keeper) AdminKeys(goCtx context.Context, req *types.QueryAdminKeysRequest) (*types.QueryAdminKeysResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	aks, err := k.GetAllAdminKeys(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	limit, offset := paginationParams(req.Pagination)
	paginated, nextKey := paginateAdminKeys(aks, offset, limit)

	return &types.QueryAdminKeysResponse{
		AdminKeys:  paginated,
		Pagination: &query.PageResponse{NextKey: nextKey, Total: uint64(len(aks))},
	}, nil
}

func (k Keeper) AdminKeysByLicenseType(goCtx context.Context, req *types.QueryAdminKeysByLicenseTypeRequest) (*types.QueryAdminKeysByLicenseTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.LicenseTypeId == "" {
		return nil, status.Error(codes.InvalidArgument, "license_type_id is required")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	allKeys, err := k.GetAllAdminKeys(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Filter admin keys that have a grant matching the license type (and optionally the permission)
	var matched []types.AdminKey
	for _, ak := range allKeys {
		for _, grant := range ak.Grants {
			if req.Permission != "" && grant.Permission != req.Permission {
				continue
			}
			for _, lt := range grant.LicenseTypes {
				if lt == req.LicenseTypeId {
					matched = append(matched, ak)
					goto nextKey
				}
			}
		}
	nextKey:
	}

	limit, offset := paginationParams(req.Pagination)
	paginated, nextPageKey := paginateAdminKeys(matched, offset, limit)

	return &types.QueryAdminKeysByLicenseTypeResponse{
		AdminKeys:  paginated,
		Pagination: &query.PageResponse{NextKey: nextPageKey, Total: uint64(len(matched))},
	}, nil
}

// paginationParams extracts limit and offset from a PageRequest (defaults: limit=100, offset=0).
func paginationParams(p *query.PageRequest) (limit, offset uint64) {
	limit = 100
	offset = 0
	if p != nil {
		if p.Limit != 0 {
			limit = p.Limit
		}
		if p.Offset != 0 {
			offset = p.Offset
		}
	}
	return
}

func paginateLicenses(all []types.License, offset, limit uint64) ([]types.License, []byte) {
	total := uint64(len(all))
	if offset >= total {
		return []types.License{}, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	var nextKey []byte
	if end < total {
		nextKey = []byte{1}
	}
	return all[offset:end], nextKey
}

func paginateLicenseTypes(all []types.LicenseType, offset, limit uint64) ([]types.LicenseType, []byte) {
	total := uint64(len(all))
	if offset >= total {
		return []types.LicenseType{}, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	var nextKey []byte
	if end < total {
		nextKey = []byte{1}
	}
	return all[offset:end], nextKey
}

func paginateAdminKeys(all []types.AdminKey, offset, limit uint64) ([]types.AdminKey, []byte) {
	total := uint64(len(all))
	if offset >= total {
		return []types.AdminKey{}, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	var nextKey []byte
	if end < total {
		nextKey = []byte{1}
	}
	return all[offset:end], nextKey
}
